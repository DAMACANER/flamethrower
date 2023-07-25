package views

import (
	"flamethrower/src/db"
	"log"

	"github.com/blockloop/scan"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ReturnClassView(app *tview.Application) *tview.Flex {
	var pageSize uint64 = 6
	var currentPageNumber uint64 = 0
	repo := &db.ClassListRepo{BaseRepo: &db.BaseRepo{}}
	rows, err := repo.Find().Paginate(currentPageNumber, pageSize).OrderBy("name COLLATE NOCASE").Query()
	if err != nil {
		log.Fatal(err)
	}
	var data []db.ClassListColumns
	err = scan.Rows(&data, rows)
	if err != nil {
		log.Fatal(err)
	}

	list := tview.NewList()

	title := tview.NewTextView()
	title.SetText("Choose a Class")
	title.SetTextAlign(tview.AlignCenter)

	for i, skill := range data {
		list.AddItem(skill.Name.String, "", rune(i), nil)
	}
	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
			if index >= len(data)-1 && len(data) < int(db.TotalClassCount) {
				currentPageNumber++
				rows, err := repo.Find().Paginate(currentPageNumber, pageSize).OrderBy("name COLLATE NOCASE").Query()
				if err != nil {
					log.Fatal(err)
				}
				var newData []db.ClassListColumns
				err = scan.Rows(&newData, rows) // newData should be here, not data
				if err != nil {
					log.Fatal(err)
				}
				for i, newSkill := range newData {
					list.AddItem(newSkill.Name.String, "", rune(len(data)+i), nil)
				}
				data = append(data, newData...)
			}
		})
	})
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		app.SetRoot(ReturnClassDetailView(data[index], app), true)
	})

	list.SetSelectedFocusOnly(true) // only change colors when in focus
	list.SetHighlightFullLine(true) // highlight the full line of selected item

	// setting color attributes
	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSecondaryTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorWhite)

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetText("[yellow]Q to quit")

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Key() returns the key code, e.g. tcell.KeyUp
		// Rune() returns the key as a rune, e.g. 'q'
		// Modifiers() returns a bitmask representing shift, ctrl, etc.

		ch := event.Rune()

		if ch == 'q' || ch == 'Q' {
			app.Stop()
		}
		// Return event to continue processing it.
		return event
	})

	vFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	hFlex := tview.NewFlex()

	// vertical centering
	vFlex.AddItem(tview.NewBox(), 0, 1, false)
	vFlex.AddItem(title, 2, 0, false)
	vFlex.AddItem(list, 0, 2, true)
	vFlex.AddItem(tview.NewBox(), 0, 1, false)
	vFlex.AddItem(footer, 1, 0, false)

	// horizontal centering
	hFlex.AddItem(tview.NewBox(), 0, 1, false)
	hFlex.AddItem(vFlex, 0, 2, true)
	hFlex.AddItem(tview.NewBox(), 0, 1, false)
	return hFlex

}
