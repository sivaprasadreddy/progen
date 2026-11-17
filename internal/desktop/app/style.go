package app

import (
	_ "embed"
	"log"
	"strings"

	"github.com/egoist/mygo/ui"
)

// The app's look: Inter, its colors and the shapes the screens share.

//go:embed assets/Inter-Medium.ttf
var interFont []byte

//go:embed assets/logo.png
var logoPNG []byte

const fontFamily = "Inter"

var (
	logo = mustBitmap(logoPNG)

	checkIcon   = mustSVG(`<svg viewBox="0 0 32 32" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M7 16.5 13 22l12-13"/></svg>`)
	tickIcon    = mustSVG(`<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3.5 8.5 6.5 11.5 12.5 4.5"/></svg>`)
	listIcon    = mustSVG(`<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M3 4.5h10M3 8h10M3 11.5h6"/></svg>`)
	slidersIcon = mustSVG(`<svg viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"><path d="M3 4.5h10M3 11.5h10"/><circle cx="6" cy="4.5" r="1.6" fill="currentColor" stroke="none"/><circle cx="10.5" cy="11.5" r="1.6" fill="currentColor" stroke="none"/></svg>`)
)

func init() {
	if err := ui.RegisterFont(interFont, fontFamily); err != nil {
		log.Fatalf("register Inter: %v", err)
	}
}

func mustBitmap(data []byte) *ui.Bitmap {
	b, err := ui.DecodeBitmap(data)
	if err != nil {
		log.Fatalf("decode logo: %v", err)
	}
	return b
}

func mustSVG(markup string) *ui.SVG {
	s, err := ui.ParseSVG([]byte(markup))
	if err != nil {
		log.Fatalf("parse icon: %v", err)
	}
	return s
}

var (
	colorText        = ui.Hex("#171828")
	colorTextStrong  = ui.Hex("#202231")
	colorLabel       = ui.Hex("#303244")
	colorMuted       = ui.Hex("#707487")
	colorSubtle      = ui.Hex("#74788b")
	colorHint        = ui.Hex("#9295a4")
	colorOptionHint  = ui.Hex("#8d91a0")
	colorAccent      = ui.Hex("#ef3e48")
	colorAccentHover = ui.Hex("#dc303b")
	colorEyebrow     = ui.Hex("#ef444d")
	colorSelected    = ui.Hex("#ef4a54")
	colorSuccess     = ui.Hex("#087f45")
	colorDanger      = ui.Hex("#a5212a")
	colorBorder      = ui.Hex("#d9dae3")
	colorButtonText  = ui.Hex("#535769")
	colorButtonEdge  = ui.Hex("#d3d5df")
	white            = ui.RGB(255, 255, 255)
	focusGlow        = ui.RGBA(239, 62, 72, 0.1)
)

// theme is the window's theme: Inter in the app's colors, the focus ring
// as the page's soft red outline.
func theme(c *ui.Context) *ui.Theme {
	t := *ui.LightTheme().WithAccent(colorAccent)
	t.Font = fontFamily
	t.Text = colorText
	t.Background = ui.Hex("#f6f4fa")
	t.Focus = ui.RGBA(239, 62, 72, 0.22)
	t.FocusRingWidth = 3
	return &t
}

// eyebrow is the small red caption above a heading.
func eyebrow(c *ui.Context, s string) ui.Element {
	return ui.Text(c, strings.ToUpper(s)).FontSize(11).FontWeight(800).LetterSpacing(1.76).TextColor(colorEyebrow)
}

// heading is a page title of the given size.
func heading(c *ui.Context, s string, size float32) ui.Element {
	return ui.Text(c, s).FontSize(size).FontWeight(700).LineHeight(1.18).LetterSpacing(-0.035 * size).TextColor(colorText)
}

// formCard is the white card a step's form sits in.
func formCard(c *ui.Context, children func()) ui.Element {
	return ui.Column(c).FillWidth().MaxWidth(760).Padding(38, 44).Radius(18).
		Border(1, ui.RGBA(205, 207, 221, 0.8)).Background(ui.RGBA(255, 255, 255, 0.68)).
		Shadow(0, 14, 45, 0, ui.RGBA(57, 57, 85, 0.08)).Children(children)
}

// formHeading is a card's eyebrow, title and description, which describe
// builds when it needs more than plain text.
func formHeading(c *ui.Context, eyebrowText, title string, describe func()) {
	ui.Column(c).Margin(0, 0, 26, 0).Children(func() {
		eyebrow(c, eyebrowText).Margin(0, 0, 8, 0)
		heading(c, title, 30)
		ui.Column(c).Margin(8, 0, 0, 0).Gap(6).Children(describe)
	})
}

// subtitle is a heading's line of description.
func subtitle(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).FontSize(13).TextColor(colorSubtle)
}

// subtitleWith is a description with one part set apart in a darker color,
// as a template's name.
func subtitleWith(c *ui.Context, before, marked, after string) ui.Element {
	return ui.RichText(c,
		ui.Span{Text: before},
		ui.Span{Text: marked, Color: colorLabel},
		ui.Span{Text: after},
	).FontSize(13).TextColor(colorSubtle)
}

