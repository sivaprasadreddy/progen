package app

import (
	"errors"
	"strings"
	"testing"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/generator"
)

// fakeProject is a prepared project of a few files, which counts its saves.
type fakeProject struct {
	name      string
	files     []generator.File
	saves     int
	saveErr   error
	discarded bool
}

func (p *fakeProject) Name() string            { return p.name }
func (p *fakeProject) Files() []generator.File { return p.files }
func (p *fakeProject) Discard()                { p.discarded = true }
func (p *fakeProject) Save() (string, error) {
	p.saves++
	if p.saveErr != nil {
		return "", p.saveErr
	}
	return "/Users/me/Downloads/" + p.name + ".zip", nil
}

// previewFixture is a wizard on the Review step whose previews are fake
// projects, and whose generations are counted.
type previewFixture struct {
	w        *Wizard
	tt       *ui.Tester
	prepared []*fakeProject
	// prepareErr fails the next preparations.
	prepareErr  error
	generations int
	// pending holds the work async was given, run by finish, when hold is
	// set.
	hold    bool
	pending []func() func()
}

func newPreviewFixture(t *testing.T) *previewFixture {
	t.Helper()
	f := &previewFixture{}
	f.w = newTestWizard(t, func(config.GeneratorConfiguration) (string, error) {
		f.generations++
		return "/Users/me/Downloads/generated.zip", nil
	})
	f.w.prepare = func(c config.GeneratorConfiguration) (PreparedProject, error) {
		if f.prepareErr != nil {
			return nil, f.prepareErr
		}
		buildFile := generator.File{Path: "pom.xml", Content: "<project>\n  <artifactId>demo</artifactId>\n</project>\n"}
		if c.Project.BuildTool == config.BuildToolGradle {
			buildFile = generator.File{Path: "build.gradle", Content: "plugins {\n    id 'java'\n}\n"}
		}
		buildFile.Size = int64(len(buildFile.Content))
		p := &fakeProject{name: strings.TrimSpace(c.Project.AppName), files: []generator.File{
			{Path: "README.md", Content: "# Demo\n", Size: 7},
			{Path: "src/main/java/dev/sivalabs/shop/Application.java", Content: "package dev.sivalabs.shop;\n", Size: 27},
			buildFile,
			{Path: ".gitignore", Content: "target/\n", Size: 8},
			{Path: ".mvn/wrapper/maven-wrapper.jar", Binary: true, Size: 62547},
		}}
		f.prepared = append(f.prepared, p)
		return p, nil
	}
	f.w.async = func(work func() func()) {
		if f.hold {
			f.pending = append(f.pending, work)
			return
		}
		work()()
	}
	// Tall enough for the Review step with an error under it.
	f.tt = ui.NewTester(f.w.View, 1200, 1000)
	click(t, f.tt, "Get started")
	for range 4 {
		click(t, f.tt, "Next")
	}
	if f.w.step != stepReview {
		t.Fatalf("not on Review: step %d", f.w.step)
	}
	return f
}

// finish runs the work async held, and applies it.
func (f *previewFixture) finish() {
	for _, work := range f.pending {
		work()()
	}
	f.pending = nil
	f.hold = false
	f.tt.Frame()
}

func (f *previewFixture) last() *fakeProject { return f.prepared[len(f.prepared)-1] }

// treeItem is what assistive technology sees of the tree's item label,
// and whether it shows.
func treeItem(tt *ui.Tester, label string) (ui.AccessInfo, bool) {
	for _, a := range tt.AccessAll() {
		if a.Role == ui.RoleTreeItem && a.Name == label {
			return a, true
		}
	}
	return ui.AccessInfo{}, false
}

func clickTreeItem(t *testing.T, tt *ui.Tester, label string) {
	t.Helper()
	item, ok := treeItem(tt, label)
	if !ok {
		t.Fatalf("no tree item %q", label)
	}
	tt.ClickAt(item.Bounds.X+item.Bounds.W/2, item.Bounds.Y+item.Bounds.H/2)
}

// shownCode is the text of the file the preview shows, "" for none.
func shownCode(w *Wizard) string {
	if w.preview == nil || w.preview.code == nil {
		return ""
	}
	return w.preview.code.Buffer.String()
}

