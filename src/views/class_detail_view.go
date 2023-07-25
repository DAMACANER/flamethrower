package views

import (
	"database/sql"
	"flamethrower/src/db"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"jaytaylor.com/html2text"
)

func ReturnClassDetailView(class db.ClassListColumns, app *tview.Application) *tview.Flex {
	//
	// views
	//
	val := reflect.ValueOf(class)
	epicView := classDetailsEpicFeat(class)
	generalStatsView := classDetailGeneralStats(val, class.Name)
	footerView := classDetailFooter()
	preReqView := classDetailsPrerequirements(val)
	spellsView := classDetailsSpells(class)
	//
	// flexes
	//
	rowBelowGeneralStatsFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(epicView, 0, 1, false).
		AddItem(preReqView, 0, 1, false).
		AddItem(spellsView, 0, 1, false)
	footbarFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(footerView, 0, 1, false)
	generalStatsFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(generalStatsView, 0, 1, false)
	mainFlex := tview.NewFlex().
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(generalStatsFlex, 0, 3, false).
			AddItem(rowBelowGeneralStatsFlex, 0, 4, false).
			AddItem(footbarFlex, 0, 1, false), 0, 1, false)
	//
	// element positioning
	//
	elements := make(map[string]*Element)
	elements[ClassDetailsFooterBox] = &Element{View: footerView, ParentFlex: footbarFlex}
	elements[ClassDetailsGeneralStatsBox] = &Element{View: generalStatsView, ParentFlex: generalStatsFlex}
	elements[ClassDetailsEpicStatsBox] = &Element{View: epicView, ParentFlex: rowBelowGeneralStatsFlex}
	elements[ClassDetailsPrereqBox] = &Element{View: preReqView, ParentFlex: rowBelowGeneralStatsFlex}
	elements[ClassDetailsSpellsBox] = &Element{View: spellsView, ParentFlex: rowBelowGeneralStatsFlex}
	//
	// navigation
	//
	var currentElement *Element

	currentElement = elements[ClassDetailsGeneralStatsBox]

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
			} else if currentElement.ParentFlex == rowBelowGeneralStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}
		case 's', 'S':
			// switch focus between generalStatsView and the columnFlex (which contains epicView1 and epicView2)
			if currentElement == elements[ClassDetailsGeneralStatsBox] {
				currentElement = elements[ClassDetailsEpicStatsBox]
				app.SetFocus(currentElement.View)
			} else if currentElement.ParentFlex == rowBelowGeneralStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}
		case 'q', 'Q':
			app.Stop()
		case 'a', 'A':
			// switch focus between generalStatsView and the columnFlex (which contains epicView1 and epicView2)
			if currentElement.ParentFlex == rowBelowGeneralStatsFlex {
				switch currentElement {
				case elements[ClassDetailsEpicStatsBox]:
					currentElement = elements[ClassDetailsSpellsBox]
					app.SetFocus(currentElement.View)
				case elements[ClassDetailsSpellsBox]:
					currentElement = elements[ClassDetailsPrereqBox]
					app.SetFocus(currentElement.View)
				case elements[ClassDetailsPrereqBox]:
					currentElement = elements[ClassDetailsEpicStatsBox]
					app.SetFocus(currentElement.View)
				}
			}
		case 'd', 'D':
			if currentElement.ParentFlex == rowBelowGeneralStatsFlex {
				switch currentElement {
				case elements[ClassDetailsEpicStatsBox]:
					currentElement = elements[ClassDetailsPrereqBox]
					app.SetFocus(currentElement.View)
				case elements[ClassDetailsPrereqBox]:
					currentElement = elements[ClassDetailsSpellsBox]
					app.SetFocus(currentElement.View)
				case elements[ClassDetailsSpellsBox]:
					currentElement = elements[ClassDetailsEpicStatsBox]
					app.SetFocus(currentElement.View)
				}
			} else if currentElement.ParentFlex == generalStatsFlex {
				currentElement = elements[ClassDetailsGeneralStatsBox]
				app.SetFocus(currentElement.View)
			}

		}
		return event
	})
	return mainFlex

}

