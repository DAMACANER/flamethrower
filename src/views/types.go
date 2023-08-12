package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xuri/excelize/v2"
)

var DefaultPageSize int64 = 5
var DefaultStartingPageNumber int64 = 0

type Application struct {
	App   *tview.Application
	Sheet *excelize.File
}
type Element struct {
	View       tview.Primitive
	ParentFlex *tview.Flex
	Position   struct {
		Row    int
		Column int
	}
}

func TraverseBoxes(currentElement *Element, elements [][]*Element, app *tview.Application) func(event *tcell.EventKey) *tcell.EventKey {
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
