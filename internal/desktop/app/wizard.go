// Package app is Progen Desktop's window: the wizard's state and actions, and the
// native view MyGo draws from them.
package app

import (
	"fmt"
	"slices"
	"strings"

	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/store"
)

// The wizard's steps, in order. stepWelcome and stepComplete bracket the
// steps that make a selection.
const (
	stepWelcome = iota
	stepProject
	stepFrontend
	stepDatabase
	stepFeatures
	stepReview
	stepComplete
)

var stepNames = []string{"Welcome", "Project", "Frontend", "Database", "Features", "Review", "Complete"}

// page is a page the sidebar opens in place of the current step, which it
// returns to.
type page int

const (
	pageNone page = iota
	pageDefaults
	pageTemplates
)

// message is a line of feedback under a form, a success or an error.
type message struct {
	text  string
	error bool
}

// projectField names a field of the Project step.
type projectField int

const (
	fieldAppName projectField = iota
	fieldAppVersion
	fieldGroupID
	fieldArtifactID
	fieldPackageName
)

// Wizard is the state of the window: the selection being made, the step it
// is on, and the user's templates and defaults.
type Wizard struct {
	templates TemplateStore
	defaults  DefaultsStore
	// generate generates a project and returns the archive's path and the
	// generation's warnings.
	generate func(config.GeneratorConfiguration) (string, []string, error)
	// prepare generates a project to preview before it is saved.
	prepare func(config.GeneratorConfiguration) (PreparedProject, error)
	// async runs work off the UI thread, then apply on it: in a goroutine
	// and Window.Update in the app, at once in tests.
	async func(work func() (apply func()))

	step int
	page page

	// selection is what the steps chose so far.
	selection config.GeneratorConfiguration
	// projectConfirmed is whether the user went past the Project step, or
	// loaded a template, after which new defaults leave the project alone.
	projectConfirmed bool
	// showProjectErrors marks the empty Project fields once the user tried
	// to go on with them.
	showProjectErrors bool

	downloadedPath  string
	generationError string
	generating      bool
	// generationWarnings are the last generation's warnings, as a failed
	// code formatting, shown on the Complete step.
	generationWarnings []string

	// preview is the project last prepared for the preview page, kept
	// while the selection it was prepared from stays the same; preparing
	// is whether one is being prepared, and previewing whether the page
	// shows in place of the Review step.
	preview    *preview
	preparing  bool
	previewing bool

	templateList   []store.Template
	templatesError string
	// activeTemplate is the template the selection was loaded from or last
	// saved as.
	activeTemplate *store.Template
	templatePanel  bool
	templateName   string
	templateMsg    *message

	projectDefaults store.ProjectDefaults
	defaultsDraft   store.ProjectDefaults
	defaultsMsg     *message

	// Confirmations shown with an alert: discarding the selection, and
	// deleting deleteTarget.
	confirmingCancel bool
	confirmingDelete bool
	deleteTarget     *store.Template
}

// TemplateStore is what the wizard needs of store.TemplateService.
type TemplateStore interface {
	List() ([]store.Template, error)
	Save(store.Template) (store.Template, error)
	Delete(id string) error
}

// DefaultsStore is what the wizard needs of store.DefaultsService.
type DefaultsStore interface {
	Get() (store.ProjectDefaults, error)
	Save(store.ProjectDefaults) (store.ProjectDefaults, error)
}

// NewWizard returns a wizard on the Welcome step. generate generates a
// project and returns the archive's path and warnings; prepare generates one to preview
// and save; async runs work off the UI thread and then the func it returns
// on it.
func NewWizard(templates TemplateStore, defaults DefaultsStore, generate func(config.GeneratorConfiguration) (string, []string, error),
	prepare func(config.GeneratorConfiguration) (PreparedProject, error), async func(func() func())) *Wizard {
	w := &Wizard{templates: templates, defaults: defaults, generate: generate, prepare: prepare, async: async}
	w.projectDefaults = store.BuiltInDefaults
	if saved, err := defaults.Get(); err == nil {
		w.projectDefaults = saved
	} else {
		w.defaultsMsg = &message{text: "Could not load defaults: " + err.Error(), error: true}
	}
	w.resetSelection()
	w.refreshTemplates()
	return w
}

// newSelection is a fresh selection: the Project step from the defaults,
// and the first choice of every other step.
func newSelection(d store.ProjectDefaults) config.GeneratorConfiguration {
	return config.GeneratorConfiguration{
		Project:  projectFromDefaults(d),
		Frontend: config.FrontendSettings{Type: config.FrontendNone},
		Database: config.DatabaseSettings{Engine: config.DatabasePostgreSQL, MigrationTool: config.MigrationFlyway, PersistenceFramework: config.PersistenceSpringDataJPA},
		Features: []config.Feature{},
		AIModels: []config.AIModel{},
	}
}

// projectFromDefaults is the Project step's starting values for a new
// application.
func projectFromDefaults(d store.ProjectDefaults) config.ProjectSettings {
	return config.ProjectSettings{
		AppName:     "demo",
		AppVersion:  d.AppVersion,
		GroupID:     d.GroupID,
		ArtifactID:  "demo",
		PackageName: deriveBasePackage(d.GroupID, "demo"),
		BuildTool:   d.BuildTool,
		TaskRunner:  d.TaskRunner,
	}
}

