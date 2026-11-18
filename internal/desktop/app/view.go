package app

import (
	"fmt"
	"slices"
	"strings"

	"github.com/egoist/mygo/ui"
	"github.com/sivaprasadreddy/progen/internal/desktop/config"
	"github.com/sivaprasadreddy/progen/internal/desktop/store"
)

const (
	authorURL     = "https://sivalabs.in"
	repositoryURL = "https://github.com/sivaprasadreddy/progen"
)

// View builds the window: the sidebar beside the current step or page, and
// the confirmation alerts over them.
func (w *Wizard) View(c *ui.Context) {
	c.SetTheme(theme(c))
	ui.Row(c).Fill().AlignItems(ui.Stretch).Children(func() {
		w.sidebar(c)
		background := ui.LinearGradient{Angle: 120, Stops: []ui.GradientStop{
			{Color: ui.Hex("#f4f3f8"), At: 0}, {Color: ui.Hex("#faf2f7"), At: 0.48}, {Color: ui.Hex("#eef2fb"), At: 1},
		}}
		if w.showsPreview() {
			// The tree and the file scroll on their own, in the whole window.
			ui.Column(c.Key("preview")).Grow(1).Shrink(1).MinWidth(0).AlignItems(ui.Stretch).LinearGradient(background).
				Children(func() { w.previewPage(c) })
			return
		}
		if w.page == pageNone && w.step == stepConfigure {
			// The section links stay over the cards as they scroll.
			ui.Column(c.Key("configure")).Grow(1).Shrink(1).MinWidth(0).AlignItems(ui.Stretch).LinearGradient(background).
				Children(func() { w.configureStep(c) })
			return
		}
		content := ui.Scroll(c.Key("content")).Grow(1).Padding(48).AlignItems(ui.Center)
		content.LinearGradient(background)
		content.Children(func() { w.content(c) })
	})

	// The alerts show in the system's accent, as macOS's own do: in the app's
	// red, their default button, Cancel, would look like the action that
	// discards. Cancel is last, so Enter cancels, and Escape dismisses.
	ui.Themed(c, alertTheme(c)).Absolute().Background(ui.Transparent).Children(func() { w.alerts(c) })
}

func alertTheme(c *ui.Context) *ui.Theme {
	accent := c.Preferences().Accent
	if accent.A == 0 {
		accent = ui.Hex("#007aff")
	}
	return theme(c).WithAccent(accent)
}

func (w *Wizard) alerts(c *ui.Context) {
	switch ui.AlertDialog(c, &w.confirmingCancel, "Discard this project’s selections?",
		"Every choice you made for this project is cleared and you return to Welcome. Your templates and defaults are kept.",
		"Discard", "Cancel") {
	case 0:
		w.cancel()
	}
	if w.deleteTarget != nil {
		target := *w.deleteTarget
		switch ui.AlertDialog(c, &w.confirmingDelete, fmt.Sprintf("Delete the template “%s”?", target.Name),
			"It is removed from your templates. This can’t be undone.", "Delete", "Cancel") {
		case 0:
			w.deleteTemplate(target)
			w.deleteTarget = nil
		case 1, ui.AlertDismissed:
			w.deleteTarget = nil
		}
	}
}

// showsPreview is whether the preview page shows in place of the Review
// step.
func (w *Wizard) showsPreview() bool {
	return w.page == pageNone && w.step == stepReview && w.previewing && w.preview != nil
}

// returnName is the name of what a page's Back returns to.
func (w *Wizard) returnName() string {
	if w.step == stepReview && w.previewing && w.preview != nil {
		return "Preview"
	}
	return stepNames[w.step]
}

func (w *Wizard) content(c *ui.Context) {
	switch {
	case w.page == pageTemplates:
		w.templatesPage(c)
	case w.page == pageDefaults:
		w.defaultsPage(c)
	case w.step == stepWelcome:
		w.welcome(c)
	case w.step == stepReview:
		w.reviewStep(c)
	case w.step == stepComplete:
		w.completeStep(c)
	}
}

