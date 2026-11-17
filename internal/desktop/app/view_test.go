package app

import (
	"errors"
	"image/png"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/generator"
	"github.com/sivaprasadreddy/progen/internal/desktop/store"
)

// newTestWizard is a wizard whose templates and defaults are in a temporary
// directory, generating with generate and running work at once.
func newTestWizard(t *testing.T, generate func(config.GeneratorConfiguration) (string, error)) *Wizard {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "Progen Desktop")
	userData := func() (string, error) { return dir, nil }
	if generate == nil {
		generate = func(config.GeneratorConfiguration) (string, error) { return "/tmp/demo.zip", nil }
	}
	prepare := func(config.GeneratorConfiguration) (PreparedProject, error) {
		return nil, errors.New("no preview in this test")
	}
	return NewWizard(store.NewTemplateService(userData), store.NewDefaultsService(userData), generate, prepare, func(work func() func()) { work()() })
}

// TestSnapshots saves a picture of every screen in DESKTOP_SNAPSHOT_DIR, to
// compare the look by eye. It is skipped without it.
func TestSnapshots(t *testing.T) {
	dir := os.Getenv("DESKTOP_SNAPSHOT_DIR")
	if dir == "" {
		t.Skip("DESKTOP_SNAPSHOT_DIR not set")
	}
	w := newTestWizard(t, nil)
	w.prepare = func(c config.GeneratorConfiguration) (PreparedProject, error) {
		pom := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<project xmlns=\"http://maven.apache.org/POM/4.0.0\">\n    <modelVersion>4.0.0</modelVersion>\n" +
			"    <groupId>" + c.Project.GroupID + "</groupId>\n    <artifactId>" + c.Project.ArtifactID + "</artifactId>\n    <!-- Spring Boot -->\n</project>\n"
		return &fakeProject{name: c.Project.AppName, files: []generator.File{
			{Path: "pom.xml", Content: pom, Size: int64(len(pom))},
			{Path: "src/main/resources/application.properties", Content: "spring.application.name=demo\n", Size: 29},
			{Path: "README.md", Content: "# demo\n", Size: 7},
			{Path: ".gitignore", Content: "target/\n", Size: 8},
			{Path: ".mvn/wrapper/maven-wrapper.properties", Content: "distributionUrl=x\n", Size: 18},
			{Path: "compose.yml", Content: "services:\n", Size: 10},
		}}, nil
	}
	tt := ui.NewTester(w.View, 1200, 800)
	save := func(name string) {
		t.Helper()
		f, err := os.Create(filepath.Join(dir, name+".png"))
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := png.Encode(f, tt.Image()); err != nil {
			t.Fatal(err)
		}
	}
	click := func(label string) {
		t.Helper()
		if err := tt.Click(label); err != nil {
			t.Fatalf("click %q: %v", label, err)
		}
	}
	save("1-welcome")
	click("Get started")
	save("2-project")
	click("Next")
	click("Thymeleaf")
	save("3-frontend")
	click("Next")
	save("4-database")
	click("Next")
	click("Spring Cloud AWS")
	save("5-features")
	click("Next")
	save("6-review")
	click("Preview")
	save("6-preview")
	click("Back to Review")
	click("Cancel")
	save("9-cancel-alert")
	tt.Key(0, ui.KeyEscape)
	click("Templates")
	save("7-templates")
	click("Defaults")
	save("8-defaults")
}

// viewTester runs w's view in a window of the app's default size.
func viewTester(w *Wizard) *ui.Tester { return ui.NewTester(w.View, 1200, 800) }

func click(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	if err := tt.Click(label); err != nil {
		t.Fatalf("click %q: %v\ntexts: %q", label, err, tt.Texts())
	}
}

