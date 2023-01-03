package main

import (
	"fmt"
	"github.com/mbpolan/lull/internal/parsers"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/ui"
	"github.com/rivo/tview"
	"os"
)

func main() {
	// initialize supporting modules
	parsers.Setup()

	// determine the user's home directory to save app state file in
	stateSavePath, err := os.UserHomeDir()
	if err != nil {
		stateSavePath = "."
	}

	stateSavePath = fmt.Sprintf("%s/.lull", stateSavePath)

	// attempt to read existing app state from file
	var st *state.AppState
	initialSave := true
	data, err := os.ReadFile(stateSavePath)
	if err != nil {
		st = state.NewAppState()
	} else {
		st, err = state.DeserializeAppState(data)
		if err != nil {
			st = state.NewAppState()
		} else {
			initialSave = false
		}
	}

	// create a state manager and flag the state as dirty to force an initial save if needed
	stateManager := state.NewStateManager(st, stateSavePath)
	if initialSave {
		stateManager.SetDirty()
	}

	app := tview.NewApplication()
	root := ui.NewRoot(app, stateManager)

	app.SetRoot(root.Widget(), true)
	app.SetFocus(root.Widget())

	if err := app.Run(); err != nil {
		panic(err)
	}

	// save the app state to file
	if err := stateManager.Shutdown(); err != nil {
		fmt.Printf("Failed to save data: %+v\n", err)
	}
}
