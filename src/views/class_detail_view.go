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
	elements := make([][]*Element, 3)
	elements[0] = []*Element{
		{View: generalStatsView, ParentFlex: generalStatsFlex, Position: struct{ Row, Column int }{0, 0}},
	}
	elements[1] = []*Element{
		{View: epicView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 0}},
		{View: preReqView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 1}},
		{View: spellsView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 2}},
	}
	elements[2] = []*Element{
		{View: footerView, ParentFlex: footbarFlex, Position: struct{ Row, Column int }{2, 0}},
	}
	//
	// navigation
	//
	currentElement := elements[0][0]

	mainFlex.SetInputCapture(TraverseBoxes(currentElement, elements, app))
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyHome:
			app.SetRoot(ReturnClassView(app), true)
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

	footer.SetBorder(true).SetBorderAttributes(tcell.AttrDim)

	fmt.Fprintf(footer, "[white]W A S D: [red]Change Focus [white]↑↓: [red]Scroll Up/Down in Focused Box \n")
	fmt.Fprintf(footer, "[white]↵ Enter: [red]Confirm Class [white]Q: [red]Quit [white]Home: [red]Previous Page \n")
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
