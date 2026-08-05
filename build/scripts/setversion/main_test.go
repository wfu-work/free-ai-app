package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "stable", input: "1.2.3", want: "1.2.3"},
		{name: "tag", input: "v2.0.1", want: "2.0.1"},
		{name: "prerelease", input: "1.5.0-beta.2", want: "1.5.0"},
		{name: "invalid", input: "release-1", wantErr: true},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := normalizeVersion(test.input)
			if test.wantErr {
				if err == nil {
					t.Fatalf("normalizeVersion(%q) 未返回错误", test.input)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("normalizeVersion(%q) = %q, %v; want %q", test.input, got, err, test.want)
			}
		})
	}
}

func TestUpdateVersionFiles(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	files := map[string]string{
		"build/config.yml": `info:
  version: "1.0.0" # The application version
`,
		"build/darwin/Info.plist": `<key>CFBundleVersion</key>
<string>1.0.0</string>
<key>CFBundleShortVersionString</key>
<string>1.0.0</string>
`,
		"build/windows/info.json": `{"fixed":{"file_version":"1.0.0"},"info":{"0000":{"ProductVersion":"1.0.0"}}}`,
		"build/windows/nsis/wails_tools.nsh": `    !define INFO_PRODUCTVERSION "1.0.0"
`,
		"build/linux/nfpm/nfpm.yaml": `version: "1.0.0"
`,
	}

	for name, contents := range files {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("创建测试目录失败: %v", err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatalf("写入测试文件失败: %v", err)
		}
	}

	if err := updateVersionFiles(root, "v2.3.4-rc.1"); err != nil {
		t.Fatalf("updateVersionFiles() error = %v", err)
	}
	for name := range files {
		contents, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatalf("读取结果失败: %v", err)
		}
		if !strings.Contains(string(contents), "2.3.4") || strings.Contains(string(contents), "1.0.0") {
			t.Fatalf("%s 未正确更新版本: %s", name, contents)
		}
	}
}