func TestPreviewShowsTheProjectTreeWithTheBuildFileChosen(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	if !f.w.showsPreview() || !f.tt.HasText("Back to Review") || len(f.prepared) != 1 {
		t.Fatalf("preview page not shown: %q", f.tt.Texts())
	}

	// The project's folder is open; the folders in it are closed.
	for _, label := range []string{"demo", ".mvn", "src", ".gitignore", "pom.xml", "README.md"} {
		if _, ok := treeItem(f.tt, label); !ok {
			t.Errorf("tree item %q not shown", label)
		}
	}
	if _, ok := treeItem(f.tt, "main"); ok {
		t.Error("src shows its children before it is opened")
	}
	if root, _ := treeItem(f.tt, "demo"); !root.Expanded {
		t.Error("the project's folder is closed")
	}

	// The build file is chosen, read-only, with line numbers.
	if pom, _ := treeItem(f.tt, "pom.xml"); !pom.Selected {
		t.Error("pom.xml not selected")
	}
	if !strings.Contains(shownCode(f.w), "<artifactId>demo</artifactId>") || !f.tt.HasText("demo/pom.xml") {
		t.Errorf("pom.xml not shown: %q", shownCode(f.w))
	}
	code := f.w.preview.code
	if !code.ReadOnly || code.LineNumbers != editor.LineNumbersOn || code.Buffer.Language() == nil || code.Buffer.Language().Name != "XML" {
		t.Errorf("pom.xml shown as %+v", code)
	}
	field, ok := f.tt.AccessByRole(ui.RoleTextField, "pom.xml")
	if !ok || !field.ReadOnly || !strings.Contains(field.Value, "<artifactId>demo</artifactId>") {
		t.Errorf("assistive technology sees %+v", field)
	}

	// Folders open and close; choosing a file shows it.
	clickTreeItem(t, f.tt, "src")
	clickTreeItem(t, f.tt, "main")
	clickTreeItem(t, f.tt, "java")
	clickTreeItem(t, f.tt, "dev")
	clickTreeItem(t, f.tt, "sivalabs")
	clickTreeItem(t, f.tt, "shop")
	if _, ok := treeItem(f.tt, "Application.java"); !ok {
		t.Fatalf("src/main/java/dev/sivalabs/shop not opened: %q", f.tt.Texts())
	}
	clickTreeItem(t, f.tt, "shop")
	if _, ok := treeItem(f.tt, "Application.java"); ok {
		t.Fatal("shop did not close")
	}
	clickTreeItem(t, f.tt, "shop")
	clickTreeItem(t, f.tt, "Application.java")
	if shownCode(f.w) != "package dev.sivalabs.shop;\n" || f.w.preview.code.Buffer.Language().Name != "Java" {
		t.Errorf("Application.java not shown: %q", shownCode(f.w))
	}

	// A binary file shows its size, not its contents.
	clickTreeItem(t, f.tt, ".mvn")
	clickTreeItem(t, f.tt, "wrapper")
	clickTreeItem(t, f.tt, "maven-wrapper.jar")
	if !f.tt.HasText("Binary file, preview not available") || !f.tt.HasText("61.1 KB") || f.w.preview.code != nil {
		t.Errorf("binary file: %q", f.tt.Texts())
	}
}

func TestPreviewOfAGradleProjectChoosesBuildGradle(t *testing.T) {
	f := newPreviewFixture(t)
	f.w.selection.Project.BuildTool = config.BuildToolGradle
	click(t, f.tt, "Preview")
	if f.w.preview.selected != "build.gradle" || !strings.Contains(shownCode(f.w), "id 'java'") ||
		f.w.preview.code.Buffer.Language().Name != "Gradle" {
		t.Errorf("selected %q: %q", f.w.preview.selected, shownCode(f.w))
	}
}

func TestGenerateFromThePreviewSavesThePreviewedFiles(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	click(t, f.tt, "Generate")
	if f.last().saves != 1 || f.generations != 0 || len(f.prepared) != 1 {
		t.Fatalf("saves %d, generations %d, preparations %d", f.last().saves, f.generations, len(f.prepared))
	}
	if f.w.step != stepComplete || !f.tt.HasText("/Users/me/Downloads/demo.zip") {
		t.Fatalf("Complete not shown: %q", f.tt.Texts())
	}
	click(t, f.tt, "Start New")
	if !f.last().discarded || f.w.preview != nil {
		t.Error("Start New kept the preview")
	}
}

func TestAnUnchangedPreviewIsReused(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	click(t, f.tt, "Back to Review")
	if f.w.step != stepReview || f.w.showsPreview() || !f.tt.HasText("Spring AI Models") {
		t.Fatalf("Back to Review: %q", f.tt.Texts())
	}
	click(t, f.tt, "Preview")
	if len(f.prepared) != 1 || !f.w.showsPreview() {
		t.Fatalf("prepared again: %d", len(f.prepared))
	}
	click(t, f.tt, "Back to Review")

	// Going back and on without changes keeps it too.
	click(t, f.tt, "Back")
	click(t, f.tt, "Next")
	if f.last().discarded {
		t.Fatal("an unchanged preview was discarded")
	}

	// Generate on the Review step saves the prepared project.
	click(t, f.tt, "Generate")
	if f.last().saves != 1 || f.generations != 0 || f.w.step != stepComplete {
		t.Errorf("saves %d, generations %d, step %d", f.last().saves, f.generations, f.w.step)
	}
}

