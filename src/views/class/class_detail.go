package class

import (
	"flamethrower/src/db/model"
	"flamethrower/src/helpers"
	"flamethrower/src/types"
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
	elements := make([][]*types.Element, 3)
	elements[0] = []*types.Element{
		{View: generalStatsView, ParentFlex: generalStatsFlex, Position: struct{ Row, Column int }{0, 0}},
	}
	elements[1] = []*types.Element{
		{View: epicView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 0}},
		{View: preReqView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 1}},
		{View: spellsView, ParentFlex: rowBelowGeneralStatsFlex, Position: struct{ Row, Column int }{1, 2}},
	}
	elements[2] = []*types.Element{
		{View: footerView, ParentFlex: footbarFlex, Position: struct{ Row, Column int }{2, 0}},
	}
	//
	// navigation
	//
	currentElement := elements[0][0]

	mainFlex.SetInputCapture(helpers.TraverseBoxes(currentElement, elements, app))
	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		switch event.Key() {
		case tcell.KeyEsc:
			app.SetRoot(ReturnClassView(app), true)
		case tcell.KeyEnter:
			app.SetRoot(ReturnClassLevelDetailView(class, app), true)
		}
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()

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

	fmt.Fprintf(footer, "[white]W A S D (or mouse click): [red]Change Focus [white]↑↓: [red]Scroll Up/Down in Focused Box \n")
	fmt.Fprintf(footer, "[white]↵ Enter: [red]Confirm Class [white]Q: [red]Quit [white]Esc: [red]Previous Page \n")
	return footer
}

func classDetailGeneralStats(class model.Class) *tview.TextView {
	lines := []types.TextViewLines{
		{Line: fmt.Sprintf("[white]Hit Die [skyblue]%s", *class.HitDie)},
		{Line: fmt.Sprintf("[white]Skill Points [skyblue]%s", *class.HitDie)},
		{Line: fmt.Sprintf("[white]Skill Points Ability [skyblue]%s", *class.SkillPoints)},
		{Line: fmt.Sprintf("[white]Class Skills [skyblue]%s", *class.SkillPoints)},
		{Line: fmt.Sprintf("[white]Weapon and Armor Proficiencies [skyblue]%s", tview.WordWrap(*class.Proficiencies, 50))},
		{Line: fmt.Sprintf("[white]Alignment [skyblue]%s", *class.Alignment)},
		{Line: fmt.Sprintf("[white]Source [skyblue]%s", *class.Reference)}}
	lines = append(lines, types.TextViewLines{Line: "[white]Specials \n"})
	profs := tview.WordWrap(*class.Proficiencies, 160)

	for _, prof := range profs {
		lines = append(lines, types.TextViewLines{Line: fmt.Sprintf("[skyblue]%s", prof)})
	}
	return helpers.CreateTextView(
		fmt.Sprintf("Stat Sheet of [green] %s", class.Name),
		lines,
		types.TextViewAppearanceOptions{
			EnableDynamicColors: true,
			SetRegions:          true,
			SetBorder:           true,
			SetWordWrap:         true,
		})
}
func classDetailsEpicFeat(class model.Class) *tview.TextView {
	return helpers.CreateTextView(
		"Epic Stat List",
		[]types.TextViewLines{
			{Line: fmt.Sprintf("[white]Epic Feat Base Level: [skyblue]%v\n\n", *class.EpicFeatBaseLevel)},
			{Line: fmt.Sprintf("[white]Epic Feat Interval: [skyblue]%v\n\n", *class.EpicFeatInterval)},
			{Line: fmt.Sprintf("[white]Epic Feat List: [skyblue]%v\n\n", *class.EpicFeatList)},
			{Line: *class.EpicFullText, ShouldFormatHTML: true}},
		types.TextViewAppearanceOptions{
			EnableDynamicColors: true,
			SetRegions:          true,
			SetWrap:             true,
			SetWordWrap:         true,
			SetBorder:           true,
		})
}

func classDetailsPrerequirements(class model.Class) *tview.TextView {
	return helpers.CreateTextView(
		"Prerequisities",
		[]types.TextViewLines{
			{Line: fmt.Sprintf("[white]Required Race: [skyblue]%v\n\n", *class.ReqRace)},
			{Line: fmt.Sprintf("[white]Required Skill: [skyblue]%v\n\n", *class.ReqSkill)},
			{Line: fmt.Sprintf("[white]Required Weapon Proficiency: [skyblue]%v\n\n", *class.ReqWeaponProficiency)},
			{Line: fmt.Sprintf("[white]Required Base Attack Bonus: [skyblue]%v\n\n", *class.ReqBaseAttackBonus)},
			{Line: fmt.Sprintf("[white]Required Psionics: [skyblue]%v\n\n", *class.ReqPsionics)},
			{Line: fmt.Sprintf("[white]Required Spells: [skyblue]%v\n\n", *class.ReqSpells)},
			{Line: fmt.Sprintf("[white]Required Epic Feat: [skyblue]%v\n\n", *class.ReqEpicFeat)},
			{Line: fmt.Sprintf("[white]Required Special: [skyblue]%v\n\n", *class.ReqSpecial)},
			{Line: fmt.Sprintf("[white]Required Psionic Power: [skyblue]%v\n\n", *class.ReqPsionics)}},
		types.TextViewAppearanceOptions{
			EnableDynamicColors: true,
			SetWrap:             true,
			SetWordWrap:         true,
			SetRegions:          true,
			SetBorder:           true,
		})

}

func classDetailsSpells(class model.Class) *tview.TextView {
	return helpers.CreateTextView(
		"Spell List",
		[]types.TextViewLines{
			{Line: fmt.Sprintf("[white]Spell Stats: [skyblue]%v\n\n", *class.SpellStat)},
			{Line: fmt.Sprintf("[white]Spell Type: [skyblue]%v\n\n", *class.SpellType)},
			{Line: fmt.Sprintf("[white]Spells List - 1: [skyblue]%v\n\n", *class.SpellList1)},
			{Line: fmt.Sprintf("[white]Spells List - 2: [skyblue]%v\n\n", *class.SpellList2)},
			{Line: fmt.Sprintf("[white]Spells List - 3: [skyblue]%v\n\n", *class.SpellList3)},
			{Line: fmt.Sprintf("[white]Spells List - 4: [skyblue]%v\n\n", *class.SpellList4)},
			{Line: fmt.Sprintf("[white]Spells List - 5: [skyblue]%v\n\n", *class.SpellList5)}},
		types.TextViewAppearanceOptions{
			EnableDynamicColors: true,
			SetWrap:             true,
			SetWordWrap:         true,
			SetRegions:          true,
			SetBorder:           true,
		})

}
