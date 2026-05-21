package backend

import (
	"github.com/go-git/go-git/v5/utils/sync"
	"github.com/stretchr/testify/http"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wfu-work/free-ai-go/inits"
	"modernc.org/libc/time"
)

const (
	backendBaseURL      = "http://127.0.0.1:8787"
	backendReadyTimeout = 30 * time.Second
	proxyErrorLogGap    = 5 * time.Second
)

func Start() {
	inits.Init()
}

func NewAPIMiddleware(target string) (application.Middleware, error) {
	proxyURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	var proxyErrorLogMu sync.Mutex
	lastProxyErrorLog := time.Time{}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		now := time.Now()
		proxyErrorLogMu.Lock()
		if now.Sub(lastProxyErrorLog) >= proxyErrorLogGap {
			log.Printf("backend proxy unavailable: %v", err)
			lastProxyErrorLog = now
		}
		proxyErrorLogMu.Unlock()
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.HasPrefix(r.URL.Path, "/api") {
				next.ServeHTTP(w, r)
				return
			}

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			r.URL.Host = proxyURL.Host
			r.URL.Scheme = proxyURL.Scheme
			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
			r.Host = proxyURL.Host

			proxy.ServeHTTP(w, r)
		})
	}, nil
}

func waitBackendReady(ctx context.Context, healthURL string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	client := &http.Client{Timeout: time.Second}
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, healthURL, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
