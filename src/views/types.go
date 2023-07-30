package views

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/xuri/excelize/v2"
)

var DefaultPageSize uint64 = 5
var DefaultStartingPageNumber uint64 = 0

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
	return func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'w', 'W':
			switch currentElement.Position.Row {
			case 0:
				break
			default:
				row := elements[currentElement.Position.Row-1]
				if len(row)-1 < currentElement.Position.Column {
					currentElement = row[len(row)-1]
				} else {
					currentElement = row[currentElement.Position.Column]
				}
				app.SetFocus(currentElement.View)
			}
		case 'a', 'A':
			switch currentElement.Position.Column {
			case 0:
				break
			default:
				currentElement = elements[currentElement.Position.Row][currentElement.Position.Column-1]
				app.SetFocus(currentElement.View)
			}
		case 's', 'S':
			switch currentElement.Position.Row {
			case len(elements) - 1:
				break
			default:
				row := elements[currentElement.Position.Row+1]
				if len(row)-1 < currentElement.Position.Column {
					currentElement = row[len(row)-1]
				} else {
					currentElement = row[currentElement.Position.Column]
				}

				app.SetFocus(currentElement.View)
			}
		case 'd', 'D':
			switch currentElement.Position.Column {
			case len(elements[currentElement.Position.Row]) - 1:
				break
			default:
				currentElement = elements[currentElement.Position.Row][currentElement.Position.Column+1]
				app.SetFocus(currentElement.View)
			}
		case 'q', 'Q':
			app.Stop()
		}
		return event
	}
}