// sidebar is the column of steps on the gradient, with the Defaults and
// Templates links and the author, which drags the window.
func (w *Wizard) sidebar(c *ui.Context) {
	side := ui.Column(c).Width(192).Shrink(0).Padding(52, 17, 20, 17).DragWindow()
	side.LinearGradient(ui.LinearGradient{Angle: 180, Stops: []ui.GradientStop{
		{Color: ui.Hex("#b9c3e2"), At: 0}, {Color: ui.Hex("#b9b9e9"), At: 0.35}, {Color: ui.Hex("#efb5d2"), At: 0.7}, {Color: ui.Hex("#d7bdd7"), At: 1},
	}})
	side.Children(func() {
		ui.Column(c).AlignItems(ui.Center).MinHeight(126).Children(func() {
			ui.Image(c, logo).Size(116, 88).Fit(ui.Contain)
			ui.Text(c, "SPRING BOOT GENERATOR").FontSize(8).FontWeight(700).LetterSpacing(0.64).
				TextColor(ui.RGBA(28, 30, 44, 0.56)).Margin(-3, 0, 0, 0)
		})
		ui.Column(c).Gap(3).Margin(18, 0, 0, 0).Children(func() {
			shown := w.step
			if w.page != pageNone {
				shown = -1
			}
			for i, name := range stepNames {
				complete := i < w.step || (i == w.step && i == stepComplete && w.downloadedPath != "")
				w.stepItem(c, i, name, complete, !complete && i == shown)
			}
		})
		ui.Box(c).Height(1).Background(ui.RGBA(28, 30, 44, 0.12)).Margin(14, 0, 10, 0)
		w.sidebarLink(c, pageDefaults, slidersIcon, "Defaults", 0)
		w.sidebarLink(c, pageTemplates, listIcon, "Templates", len(w.templateList)).Margin(3, 0, 0, 0)
		ui.Spacer(c)
		ui.Column(c).Gap(6).Padding(14, 8, 0, 8).BorderWidth(1, 0, 0, 0).BorderColor(ui.RGBA(28, 30, 44, 0.12)).Children(func() {
			ui.Row(c).Gap(3).Children(func() {
				ui.Text(c, "Author:").FontSize(11).TextColor(ui.Hex("#586276"))
				footerLink(c, "SivaLabs", authorURL)
			})
			footerLink(c, "GitHub repository ↗", repositoryURL)
		})
	})
}

func (w *Wizard) stepItem(c *ui.Context, index int, name string, complete, active bool) {
	row := ui.Row(c.Key("step:"+name)).Height(50).Gap(13).Padding(0, 14).Radius(9)
	text, numberText, numberFill := ui.Hex("#445064"), ui.Hex("#485164"), ui.RGBA(222, 228, 236, 0.86)
	switch {
	case active:
		row.Background(ui.RGBA(255, 255, 255, 0.95)).Shadow(0, 2, 8, 0, ui.RGBA(55, 57, 88, 0.05))
		text, numberText, numberFill = ui.Hex("#151626"), white, colorEyebrow
	case complete:
		text, numberText, numberFill = colorSuccess, white, ui.Hex("#18ba66")
	}
	row.Children(func() {
		ui.Box(c).Size(28, 28).Shrink(0).Radius(14).Background(numberFill).Center().Children(func() {
			if complete {
				ui.Icon(c, tickIcon).Size(14, 14).TextColor(numberText)
			} else {
				ui.Textf(c, "%d", index+1).FontSize(13).TextColor(numberText)
			}
		})
		ui.Text(c, name).FontSize(15).TextColor(text)
	})
}

// sidebarLink opens a page in place of the step, or closes it again.
func (w *Wizard) sidebarLink(c *ui.Context, p page, icon *ui.SVG, name string, count int) ui.Element {
	active := w.page == p
	link := ui.ButtonBase(c.Key("link:"+name)).Label(name).Disabled(w.busy()).Height(50).Gap(13).
		Padding(0, 6, 0, 14).Radius(9).Justify(ui.Start).Cursor(ui.CursorPointer)
	iconText, iconFill, text := ui.Hex("#485164"), ui.RGBA(222, 228, 236, 0.86), ui.Hex("#445064")
	switch {
	case active:
		link.Background(ui.RGBA(255, 255, 255, 0.95)).Shadow(0, 2, 8, 0, ui.RGBA(55, 57, 88, 0.05))
		iconText, iconFill, text = white, colorEyebrow, ui.Hex("#151626")
	case link.Hovered():
		link.Background(ui.RGBA(255, 255, 255, 0.5))
	}
	if link.Clicked() {
		if active {
			w.page = pageNone
		} else {
			w.openPage(p)
		}
	}
	link.Children(func() {
		ui.Box(c).Size(28, 28).Shrink(0).Radius(14).Background(iconFill).Center().Children(func() {
			ui.Icon(c, icon).Size(14, 14).TextColor(iconText)
		})
		ui.Text(c, name).FontSize(15).TextColor(text)
		if count > 0 {
			ui.Textf(c, "%d", count).FontSize(10).FontWeight(700).TextColor(ui.Hex("#485164")).Margin(0, 0, 0, -9).
				MinWidth(18).Padding(1, 5).Radius(999).Background(ui.RGBA(222, 228, 236, 0.86)).TextAlign(ui.Center)
		}
	})
	return link
}