// optionGroup is a legend over a grid of option cards, three to a row.
func optionGroup(c *ui.Context, legend string, first bool, cards func()) {
	top := float32(25)
	if first {
		top = 0
	}
	ui.Column(c).Margin(top, 0, 0, 0).Children(func() {
		ui.Text(c, legend).FontSize(12).FontWeight(700).TextColor(colorLabel).Margin(0, 0, 10, 0)
		ui.Grid(c).Columns(3).Gap(10).Children(cards)
	})
}

// cardLook styles an option card, plain or selected.
func cardLook(card ui.Element, selected bool, radius float32) ui.Element {
	card.Radius(radius).Cursor(ui.CursorPointer)
	if selected {
		return card.Border(1, colorSelected).Background(ui.Hex("#fff7f7")).Shadow(0, 0, 0, 2, ui.RGBA(239, 62, 72, 0.07))
	}
	if card.Hovered() {
		return card.Border(1, ui.Hex("#c9cbd8")).Background(ui.RGBA(255, 255, 255, 0.8))
	}
	return card.Border(1, colorBorder).Background(ui.RGBA(255, 255, 255, 0.55))
}

// radioDot is a radio button's circle, filled in the accent when on.
func radioDot(c *ui.Context, on bool) {
	dot := ui.Box(c).Size(16, 16).Shrink(0).Radius(8).Center()
	if !on {
		dot.Background(white).Border(1, ui.Hex("#b9bcc8"))
		return
	}
	dot.Background(colorAccent).Children(func() {
		ui.Box(c).Size(6, 6).Radius(3).Background(white)
	})
}

// checkMark is a check box's square, filled with a tick when on.
func checkMark(c *ui.Context, on bool) {
	box := ui.Box(c).Size(16, 16).Shrink(0).Radius(4).Center()
	if !on {
		box.Background(white).Border(1, ui.Hex("#b9bcc8"))
		return
	}
	box.Background(colorAccent).Children(func() {
		ui.Icon(c, tickIcon).Size(12, 12).TextColor(white)
	})
}

// optionText is an option card's name over its description.
func optionText(c *ui.Context, name, description string, nameSize, gap, lineHeight float32) {
	ui.Column(c).Grow(1).Shrink(1).Gap(gap).Children(func() {
		ui.Text(c, name).FontSize(nameSize).FontWeight(700).TextColor(colorTextStrong)
		ui.Text(c, description).FontSize(11).LineHeight(lineHeight).TextColor(colorOptionHint)
	})
}

// radioCard is an option card choosing o.value into a copy of *selected;
// choose runs when the user picks it.
func radioCard[T comparable](c *ui.Context, selected T, o option[T], choose func(T)) {
	value := selected
	card := ui.RadioBase(c, &value, o.value).Label(o.name)
	if value != selected {
		choose(value)
	}
	on := value == o.value
	cardLook(card, on, 11).Gap(11).MinHeight(72).Padding(14, 16).AlignItems(ui.Center).Children(func() {
		radioDot(c, on)
		optionText(c, o.name, o.description, 13, 4, 1.4)
	})
}

// checkCard is an option card toggling a choice; toggle runs on a click.
// large is the Features step's taller card.
func checkCard(c *ui.Context, checked bool, name, description string, large bool, toggle func()) {
	value := checked
	card := ui.CheckboxBase(c, &value).Label(name)
	if value != checked {
		toggle()
	}
	if large {
		cardLook(card, value, 12).Gap(14).MinHeight(92).Padding(18, 20)
	} else {
		cardLook(card, value, 11).Gap(11).MinHeight(72).Padding(14, 16)
	}
	card.AlignItems(ui.Center).Children(func() {
		checkMark(c, value)
		if large {
			optionText(c, name, description, 14, 5, 1.45)
		} else {
			optionText(c, name, description, 13, 4, 1.4)
		}
	})
}

// textField is a labeled text input with a hint below it. invalid marks it
// empty when the user tried to go on.
func textField(c *ui.Context, label, hint, placeholder string, value string, invalid bool, change func(string)) ui.Element {
	return textFieldWith(c, label, hint, placeholder, value, invalid, false, change)
}

// textFieldWith is textField, taking the focus as it appears when focus is
// true, as a field opened to be filled in.
func textFieldWith(c *ui.Context, label, hint, placeholder string, value string, invalid, focus bool, change func(string)) ui.Element {
	return ui.Column(c).Gap(7).Children(func() {
		ui.Text(c, label).FontSize(12).FontWeight(700).TextColor(colorLabel)
		v := value
		input := ui.TextInput(c.Key(label), &v).Placeholder(placeholder).Label(label).
			Padding(11, 13).Radius(8).FontSize(13).TextColor(colorTextStrong).FocusRing(false)
		if focus {
			input.AutoFocus()
		}
		if v != value {
			change(v)
		}
		inputLook(input, invalid && strings.TrimSpace(v) == "", ui.RGBA(255, 255, 255, 0.85), ui.Hex("#d7d8e2"))
		if invalid && strings.TrimSpace(v) == "" {
			ui.Text(c, "Required.").FontSize(9).TextColor(colorDanger)
		} else {
			ui.Text(c, hint).FontSize(9).TextColor(colorHint)
		}
	})
}

