package views

import (
	"flamethrower/src/db/model"
	"flamethrower/src/db/table"
	"flamethrower/src/engine"
	"fmt"
	"strconv"
	"strings"

	. "github.com/go-jet/jet/v2/sqlite"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LevelChangeEvent struct {
	Level int
}

var levelChangeChannel = make(chan LevelChangeEvent)

func ReturnClassLevelDetailView(class model.Class, app *tview.Application) *tview.Flex {
	// set up selected class levels
	// join the classleveltable with class based on the class name
	//
	// we will not match based on the id
	levels := fetchClassLevelData(app, class)
	// set up main flex
	mainFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	// put a class field in the main flex
	// it will show first level details initially
	classFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	// first part: general details
	generalDetailsFlex := returnGeneralClassDetailsFlex(levels[0])
	// third part: spell slots
	slotsDetailFlex := returnSlotsFlex(levels[0])
	// fourth part: spells known
	spellsKnownFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	// add all the parts to the class flex
	// it will be left to right in order
	classFlex.AddItem(generalDetailsFlex, 0, 1, false)
	classFlex.AddItem(slotsDetailFlex, 0, 1, false)
	classFlex.AddItem(spellsKnownFlex, 0, 1, false)
	// put up the input field
	inputFlex := returnLevelInputBox(app, len(levels))
	// get the footbar flex for commands
	footbarFlex := classLevelDetailFooter()

	// add all the components to the main flex
	mainFlex.AddItem(inputFlex, 3, 0, false)
	mainFlex.AddItem(classFlex, 0, 1, false)
	mainFlex.AddItem(footbarFlex, 5, 0, false)

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ch := event.Rune()
		if ch == 'q' || ch == 'Q' {
			app.Stop()
		}
		return event
	})
	go eventHandler(app, levels, &generalDetailsFlex, &slotsDetailFlex)
	// return the centered main flex
	return mainFlex
}

func classLevelDetailFooter() *tview.TextView {

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetText("[yellow]Commands: \n").SetTextAlign(tview.AlignCenter)

	footer.SetBorder(true).SetBorderAttributes(tcell.AttrDim)
	fmt.Fprintf(footer, "[white]↵ Enter: [red]Confirm Class [white]Q: [red]Quit [white]Home: [red]Previous Page \n")
	return footer
}

func fetchClassLevelData(app *tview.Application, class model.Class) []model.ClassLevelTable {
	stmt := SELECT(table.ClassLevelTable.AllColumns).
		FROM(table.ClassLevelTable.INNER_JOIN(table.Class, table.ClassLevelTable.Name.EQ(String(class.Name)))).
		WHERE(table.ClassLevelTable.Name.EQ(String(class.Name))).
		ORDER_BY(table.ClassLevelTable.Level.ASC())
	var data []model.ClassLevelTable
	err := stmt.Query(engine.DB, &data)
	HandleError(err, app)
	return data
}

func returnLevelInputBox(app *tview.Application, max_levels int) *tview.InputField {
	inputField := createInputField(max_levels)

	// add a callback for when the user writes anything
	inputField.SetChangedFunc(func(text string) {
		handleInputTextChanged(text, app, inputField, max_levels)
	})

	return inputField
}

func createInputField(max_levels int) *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(fmt.Sprintf("Enter a level for seeing class details (min 1 / max %d):", max_levels)).
		SetFieldWidth(2).
		SetAcceptanceFunc(tview.InputFieldInteger)
	inputField.SetTitleAlign(tview.AlignCenter)
	inputField.SetBorderColor(tcell.ColorLightBlue)
	inputField.SetBorderAttributes(tcell.AttrBold)
	inputField.SetBorder(true)
	return inputField
}

func handleInputTextChanged(text string, app *tview.Application, inputField *tview.InputField, max_levels int) {
	if len(text) >= 3 {
		inputField.SetText("")
		return
	}

	if len(text) > 0 && len(text) < 3 {
		level, err := strconv.Atoi(text)
		if err != nil {
			HandleError(err, app)
		}

		if level < 1 || level > max_levels {
			inputField.SetText("")
			levelChangeChannel <- LevelChangeEvent{Level: 1}
			return
		}

		// Send the event instead of directly updating
		levelChangeChannel <- LevelChangeEvent{Level: level}
	}
}
func eventHandler(app *tview.Application, levels []model.ClassLevelTable, generalDetails **tview.Flex, slotDetails **tview.Flex) {
	for {
		select {
		case event := <-levelChangeChannel:
			app.QueueUpdateDraw(func() {
				updateGeneralDetails(levels[event.Level-1], *generalDetails)
				updateSlotDetails(levels[event.Level-1], *slotDetails)
			})
		}
	}
}

func updateGeneralDetails(level model.ClassLevelTable, generalDetails *tview.Flex) {
	// Clear existing items
	(*generalDetails).Clear()
	(*generalDetails).AddItem(returnGeneralClassDetailsFlex(level), 0, 1, false)
}

func updateSlotDetails(level model.ClassLevelTable, slotsFlex *tview.Flex) {
	// Clear existing items
	(*slotsFlex).Clear()
	(*slotsFlex).AddItem(returnSlotsFlex(level), 0, 1, false)
}

