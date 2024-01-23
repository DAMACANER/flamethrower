package class

import (
	"flamethrower/src/db/model"
	"flamethrower/src/db/table"
	"flamethrower/src/engine"
	"flamethrower/src/handlers"
	"flamethrower/src/helpers"
	"flamethrower/src/types"
	"fmt"
	"strconv"
	"strings"
	"sync"

	. "github.com/go-jet/jet/v2/sqlite"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type LevelChangeEvent struct {
	Level int
}

var levelChangeChannel = make(chan LevelChangeEvent)
var levelChangeMutex = &sync.Mutex{}
var levelIncrementMutex = &sync.Mutex{}

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
	initialLevel := levels[0]
	// first part: general details
	generalDetailsFlex := returnGeneralClassDetailsFlex(initialLevel)
	// third part: spell slots
	slotsDetailFlex := returnSlotsFlex(initialLevel)
	// fourth part: spells known
	spellsKnownFlex := returnSpellsKnownFlex(initialLevel)
	// add all the parts to the class flex
	// it will be left to right in order
	classFlex.AddItem(generalDetailsFlex, 0, 1, false)
	classFlex.AddItem(slotsDetailFlex, 0, 1, false)
	classFlex.AddItem(spellsKnownFlex, 0, 1, false)
	// put up the input field
	inputFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	inputTextField := returnLevelInputBox(app, len(levels))
	inputFlex.AddItem(inputTextField, 0, 1, false)
	// get the footbar flex for commands
	footbarFlex := classLevelDetailFooter()

	// add all the components to the main flex
	mainFlex.AddItem(inputFlex, 3, 0, false)
	mainFlex.AddItem(classFlex, 0, 1, false)
	mainFlex.AddItem(footbarFlex, 5, 0, false)
	go eventHandler(app, levels, &generalDetailsFlex, &slotsDetailFlex, &spellsKnownFlex)
	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q', 'Q':
			app.Stop()
		}
		switch event.Key() {
		case tcell.KeyEnter:
			// get the level from the input field
		case tcell.KeyEsc:
			app.SetRoot(ReturnClassView(app), true)
		}
		switch event.Key() {
		case tcell.KeyUp:
			// do not handle the updating of the detail boxes here
			// or do not handle the app drawing
			//
			// just increment the level and event handler of input
			// text field will handle the rest
			level, err := strconv.Atoi(inputTextField.GetText())
			if inputTextField.GetText() == "" {
				// initial input
				level = 1
			} else {
				handlers.Error(err, app)
			}
			levelIncrementMutex.Lock()
			inputTextField.SetText(strconv.Itoa(level + 1))
			levelIncrementMutex.Unlock()
		case tcell.KeyDown:
			// same as above
			level, err := strconv.Atoi(inputTextField.GetText())
			if inputTextField.GetText() == "" {
				level = 20
			} else {
				handlers.Error(err, app)
			}
			levelIncrementMutex.Lock()
			inputTextField.SetText(strconv.Itoa(level - 1))
			levelIncrementMutex.Unlock()
		}
		return event
	}) // return the centered main flex
	return mainFlex
}

func classLevelDetailFooter() *tview.TextView {

	footer := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetText("[yellow]Commands: \n").SetTextAlign(tview.AlignCenter)

	footer.SetBorder(true).SetBorderAttributes(tcell.AttrDim)
	fmt.Fprintf(footer, "[white]Up/Down Keys: [red] See Details of Next/Previous Level\n")
	fmt.Fprintf(footer, "[white]↵ Enter: [red]Confirm Class [white]Q: [red]Quit [white]Esc: [red]Previous Page \n")
	return footer
}