// replace types value into the text field labeled label, over its text.
func replace(t *testing.T, tt *ui.Tester, label, value string) {
	t.Helper()
	field, ok := tt.AccessByRole(ui.RoleTextField, label)
	if !ok {
		t.Fatalf("no text field %q", label)
	}
	if !field.Focused {
		tt.ClickAt(field.Bounds.X+field.Bounds.W/2, field.Bounds.Y+field.Bounds.H/2)
	}
	if !tt.Focused(label) {
		t.Fatalf("text field %q did not take the focus", label)
	}
	tt.Key(ui.Cmd, ui.KeyA)
	if value == "" {
		tt.Key(0, ui.KeyBackspace)
		return
	}
	tt.Type(value)
}

func TestProjectStepStartsFromDefaultsAndDerivesThePackage(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Get started")
	if p := w.selection.Project; p != projectFromDefaults(store.BuiltInDefaults) {
		t.Fatalf("Project starts at %+v", p)
	}
	replace(t, tt, "Group ID", "com.acme")
	replace(t, tt, "Artifact ID", "orders")
	if got := w.selection.Project.PackageName; got != "com.acme.orders" {
		t.Errorf("package follows the IDs: got %q", got)
	}
	replace(t, tt, "Package Name", "com.acme.api")
	replace(t, tt, "Group ID", "com.shop")
	if got := w.selection.Project.PackageName; got != "com.acme.api" {
		t.Errorf("a customized package stays: got %q", got)
	}
}

func TestProjectStepNeedsEveryField(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Get started")
	replace(t, tt, "Application Name", "")
	click(t, tt, "Next")
	if w.step != stepProject || !tt.HasText("Required.") {
		t.Fatalf("Next went on without a name: step %d, texts %q", w.step, tt.Texts())
	}
	replace(t, tt, "Application Name", "orders")
	click(t, tt, "Next")
	if w.step != stepFrontend || tt.HasText("Required.") {
		t.Errorf("Next with a name: step %d", w.step)
	}
}

func TestHTMXGoesWithServerRenderedFrontends(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Get started")
	click(t, tt, "Next")
	if tt.HasText("HTMX") {
		t.Error("HTMX offered without a server-rendered frontend")
	}
	click(t, tt, "Thymeleaf")
	click(t, tt, "HTMX")
	if !w.selection.Frontend.HTMX {
		t.Fatal("HTMX not chosen")
	}
	click(t, tt, "React")
	if w.selection.Frontend.HTMX || tt.HasText("HTMX") {
		t.Errorf("HTMX kept with React: %+v", w.selection.Frontend)
	}
}

func TestGeneratesTheSelection(t *testing.T) {
	var generated config.GeneratorConfiguration
	w := newTestWizard(t, func(c config.GeneratorConfiguration) (string, error) {
		generated = c
		return "/Users/me/Downloads/orders.zip", nil
	})
	tt := viewTester(w)
	click(t, tt, "Get started")
	replace(t, tt, "Application Name", "orders")
	click(t, tt, "Gradle")
	click(t, tt, "Next")
	click(t, tt, "JTE")
	click(t, tt, "HTMX")
	click(t, tt, "Next")
	click(t, tt, "MySQL")
	click(t, tt, "jOOQ")
	click(t, tt, "Liquibase")
	click(t, tt, "Next")
	click(t, tt, "Email")
	click(t, tt, "Anthropic")
	click(t, tt, "Ollama")
	click(t, tt, "Next")
	if w.step != stepReview || !tt.HasText("Email") {
		t.Fatalf("not on Review: step %d", w.step)
	}
	click(t, tt, "Generate")
	want := config.GeneratorConfiguration{
		Project:  config.ProjectSettings{AppName: "orders", AppVersion: "1.0.0", GroupID: "dev.sivalabs", ArtifactID: "demo", PackageName: "dev.sivalabs.demo", BuildTool: config.BuildToolGradle, TaskRunner: config.TaskRunnerMise},
		Frontend: config.FrontendSettings{Type: config.FrontendJTE, HTMX: true},
		Database: config.DatabaseSettings{Engine: config.DatabaseMySQL, MigrationTool: config.MigrationLiquibase, PersistenceFramework: config.PersistenceJOOQ},
		Features: []config.Feature{config.FeatureEmail},
		AIModels: []config.AIModel{config.AIModelOllama, config.AIModelAnthropic}, // canonical order
	}
	if generated.Project != want.Project || generated.Frontend != want.Frontend || generated.Database != want.Database ||
		!slices.Equal(generated.Features, want.Features) || !slices.Equal(generated.AIModels, want.AIModels) {
		t.Errorf("generated %+v\nwant      %+v", generated, want)
	}
	if w.step != stepComplete || !tt.HasText("/Users/me/Downloads/orders.zip") || !tt.HasText("Your Spring Boot Gradle project was generated and saved as a ZIP archive.") {
		t.Fatalf("Complete step not shown: %q", tt.Texts())
	}
	click(t, tt, "Start New")
	if w.step != stepWelcome || w.selection.Project.AppName != "demo" {
		t.Errorf("Start New: step %d, project %+v", w.step, w.selection.Project)
	}
}

