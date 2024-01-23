package types

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Element struct {
	View       tview.Primitive
	ParentFlex *tview.Flex
	Position   struct {
		Row    int
		Column int
	}
}

type ListItem struct {
	Label      string
	Secondary  string
	Shortcut   rune
	Color      tcell.Color
	SelectFunc func()
}

type BorderPadding struct {
	Top    int
	Bottom int
	Left   int
	Right  int
}

type ListAppearanceOptions struct {
	ShouldDrawBorder        bool
	ShouldShowSecondaryText bool
	BorderColor             tcell.Color
	BackgroundColor         tcell.Color
	// please use tview.FlexColumn etc.
	ListDirection int
	Title         string
	// please use tview.AlignLeft, or tview.AlignRight, dont just use 0,1, etc.
	TitleAlignment int
	BorderPadding  BorderPadding
}

type TextViewLines struct {
	Line string
	// Pass this as true if line needs to be formatted from HTML to Text.
	ShouldFormatHTML bool
}

type TextViewAppearanceOptions struct {
	EnableDynamicColors bool
	SetRegions          bool
	SetWrap             bool // wrap of text view
	SetWordWrap         bool // word wrap of text within the view
	SetBorder           bool
}
