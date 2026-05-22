package main

import (
	"embed"
	_ "embed"
	"freeai-app/backend"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist/freeai-web/browser
var assets embed.FS

func main() {
	go backend.Start()
	apiMiddleware, err := backend.NewAPIMiddleware()

	app := application.New(application.Options{
		Name:        "FreeAi",
		Description: "FreeAI desktop application",
		Services:    []application.Service{},
		Assets: application.AssetOptions{
			Handler:    application.AssetFileServerFS(assets),
			Middleware: application.ChainMiddleware(apiMiddleware),
		},
		Server: application.ServerOptions{
			Host: "127.0.0.1",
			Port: 8787,
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		OnShutdown: func() {

		},
	})

	setupApplicationMenu(app)

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:              "FreeAi",
		Width:              1280,
		Height:             860,
		UseApplicationMenu: true,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "http://127.0.0.1:8787",
	})

	setupSystemTray(app, window)
	setupApplicationEvents(app, window)

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
