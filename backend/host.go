package backend

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/wfu-work/free-ai-go/inits"
)

const (
	BackendBaseURL   = "http://127.0.0.1:8787"
	backendHealthURL = BackendBaseURL + "/api/health"
)

// Host 管理嵌入式 FreeAI 后端的启动和就绪探测。
type Host struct {
	dataDirectory string
	startOnce     sync.Once
	done          chan struct{}
}

// NewHost 创建使用指定数据目录的后端宿主。
func NewHost(dataDirectory string) *Host {
	return &Host{dataDirectory: dataDirectory, done: make(chan struct{})}
}

// Start 在后台启动 FreeAI HTTP 服务；重复调用不会重复启动。
func (h *Host) Start() {
	h.startOnce.Do(func() {
		go func() {
			defer close(h.done)
			inits.Init()
		}()
	})
}

// BaseURL 返回桌面前端访问后台管理 API 的回环地址。
func (h *Host) BaseURL() string {
	return BackendBaseURL
}

// DataDirectory 返回当前桌面实例的数据根目录。
func (h *Host) DataDirectory() string {
	return h.dataDirectory
}

// WaitReady 等待健康检查成功，或在后端提前退出和上下文取消时返回错误。
func (h *Host) WaitReady(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	client := &http.Client{Timeout: time.Second}
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		request, err := http.NewRequestWithContext(ctx, http.MethodGet, backendHealthURL, nil)
		if err != nil {
			return err
		}
		response, requestErr := client.Do(request)
		if requestErr == nil {
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("等待后台服务就绪: %w", ctx.Err())
		case <-h.done:
			return fmt.Errorf("后台服务在就绪前退出")
		case <-ticker.C:
		}
	}
}
