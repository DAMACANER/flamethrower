package main

import (
	"flamethrower/src/db"
	"flamethrower/src/views"
	"fmt"
	"log"
	"os"

	"github.com/blockloop/scan"
	"github.com/rivo/tview"
)

func init() {
	db.InitDB("dnd35.db")
	repo := &db.ClassListRepo{BaseRepo: &db.BaseRepo{}}
	rows, err := repo.Count().Query()
	if err != nil {
		log.Fatal(err)
	}
	err = scan.Row(&db.TotalClassCount, rows)
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	app := tview.NewApplication().EnableMouse(true)
	app.SetRoot(views.ReturnClassView(app), true)
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

}