// footerLink opens a URL in the default browser.
func footerLink(c *ui.Context, label, url string) {
	link := ui.Link(c, label, url).FontSize(11).FontWeight(700)
	if link.Hovered() {
		link.TextColor(colorAccent)
	} else {
		link.TextColor(colorLabel)
	}
}

// welcome is the first screen: the app's pitch, Get started, and the
// templates to start from.
func (w *Wizard) welcome(c *ui.Context) {
	width, _ := c.Size()
	ui.Column(c).WidthPercent(78).MaxWidth(560).AlignItems(ui.Center).Padding(38, 0, 0, 0).Children(func() {
		successIcon(c)
		eyebrow(c, "Welcome").Margin(0, 0, 10, 0)
		heading(c, "Build production-ready\nSpring Boot applications", clamp(width*0.031, 27, 42)).TextAlign(ui.Center)
		ui.Text(c, "Create a thoughtfully structured Spring Boot project with sensible defaults and the dependencies your application needs—all through one simple, guided setup.").
			FontSize(15).LineHeight(1.65).TextColor(colorMuted).TextAlign(ui.Center).MaxWidth(525).Margin(16, 0, 29, 0)
		if outlineButton(c, "Get started").Clicked() {
			w.getStarted()
		}
		if w.templatesError != "" {
			errorBox(c, w.templatesError).FillWidth().Margin(24, 0, 0, 0)
		}
		if len(w.templateList) > 0 {
			ui.Column(c).FillWidth().Margin(40, 0, 0, 0).Children(func() {
				ui.Text(c, "Or start from a template").FontSize(13).FontWeight(700).TextColor(colorLabel).
					TextAlign(ui.Center).FillWidth().Margin(0, 0, 12, 0)
				w.templateItems(c)
			})
		}
	})
}

func successIcon(c *ui.Context) {
	ui.Box(c).Size(65, 65).Radius(16).Background(ui.RGBA(83, 203, 139, 0.22)).Center().Margin(0, 0, 22, 0).Children(func() {
		ui.Icon(c, checkIcon).Size(37, 37).TextColor(ui.Hex("#19bd68"))
	})
}

func clamp(v, low, high float32) float32 {
	return max(low, min(v, high))
}

// templateItems lists the templates, each opening at Review, with Update
// and Delete.
func (w *Wizard) templateItems(c *ui.Context) {
	ui.Column(c).Gap(10).Children(func() {
		for _, t := range w.templateList {
			row := ui.Row(c.Key(t.ID)).Gap(12).Padding(6, 12, 6, 6).Radius(12)
			if row.Hovered() {
				row.Border(1, colorSelected).Background(ui.Hex("#fff7f7"))
			} else {
				row.Border(1, colorBorder).Background(ui.RGBA(255, 255, 255, 0.62))
			}
			row.Children(func() {
				open := ui.ButtonBase(c.Key("open")).Label(t.Name).Grow(1).Shrink(1).Column().AlignItems(ui.Start).Gap(4).
					Padding(10, 12).Radius(8).Cursor(ui.CursorPointer).Tooltip(fmt.Sprintf("Review “%s”", t.Name))
				open.OnClick(func() { w.applyTemplate(t, stepReview) })
				open.Children(func() {
					ui.Text(c, t.Name).FontSize(14).FontWeight(700).TextColor(colorTextStrong)
					ui.Text(c, templateSummary(t.Configuration)).FontSize(12).LineHeight(1.4).TextColor(ui.Hex("#6b6f82"))
				})
				ui.Row(c).Gap(6).Shrink(0).Children(func() {
					smallButton(c, "Update").Tooltip("Change this template’s selections on the Configure step").
						OnClick(func() { w.applyTemplate(t, stepConfigure) })
					smallButton(c, "Delete").OnClick(func() {
						target := t
						w.deleteTarget = &target
						w.confirmingDelete = true
					})
				})
			})
		}
	})
}

