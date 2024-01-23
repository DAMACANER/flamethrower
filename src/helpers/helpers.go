package helpers

import (
	"flamethrower/src/types"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"jaytaylor.com/html2text"
)

func NotEmpty(v interface{}) bool {
	if v == nil {
		return false
	}

	switch v := v.(type) {
	case string:
		trimmed := strings.TrimSpace(v)
		return trimmed != "" && trimmed != "None"
	case *string:
		if v == nil {
			return false
		}
		trimmed := strings.TrimSpace(*v)
		return trimmed != "" && trimmed != "None"
	case *int, *int32, *int64:
		return v != nil
	case *float32, *float64:
		return v != nil
	default:
		return true
	}
}

func TraverseBoxes(currentElement *types.Element, elements [][]*types.Element, app *tview.Application) func(event *tcell.EventKey) *tcell.EventKey {
	moveFocus := func(rowChange, colChange int) {
		newRow := currentElement.Position.Row + rowChange
		newCol := currentElement.Position.Column + colChange

		if newRow < 0 || newRow >= len(elements) {
			return // Out of row bounds
		}

		row := elements[newRow]
		if newCol < 0 || newCol >= len(row) {
			return // Out of column bounds
		}

		// Adjust for columns out of index
		if len(row)-1 < newCol {
			newCol = len(row) - 1
		}

		currentElement = row[newCol]
		app.SetFocus(currentElement.View)
	}

	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'w', 'W':	
			moveFocus(-1, 0)
		case 'a', 'A':
			moveFocus(0, -1)
		case 's', 'S':
			moveFocus(1, 0)
		case 'd', 'D':
			moveFocus(0, 1)
		case 'q', 'Q':
			app.Stop()
		}
		return event
	}
}

func CreateListFlex(title string, items []types.ListItem, opts types.ListAppearanceOptions) *tview.Flex {
	listFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	list := tview.NewList()
	list.ShowSecondaryText(opts.ShouldShowSecondaryText)
	list.SetBorder(opts.ShouldDrawBorder)
	list.SetTitle(opts.Title)
	list.SetTitleAlign(opts.TitleAlignment)
	list.SetBorderColor(opts.BorderColor)
	list.SetBorderPadding(opts.BorderPadding.Top, opts.BorderPadding.Bottom, opts.BorderPadding.Left, opts.BorderPadding.Right)
	list.SetBackgroundColor(opts.BackgroundColor)
	// Add items to the list
	for _, item := range items {
		list.AddItem(item.Label, item.Secondary, item.Shortcut, item.SelectFunc).SetMainTextColor(item.Color).SetSecondaryTextColor(tcell.ColorWhite)
	}
	listFlex.AddItem(list, 0, 1, false)
	return listFlex
}
func CreateTextView(title string, lines []types.TextViewLines, opts types.TextViewAppearanceOptions) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(opts.EnableDynamicColors).
		SetRegions(opts.SetRegions).
		SetWrap(opts.SetWrap).
		SetWordWrap(opts.SetWordWrap)
	textView.SetBorder(opts.SetBorder).SetTitle(title)
	for _, line := range lines {
		if NotEmpty(line) {
			if line.ShouldFormatHTML {
				plainText, err := html2text.FromString(fmt.Sprintf("%v", &line.Line))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprintf(textView, "%s \n", plainText)
			} else {
				fmt.Fprint(textView, line.Line)
				fmt.Fprint(textView, "\n")
			}
		}
	}
	return textView
}
