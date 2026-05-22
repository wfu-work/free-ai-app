package backend

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wfu-work/free-ai-go/inits"
)

const (
	BackendBaseURL      = "http://127.0.0.1:8787"
	backendHealthURL    = BackendBaseURL + "/api/health"
	backendReadyTimeout = 30 * time.Second
	proxyErrorLogGap    = 5 * time.Second
)

func Start() {
	inits.Init()
}

func StartAndWait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		Start()
	}()

	readyCh := make(chan error, 1)
	go func() {
		readyCh <- waitBackendReady(ctx, backendHealthURL, backendReadyTimeout)
	}()

	select {
	case err := <-readyCh:
		return err
	case <-done:
		return errors.New("backend exited before it became ready")
	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewAPIMiddleware() (application.Middleware, error) {
	proxyURL, err := url.Parse(BackendBaseURL)
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

			if r.Body != nil && r.Method != http.MethodGet && r.Method != http.MethodHead {
				body, readErr := io.ReadAll(r.Body)
				if readErr != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				_ = r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(body))
				r.GetBody = func() (io.ReadCloser, error) {
					return io.NopCloser(bytes.NewReader(body)), nil
				}
				r.ContentLength = int64(len(body))
				r.Header.Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
			}

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