// cancelButton asks before it discards the selection.
func (w *Wizard) cancelButton(c *ui.Context) {
	b := secondaryButton(c, "Cancel", w.busy()).Tooltip("Clear all selections and return to Welcome")
	if b.Clicked() {
		w.confirmingCancel = true
	}
}

// sectionBarHeight is the height of the bar of section links over the
// Configure step.
const sectionBarHeight = 60

// configureStep is every choice of a project in one scrolling page, a card
// for each section, under a bar of links that scroll to them and mark the
// one in view.
func (w *Wizard) configureStep(c *ui.Context) {
	// targets are where the links scroll to, from where the previous frame
	// laid the cards out: each card's top in the page, which puts it just
	// under the bar.
	var targets [sectionCount]float32
	var view float32
	scroll := ui.Scroll(c.Key("configure-scroll")).TrackScroll(&w.configureScroll).Grow(1).Shrink(1).
		Padding(sectionBarHeight+28, 48, 48, 48).ScrollbarInsets(sectionBarHeight, 0, 0).AlignItems(ui.Center)
	scroll.Children(func() {
		page := ui.Column(c).FillWidth().MaxWidth(760).Gap(18)
		page.Children(func() {
			ui.Column(c).Padding(0, 4, 8, 4).Children(func() {
				eyebrow(c, "Step 2 of 4").Margin(0, 0, 8, 0)
				heading(c, "Configure Your Project", 30)
				subtitle(c, "Choose everything your Spring Boot application needs on one page. Jump to a section with the links above.").Margin(8, 0, 0, 0)
			})
			cards := [sectionCount]func(*ui.Context){w.projectSection, w.frontendSection, w.databaseSection, w.featuresSection}
			for s, card := range cards {
				e := formCard(c.Key("section:"+sectionNames[s]), func() { card(c) })
				if s > 0 {
					targets[s] = e.Bounds().Y - page.Bounds().Y
				}
			}
			formActions(c, func() { w.cancelButton(c) }, func() {
				if secondaryButton(c, "Back", false).Clicked() {
					w.back()
				}
				if primaryButton(c, "Next", "→", false).Clicked() {
					w.next()
				}
			})
		})
		view = scroll.Bounds().H
	})

	// The section in view is the last whose top passed a third of the way
	// down the page under the bar.
	in := sectionProject
	line := w.configureScroll.Y + max(view-sectionBarHeight, 0)/3
	for s := sectionFrontend; s < sectionCount; s++ {
		if targets[s] > 0 && targets[s] <= line {
			in = s
		}
	}
	bar := ui.Row(c).Absolute().Top(0).Left(0).Right(0).Height(sectionBarHeight).Padding(0, 48).Justify(ui.Center).
		AlignItems(ui.Center).Background(ui.RGBA(248, 246, 251, 0.97)).BorderWidth(0, 0, 1, 0).BorderColor(ui.RGBA(205, 207, 221, 0.7))
	bar.Role(ui.RoleNavigation).Label("Sections").Children(func() {
		ui.Row(c).FillWidth().MaxWidth(760).Gap(6).AlignItems(ui.Center).Children(func() {
			for s := range sectionCount {
				w.sectionLink(c, s, s == in, targets[s])
			}
		})
	})
}

// sectionLink scrolls the Configure step to a section's card.
func (w *Wizard) sectionLink(c *ui.Context, s section, active bool, target float32) {
	name := sectionNames[s]
	link := ui.ButtonBase(c.Key("section-link:"+name)).Label(name).Padding(8, 16).Radius(999).Cursor(ui.CursorPointer).
		Tooltip("Go to " + name)
	text := colorButtonText
	current := ui.CurrentNone
	switch {
	case active:
		current = ui.CurrentLocation
		link.Background(white).Border(1, ui.RGBA(239, 74, 84, 0.45)).Shadow(0, 2, 8, 0, ui.RGBA(55, 57, 88, 0.06))
		text = colorAccent
	case link.Hovered():
		link.Background(ui.RGBA(255, 255, 255, 0.7)).Border(1, colorButtonEdge)
	default:
		link.Border(1, ui.Transparent)
	}
	link.Current(current).OnClick(func() { w.configureScroll.ScrollTo(0, target, ui.Smooth) })
	link.Children(func() {
		ui.Text(c, name).FontSize(13).FontWeight(600).TextColor(text)
	})
}

