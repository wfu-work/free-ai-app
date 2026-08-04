package backend

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

const proxyErrorLogGap = 5 * time.Second

// NewAPIMiddleware 创建只代理 /api 路径的 Wails 资源中间件。
func NewAPIMiddleware(target string) (application.Middleware, error) {
	return newAPIMiddleware(target, nil)
}

func newAPIMiddleware(target string, transport http.RoundTripper) (application.Middleware, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("解析后台地址: %w", err)
	}
	if targetURL.Scheme == "" || targetURL.Host == "" {
		return nil, fmt.Errorf("后台地址必须包含协议和主机")
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	if transport != nil {
		proxy.Transport = transport
	}
	var errorLogMu sync.Mutex
	lastErrorLog := time.Time{}
	proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, proxyErr error) {
		now := time.Now()
		errorLogMu.Lock()
		if now.Sub(lastErrorLog) >= proxyErrorLogGap {
			log.Printf("后台代理 %s 不可用: %v\n", request.URL.Path, proxyErr)
			lastErrorLog = now
		}
		errorLogMu.Unlock()
		http.Error(writer, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path != "/api" && !strings.HasPrefix(request.URL.Path, "/api/") {
				next.ServeHTTP(writer, request)
				return
			}
			proxy.ServeHTTP(writer, request)
		})
	}, nil
}
