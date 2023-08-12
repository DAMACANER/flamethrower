package main

import (
	"flag"
	"flamethrower/src/engine"
	"flamethrower/src/views"
	"fmt"
	"os"

	"github.com/rivo/tview"
)

var DBLocation = flag.String("db", "dnd35.db", "Location of the database file")

func init() {
	flag.Parse()
	engine.InitDB(DBLocation)
}
func main() {
	app := tview.NewApplication().EnableMouse(true)
	app.SetRoot(views.ReturnClassView(app), true)
	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

}
