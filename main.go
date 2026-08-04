package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"sync"
	"time"

	"freeai-app/backend"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const (
	applicationID      = "com.xiaoxi.freeai"
	backendReadyPeriod = 30 * time.Second
)

// embeddedAssets 保存由 free-ai-web 构建并同步过来的 Angular 静态资源。
//
//go:embed all:frontend/dist
var embeddedAssets embed.FS

//go:embed build/appicon.png
var applicationIcon []byte

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run 组装桌面资源、后台服务、单实例窗口和系统托盘生命周期。
func run() error {
	// 桌面开发模式始终使用同步到 frontend/dist 的资源，避免意外加载外部开发服务器。
	if err := os.Unsetenv("FRONTEND_DEVSERVER_URL"); err != nil {
		return fmt.Errorf("关闭外部前端开发服务器: %w", err)
	}

	assets, err := fs.Sub(embeddedAssets, "frontend/dist")
	if err != nil {
		return fmt.Errorf("打开内嵌前端资源: %w", err)
	}
	dataDirectory, err := backend.ResolveDataDirectory()
	if err != nil {
		return err
	}
	host := backend.NewHost(dataDirectory)
	apiMiddleware, err := backend.NewAPIMiddleware(host.BaseURL())
	if err != nil {
		return err
	}

	var windowMu sync.RWMutex
	var mainWindow *application.WebviewWindow
	app := application.New(application.Options{
		Name:        "FreeAI",
		Description: "本地官方 AI 账号池与 OpenAI 兼容网关",
		Icon:        applicationIcon,
		Services: []application.Service{
			application.NewService(backend.NewDesktopService(dataDirectory, host.BaseURL())),
		},
		Assets: application.AssetOptions{
			Handler:    newDesktopAssetHandler(assets),
			Middleware: apiMiddleware,
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: applicationID,
			OnSecondInstanceLaunch: func(application.SecondInstanceData) {
				windowMu.RLock()
				window := mainWindow
				windowMu.RUnlock()
				if window != nil {
					showDesktopWindow(window)
				}
			},
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Windows: application.WindowsOptions{
			DisableQuitOnLastWindowClosed: true,
		},
		Linux: application.LinuxOptions{
			DisableQuitOnLastWindowClosed: true,
		},
	})

	host.Start()
	readyContext, cancelReady := context.WithTimeout(context.Background(), backendReadyPeriod)
	defer cancelReady()
	if err = host.WaitReady(readyContext); err != nil {
		return fmt.Errorf("启动内嵌后台服务: %w", err)
	}

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:      "main",
		Title:     "FreeAI",
		Width:     1440,
		Height:    900,
		MinWidth:  1024,
		MinHeight: 680,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})
	windowMu.Lock()
	mainWindow = window
	windowMu.Unlock()

	setupDesktopMenus(app, window, dataDirectory)
	if err = app.Run(); err != nil {
		return fmt.Errorf("运行 Wails 桌面应用: %w", err)
	}
	return nil
}
