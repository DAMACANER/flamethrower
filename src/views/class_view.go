package views

import (
	"flamethrower/src/db"
	"fmt"
	"log"

	"github.com/blockloop/scan"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"jaytaylor.com/html2text"
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
		list.AddItem(skill.Name, "", rune(i), nil)
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
					list.AddItem(newSkill.Name, "", rune(len(data)+i), nil)
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

func ReturnClassDetailView(class db.ClassListColumns, app *tview.Application) *tview.Flex {
	epicView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		}).SetWrap(true).SetWordWrap(true)
	epicView.SetBorder(true)
	epicView.SetTitle("Epic Stat List")
	fmt.Fprintf(epicView, "[white]Epic Feat Base Level: [red]%s\n\n", class.EpicFeatBaseLevel.String)
	fmt.Fprintf(epicView, "[white]Epic Feat Interval: [red]%s\n\n", class.EpicFeatInterval.String)
	fmt.Fprintf(epicView, "[white]Epic Feat List: [red]%s\n\n", class.EpicFeatList.String)
	epicFullTextPlain, err := html2text.FromString(class.EpicFullText.String)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(epicView, "[white]Epic Full Text: \n [red]%s\n\n", epicFullTextPlain)

	generalStatsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	generalStatsView.SetBorder(true)
	generalStatsView.SetTitle(fmt.Sprintf("Stat Sheet of [green]%s", class.Name))
	fmt.Fprintf(generalStatsView, "[white]Alignment: [blue]%s\n\n", class.Alignment.String)
	fmt.Fprintf(generalStatsView, "[white]Skill Points: [blue]%s\n\n", class.SkillPoints.String)
	fmt.Fprintf(generalStatsView, "[white]Class Type: [blue]%s\n\n", class.Type.String)
	fmt.Fprintf(generalStatsView, "[white]Hit Die: [blue]%s\n\n", class.HitDie.String)
	fmt.Fprintf(generalStatsView, "[white]Class Skills: [blue]%s\n\n", class.ClassSkills.String)
	fmt.Fprintf(generalStatsView, "[white]Skill Points Ability: [blue]%s\n\n", class.SkillPointsAbility.String)
	fmt.Fprintf(generalStatsView, "[white]Spell Stat: [blue]%s\n\n", class.SpellStat.String)
	fmt.Fprintf(generalStatsView, "[white]Spell Type: [blue]%s\n\n", class.SpellType.String)
	fmt.Fprintf(generalStatsView, "[white]Proficencies: [blue]%s\n\n", class.Proficiencies.String)

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetText("[yellow]Commands: \n").SetTextAlign(tview.AlignCenter)

	fmt.Fprintf(footer, "[white]W: [red]Change Focus To Upper Box\n")
	fmt.Fprintf(footer, "[white]A: [red]Change Focus To Left Box\n")
	fmt.Fprintf(footer, "[white]S: [red]Change Focus To Below Box\n")
	fmt.Fprintf(footer, "[white]D: [red]Change Focus To Right Box\n")
	fmt.Fprintf(footer, "[white]Q: [red]Quit\n")

	elements := make(map[string]*Element)
	elements[ClassDetailsGeneralStatsBox] = &Element{View: generalStatsView, ParentFlex: nil}
	columnBelowGeneralStatsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(epicView, 0, 1, false)
	elements[ClassDetailsEpicStatsBox] = &Element{View: epicView, ParentFlex: columnBelowGeneralStatsFlex}
	columnFootbar := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(footer, 0, 1, false)
	elements[ClassDetailsFooterBox] = &Element{View: footer, ParentFlex: columnFootbar}
	flex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(generalStatsView, 0, 3, false).
			AddItem(columnBelowGeneralStatsFlex, 0, 4, false).
			AddItem(columnFootbar, 0, 1, false), 0, 1, false)

	var currentElement *Element

	currentElement = elements[ClassDetailsGeneralStatsBox]

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft, tcell.KeyRight:
			print("a")
		}
		switch event.Rune() {
		case 'w', 'W':
			// switch focus between generalStatsView and the columnFlex (which contains epicView1 and epicView2)
			if currentElement == elements[ClassDetailsGeneralStatsBox] {
				currentElement = elements[ClassDetailsEpicStatsBox]
				app.SetFocus(currentElement.View)
			} else if currentElement.ParentFlex == columnBelowGeneralStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}
		case 's', 'S':
			// switch focus between generalStatsView and the columnFlex (which contains epicView1 and epicView2)
			if currentElement == elements[ClassDetailsGeneralStatsBox] {
				currentElement = elements[ClassDetailsEpicStatsBox]
				app.SetFocus(currentElement.View)
			} else if currentElement.ParentFlex == columnBelowGeneralStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}
		case 'q', 'Q':
			app.Stop()
		case 'a', 'A':
			// switch focus between generalStatsView and the columnFlex (which contains epicView1 and epicView2)
			if currentElement == elements[ClassDetailsGeneralStatsBox] {
				currentElement = elements[ClassDetailsEpicStatsBox]
				app.SetFocus(currentElement.View)
			}
		case 'd', 'D':
			if currentElement.ParentFlex == columnBelowGeneralStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}

		}
		return event
	})
	return flex

}
