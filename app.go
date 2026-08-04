package main

import (
	_ "embed"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

//go:embed build/tray/freeai-template.png
var trayIcon []byte

const (
	dashboardRoute   = "/dashboard"
	accountsRoute    = "/accounts/list"
	modelsRoute      = "/models/list"
	apiKeysRoute     = "/access/keys"
	integrationRoute = "/access/guide"
	usageRoute       = "/usage"
	requestLogsRoute = "/request-logs/list"
	tasksRoute       = "/ops/tasks"
	settingsRoute    = "/settings/gateway"
)

var quitting atomic.Bool

type desktopMenuActions struct {
	app           *application.App
	window        *application.WebviewWindow
	dataDirectory string
}

// setupDesktopMenus 统一配置原生菜单、托盘菜单和关闭到后台行为。
func setupDesktopMenus(app *application.App, window *application.WebviewWindow, dataDirectory string) {
	actions := desktopMenuActions{app: app, window: window, dataDirectory: dataDirectory}

	window.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if quitting.Load() {
			return
		}
		window.Hide()
		event.Cancel()
	})
	app.Event.OnApplicationEvent(events.Mac.ApplicationShouldHandleReopen, func(*application.ApplicationEvent) {
		actions.showWindow()
	})

	app.Menu.Set(buildApplicationMenu(app, actions))
	tray := app.SystemTray.New()
	tray.SetTooltip("FreeAI · 官方账号池网关")
	if runtime.GOOS == "darwin" {
		tray.SetTemplateIcon(trayIcon)
	} else {
		tray.SetIcon(applicationIcon)
	}
	tray.SetMenu(buildTrayMenu(app, actions))
	tray.OnClick(actions.showWindow)
}

// buildApplicationMenu 创建适合 FreeAI 管理工作流的原生应用菜单。
func buildApplicationMenu(app *application.App, actions desktopMenuActions) *application.Menu {
	menu := app.NewMenu()

	freeAIMenu := addTopLevelMenu(menu, "FreeAI", application.AppMenu)
	freeAIMenu.Add("关于 FreeAI").OnClick(func(*application.Context) { app.Menu.ShowAbout() })
	freeAIMenu.Add("网关设置…").SetAccelerator("CmdOrCtrl+,").OnClick(func(*application.Context) {
		actions.navigate(settingsRoute)
	})
	freeAIMenu.AddSeparator()
	freeAIMenu.Add("关闭窗口").SetAccelerator("CmdOrCtrl+W").OnClick(func(*application.Context) {
		actions.window.Hide()
	})
	freeAIMenu.Add("隐藏 FreeAI").SetAccelerator("CmdOrCtrl+H").OnClick(func(*application.Context) {
		actions.window.Hide()
	})
	freeAIMenu.AddSeparator()
	freeAIMenu.Add("退出 FreeAI").SetAccelerator("CmdOrCtrl+Q").OnClick(func(*application.Context) {
		quitting.Store(true)
		app.Quit()
	})

	editMenu := addTopLevelMenu(menu, "编辑", application.EditMenu)
	addRoleMenuItem(editMenu, application.Undo, "撤销")
	addRoleMenuItem(editMenu, application.Redo, "重做")
	editMenu.AddSeparator()
	addRoleMenuItem(editMenu, application.Cut, "剪切")
	addRoleMenuItem(editMenu, application.Copy, "复制")
	addRoleMenuItem(editMenu, application.Paste, "粘贴")
	addRoleMenuItem(editMenu, application.Delete, "删除")
	editMenu.AddSeparator()
	addRoleMenuItem(editMenu, application.SelectAll, "全选")

	navigationMenu := addTopLevelMenu(menu, "导航", application.NoRole)
	addNavigationMenuItem(navigationMenu, actions, "工作台", dashboardRoute)
	navigationMenu.AddSeparator()
	addNavigationMenuItem(navigationMenu, actions, "官方账号", accountsRoute)
	addNavigationMenuItem(navigationMenu, actions, "模型目录", modelsRoute)
	addNavigationMenuItem(navigationMenu, actions, "API 密钥", apiKeysRoute)
	addNavigationMenuItem(navigationMenu, actions, "接入指南", integrationRoute)
	navigationMenu.AddSeparator()
	addNavigationMenuItem(navigationMenu, actions, "用量分析", usageRoute)
	addNavigationMenuItem(navigationMenu, actions, "请求记录", requestLogsRoute)

	viewMenu := addTopLevelMenu(menu, "视图", application.ViewMenu)
	addRoleMenuItem(viewMenu, application.Reload, "刷新页面")
	viewMenu.AddSeparator()
	addRoleMenuItem(viewMenu, application.ResetZoom, "恢复默认缩放")
	addRoleMenuItem(viewMenu, application.ZoomIn, "放大")
	addRoleMenuItem(viewMenu, application.ZoomOut, "缩小")
	viewMenu.AddSeparator()
	addRoleMenuItem(viewMenu, application.ToggleFullscreen, "进入/退出全屏")

	windowMenu := addTopLevelMenu(menu, "窗口", application.WindowMenu)
	addRoleMenuItem(windowMenu, application.Minimise, "最小化")
	windowMenu.Add("显示主窗口").OnClick(func(*application.Context) { actions.showWindow() })

	helpMenu := addTopLevelMenu(menu, "帮助", application.HelpMenu)
	addNavigationMenuItem(helpMenu, actions, "任务中心", tasksRoute)
	helpMenu.Add("打开数据目录").OnClick(func(*application.Context) { actions.openDataDirectory() })
	helpMenu.Add("打开日志目录").OnClick(func(*application.Context) { actions.openLogDirectory() })
	helpMenu.AddSeparator()
	helpMenu.Add("开发者工具").OnClick(func(*application.Context) { actions.openDeveloperTools() })
	helpMenu.Add("关于 FreeAI").OnClick(func(*application.Context) { app.Menu.ShowAbout() })

	actions.window.RegisterKeyBinding("F12", func(currentWindow application.Window) {
		currentWindow.OpenDevTools()
	})
	return menu
}

