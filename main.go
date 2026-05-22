package main

import (
	"embed"
	_ "embed"
	"freeai-app/backend"
	"log"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed all:frontend/dist/freeai-web/browser
var assets embed.FS

//go:embed build/tray/freeai-template.png
var trayIcon []byte

var quitting atomic.Bool

const (
	shortcutShow   = "CmdOrCtrl+Shift+F"
	shortcutReload = "CmdOrCtrl+R"
	shortcutQuit   = "CmdOrCtrl+Q"
)

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

func setupApplicationMenu(app *application.App) {
	menu := app.NewMenu()
	if runtime.GOOS == "darwin" {
		menu.AddRole(application.AppMenu)
	}
	menu.AddRole(application.FileMenu)
	menu.AddRole(application.EditMenu)
	menu.AddRole(application.WindowMenu)

	appMenu := menu.AddSubmenu("Setting")
	appMenu.Add("Show").SetAccelerator(shortcutShow).OnClick(func(ctx *application.Context) {
		showMainWindow(app)
	})
	appMenu.Add("Reload").SetAccelerator(shortcutReload).OnClick(func(ctx *application.Context) {
		app.Window.Current().Reload()
	})
	appMenu.AddSeparator()
	appMenu.Add("Quit").SetAccelerator(shortcutQuit).OnClick(func(ctx *application.Context) {
		quitting.Store(true)
		app.Quit()
	})

	menu.AddRole(application.HelpMenu)
	app.Menu.Set(menu)
}

func setupSystemTray(app *application.App, window *application.WebviewWindow) {
	tray := app.SystemTray.New()
	tray.SetTooltip("FreeAi")
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(trayIcon)
	} else {
		tray.SetIcon(trayIcon)
	}

	trayMenu := app.NewMenu()
	trayMenu.Add(menuLabel("显示主窗口", shortcutShow)).SetAccelerator(shortcutShow).OnClick(func(ctx *application.Context) {
		showWindow(window)
	})
	trayMenu.Add("隐藏到托盘").OnClick(func(ctx *application.Context) {
		window.Hide()
	})
	trayMenu.Add(menuLabel("重新加载", shortcutReload)).SetAccelerator(shortcutReload).OnClick(func(ctx *application.Context) {
		window.Reload()
	})
	trayMenu.AddSeparator()
	trayMenu.Add("打开后保持后台运行").SetEnabled(false)
	trayMenu.Add(menuLabel("退出 FreeAi", shortcutQuit)).SetAccelerator(shortcutQuit).OnClick(func(ctx *application.Context) {
		quitting.Store(true)
		app.Quit()
	})

	tray.SetMenu(trayMenu)
	tray.AttachWindow(window).WindowOffset(8)
	tray.OnClick(func() {
		showWindow(window)
	})
}

func setupApplicationEvents(app *application.App, window *application.WebviewWindow) {
	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting.Load() {
			return
		}
		window.Hide()
		event.Cancel()
	})

	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(event *application.ApplicationEvent) {
		showWindow(window)
	})
}

func showMainWindow(app *application.App) {
	window := app.Window.Current()
	window.Show().Focus()
	window.Center()
}

func showWindow(window *application.WebviewWindow) {
	window.Show().Focus()
	window.Center()
}

func menuLabel(label string, shortcut string) string {
	return label + " (" + shortcutDisplay(shortcut) + ")"
}

func shortcutDisplay(shortcut string) string {
	if runtime.GOOS != "darwin" {
		return strings.ReplaceAll(shortcut, "CmdOrCtrl", "Ctrl")
	}

	display := strings.ReplaceAll(shortcut, "CmdOrCtrl", "⌘")
	display = strings.ReplaceAll(display, "Shift", "⇧")
	display = strings.ReplaceAll(display, "+", "")
	return display
}
