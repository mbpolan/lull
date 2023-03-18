package main

import (
	"flag"
	"fmt"
	"github.com/mbpolan/lull/internal/events"
	"github.com/mbpolan/lull/internal/logger"
	"github.com/mbpolan/lull/internal/parsers"
	"github.com/mbpolan/lull/internal/state"
	"github.com/mbpolan/lull/internal/system"
	"github.com/mbpolan/lull/internal/ui"
	"github.com/mbpolan/lull/internal/util"
	"github.com/rivo/tview"
	"os"
	"time"
)

var (
	version = "dev"
	commit  = "local"
	date    = time.Now().Format(time.RFC3339)
)

func main() {
	var logLevel string
	flag.StringVar(&logLevel, "log-level", "error", "sets the verbosity for logging (debug, info, error)")
	flag.Parse()

	// initialize supporting modules
	if err := system.Setup(); err != nil {
		fmt.Printf("Could not initialize system: %s\n", err)
		os.Exit(1)
	}

	if err := logger.Setup(logLevel); err != nil {
		fmt.Printf("Could not initialize logging: %s\n", err)
		os.Exit(1)
	}

	events.Setup()
	parsers.Setup()

	// populate build information
	buildMeta := util.NewBuildMeta(version, commit, date)

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
		logger.Infof("initializing new app state due to error: %s", err)
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
	root := ui.NewRoot(app, stateManager, buildMeta)

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