func fetchClassLevelData(app *tview.Application, class model.Class) []model.ClassLevelTable {
	stmt := SELECT(table.ClassLevelTable.AllColumns).
		FROM(table.ClassLevelTable.INNER_JOIN(table.Class, table.ClassLevelTable.Name.EQ(String(class.Name)))).
		WHERE(table.ClassLevelTable.Name.EQ(String(class.Name))).
		ORDER_BY(table.ClassLevelTable.Level.ASC())
	var data []model.ClassLevelTable
	err := stmt.Query(engine.DB, &data)
	handlers.Error(err, app)
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
		SetLabel(fmt.Sprintf("Enter a level (or press Up / Down arrow keys) for seeing class details (min 1 / max %d):", max_levels)).
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
		levelChangeChannel <- LevelChangeEvent{Level: 1}
		return
	}

	if len(text) > 0 && len(text) < 3 {
		level, err := strconv.Atoi(text)
		if err != nil {
			handlers.Error(err, app)
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
func eventHandler(app *tview.Application, levels []model.ClassLevelTable, generalDetails **tview.Flex, slotDetails **tview.Flex, spellsKnown **tview.Flex) {
	for v := range levelChangeChannel {
		levelChangeMutex.Lock()
		app.QueueUpdateDraw(func() {
			level := levels[v.Level-1]
			updateGeneralDetails(level, *generalDetails)
			updateSlotDetails(level, *slotDetails)
			updateSpellsKnownDetails(level, *spellsKnown)
		})
		levelChangeMutex.Unlock()
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

func updateSpellsKnownDetails(level model.ClassLevelTable, spellsKnownFlex *tview.Flex) {
	// Clear existing items
	(*spellsKnownFlex).Clear()
	(*spellsKnownFlex).AddItem(returnSpellsKnownFlex(level), 0, 1, false)
}

func returnGeneralClassDetailsFlex(level model.ClassLevelTable) *tview.Flex {
	generalDetailsItems := []types.ListItem{
		{Label: "", Secondary: "", Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Name: [skyblue]%s", level.Name), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Level: [skyblue]%s", *level.Level), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Base Attack Bonus:: [skyblue]%s", *level.BaseAttackBonus), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("AC Bonus: [skyblue]%s", *level.AcBonus), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Fort Save: [skyblue]%s", *level.FortSave), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Ref Save: [skyblue]%s", *level.RefSave), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Will Save: [skyblue]%s", *level.WillSave), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Specials: \n %s", formatSpecial(level.Special)), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Points Per Day: [skyblue]%s", *level.PointsPerDay), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Powers Known: [skyblue]%s", *level.PowersKnown), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Unarmed Damage: [skyblue]%s", *level.UnarmedDamage), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Flurry of Blows: [skyblue]%s", *level.FlurryOfBlows), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("Unarmored Speed Bonus: [skyblue]%s", *level.UnarmoredSpeedBonus), Shortcut: 0, Color: tcell.ColorWhite},
	}

	options := types.ListAppearanceOptions{
		ShouldDrawBorder:        true,
		ShouldShowSecondaryText: true,
		BorderColor:             tcell.ColorLightBlue,
		BackgroundColor:         tcell.ColorBlack,
		TitleAlignment:          tview.AlignLeft,
		Title:                   "General Details",
		ListDirection:           tview.FlexColumn,
		BorderPadding: types.BorderPadding{
			Top:    1,
			Bottom: 1,
			Left:   2,
			Right:  2,
		},
	}

	return helpers.CreateListFlex("General Details", generalDetailsItems, options)
}

func returnSlotsFlex(level model.ClassLevelTable) *tview.Flex {
	slotItems := []types.ListItem{
		{Label: "", Secondary: "", Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("0: [skyblue]%s", *level.Slots0), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("1: [skyblue]%s", *level.Slots1), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("2: [skyblue]%s", *level.Slots2), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("3: [skyblue]%s", *level.Slots3), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("4: [skyblue]%s", *level.Slots4), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("5: [skyblue]%s", *level.Slots5), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("6: [skyblue]%s", *level.Slots6), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("7: [skyblue]%s", *level.Slots7), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("8: [skyblue]%s", *level.Slots8), Shortcut: 0, Color: tcell.ColorWhite},
		{Label: fmt.Sprintf("9: [skyblue]%s", *level.Slots9), Shortcut: 0, Color: tcell.ColorWhite},
	}

	options := types.ListAppearanceOptions{
		ShouldDrawBorder:        true,
		ShouldShowSecondaryText: false,
		BorderColor:             tcell.ColorLightBlue,
		BackgroundColor:         tcell.ColorBlack,
		TitleAlignment:          tview.AlignLeft,
		Title:                   "Slots",
		ListDirection:           tview.FlexColumn,
		BorderPadding: types.BorderPadding{
			Top:    1,
			Bottom: 1,
			Left:   2,
			Right:  2,
		},
	}

	return helpers.CreateListFlex("Slots", slotItems, options)
}
func returnSpellsKnownFlex(level model.ClassLevelTable) *tview.Flex {
	return helpers.CreateListFlex(
		"Spells Known",
		[]types.ListItem{
			{Label: "", Secondary: "", Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("0: [skyblue]%s", *level.SpellsKnown0), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("1: [skyblue]%s", *level.SpellsKnown1), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("2: [skyblue]%s", *level.SpellsKnown2), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("3: [skyblue]%s", *level.SpellsKnown3), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("4: [skyblue]%s", *level.SpellsKnown4), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("5: [skyblue]%s", *level.SpellsKnown5), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("6: [skyblue]%s", *level.SpellsKnown6), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("7: [skyblue]%s", *level.SpellsKnown7), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("8: [skyblue]%s", *level.SpellsKnown8), Shortcut: 0, Color: tcell.ColorWhite},
			{Label: fmt.Sprintf("9: [skyblue]%s", *level.SpellsKnown9), Shortcut: 0, Color: tcell.ColorWhite},
		},
		types.ListAppearanceOptions{
			ShouldDrawBorder:        true,
			ShouldShowSecondaryText: false,
			BorderColor:             tcell.ColorLightBlue,
			BackgroundColor:         tcell.ColorBlack,
			TitleAlignment:          tview.AlignLeft,
			Title:                   "Spells Known",
			ListDirection:           tview.FlexColumn,
			BorderPadding: types.BorderPadding{
				Top:    1,
				Bottom: 1,
				Left:   2,
				Right:  2,
			},
		})
}

func formatSpecial(special *string) string {
	if !helpers.NotEmpty(special) {
		return ""
	}

	spec := strings.Split(*special, ",")
	if len(spec) == 0 {
		return fmt.Sprintf("[skyblue]%s", strings.TrimSpace(*special))
	}

	var formattedSpecial []string
	for i, s := range spec {
		formattedSpecial = append(formattedSpecial, fmt.Sprintf("Special %d: [skyblue]%s", i+1, strings.TrimSpace(s)))
	}

	return strings.Join(formattedSpecial, "\n")
}
