package backend

// DesktopService 向前端暴露少量桌面环境信息，业务操作仍统一使用 REST API。
type DesktopService struct {
	dataDirectory string
	backendURL    string
}

// NewDesktopService 创建桌面环境信息服务。
func NewDesktopService(dataDirectory, backendURL string) *DesktopService {
	return &DesktopService{dataDirectory: dataDirectory, backendURL: backendURL}
}

// DataDirectory 返回本地数据库、密钥和日志使用的数据根目录。
func (s *DesktopService) DataDirectory() string {
	return s.dataDirectory
}

// BackendURL 返回对外提供 OpenAI 兼容接口的本地网关地址。
func (s *DesktopService) BackendURL() string {
	return s.backendURL
}