// deriveBasePackage is the package name the group and artifact IDs give.
func deriveBasePackage(groupID, artifactID string) string {
	group := strings.TrimSpace(groupID)
	segment := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(artifactID), "-", ""))
	switch {
	case group == "":
		return segment
	case segment == "":
		return group
	}
	return group + "." + segment
}

// projectValue returns a pointer to one of p's text fields.
func projectValue(p *config.ProjectSettings, f projectField) *string {
	switch f {
	case fieldAppName:
		return &p.AppName
	case fieldAppVersion:
		return &p.AppVersion
	case fieldGroupID:
		return &p.GroupID
	case fieldArtifactID:
		return &p.ArtifactID
	}
	return &p.PackageName
}

// setProjectField changes a Project field. The package name follows the
// group and artifact IDs until the user customizes it.
func (w *Wizard) setProjectField(f projectField, value string) {
	p := &w.selection.Project
	derived := p.PackageName == deriveBasePackage(p.GroupID, p.ArtifactID)
	*projectValue(p, f) = value
	if derived && (f == fieldGroupID || f == fieldArtifactID) {
		p.PackageName = deriveBasePackage(p.GroupID, p.ArtifactID)
	}
	w.generationError = ""
	w.dropStalePreview()
}

// projectComplete reports whether every Project text field is filled in.
func (w *Wizard) projectComplete() bool {
	for f := fieldAppName; f <= fieldPackageName; f++ {
		if strings.TrimSpace(*projectValue(&w.selection.Project, f)) == "" {
			return false
		}
	}
	return true
}

// setFrontend chooses a frontend. HTMX stays only with a server-rendered one.
func (w *Wizard) setFrontend(t config.FrontendType) {
	w.selection.Frontend = config.FrontendSettings{Type: t, HTMX: config.IsServerRendered(t) && w.selection.Frontend.HTMX}
	w.dropStalePreview()
}

func (w *Wizard) toggleFeature(f config.Feature) {
	if i := slices.Index(w.selection.Features, f); i >= 0 {
		w.selection.Features = slices.Delete(slices.Clone(w.selection.Features), i, i+1)
	} else {
		w.selection.Features = append(slices.Clone(w.selection.Features), f)
	}
	w.dropStalePreview()
}

// toggleAIModel keeps the models in their canonical order.
func (w *Wizard) toggleAIModel(m config.AIModel) {
	chosen := slices.Contains(w.selection.AIModels, m)
	models := []config.AIModel{}
	for _, o := range aiModelOptions {
		if (o.value == m && !chosen) || (o.value != m && slices.Contains(w.selection.AIModels, o.value)) {
			models = append(models, o.value)
		}
	}
	w.selection.AIModels = models
	w.dropStalePreview()
}

// next goes on from the current step. The Project step goes on only with
// its fields filled in. A preview of another selection is discarded.
func (w *Wizard) next() {
	w.dropStalePreview()
	if w.step == stepProject {
		if !w.projectComplete() {
			w.showProjectErrors = true
			return
		}
		w.showProjectErrors = false
		w.projectConfirmed = true
	}
	if w.step < stepReview {
		w.step++
	}
}

func (w *Wizard) back() {
	w.dropStalePreview()
	if w.step == stepReview {
		// A message about the selection goes stale once it can change.
		w.templatePanel = false
		w.templateMsg = nil
	}
	if w.step > stepWelcome {
		w.step--
	}
}

// getStarted begins a project. Without a template it starts from the
// defaults, not from a template looked at earlier.
func (w *Wizard) getStarted() {
	if w.activeTemplate != nil {
		w.resetSelection()
	}
	w.step = stepProject
}

// resetSelection drops every choice, the loaded template and the preview.
func (w *Wizard) resetSelection() {
	w.discardPreview()
	w.selection = newSelection(w.projectDefaults)
	w.projectConfirmed = false
	w.showProjectErrors = false
	w.activeTemplate = nil
	w.templateName = ""
	w.templatePanel = false
	w.templateMsg = nil
}

// cancel is Cancel's action once the user confirmed it.
func (w *Wizard) cancel() {
	w.resetSelection()
	w.generationError = ""
	w.step = stepWelcome
}

func (w *Wizard) startNew() {
	w.resetSelection()
	w.downloadedPath = ""
	w.generationError = ""
	w.generationWarnings = nil
	w.generating = false
	w.step = stepWelcome
}

// applyTemplate loads a template's selection into every step, then shows
// step.
func (w *Wizard) applyTemplate(t store.Template, step int) {
	w.selection = cloneConfiguration(t.Configuration)
	if w.selection.Features == nil {
		w.selection.Features = []config.Feature{}
	}
	ordered := []config.AIModel{}
	for _, o := range aiModelOptions {
		if slices.Contains(t.Configuration.AIModels, o.value) {
			ordered = append(ordered, o.value)
		}
	}
	w.selection.AIModels = ordered
	template := t
	w.activeTemplate = &template
	w.templateName = t.Name
	w.templatePanel = false
	w.templateMsg = nil
	w.generationError = ""
	w.projectConfirmed = true
	w.showProjectErrors = false
	w.page = pageNone
	w.step = step
	w.dropStalePreview()
}

