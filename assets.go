package main

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// newDesktopAssetHandler 返回支持 Angular 客户端路由回退的静态资源处理器。
func newDesktopAssetHandler(assets fs.FS) http.Handler {
	staticAssets := application.BundledAssetFileServer(assets)
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if isSPAClientRoute(request) {
			request = request.Clone(request.Context())
			request.URL.Path = "/"
			request.URL.RawPath = ""
		}
		staticAssets.ServeHTTP(response, request)
	})
}

// isSPAClientRoute 判断请求是否应回退到 Angular 的 index.html。
func isSPAClientRoute(request *http.Request) bool {
	if request.Method != http.MethodGet && request.Method != http.MethodHead {
		return false
	}
	requestPath := request.URL.Path
	if requestPath == "" || requestPath == "/" {
		return false
	}
	if requestPath == "/api" || strings.HasPrefix(requestPath, "/api/") || strings.HasPrefix(requestPath, "/wails/") {
		return false
	}
	return path.Ext(requestPath) == ""
}