// sectionHeading is a Configure card's title and description.
func sectionHeading(c *ui.Context, title, description string) {
	ui.Column(c).Margin(0, 0, 22, 0).Gap(6).Children(func() {
		heading(c, title, 21)
		subtitle(c, description)
	})
}

func (w *Wizard) projectSection(c *ui.Context) {
	p := w.selection.Project
	sectionHeading(c, "Project Metadata", "Define the coordinates, build system, and task runner for your new Spring Boot application.")
	ui.Grid(c).Columns(2).Gap(20).Children(func() {
		w.projectField(c, fieldAppName, "Application Name", "The generated project directory name.", "boot-demo")
		w.projectField(c, fieldAppVersion, "Application Version", "The initial version of your application.", "1.0.0")
		w.projectField(c, fieldGroupID, "Group ID", "Usually your organization’s reversed domain.", "com.example")
		w.projectField(c, fieldArtifactID, "Artifact ID", "The Maven/Gradle artifact name.", "my-application")
		w.projectField(c, fieldPackageName, "Package Name", "The base package for generated source code.", "com.example.demo").ColumnSpan(-1)
	})
	optionGroup(c, "Build Tool", false, func() {
		for _, o := range buildToolOptions {
			radioCard(c.Key(o.value), p.BuildTool, o, func(v config.BuildTool) { w.selection.Project.BuildTool = v })
		}
	})
	optionGroup(c, "Task Runner", false, func() {
		for _, o := range taskRunnerOptions {
			radioCard(c.Key(o.value), p.TaskRunner, o, func(v config.TaskRunner) { w.selection.Project.TaskRunner = v })
		}
	})
}

func (w *Wizard) projectField(c *ui.Context, f projectField, label, hint, placeholder string) ui.Element {
	return textField(c, label, hint, placeholder, *projectValue(&w.selection.Project, f), w.showProjectErrors,
		func(v string) { w.setProjectField(f, v) })
}

func (w *Wizard) frontendSection(c *ui.Context) {
	f := w.selection.Frontend
	sectionHeading(c, "Frontend", "Choose how the user interface for "+w.selection.Project.AppName+" should be built.")
	optionGroup(c, "Frontend", true, func() {
		for _, o := range frontendOptions {
			radioCard(c.Key(o.value), f.Type, o, w.setFrontend)
		}
	})
	if config.IsServerRendered(f.Type) {
		optionGroup(c, "Additional Support", false, func() {
			checkCard(c.Key("htmx"), f.HTMX, "HTMX", "Add HTMX for dynamic server-rendered pages", false,
				func() { w.selection.Frontend.HTMX = !w.selection.Frontend.HTMX })
		})
	}
}

func (w *Wizard) databaseSection(c *ui.Context) {
	d := w.selection.Database
	sectionHeading(c, "Database", "Choose the database, persistence framework, and schema migration tool for "+w.selection.Project.AppName+".")
	optionGroup(c, "Database", true, func() {
		for _, o := range databaseOptions {
			radioCard(c.Key(o.value), d.Engine, o, func(v config.DatabaseEngine) { w.selection.Database.Engine = v })
		}
	})
	optionGroup(c, "Persistence Framework", false, func() {
		for _, o := range persistenceOptions {
			radioCard(c.Key(o.value), d.PersistenceFramework, o, func(v config.PersistenceFramework) { w.selection.Database.PersistenceFramework = v })
		}
	})
	optionGroup(c, "Database Migration Tool", false, func() {
		for _, o := range migrationOptions {
			radioCard(c.Key(o.value), d.MigrationTool, o, func(v config.MigrationTool) { w.selection.Database.MigrationTool = v })
		}
	})
}