func TestShowsAGenerationError(t *testing.T) {
	w := newTestWizard(t, func(config.GeneratorConfiguration) (string, error) {
		return "", errors.New("generate project: exit status 1")
	})
	tt := viewTester(w)
	click(t, tt, "Get started")
	for range 4 {
		click(t, tt, "Next")
	}
	click(t, tt, "Generate")
	if w.step != stepReview || !tt.HasText("generate project: exit status 1") {
		t.Errorf("error not shown: step %d, texts %q", w.step, tt.Texts())
	}
}

func TestTemplatesOpenAtReviewAndSaveTheEditedSelection(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Thymeleaf + HTMX")
	if w.step != stepReview || !tt.HasText("Thymeleaf + HTMX") || w.activeTemplate == nil {
		t.Fatalf("template not opened at Review: step %d", w.step)
	}
	replace(t, tt, "Name", "orders")
	click(t, tt, "Save as template")
	replace(t, tt, "Template Name", "Orders")
	click(t, tt, "Save as new")
	if !tt.HasText("Saved template “Orders”.") {
		t.Fatalf("not saved: %q", tt.Texts())
	}
	list, _ := w.templates.List()
	names := []string{}
	for _, tpl := range list {
		names = append(names, tpl.Name)
	}
	if !slices.Equal(names, []string{"REST API", "Thymeleaf + HTMX", "Angular", "Orders"}) {
		t.Fatalf("templates %q", names)
	}
	if got := list[3].Configuration; got.Project.AppName != "orders" || !got.Frontend.HTMX || got.Frontend.Type != config.FrontendThymeleaf {
		t.Errorf("saved %+v", got)
	}
	if list[1].Configuration.Project.AppName != "demo" {
		t.Error("the starter template changed")
	}

	// Saved as Orders, the Review step now updates Orders.
	replace(t, tt, "Version", "2.0.0")
	click(t, tt, "Update “Orders”")
	list, _ = w.templates.List()
	if !tt.HasText("Updated template “Orders”.") || len(list) != 4 || list[3].Configuration.Project.AppVersion != "2.0.0" {
		t.Errorf("not updated in place: %+v", list)
	}
}

func TestUpdateOpensATemplateAtTheProjectStep(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Templates")
	if !tt.HasText("3 TEMPLATES") {
		t.Fatalf("Templates page: %q", tt.Texts())
	}
	if err := tt.Click("Update"); err != nil {
		t.Fatal(err)
	}
	if w.page != pageNone || w.step != stepProject || w.activeTemplate == nil || w.activeTemplate.Name != "REST API" {
		t.Errorf("Update: page %d, step %d, template %v", w.page, w.step, w.activeTemplate)
	}
}