func returnGeneralClassDetailsFlex(level model.ClassLevelTable) *tview.Flex {
	generalDetailsFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	generalDetailsList := tview.NewList()
	generalDetailsList.ShowSecondaryText(false)
	generalDetailsList.SetBorder(true)
	generalDetailsList.SetTitle("General Details")
	generalDetailsList.SetTitleAlign(tview.AlignLeft)
	generalDetailsList.SetBorderColor(tcell.ColorLightBlue)
	generalDetailsList.SetBorderPadding(1, 1, 2, 2)
	generalDetailsList.SetBackgroundColor(tcell.ColorBlack)
	generalDetailsList.AddItem("", "", 0, nil)
	generalDetailsList.AddItem(fmt.Sprintf("Name: [skyblue]%s", level.Name), "", 0, nil)
	if notEmpty(level.Level) {
		generalDetailsList.AddItem(fmt.Sprintf("Level: [skyblue]%s", *level.Level), "", 0, nil)
	}
	if notEmpty(level.BaseAttackBonus) {
		generalDetailsList.AddItem(fmt.Sprintf("Base Attack Bonus: [skyblue]%s", *level.BaseAttackBonus), "", 0, nil)
	}
	if notEmpty(level.AcBonus) {
		generalDetailsList.AddItem(fmt.Sprintf("AC Bonus: [skyblue]%s", *level.AcBonus), "", 0, nil)
	}
	if notEmpty(level.FortSave) {
		generalDetailsList.AddItem(fmt.Sprintf("Fort Save: [skyblue]%s", *level.FortSave), "", 0, nil)
	}
	if notEmpty(level.RefSave) {
		generalDetailsList.AddItem(fmt.Sprintf("Ref Save: [skyblue]%s", *level.RefSave), "", 0, nil)
	}
	if notEmpty(level.WillSave) {
		generalDetailsList.AddItem(fmt.Sprintf("Will Save: [skyblue]%s", *level.WillSave), "", 0, nil)
	}
	if notEmpty(level.Special) {
		// split special by , and make it a list
		spec := strings.Split(*level.Special, ",")
		for i, s := range spec {
			generalDetailsList.AddItem(fmt.Sprintf("Special %d: [skyblue]%s", i+1, strings.TrimSpace(s)), "", 0, nil)
		}
	}
	if notEmpty(level.PointsPerDay) {
		generalDetailsList.AddItem(fmt.Sprintf("Points Per Day: [skyblue]%s", *level.PointsPerDay), "", 0, nil)
	}
	if notEmpty(level.PowersKnown) {
		generalDetailsList.AddItem(fmt.Sprintf("Powers Known: [skyblue]%s", *level.PowersKnown), "", 0, nil)
	}
	if notEmpty(level.UnarmedDamage) {
		generalDetailsList.AddItem(fmt.Sprintf("Unarmed Damage: [skyblue]%s", *level.UnarmedDamage), "", 0, nil)
	}
	if notEmpty(level.FlurryOfBlows) {
		generalDetailsList.AddItem(fmt.Sprintf("Flurry of Blows: [skyblue]%s", *level.FlurryOfBlows), "", 0, nil)
	}
	if notEmpty(level.UnarmoredSpeedBonus) {
		generalDetailsList.AddItem(fmt.Sprintf("Unarmored Speed Bonus: [skyblue]%s", *level.UnarmoredSpeedBonus), "", 0, nil)
	}
	generalDetailsFlex.AddItem(generalDetailsList, 0, 1, false)
	return generalDetailsFlex
}

func returnSlotsFlex(level model.ClassLevelTable) *tview.Flex {
	slotFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	slotList := tview.NewList()
	slotList.ShowSecondaryText(false)
	slotList.SetBorder(true)
	slotList.SetTitle("Slots")
	slotList.SetTitleAlign(tview.AlignLeft)
	slotList.SetBorderColor(tcell.ColorLightBlue)
	slotList.SetBorderPadding(1, 1, 2, 2)
	slotList.SetBackgroundColor(tcell.ColorBlack)
	slotList.AddItem("", "", 0, nil)
	if notEmpty(level.Slots0) {
		slotList.AddItem(fmt.Sprintf("0: [skyblue]%s", *level.Slots0), "", 0, nil)
	}
	if notEmpty(level.Slots1) {
		slotList.AddItem(fmt.Sprintf("1: [skyblue]%s", *level.Slots1), "", 0, nil)
	}
	if notEmpty(level.Slots2) {
		slotList.AddItem(fmt.Sprintf("2: [skyblue]%s", *level.Slots2), "", 0, nil)
	}
	if notEmpty(level.Slots3) {
		slotList.AddItem(fmt.Sprintf("3: [skyblue]%s", *level.Slots3), "", 0, nil)
	}
	if notEmpty(level.Slots4) {
		slotList.AddItem(fmt.Sprintf("4: [skyblue]%s", *level.Slots4), "", 0, nil)
	}
	if notEmpty(level.Slots5) {
		slotList.AddItem(fmt.Sprintf("5: [skyblue]%s", *level.Slots5), "", 0, nil)
	}
	if notEmpty(level.Slots6) {
		slotList.AddItem(fmt.Sprintf("6: [skyblue]%s", *level.Slots6), "", 0, nil)
	}
	if notEmpty(level.Slots7) {
		slotList.AddItem(fmt.Sprintf("7: [skyblue]%s", *level.Slots7), "", 0, nil)
	}
	if notEmpty(level.Slots8) {
		slotList.AddItem(fmt.Sprintf("8: [skyblue]%s", *level.Slots8), "", 0, nil)
	}
	if notEmpty(level.Slots9) {
		slotList.AddItem(fmt.Sprintf("9: %s", *level.Slots9), "", 0, nil)
	}
	slotFlex.AddItem(slotList, 0, 1, false)
	return slotFlex
}
