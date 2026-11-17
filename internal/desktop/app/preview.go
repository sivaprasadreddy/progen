package app

import (
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/generator"
)

// PreparedProject is a generated project waiting to be previewed and saved:
// generator.Prepared in the app.
type PreparedProject interface {
	// Name is the project directory's name.
	Name() string
	// Files are the project's files, their paths in the project directory.
	Files() []generator.File
	// Warnings describe steps that failed without failing generation, as
	// code formatting.
	Warnings() []string
	// Save archives the project in the Downloads directory and returns the
	// archive's path.
	Save() (string, error)
	// Discard removes the project's files.
	Discard()
}

// preview is a prepared project shown on the preview page: its file tree,
// the file chosen in it, and the selection it was prepared from.
type preview struct {
	project PreparedProject
	// configuration is the selection the project was generated from: the
	// preview is the project's only while the selection stays the same.
	configuration config.GeneratorConfiguration
	root          *previewNode
	files         map[string]generator.File
	selected      string
	// code shows the selected file, nil for a binary one.
	code *editor.State
}

// previewNode is a folder or a file of the tree.
type previewNode struct {
	name, path string
	// open is whether a folder shows its children; nil for a file.
	open     *bool
	children []*previewNode
}

// newPreview shows project's files: the project's folder open, the folders
// in it closed, and the build file chosen.
func newPreview(project PreparedProject, configuration config.GeneratorConfiguration) *preview {
	open := true
	p := &preview{
		project:       project,
		configuration: configuration,
		root:          &previewNode{name: project.Name(), open: &open},
		files:         map[string]generator.File{},
	}
	for _, f := range project.Files() {
		p.files[f.Path] = f
		p.root.add(f.Path, strings.Split(f.Path, "/"))
	}
	p.root.sort()

	buildFile := "pom.xml"
	if configuration.Project.BuildTool == config.BuildToolGradle {
		buildFile = "build.gradle"
	}
	if _, ok := p.files[buildFile]; !ok {
		buildFile = p.root.firstFile()
	}
	p.selectFile(buildFile)
	return p
}

// add adds the file at filePath, whose segments below n are segments.
func (n *previewNode) add(filePath string, segments []string) {
	if len(segments) == 1 {
		n.children = append(n.children, &previewNode{name: segments[0], path: filePath})
		return
	}
	folder := path.Join(n.path, segments[0])
	var child *previewNode
	for _, c := range n.children {
		if c.open != nil && c.name == segments[0] {
			child = c
		}
	}
	if child == nil {
		child = &previewNode{name: segments[0], path: folder, open: new(bool)}
		n.children = append(n.children, child)
	}
	child.add(filePath, segments[1:])
}

// sort orders each folder's children: folders, then files, by name.
func (n *previewNode) sort() {
	slices.SortFunc(n.children, func(a, b *previewNode) int {
		if (a.open == nil) != (b.open == nil) {
			if a.open != nil {
				return -1
			}
			return 1
		}
		return strings.Compare(strings.ToLower(a.name), strings.ToLower(b.name))
	})
	for _, c := range n.children {
		c.sort()
	}
}

// firstFile is the path of the first file the tree shows, "" for none.
func (n *previewNode) firstFile() string {
	for _, c := range n.children {
		if c.open == nil {
			return c.path
		}
	}
	for _, c := range n.children {
		if f := c.firstFile(); f != "" {
			return f
		}
	}
	return ""
}

// selectFile shows the file at filePath, read-only and highlighted.
func (p *preview) selectFile(filePath string) {
	p.selected = filePath
	p.code = nil
	f, ok := p.files[filePath]
	if !ok || f.Binary {
		return
	}
	buffer := editor.NewBuffer(f.Content)
	buffer.SetLanguage(languageFor(path.Base(filePath), f.Content))
	p.code = &editor.State{Buffer: buffer, ReadOnly: true, Label: filePath, NoScrollPastEnd: true}
}

// sameSelection reports whether two selections generate the same project:
// the features in any order.
func sameSelection(a, b config.GeneratorConfiguration) bool {
	features := func(c config.GeneratorConfiguration) []config.Feature {
		f := slices.Clone(c.Features)
		slices.Sort(f)
		return f
	}
	return a.Project == b.Project && a.Frontend == b.Frontend && a.Database == b.Database &&
		slices.Equal(features(a), features(b)) && slices.Equal(a.AIModels, b.AIModels)
}

// currentPreview is the prepared preview when it is still the selection's,
// else nil.
func (w *Wizard) currentPreview() *preview {
	if w.preview != nil && sameSelection(w.preview.configuration, w.configuration()) {
		return w.preview
	}
	return nil
}

// discardPreview drops the prepared project and removes its files.
func (w *Wizard) discardPreview() {
	if w.preview != nil {
		w.preview.project.Discard()
		w.preview = nil
	}
	w.previewing = false
}

// dropStalePreview discards the prepared project once the selection
// changed from the one it was generated from.
func (w *Wizard) dropStalePreview() {
	if w.preview != nil && w.currentPreview() == nil {
		w.discardPreview()
	}
}

// Close discards what the wizard keeps on disk: the prepared preview. The
// app calls it as it quits.
func (w *Wizard) Close() {
	w.discardPreview()
}

// busy is whether a project is being generated, saved or prepared, while
// the wizard's actions wait.
func (w *Wizard) busy() bool { return w.generating || w.preparing }

