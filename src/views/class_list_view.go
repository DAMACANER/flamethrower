package views

import (
	"flamethrower/src/db"
	"flamethrower/src/db/model"
	"flamethrower/src/db/table"
	"log"
	"sync"

	"github.com/gdamore/tcell/v2"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/rivo/tview"
)

func ReturnClassView(app *tview.Application) *tview.Flex {
	var PageSize int64 = DefaultPageSize
	var CurrentPageNumber int64 = DefaultStartingPageNumber
	stmt := SELECT(table.Class.AllColumns).
		FROM(table.Class).
		LIMIT(PageSize).
		OFFSET(CurrentPageNumber * PageSize).
		ORDER_BY(table.Class.Name.ASC())
	var data []model.Class
	err := stmt.Query(db.DB, &data)
	if err != nil {
		return ReturnErrorView(app, err.Error())
	}	
	list := tview.NewList()

	title := tview.NewTextView()
	title.SetText("Choose a Class")
	title.SetTextAlign(tview.AlignCenter)

	for i, skill := range data {
		list.AddItem(skill.Name, "", rune(i), nil)
	}
	var cnt int
	stmt = SELECT(COUNT(table.Class.ID)).
		FROM(table.Class)
	err = stmt.Query(db.DB, &cnt)
	if err != nil {
		log.Fatal(err)
	}
	fetchingNewData := false
	var fetchingMutex sync.Mutex

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		fetchingMutex.Lock()
		nearEnd := index >= len(data)-int(PageSize)/5
		fetchingMutex.Unlock()

		if nearEnd && len(data) < cnt {
			fetchingMutex.Lock()
			if fetchingNewData {
				fetchingMutex.Unlock()
				return
			}
			fetchingNewData = true
			fetchingMutex.Unlock()

			go func() {
				CurrentPageNumber++
				stmt = SELECT(table.Class.Name).
					FROM(table.Class).
					LIMIT(PageSize).
					OFFSET(CurrentPageNumber * PageSize).
					ORDER_BY(table.Class.Name.ASC())
				var newData []model.Class
				err = stmt.Query(db.DB, &newData)
				if err != nil {
					app.QueueUpdateDraw(func() {
						app.SetRoot(ReturnErrorView(app, err.Error()), true)
					})
					return
				}
				app.QueueUpdateDraw(func() {
					for i, newSkill := range newData {
						list.AddItem(newSkill.Name, "", rune(len(data)+i), nil)
					}
					data = append(data, newData...)
					fetchingMutex.Lock()
					fetchingNewData = false
					fetchingMutex.Unlock()
				})
			}()
		}
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
		ch := event.Rune()
		if ch == 'q' || ch == 'Q' {
			app.Stop()
		}
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