func cloneConfiguration(c config.GeneratorConfiguration) config.GeneratorConfiguration {
	c.Features = slices.Clone(c.Features)
	c.AIModels = slices.Clone(c.AIModels)
	return c
}

// configuration is the selection as the generator and templates take it.
func (w *Wizard) configuration() config.GeneratorConfiguration {
	c := cloneConfiguration(w.selection)
	c.Frontend.HTMX = config.IsServerRendered(c.Frontend.Type) && c.Frontend.HTMX
	return c
}

func (w *Wizard) refreshTemplates() {
	list, err := w.templates.List()
	if err != nil {
		w.templatesError = "Could not load templates: " + err.Error()
		return
	}
	w.templateList = list
	w.templatesError = ""
}

// saveTemplate saves the selection as a new template, or over the active
// one when update is true.
func (w *Wizard) saveTemplate(update bool) {
	if !w.projectComplete() {
		w.showProjectErrors = true
		return
	}
	id := ""
	if update && w.activeTemplate != nil {
		id = w.activeTemplate.ID
	}
	saved, err := w.templates.Save(store.Template{ID: id, Name: w.templateName, Configuration: w.configuration()})
	if err != nil {
		w.templateMsg = &message{text: err.Error(), error: true}
		return
	}
	w.activeTemplate = &saved
	w.templateName = saved.Name
	verb := "Saved"
	if update {
		verb = "Updated"
	}
	w.templateMsg = &message{text: fmt.Sprintf("%s template “%s”.", verb, saved.Name)}
	w.refreshTemplates()
}

func (w *Wizard) deleteTemplate(t store.Template) {
	if err := w.templates.Delete(t.ID); err != nil {
		w.templatesError = fmt.Sprintf("Could not delete “%s”: %v", t.Name, err)
		return
	}
	if w.activeTemplate != nil && w.activeTemplate.ID == t.ID {
		w.activeTemplate = nil
	}
	w.refreshTemplates()
}

// startGenerating generates the project off the UI thread and shows the
// Complete step, or the error, when it is done. A preview of the selection
// is saved as it is, without generating it again.
func (w *Wizard) startGenerating() {
	if w.busy() {
		return
	}
	if !w.projectComplete() {
		w.showProjectErrors = true
		return
	}
	w.generating = true
	w.generationError = ""
	generate := w.generate
	if p := w.currentPreview(); p != nil {
		generate = func(config.GeneratorConfiguration) (string, []string, error) {
			path, err := p.project.Save()
			return path, p.project.Warnings(), err
		}
	}
	configuration := w.configuration()
	w.async(func() func() {
		path, warnings, err := generate(configuration)
		return func() {
			w.generating = false
			if err != nil {
				w.generationError = err.Error()
				return
			}
			w.downloadedPath = path
			w.generationWarnings = warnings
			w.previewing = false
			w.step = stepComplete
		}
	})
}

func (w *Wizard) openPage(p page) {
	switch p {
	case pageTemplates:
		w.refreshTemplates()
	case pageDefaults:
		w.defaultsDraft = w.projectDefaults
		w.defaultsMsg = nil
	}
	w.page = p
}

// saveDefaults saves the Defaults page. A Project step that still shows the
// previous defaults follows the new ones; one the user changed, confirmed
// or filled from a template keeps its values.
func (w *Wizard) saveDefaults() {
	saved, err := w.defaults.Save(w.defaultsDraft)
	if err != nil {
		w.defaultsMsg = &message{text: err.Error(), error: true}
		return
	}
	if !w.projectConfirmed && w.activeTemplate == nil && w.selection.Project == projectFromDefaults(w.projectDefaults) {
		w.selection.Project = projectFromDefaults(saved)
	}
	w.projectDefaults = saved
	w.defaultsDraft = saved
	w.defaultsMsg = &message{text: "Saved. New projects start with these values."}
}

// templateSummary is a one-line description of a template's main choices,
// e.g. "Maven · JTE + HTMX · PostgreSQL · Spring Data JPA · 2 features".
func templateSummary(c config.GeneratorConfiguration) string {
	frontend := "REST API"
	if c.Frontend.Type != config.FrontendNone {
		frontend = label(c.Frontend.Type)
		if c.Frontend.HTMX {
			frontend += " + HTMX"
		}
	}
	parts := []string{label(c.Project.BuildTool), frontend, label(c.Database.Engine), label(c.Database.PersistenceFramework)}
	count := func(n int, noun string) {
		switch {
		case n == 1:
			parts = append(parts, fmt.Sprintf("1 %s", noun))
		case n > 1:
			parts = append(parts, fmt.Sprintf("%d %ss", n, noun))
		}
	}
	count(len(c.Features), "feature")
	count(len(c.AIModels), "AI model")
	return strings.Join(parts, " · ")
}