// buildTrayMenu 创建常驻后台时可用的快速操作菜单。
func buildTrayMenu(app *application.App, actions desktopMenuActions) *application.Menu {
	menu := app.NewMenu()
	menu.Add("打开 FreeAI").OnClick(func(*application.Context) { actions.showWindow() })
	menu.AddSeparator()
	menu.Add("● 本地网关运行中 · 127.0.0.1:8787").SetEnabled(false)

	quickMenu := menu.AddSubmenu("快速入口")
	addNavigationMenuItem(quickMenu, actions, "官方账号", accountsRoute)
	addNavigationMenuItem(quickMenu, actions, "模型目录", modelsRoute)
	addNavigationMenuItem(quickMenu, actions, "API 密钥", apiKeysRoute)
	addNavigationMenuItem(quickMenu, actions, "用量分析", usageRoute)

	menu.Add("打开数据目录").OnClick(func(*application.Context) { actions.openDataDirectory() })
	menu.Add("打开日志目录").OnClick(func(*application.Context) { actions.openLogDirectory() })

	autostartEnabled, err := app.Autostart.IsEnabled()
	if err != nil {
		log.Printf("读取自动启动状态失败: %v\n", err)
	}
	menu.AddSeparator()
	menu.AddCheckbox("登录时自动启动", autostartEnabled).OnClick(func(ctx *application.Context) {
		var toggleErr error
		if ctx.IsChecked() {
			toggleErr = app.Autostart.EnableWithOptions(application.AutostartOptions{Identifier: applicationID})
		} else {
			toggleErr = app.Autostart.Disable()
		}
		if toggleErr != nil {
			ctx.ClickedMenuItem().SetChecked(!ctx.IsChecked())
			actions.showError("自动启动设置失败", toggleErr)
		}
	})

	menu.AddSeparator()
	menu.Add("关于 FreeAI").OnClick(func(*application.Context) { app.Menu.ShowAbout() })
	menu.Add("退出 FreeAI").OnClick(func(*application.Context) {
		quitting.Store(true)
		app.Quit()
	})
	return menu
}

func addTopLevelMenu(menu *application.Menu, label string, role application.Role) *application.Menu {
	item := application.NewSubMenuItem(label).SetRole(role)
	menu.Append(application.NewMenuFromItems(item))
	return item.GetSubmenu()
}

func addRoleMenuItem(menu *application.Menu, role application.Role, label string) {
	menu.AddRole(role)
	if item := menu.FindByRole(role); item != nil {
		item.SetLabel(label)
	}
}

func addNavigationMenuItem(menu *application.Menu, actions desktopMenuActions, label, route string) {
	menu.Add(label).OnClick(func(*application.Context) { actions.navigate(route) })
}

func showDesktopWindow(window *application.WebviewWindow) {
	window.Show()
	if window.IsMinimised() {
		window.UnMinimise()
	}
	window.Focus()
}

func (a desktopMenuActions) showWindow() {
	showDesktopWindow(a.window)
}

func (a desktopMenuActions) navigate(route string) {
	a.showWindow()
	// FreeAI 使用 HashLocationStrategy，原生菜单直接更新 hash 即可交给 Angular 路由处理。
	a.window.ExecJS(fmt.Sprintf(`window.location.hash = %s;`, strconv.Quote("#"+route)))
}

func (a desktopMenuActions) openDataDirectory() {
	a.openDirectory(filepath.Join(a.dataDirectory, "data"), "无法打开数据目录")
}

func (a desktopMenuActions) openLogDirectory() {
	a.openDirectory(filepath.Join(a.dataDirectory, "logback"), "无法打开日志目录")
}

func (a desktopMenuActions) openDirectory(directory, title string) {
	if err := a.app.Env.OpenFileManager(directory, false); err != nil {
		a.showError(title, err)
	}
}

func (a desktopMenuActions) openDeveloperTools() {
	a.showWindow()
	a.window.OpenDevTools()
}

func (a desktopMenuActions) showError(title string, err error) {
	log.Printf("%s: %v\n", title, err)
	a.showWindow()
	dialog := a.app.Dialog.Error().SetTitle(title).SetMessage(err.Error()).SetIcon(applicationIcon).AttachToWindow(a.window)
	dialog.AddButton("确定").SetAsDefault()
	dialog.Show()
}
