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

func init() {
	application.RegisterEvent[string]("time")
}

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
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		OnShutdown: func() {

		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "FreeAi",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "http://127.0.0.1:8787",
	})

	// Run the application. This blocks until the application has been exited.
	err = app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
