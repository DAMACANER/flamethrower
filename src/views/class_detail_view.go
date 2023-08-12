package views

import (
	"flamethrower/src/db/model"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"jaytaylor.com/html2text"
)

func ReturnClassDetailView(class model.Class, app *tview.Application) *tview.Flex {
	//
	// views
	//
	epicView := classDetailsEpicFeat(class)
	generalStatsView := classDetailGeneralStats(class)
	footerView := classDetailFooter()
	preReqView := classDetailsPrerequirements(class)
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

func classDetailGeneralStats(class model.Class) *tview.TextView {
	generalStatsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	generalStatsView.SetBorder(true)
	generalStatsView.SetTitle(fmt.Sprintf("Stat Sheet of [green] %s", class.Name))
	fmt.Fprintf(generalStatsView, "[white]Name: [skyblue]%v\n\n", class.Name)
	fmt.Fprintf(generalStatsView, "[white]Hit Die: [skyblue]%v\n\n", *class.HitDie)
	fmt.Fprintf(generalStatsView, "[white]Skill Points: [skyblue]%v\n\n", *class.SkillPoints)
	fmt.Fprintf(generalStatsView, "[white]Skill Points Ability: [skyblue]%v\n\n", *class.SkillPointsAbility)
	fmt.Fprintf(generalStatsView, "[white]Class Skills: [skyblue]%v\n\n", *class.ClassSkills)
	fmt.Fprintf(generalStatsView, "[white]Weapon and Armor Proficiencies: [skyblue]%v\n\n", *class.Proficiencies)
	fmt.Fprintf(generalStatsView, "[white]Alignment: [skyblue]%v\n\n", *class.Alignment)
	fmt.Fprintf(generalStatsView, "[white]Source: [skyblue]%v\n\n", *class.Reference)

	return generalStatsView
}

func classDetailsEpicFeat(class model.Class) *tview.TextView {
	epicView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	epicView.SetBorder(true)
	epicView.SetTitle("Epic Stat List")
	fmt.Fprintf(epicView, "[white]Epic Feat Base Level: [skyblue]%v\n\n", *class.EpicFeatBaseLevel)
	fmt.Fprintf(epicView, "[white]Epic Feat Interval: [skyblue]%v\n\n", *class.EpicFeatInterval)
	fmt.Fprintf(epicView, "[white]Epic Feat List: [skyblue]%v\n\n", *class.EpicFeatList)
	epicFullTextPlain, err := html2text.FromString(fmt.Sprintf("%v", *class.EpicFullText))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(epicView, "[white]Epic Full Text: \n [skyblue]%s\n\n", epicFullTextPlain)
	return epicView
}

func classDetailsPrerequirements(class model.Class) *tview.TextView {
	preReqView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	preReqView.SetBorder(true)
	preReqView.SetTitle("Prerequisites")
	fmt.Fprintf(preReqView, "[white]Required Race: [skyblue]%v\n\n", *class.ReqRace)
	fmt.Fprintf(preReqView, "[white]Required Feat: [skyblue]%v\n\n", *class.ReqFeat)
	fmt.Fprintf(preReqView, "[white]Required Skill: [skyblue]%v\n\n", *class.ReqSkill)
	fmt.Fprintf(preReqView, "[white]Required Weapon Proficiency: [skyblue]%v\n\n", *class.ReqWeaponProficiency)
	fmt.Fprintf(preReqView, "[white]Required Base Attack Bonus: [skyblue]%v\n\n", *class.ReqBaseAttackBonus)
	fmt.Fprintf(preReqView, "[white]Required Psionics: [skyblue]%v\n\n", *class.ReqPsionics)
	fmt.Fprintf(preReqView, "[white]Required Spells: [skyblue]%v\n\n", *class.ReqSpells)
	return preReqView
}

func classDetailsSpells(class model.Class) *tview.TextView {
	spellsView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetWordWrap(true)
	spellsView.SetBorder(true).SetTitle("Spell List")
	fmt.Fprintf(spellsView, "[white]Spell Statistics: [skyblue]%v\n\n", *class.SpellStat)
	fmt.Fprintf(spellsView, "[white]Spells List - 1: [skyblue]%v\n\n", *class.SpellList1)
	fmt.Fprintf(spellsView, "[white]Spells List - 2: [skyblue]%v\n\n", *class.SpellList2)
	fmt.Fprintf(spellsView, "[white]Spells List - 3: [skyblue]%v\n\n", *class.SpellList3)
	fmt.Fprintf(spellsView, "[white]Spells List - 4: [skyblue]%v\n\n", *class.SpellList4)
	fmt.Fprintf(spellsView, "[white]Spells List - 5: [skyblue]%v\n\n", *class.SpellList5)
	return spellsView

}