func TestDeleteAsksFirst(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Delete")
	if !tt.HasText("Delete the template “REST API”?") {
		t.Fatalf("no confirmation: %q", tt.Texts())
	}
	click(t, tt, "Cancel")
	if len(w.templateList) != 3 || tt.HasText("Delete the template “REST API”?") {
		t.Fatalf("Cancel deleted or kept the alert: %d templates", len(w.templateList))
	}
	click(t, tt, "Delete")
	// The alert's Delete is the last match: the list's buttons come first.
	tt.ClickAt(alertButton(t, tt, "Delete"))
	if len(w.templateList) != 2 || w.templateList[0].Name != "Thymeleaf + HTMX" {
		t.Errorf("not deleted: %+v", w.templateList)
	}
}

func TestCancelAsksBeforeDiscardingTheSelection(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Angular")
	replace(t, tt, "Name", "orders")
	click(t, tt, "Cancel")
	if !tt.HasText("Discard this project’s selections?") {
		t.Fatalf("no confirmation: %q", tt.Texts())
	}
	tt.ClickAt(alertButton(t, tt, "Cancel"))
	if w.step != stepReview || w.selection.Project.AppName != "orders" {
		t.Fatalf("declining changed the selection: step %d, %+v", w.step, w.selection.Project)
	}
	click(t, tt, "Cancel")
	click(t, tt, "Discard")
	if w.step != stepWelcome || w.activeTemplate != nil || w.selection.Frontend.Type != config.FrontendNone || w.selection.Project.AppName != "demo" {
		t.Errorf("not discarded: step %d, %+v", w.step, w.selection)
	}
}

func TestDefaultsPrefillNewProjects(t *testing.T) {
	w := newTestWizard(t, nil)
	tt := viewTester(w)
	click(t, tt, "Defaults")
	replace(t, tt, "Group ID", "com acme")
	click(t, tt, "Save defaults")
	if !tt.HasText(`invalid group ID "com acme": use letters, digits, '.', '-' and '_'`) {
		t.Fatalf("invalid group ID accepted: %q", tt.Texts())
	}
	replace(t, tt, "Group ID", "com.acme")
	replace(t, tt, "Application Version", "0.1.0")
	click(t, tt, "Gradle")
	click(t, tt, "Taskfile")
	click(t, tt, "Save defaults")
	if !tt.HasText("Saved. New projects start with these values.") {
		t.Fatalf("not saved: %q", tt.Texts())
	}
	click(t, tt, "Back to Welcome")
	click(t, tt, "Get started")
	want := config.ProjectSettings{AppName: "demo", AppVersion: "0.1.0", GroupID: "com.acme", ArtifactID: "demo", PackageName: "com.acme.demo", BuildTool: config.BuildToolGradle, TaskRunner: config.TaskRunnerTaskfile}
	if w.selection.Project != want {
		t.Errorf("Project %+v, want %+v", w.selection.Project, want)
	}

	// A Project step the user changed keeps its values when the defaults change.
	replace(t, tt, "Application Name", "orders")
	click(t, tt, "Defaults")
	replace(t, tt, "Application Version", "2.0.0")
	click(t, tt, "Save defaults")
	click(t, tt, "Back to Project")
	if p := w.selection.Project; p.AppName != "orders" || p.AppVersion != "0.1.0" {
		t.Errorf("edited Project changed: %+v", p)
	}

	// A new wizard reads the saved defaults.
	again := NewWizard(w.templates, w.defaults, w.generate, w.prepare, w.async)
	if again.selection.Project.AppVersion != "2.0.0" {
		t.Errorf("defaults not kept: %+v", again.selection.Project)
	}
}

// alertButton returns the center of the last element labeled label, the
// alert's button, which shows over the window's own.
func alertButton(t *testing.T, tt *ui.Tester, label string) (float32, float32) {
	t.Helper()
	var found *ui.AccessInfo
	for _, a := range tt.AccessAll() {
		if a.Role == ui.RoleButton && a.Name == label {
			found = &a
		}
	}
	if found == nil {
		t.Fatalf("no button %q", label)
	}
	b := found.Bounds
	return b.X + b.W/2, b.Y + b.H/2
}