// inputLook styles a text input: its border reddens while it has the focus,
// and while it is invalid.
func inputLook(input ui.Element, invalid bool, background, border ui.Color) {
	switch {
	case input.Focused():
		input.Background(white).Border(1, ui.Hex("#ef5660")).Shadow(0, 0, 0, 3, focusGlow)
	case invalid:
		input.Background(ui.Hex("#fff2f3")).Border(1, ui.Hex("#f0b4b8"))
	case input.Hovered():
		input.Background(white).Border(1, ui.Hex("#cfd1dc"))
	default:
		input.Background(background).Border(1, border)
	}
}

// secondaryButton is an outlined button, as Back.
func secondaryButton(c *ui.Context, label string, disabled bool) ui.Element {
	b := ui.ButtonBase(c.Key("btn:"+label)).Label(label).Disabled(disabled).Padding(10, 20).Radius(8).Cursor(ui.CursorPointer)
	if b.Hovered() && !disabled {
		b.Background(white)
	} else {
		b.Background(ui.RGBA(255, 255, 255, 0.45))
	}
	b.Border(1, colorButtonEdge).Children(func() {
		ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(colorButtonText)
	})
	return b
}

// primaryButton is the red button that goes on, with an arrow after its
// label when arrow is not empty.
func primaryButton(c *ui.Context, label, arrow string, disabled bool) ui.Element {
	b := ui.ButtonBase(c.Key("btn:"+label)).Label(label).Disabled(disabled).Padding(10, 20).Radius(8).Gap(20).Cursor(ui.CursorPointer)
	fill := colorAccent
	if b.Hovered() && !disabled {
		fill = colorAccentHover
	}
	b.Background(fill).Border(1, colorAccent).Children(func() {
		ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(white)
		if arrow != "" {
			ui.Text(c, arrow).FontSize(12).FontWeight(600).TextColor(white)
		}
	})
	return b
}

// outlineButton is the outlined red button of the Welcome and Complete
// screens, filled while the pointer is over it.
func outlineButton(c *ui.Context, label string) ui.Element {
	b := ui.ButtonBase(c.Key("btn:"+label)).Label(label).MinWidth(118).Padding(11, 23).Radius(9).Justify(ui.Center).Cursor(ui.CursorPointer)
	hovered := b.Hovered()
	text := ui.Hex("#e9323d")
	if hovered {
		b.Background(colorAccent).Top(-1)
		text = white
	} else {
		b.Background(ui.RGBA(255, 255, 255, 0.35))
	}
	b.Border(1.5, colorAccent).Children(func() {
		ui.Text(c, label).FontSize(13).FontWeight(600).TextColor(text)
	})
	return b
}

// smallButton is a template's Update or Delete.
func smallButton(c *ui.Context, label string) ui.Element {
	b := ui.ButtonBase(c.Key("small:"+label)).Label(label).Padding(6, 11).Radius(7).Cursor(ui.CursorPointer)
	text := colorButtonText
	if b.Hovered() {
		b.Background(white)
		text = colorTextStrong
	} else {
		b.Background(ui.RGBA(255, 255, 255, 0.7))
	}
	b.Border(1, colorButtonEdge).Children(func() {
		ui.Text(c, label).FontSize(12).FontWeight(600).TextColor(text)
	})
	return b
}

// formActions is the row of buttons at a card's foot: start on the left,
// end on the right.
func formActions(c *ui.Context, start, end func()) {
	ui.Row(c).Margin(29, 0, 0, 0).Gap(10).Children(func() {
		if start != nil {
			start()
		}
		ui.Spacer(c)
		end()
	})
}

// feedback is a line of success or error under a form.
func feedback(c *ui.Context, m *message) {
	if m == nil {
		return
	}
	color := colorSuccess
	if m.error {
		color = colorDanger
	}
	ui.Text(c, m.text).FontSize(13).TextColor(color)
}

// errorBox is an error in a red box, as a failed generation.
func errorBox(c *ui.Context, text string) ui.Element {
	return ui.Column(c).Margin(14, 0, 0, 0).Padding(10, 12).Radius(8).
		Border(1, ui.Hex("#f0b4b8")).Background(ui.Hex("#fff2f3")).Children(func() {
		ui.Text(c, text).FontSize(13).LineHeight(1.45).TextColor(colorDanger)
	})
}

// warningBox is a warning in an amber box, as a failed code formatting.
func warningBox(c *ui.Context, text string) ui.Element {
	return ui.Column(c).Margin(14, 0, 0, 0).Padding(10, 12).Radius(8).
		Border(1, ui.Hex("#f0d58a")).Background(ui.Hex("#fff8e6")).Children(func() {
		ui.RichText(c,
			ui.Span{Text: "WARNING: ", Weight: 700},
			ui.Span{Text: text},
		).FontSize(13).LineHeight(1.45).TextColor(ui.Hex("#8a5a00")).Selectable()
	})
}
