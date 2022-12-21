package main

import (
	"github.com/mbpolan/lull/internal/ui"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	root := ui.NewRoot(app)

	app.SetRoot(root.Widget(), true)
	app.SetFocus(root.Widget())

	if err := app.Run(); err != nil {
		panic(err)
	}
}
