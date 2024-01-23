package handlers

import (
	"log"

	"github.com/rivo/tview"
)

// func returnErrorView(app *tview.Application, errorMessage string) *tview.Flex {
// 	modal := tview.NewModal().
// 		SetText(fmt.Sprintf("[red]%s \n\n Go back to [white]class selection?", errorMessage)).
// 		AddButtons([]string{"OK"}).
// 		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 			if buttonLabel == "OK" {
// 				app.SetRoot(ReturnClassView(app), true)
// 			}
// 		})
// 	modalFlex := tview.NewFlex().
// 		AddItem(modal, 0, 1, false)
// 	app.SetFocus(modal.Box)
// 	return modalFlex
// }

func Error(err error, app *tview.Application) {
	if err != nil {
		// TODO
		// app.SetRoot(returnErrorView(app, err.Error()), true)
		log.Fatal(err)
	}
}