// startPreview shows the preview page: at once with the prepared project
// when the selection is unchanged, else once a project is prepared off the
// UI thread. A failure shows on the Review step.
func (w *Wizard) startPreview() {
	if w.busy() {
		return
	}
	if !w.projectComplete() {
		w.showProjectErrors = true
		return
	}
	w.generationError = ""
	if w.currentPreview() != nil {
		w.previewing = true
		return
	}
	w.discardPreview()
	w.preparing = true
	configuration := w.configuration()
	w.async(func() func() {
		project, err := w.prepare(configuration)
		var prepared *preview
		if err == nil {
			prepared = newPreview(project, configuration)
		}
		return func() {
			w.preparing = false
			if err != nil {
				w.generationError = err.Error()
				return
			}
			w.preview = prepared
			w.previewing = true
		}
	})
}

// backToReview leaves the preview page, keeping the prepared project.
func (w *Wizard) backToReview() {
	w.previewing = false
	w.generationError = ""
}

// previewPage shows the prepared project: its file tree beside the chosen
// file, with Back to Review and Generate.
func (w *Wizard) previewPage(c *ui.Context) {
	p := w.preview
	ui.Column(c).Grow(1).Shrink(1).MinHeight(0).Padding(34, 40).Children(func() {
		ui.Row(c).AlignItems(ui.End).Margin(0, 0, 18, 0).Children(func() {
			ui.Column(c).Grow(1).Shrink(1).Children(func() {
				eyebrow(c, "Preview").Margin(0, 0, 8, 0)
				heading(c, p.project.Name(), 28)
				subtitle(c, "Every file of the project, as the ZIP will have it. Generate saves these files.").Margin(8, 0, 0, 0)
			})
			fileCount := len(p.files)
			noun := "files"
			if fileCount == 1 {
				noun = "file"
			}
			ui.Textf(c, "%d %s", fileCount, noun).FontSize(12).FontWeight(700).TextColor(colorMuted)
		})
		ui.Row(c).Grow(1).Shrink(1).MinHeight(0).Gap(14).AlignItems(ui.Stretch).Children(func() {
			previewPanel(c).Width(290).Shrink(0).Children(func() {
				ui.Scroll(c.Key("preview:tree")).Grow(1).Padding(8, 6).Children(func() {
					w.previewTree(c, p)
				})
			})
			previewPanel(c).Grow(1).Shrink(1).MinWidth(0).Children(func() {
				w.previewFile(c, p)
			})
		})
		for _, warning := range p.project.Warnings() {
			warningBox(c, warning)
		}
		if w.generationError != "" {
			errorBox(c, w.generationError)
		}
		formActions(c, nil, func() {
			if secondaryButton(c, "Back to Review", w.busy()).Clicked() {
				w.backToReview()
			}
			w.generateButton(c)
		})
	})
}

// previewPanel is the white panel the tree and the file sit in.
func previewPanel(c *ui.Context) ui.Element {
	return ui.Column(c).Radius(12).Border(1, ui.RGBA(205, 207, 221, 0.8)).Background(ui.RGBA(255, 255, 255, 0.82)).
		Shadow(0, 10, 30, 0, ui.RGBA(57, 57, 85, 0.06)).Clip()
}

func (w *Wizard) previewTree(c *ui.Context, p *preview) {
	var item func(n *previewNode)
	item = func(n *previewNode) {
		if n.open == nil {
			if ui.TreeItem(c.Key("file:"+n.path), n.name, nil, nil).Selected(p.selected == n.path).Clicked() {
				p.selectFile(n.path)
			}
			return
		}
		folder := ui.TreeItem(c.Key("folder:"+n.path), n.name, n.open, func() {
			for _, child := range n.children {
				item(child)
			}
		})
		if folder.Clicked() {
			*n.open = !*n.open
		}
	}
	ui.Tree(c, func() { item(p.root) }).Label("Project files")
}

// previewFile is the chosen file's path over its contents, or over a note
// for a binary file.
func (w *Wizard) previewFile(c *ui.Context, p *preview) {
	ui.Row(c).Padding(11, 16).Gap(10).BorderWidth(0, 0, 1, 0).BorderColor(ui.Hex("#e4e5ec")).Children(func() {
		ui.Text(c, p.project.Name()+"/"+p.selected).Font("monospace").FontSize(12).TextColor(colorLabel).Grow(1).Shrink(1)
		if f, ok := p.files[p.selected]; ok {
			ui.Text(c, formatSize(f.Size)).FontSize(11).TextColor(colorHint).Shrink(0)
		}
	})
	f, ok := p.files[p.selected]
	switch {
	case !ok:
		ui.Text(c, "The project has no files.").FontSize(13).TextColor(colorMuted).Padding(20)
	case f.Binary:
		ui.Column(c).Grow(1).Center().Gap(6).Children(func() {
			ui.Text(c, "Binary file, preview not available").FontSize(14).FontWeight(700).TextColor(colorLabel)
			ui.Text(c, formatSize(f.Size)).FontSize(12).TextColor(colorMuted)
		})
	default:
		editor.View(c.Key("code:"+p.selected), p.code).Grow(1).FillWidth()
	}
}

// formatSize is a file's size for people, as "512 bytes" or "43.8 KB".
func formatSize(size int64) string {
	switch {
	case size == 1:
		return "1 byte"
	case size < 1024:
		return fmt.Sprintf("%d bytes", size)
	case size < 1024*1024:
		return fmt.Sprintf("%.1f KB", float64(size)/1024)
	}
	return fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
}
