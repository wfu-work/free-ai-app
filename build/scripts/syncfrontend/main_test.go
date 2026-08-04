package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestSyncDirectory 验证前端同步会移除旧产物并复制完整的新构建。
func TestSyncDirectory(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	target := filepath.Join(root, "target")
	if err := os.MkdirAll(filepath.Join(source, "assets"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "index.html"), []byte("new-index"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "assets", "logo.svg"), []byte("logo"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "stale.js"), []byte("stale"), 0644); err != nil {
		t.Fatal(err)
	}

	if err := syncDirectory(source, target); err != nil {
		t.Fatalf("同步目录失败: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "stale.js")); !os.IsNotExist(err) {
		t.Fatalf("旧资源未被清理: %v", err)
	}
	if content, err := os.ReadFile(filepath.Join(target, "assets", "logo.svg")); err != nil || string(content) != "logo" {
		t.Fatalf("新资源未被复制: content=%q err=%v", content, err)
	}
}