func (w *Wizard) featuresSection(c *ui.Context) {
	sectionHeading(c, "Features", "Select any additional capabilities to include in your generated application.")
	ui.Column(c).Gap(25).Children(func() {
		featureGroup(c, "Optional Features", func() {
			for _, o := range featureOptions {
				checkCard(c.Key(o.value), slices.Contains(w.selection.Features, o.value), o.name, o.description, true,
					func() { w.toggleFeature(o.value) })
			}
		})
		featureGroup(c, "Spring AI Models (optional)", func() {
			for _, o := range aiModelOptions {
				checkCard(c.Key(o.value), slices.Contains(w.selection.AIModels, o.value), o.name, o.description, true,
					func() { w.toggleAIModel(o.value) })
			}
		})
	})
}

func featureGroup(c *ui.Context, legend string, cards func()) {
	ui.Column(c).Children(func() {
		ui.Text(c, legend).FontSize(12).FontWeight(700).TextColor(colorLabel).Margin(0, 0, 10, 0)
		ui.Grid(c).Columns(3).Gap(12).Children(cards)
	})
}

func (w *Wizard) reviewStep(c *ui.Context) {
	s := w.selection
	formCard(c, func() {
		formHeading(c, "Step 3 of 4", "Review", func() {
			subtitle(c, "Review your selections and adjust the project’s name and coordinates before generating it.")
			if w.activeTemplate != nil {
				subtitleWith(c, "From template ", w.activeTemplate.Name, "")
			}
		})
		ui.Grid(c).Columns(2).Gap(14).Children(func() {
			reviewSection(c, "Project", func() {
				w.reviewField(c, fieldAppName, "Name", "boot-demo")
				w.reviewField(c, fieldAppVersion, "Version", "1.0.0")
				w.reviewField(c, fieldGroupID, "Group ID", "com.example")
				w.reviewField(c, fieldArtifactID, "Artifact ID", "my-application")
				w.reviewField(c, fieldPackageName, "Package", "com.example.demo")
				ui.Box(c).Height(4)
				reviewRow(c, "Build Tool", label(s.Project.BuildTool))
				reviewRow(c, "Task Runner", label(s.Project.TaskRunner))
			}).RowSpan(2)
			reviewSection(c, "Frontend", func() {
				reviewRow(c, "Frontend", label(s.Frontend.Type))
				htmx := "Disabled"
				if config.IsServerRendered(s.Frontend.Type) && s.Frontend.HTMX {
					htmx = "Enabled"
				}
				reviewRow(c, "HTMX Support", htmx)
			})
			reviewSection(c, "Database", func() {
				reviewRow(c, "Database", label(s.Database.Engine))
				reviewRow(c, "Persistence", label(s.Database.PersistenceFramework))
				reviewRow(c, "Migration", label(s.Database.MigrationTool))
			})
			reviewSection(c, "Features", func() {
				names := []string{}
				for _, f := range s.Features {
					names = append(names, label(f))
				}
				reviewTags(c, names)
			}).ColumnSpan(-1)
			reviewSection(c, "Spring AI Models", func() {
				names := []string{}
				for _, m := range s.AIModels {
					names = append(names, label(m))
				}
				reviewTags(c, names)
			}).ColumnSpan(-1)
		})
		if w.templatePanel {
			w.templatePanelView(c)
		}
		if w.generationError != "" {
			errorBox(c, w.generationError)
		}
		formActions(c, func() {
			w.cancelButton(c)
			toggle := "Save as template"
			if w.templatePanel {
				toggle = "Close"
			}
			if secondaryButton(c, toggle, w.busy()).Clicked() {
				w.templatePanel = !w.templatePanel
				w.templateMsg = nil
			}
		}, func() {
			if secondaryButton(c, "Back", w.busy()).Clicked() {
				w.back()
			}
			preview := "Preview"
			if w.preparing {
				preview = "Preparing preview…"
			}
			if secondaryButton(c, preview, w.busy()).Tooltip("See every generated file before saving the project").Clicked() {
				w.startPreview()
			}
			w.generateButton(c)
		})
	})
}

// generateButton is Generate, of the Review step and the preview page.
func (w *Wizard) generateButton(c *ui.Context) {
	label := "Generate"
	if w.generating {
		label = "Generating…"
	}
	if primaryButton(c, label, "↓", w.busy()).MinWidth(125).Justify(ui.SpaceBetween).Clicked() {
		w.startGenerating()
	}
}

