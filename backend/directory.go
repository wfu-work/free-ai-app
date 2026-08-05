package backend

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	freeaiutils "github.com/wfu-work/free-ai-go/utils"
)

var localConfigFileNames = []string{
	"config.debug.yaml",
	"config.release.yaml",
	"config.test.yaml",
	"config.yaml",
}

// ResolveDataDirectory 按实际配置来源解析数据库、密钥和日志所在的根目录。
func ResolveDataDirectory() (string, error) {
	if configPath := explicitConfigPath(os.Args[1:]); configPath != "" {
		return absoluteParent(configPath)
	}
	if configPath := strings.TrimSpace(os.Getenv("NAV_CONFIG")); configPath != "" {
		return absoluteParent(configPath)
	}
	workingDirectory, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("读取工作目录: %w", err)
	}
	for _, name := range localConfigFileNames {
		if _, statErr := os.Stat(filepath.Join(workingDirectory, name)); statErr == nil {
			return workingDirectory, nil
		}
	}
	configPath, err := freeaiutils.DefaultConfigPath()
	if err != nil {
		return "", fmt.Errorf("解析 FreeAI 默认配置路径: %w", err)
	}
	dataDirectory := filepath.Dir(configPath)
	if err = os.MkdirAll(dataDirectory, 0700); err != nil {
		return "", fmt.Errorf("创建 FreeAI 数据目录: %w", err)
	}
	return filepath.Clean(dataDirectory), nil
}

func explicitConfigPath(args []string) string {
	for index, arg := range args {
		if (arg == "-c" || arg == "--c") && index+1 < len(args) {
			return strings.TrimSpace(args[index+1])
		}
		if strings.HasPrefix(arg, "-c=") {
			return strings.TrimSpace(strings.TrimPrefix(arg, "-c="))
		}
		if strings.HasPrefix(arg, "--c=") {
			return strings.TrimSpace(strings.TrimPrefix(arg, "--c="))
		}
	}
	return ""
}

func absoluteParent(configPath string) (string, error) {
	absolutePath, err := filepath.Abs(filepath.FromSlash(configPath))
	if err != nil {
		return "", fmt.Errorf("解析配置路径: %w", err)
	}
	return filepath.Dir(absolutePath), nil
}
