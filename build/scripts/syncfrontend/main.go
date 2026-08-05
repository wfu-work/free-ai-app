package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	source := flag.String("source", "", "前端构建目录")
	target := flag.String("target", "", "Wails 内嵌资源目录")
	flag.Parse()
	if err := syncDirectory(*source, *target); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// syncDirectory 校验源目录后，完整替换桌面端的生成资源目录。
func syncDirectory(source, target string) error {
	sourcePath, err := safeDirectory(source)
	if err != nil {
		return fmt.Errorf("无效源目录: %w", err)
	}
	targetPath, err := safeDirectory(target)
	if err != nil {
		return fmt.Errorf("无效目标目录: %w", err)
	}
	if sourcePath == targetPath {
		return fmt.Errorf("源目录和目标目录不能相同")
	}
	if info, statErr := os.Stat(sourcePath); statErr != nil {
		return fmt.Errorf("读取源目录: %w", statErr)
	} else if !info.IsDir() {
		return fmt.Errorf("源路径不是目录: %s", sourcePath)
	}
	if _, statErr := os.Stat(filepath.Join(sourcePath, "index.html")); statErr != nil {
		return fmt.Errorf("源目录缺少 index.html: %w", statErr)
	}
	if err = os.RemoveAll(targetPath); err != nil {
		return fmt.Errorf("清理目标目录: %w", err)
	}
	if err = os.MkdirAll(targetPath, 0755); err != nil {
		return fmt.Errorf("创建目标目录: %w", err)
	}
	if err = os.CopyFS(targetPath, os.DirFS(sourcePath)); err != nil {
		return fmt.Errorf("复制前端资源: %w", err)
	}
	// 保留空目录占位文件，确保全新克隆在同步前也能通过 go:embed 编译检查。
	if err = os.WriteFile(filepath.Join(targetPath, ".keep"), nil, 0644); err != nil {
		return fmt.Errorf("写入前端目录占位文件: %w", err)
	}
	return nil
}

func safeDirectory(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("路径为空")
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	volumeRoot := filepath.VolumeName(absolutePath) + string(filepath.Separator)
	if filepath.Clean(absolutePath) == filepath.Clean(volumeRoot) {
		return "", fmt.Errorf("不允许使用文件系统根目录")
	}
	return filepath.Clean(absolutePath), nil
}