func reviewSection(c *ui.Context, title string, rows func()) ui.Element {
	return ui.Column(c).Padding(16, 18).Radius(10).Border(1, ui.Hex("#dedfe7")).Background(ui.RGBA(255, 255, 255, 0.58)).Children(func() {
		ui.Text(c, title).FontSize(15).FontWeight(700).TextColor(ui.Hex("#272938")).Margin(0, 0, 10, 0)
		rows()
	})
}

func reviewRow(c *ui.Context, name, value string) {
	ui.Row(c).Gap(16).Padding(5, 0).Children(func() {
		ui.Text(c, name).FontSize(13).LineHeight(1.45).TextColor(ui.Hex("#5f6375"))
		ui.Spacer(c)
		ui.Text(c, value).FontSize(13).LineHeight(1.45).FontWeight(700).TextColor(ui.Hex("#24263a")).TextAlign(ui.End)
	})
}

// reviewField is a Project value the Review step can change in place.
func (w *Wizard) reviewField(c *ui.Context, f projectField, name, placeholder string) {
	ui.Row(c).Gap(16).Padding(3, 0).Children(func() {
		ui.Text(c, name).FontSize(13).TextColor(ui.Hex("#5f6375"))
		ui.Spacer(c)
		value := *projectValue(&w.selection.Project, f)
		v := value
		input := ui.TextInput(c.Key("review:"+name), &v).Label(name).Placeholder(placeholder).
			WidthPercent(62).Padding(5, 8).Radius(6).FontSize(13).FontWeight(700).TextColor(ui.Hex("#24263a")).FocusRing(false).
			Disabled(w.busy())
		if v != value {
			w.setProjectField(f, v)
		}
		inputLook(input, w.showProjectErrors && strings.TrimSpace(v) == "", ui.RGBA(255, 255, 255, 0.75), ui.Hex("#e1e2ea"))
	})
}

func reviewTags(c *ui.Context, names []string) {
	ui.Row(c).Wrap().Gap(8).Children(func() {
		if len(names) == 0 {
			ui.Text(c, "None").FontSize(12).FontWeight(700).TextColor(ui.Hex("#5f6375")).
				Padding(6, 12).Radius(999).Background(ui.Hex("#eff0f4"))
			return
		}
		for _, n := range names {
			ui.Text(c, n).FontSize(12).FontWeight(700).TextColor(ui.Hex("#9f2830")).
				Padding(6, 12).Radius(999).Background(ui.Hex("#fff0f1"))
		}
	})
}

// templatePanelView names the template to save the selection as.
func (w *Wizard) templatePanelView(c *ui.Context) {
	ui.Column(c).Margin(18, 0, 0, 0).Padding(16, 18).Gap(8).Radius(10).Border(1, ui.Hex("#dedfe7")).
		Background(ui.RGBA(255, 255, 255, 0.58)).Children(func() {
		ui.Row(c).Gap(12).AlignItems(ui.Start).Children(func() {
			ui.Column(c).Grow(1).Children(func() {
				textFieldWith(c, "Template Name", "Saved templates are listed on the Welcome step.", "e.g. JTE + PostgreSQL",
					w.templateName, false, true, func(v string) { w.templateName = v })
			})
			ui.Row(c).Gap(8).Padding(21, 0, 0, 0).Children(func() {
				if w.activeTemplate != nil {
					if primaryButton(c, fmt.Sprintf("Update “%s”", w.activeTemplate.Name), "", w.busy()).Padding(11, 20).Clicked() {
						w.saveTemplate(true)
					}
					if secondaryButton(c, "Save as new", w.busy()).Padding(11, 20).Clicked() {
						w.saveTemplate(false)
					}
				} else if primaryButton(c, "Save template", "", w.busy()).Padding(11, 20).Clicked() {
					w.saveTemplate(false)
				}
			})
		})
		feedback(c, w.templateMsg)
	})
}