func TestChangingTheSelectionDiscardsThePreview(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	click(t, f.tt, "Back to Review")
	replace(t, f.tt, "Name", "orders")
	if !f.prepared[0].discarded || f.w.preview != nil {
		t.Fatal("editing the name kept the preview")
	}
	click(t, f.tt, "Preview")
	if len(f.prepared) != 2 || !f.tt.HasText("orders") {
		t.Fatalf("not prepared again: %d", len(f.prepared))
	}

	// A change in an earlier step discards it once the user goes on.
	click(t, f.tt, "Back to Review")
	click(t, f.tt, "Back")
	click(t, f.tt, "Email")
	click(t, f.tt, "Next")
	if !f.prepared[1].discarded || f.w.preview != nil {
		t.Fatal("choosing a feature kept the preview")
	}
	click(t, f.tt, "Generate")
	if f.generations != 1 || f.prepared[1].saves != 0 {
		t.Errorf("Generate saved a stale preview: generations %d", f.generations)
	}
}

func TestCancelDiscardsThePreview(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	click(t, f.tt, "Back to Review")
	click(t, f.tt, "Cancel")
	click(t, f.tt, "Discard")
	if f.w.step != stepWelcome || !f.last().discarded || f.w.preview != nil {
		t.Errorf("Cancel kept the preview: step %d", f.w.step)
	}
}

func TestCloseDiscardsThePreview(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	f.w.Close()
	if !f.last().discarded || f.w.preview != nil {
		t.Error("Close kept the preview")
	}
}

func TestActionsWaitWhileThePreviewIsPrepared(t *testing.T) {
	f := newPreviewFixture(t)
	f.hold = true
	click(t, f.tt, "Preview")
	if f.w.showsPreview() || f.w.step != stepReview {
		t.Fatal("left Review before the preview was ready")
	}
	button, ok := f.tt.AccessByRole(ui.RoleButton, "Preparing preview…")
	if !ok || !button.Disabled {
		t.Fatalf("Preview button while preparing: %+v, %v", button, ok)
	}
	for _, label := range []string{"Generate", "Back", "Cancel", "Save as template", "Templates", "Defaults"} {
		if b, ok := f.tt.AccessByRole(ui.RoleButton, label); !ok || !b.Disabled {
			t.Errorf("%s is not disabled while preparing: %+v", label, b)
		}
	}
	if field, ok := f.tt.AccessByRole(ui.RoleTextField, "Name"); !ok || !field.Disabled {
		t.Errorf("Name is editable while preparing: %+v", field)
	}
	_ = f.tt.Click("Generate")
	_ = f.tt.Click("Back")
	if f.generations != 0 || f.w.step != stepReview {
		t.Fatalf("an action ran while preparing: generations %d, step %d", f.generations, f.w.step)
	}

	f.finish()
	if !f.w.showsPreview() || len(f.prepared) != 1 {
		t.Errorf("preview not shown once ready: %q", f.tt.Texts())
	}
}

func TestAFailedPreviewStaysOnReview(t *testing.T) {
	f := newPreviewFixture(t)
	f.prepareErr = errors.New("generate project: exit status 1")
	click(t, f.tt, "Preview")
	if f.w.showsPreview() || f.w.preview != nil || !f.tt.HasText("generate project: exit status 1") || !f.tt.HasText("Preview") {
		t.Fatalf("failure: %q", f.tt.Texts())
	}

	// The user can try again.
	f.prepareErr = nil
	click(t, f.tt, "Preview")
	if !f.w.showsPreview() || f.tt.HasText("generate project: exit status 1") {
		t.Errorf("second try: %q", f.tt.Texts())
	}
}

func TestPreviewNeedsTheProjectFields(t *testing.T) {
	f := newPreviewFixture(t)
	replace(t, f.tt, "Name", "")
	click(t, f.tt, "Preview")
	if len(f.prepared) != 0 || f.w.showsPreview() || !f.w.showProjectErrors {
		t.Errorf("previewed without a name: %d prepared", len(f.prepared))
	}
	if field, _ := f.tt.AccessByRole(ui.RoleTextField, "Name"); field.Bounds.W == 0 {
		t.Error("Name field missing")
	}
}

func TestAFailedSaveStaysOnThePreview(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	f.last().saveErr = errors.New("create project archive: permission denied")
	click(t, f.tt, "Generate")
	if !f.w.showsPreview() || !f.tt.HasText("create project archive: permission denied") || f.w.preview == nil {
		t.Fatalf("save failure: %q", f.tt.Texts())
	}
	f.last().saveErr = nil
	click(t, f.tt, "Generate")
	if f.w.step != stepComplete || len(f.prepared) != 1 {
		t.Errorf("second save: step %d", f.w.step)
	}
}

func TestPagesOpenedFromThePreviewReturnToIt(t *testing.T) {
	f := newPreviewFixture(t)
	click(t, f.tt, "Preview")
	click(t, f.tt, "Defaults")
	click(t, f.tt, "Back to Preview")
	if !f.w.showsPreview() {
		t.Errorf("Defaults did not return to the preview: %q", f.tt.Texts())
	}
}