func classDetailFooter() *tview.TextView {

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetText("[yellow]Commands: \n").SetTextAlign(tview.AlignCenter)

	fmt.Fprintf(footer, "[white]↵ Enter: [red] Confirm Class \n ")
	fmt.Fprintf(footer, "[white]W A S D: [red]Change Focus ")
	fmt.Fprintf(footer, "[white]↑↓: [red]Scroll Up/Down in Focused Box ")
	fmt.Fprintf(footer, "[white]Q: [red]Quit \n")
	return footer
}

func classDetailGeneralStats(val reflect.Value, className sql.NullString) *tview.TextView {
	generalStatsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	generalStatsView.SetBorder(true)
	generalStatsView.SetTitle(fmt.Sprintf("Stat Sheet of [green] %s", className.String))
	typeOfS := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if !strings.Contains(typeOfS.Field(i).Name, "Epic") && !strings.Contains(typeOfS.Field(i).Name, "Req") && !strings.Contains(typeOfS.Field(i).Name, "FullText") {
			if reflect.TypeOf(val.Field(i).Interface()) == reflect.TypeOf(sql.NullString{}) {
				if val.Field(i).Interface().(sql.NullString).Valid {
					fmt.Fprintf(generalStatsView, "[white]%s: [blue]%s\n\n", typeOfS.Field(i).Name, val.Field(i).Interface().(sql.NullString).String)
				}
			} else if reflect.TypeOf(val.Field(i).Interface()) == reflect.TypeOf(sql.NullInt16{}) {
				if val.Field(i).Interface().(sql.NullInt16).Valid {
					fmt.Fprintf(generalStatsView, "[white]%s: [blue]%d\n\n", typeOfS.Field(i).Name, val.Field(i).Interface().(sql.NullInt16).Int16)
				}
			}
		} else if strings.Contains(typeOfS.Field(i).Name, "FullText") && !strings.Contains(typeOfS.Field(i).Name, "Epic") {
			fullTextPlain, err := html2text.FromString(val.Field(i).Interface().(sql.NullString).String)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(generalStatsView, "[white]%s: \n [skyblue]%s\n\n", typeOfS.Field(i).Name, fullTextPlain)
		}

	}
	return generalStatsView
}

func classDetailsEpicFeat(class db.ClassListColumns) *tview.TextView {
	epicView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	epicView.SetBorder(true)
	epicView.SetTitle("Epic Stat List")
	fmt.Fprintf(epicView, "[white]Epic Feat Base Level: [blue]%s\n\n", class.EpicFeatBaseLevel.String)
	fmt.Fprintf(epicView, "[white]Epic Feat Interval: [blue]%s\n\n", class.EpicFeatInterval.String)
	fmt.Fprintf(epicView, "[white]Epic Feat List: [blue]%s\n\n", class.EpicFeatList.String)
	epicFullTextPlain, err := html2text.FromString(class.EpicFullText.String)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(epicView, "[white]Epic Full Text: \n [skyblue]%s\n\n", epicFullTextPlain)
	return epicView
}

func classDetailsPrerequirements(val reflect.Value) *tview.TextView {
	preReqView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	preReqView.SetBorder(true)
	typeOfS := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if strings.Contains(typeOfS.Field(i).Name, "Req") {
			if reflect.TypeOf(val.Field(i).Interface()) == reflect.TypeOf(sql.NullString{}) {
				if val.Field(i).Interface().(sql.NullString).Valid {
					fmt.Fprintf(preReqView, "[white]%s: [blue]%s\n\n", typeOfS.Field(i).Name, val.Field(i).Interface().(sql.NullString).String)
				}
			}
		}
	}
	preReqView.SetTitle("Prerequisites")
	return preReqView
}

func classDetailsSpells(class db.ClassListColumns) *tview.TextView {
	spellsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	spellsView.SetBorder(true).SetTitle("Spell List")
	val := reflect.ValueOf(class)
	typeOfS := val.Type()
	for i := 0; i < val.NumField(); i++ {
		if strings.Contains(typeOfS.Field(i).Name, "SpellList") || strings.Contains(typeOfS.Field(i).Name, "SpellType") {
			if reflect.TypeOf(val.Field(i).Interface()) == reflect.TypeOf(sql.NullString{}) {
				if val.Field(i).Interface().(sql.NullString).Valid {
					fmt.Fprintf(spellsView, "[white]%s: [blue]%s\n\n", typeOfS.Field(i).Name, val.Field(i).Interface().(sql.NullString).String)
				}
			}
		}
	}
	return spellsView

}