func (w *Wizard) completeStep(c *ui.Context) {
	width, _ := c.Size()
	ui.Column(c).WidthPercent(90).MaxWidth(650).AlignItems(ui.Center).Padding(38, 0, 0, 0).Children(func() {
		successIcon(c)
		eyebrow(c, "Complete").Margin(0, 0, 10, 0)
		heading(c, "Project downloaded", clamp(width*0.031, 27, 42)).TextAlign(ui.Center)
		ui.Text(c, fmt.Sprintf("Your Spring Boot %s project was generated and saved as a ZIP archive.", label(w.selection.Project.BuildTool))).
			FontSize(15).LineHeight(1.65).TextColor(colorMuted).TextAlign(ui.Center).MaxWidth(525).Margin(16, 0, 29, 0)
		ui.Column(c).FillWidth().Margin(0, 0, 25, 0).Padding(15, 17).Radius(10).Border(1, colorBorder).
			Background(ui.RGBA(255, 255, 255, 0.7)).Children(func() {
			ui.Text(c, "DOWNLOADED PROJECT").FontSize(9).FontWeight(700).LetterSpacing(0.72).TextColor(ui.Hex("#8b8e9b")).Margin(0, 0, 7, 0)
			ui.Text(c, w.downloadedPath).Font("monospace").FontSize(11).LineHeight(1.5).TextColor(ui.Hex("#303242")).Selectable()
		})
		for _, warning := range w.generationWarnings {
			warningBox(c, warning).FillWidth().Margin(0, 0, 25, 0)
		}
		if outlineButton(c, "Start New").Clicked() {
			w.startNew()
		}
	})
}

func (w *Wizard) templatesPage(c *ui.Context) {
	formCard(c, func() {
		count := "1 template"
		if len(w.templateList) != 1 {
			count = fmt.Sprintf("%d templates", len(w.templateList))
		}
		formHeading(c, count, "Templates", func() {
			subtitle(c, "Start a project from a saved selection, or update and delete your templates. Save new ones from the Review step.")
		})
		if w.templatesError != "" {
			errorBox(c, w.templatesError).Margin(0, 0, 14, 0)
		}
		if len(w.templateList) > 0 {
			w.templateItems(c)
		} else {
			ui.RichText(c,
				ui.Span{Text: "No templates yet. Make your selections, then use "},
				ui.Span{Text: "Save as template", Weight: 700, Color: colorLabel},
				ui.Span{Text: " on the Review step."},
			).FontSize(13).LineHeight(1.5).TextColor(ui.Hex("#6b6f82")).Padding(18, 20).Radius(12).
				Border(1, colorButtonEdge).BorderStyle(ui.BorderDashed)
		}
		formActions(c, nil, func() {
			if secondaryButton(c, "Back to "+w.returnName(), false).Clicked() {
				w.page = pageNone
			}
		})
	})
}

func (w *Wizard) defaultsPage(c *ui.Context) {
	d := w.defaultsDraft
	formCard(c, func() {
		formHeading(c, "Defaults", "Project Defaults", func() {
			subtitle(c, "New projects start with these values on the Configure step. Templates keep their own.")
		})
		ui.Grid(c).Columns(2).Gap(20).Children(func() {
			textField(c, "Application Version", "The initial version of new applications.", "1.0.0", d.AppVersion, false,
				func(v string) { w.defaultsDraft.AppVersion = v; w.defaultsMsg = nil })
			textField(c, "Group ID", "Also the start of the package name.", "com.example", d.GroupID, false,
				func(v string) { w.defaultsDraft.GroupID = v; w.defaultsMsg = nil })
		})
		optionGroup(c, "Build Tool", false, func() {
			for _, o := range buildToolOptions {
				radioCard(c.Key(o.value), d.BuildTool, o, func(v config.BuildTool) { w.defaultsDraft.BuildTool = v; w.defaultsMsg = nil })
			}
		})
		optionGroup(c, "Task Runner", false, func() {
			for _, o := range taskRunnerOptions {
				radioCard(c.Key(o.value), d.TaskRunner, o, func(v config.TaskRunner) { w.defaultsDraft.TaskRunner = v; w.defaultsMsg = nil })
			}
		})
		if w.defaultsMsg != nil {
			ui.Column(c).Margin(18, 0, 0, 0).Children(func() { feedback(c, w.defaultsMsg) })
		}
		formActions(c, func() {
			if secondaryButton(c, "Restore built-in values", false).Tooltip("Fill in Progen Desktop's own defaults; Save to keep them").Clicked() {
				w.defaultsDraft = store.BuiltInDefaults
				w.defaultsMsg = nil
			}
		}, func() {
			if secondaryButton(c, "Back to "+w.returnName(), false).Clicked() {
				w.page = pageNone
			}
			if primaryButton(c, "Save defaults", "", false).Clicked() {
				w.saveDefaults()
			}
		})
	})
}
