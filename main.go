package main

import (
	"context"
	"freeai-app/backend"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	if err := backend.StartAndWait(context.Background()); err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name:        "FreeAi",
		Description: "FreeAI desktop application",
		Services:    []application.Service{},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
	})

	setupApplicationMenu(app)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:              "FreeAi",
		Width:              1280,
		Height:             860,
		UseApplicationMenu: true,
		DevToolsEnabled:    true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              backend.BackendBaseURL,
	})

	setupSystemTray(app, window)
	setupApplicationEvents(app, window)

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
