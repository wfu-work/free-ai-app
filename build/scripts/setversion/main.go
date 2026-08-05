package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var semanticVersionPattern = regexp.MustCompile(`^v?([0-9]+)\.([0-9]+)\.([0-9]+)(?:-[0-9A-Za-z][0-9A-Za-z.-]*)?(?:\+[0-9A-Za-z][0-9A-Za-z.-]*)?$`)

func main() {
	version := flag.String("version", "", "发布版本，例如 1.2.0 或 1.2.0-beta.1")
	root := flag.String("root", ".", "free-ai-app 工程根目录")
	flag.Parse()

	if err := updateVersionFiles(*root, *version); err != nil {
		fmt.Fprintf(os.Stderr, "更新发布版本失败: %v\n", err)
		os.Exit(1)
	}
}

// normalizeVersion 校验发布版本，并返回适合桌面安装包元数据的三段数字版本。
func normalizeVersion(value string) (string, error) {
	value = strings.TrimSpace(value)
	matches := semanticVersionPattern.FindStringSubmatch(value)
	if len(matches) != 4 {
		return "", fmt.Errorf("版本 %q 不是有效的 SemVer", value)
	}
	return strings.Join(matches[1:4], "."), nil
}

// updateVersionFiles 同步更新各桌面平台的安装包版本，避免文件名与包内版本不一致。
func updateVersionFiles(root, releaseVersion string) error {
	version, err := normalizeVersion(releaseVersion)
	if err != nil {
		return err
	}

	configPath := filepath.Join(root, "build", "config.yml")
	if err = replaceExactlyOnce(
		configPath,
		regexp.MustCompile(`(?m)^(\s{2}version:\s*)"[^"]+"(\s*# The application version.*)$`),
		`${1}"`+version+`"${2}`,
	); err != nil {
		return err
	}

	darwinPath := filepath.Join(root, "build", "darwin", "Info.plist")
	for _, key := range []string{"CFBundleVersion", "CFBundleShortVersionString"} {
		if err = replaceExactlyOnce(
			darwinPath,
			regexp.MustCompile(`(?s)(<key>`+regexp.QuoteMeta(key)+`</key>\s*<string>)[^<]+(</string>)`),
			`${1}`+version+`${2}`,
		); err != nil {
			return err
		}
	}

	windowsInfoPath := filepath.Join(root, "build", "windows", "info.json")
	for _, key := range []string{"file_version", "ProductVersion"} {
		if err = replaceExactlyOnce(
			windowsInfoPath,
			regexp.MustCompile(`("`+regexp.QuoteMeta(key)+`"\s*:\s*")[^"]+("\s*[,}])`),
			`${1}`+version+`${2}`,
		); err != nil {
			return err
		}
	}

	if err = replaceExactlyOnce(
		filepath.Join(root, "build", "windows", "nsis", "wails_tools.nsh"),
		regexp.MustCompile(`(?m)^(\s*!define INFO_PRODUCTVERSION\s+")[^"]+(".*)$`),
		`${1}`+version+`${2}`,
	); err != nil {
		return err
	}

	return replaceExactlyOnce(
		filepath.Join(root, "build", "linux", "nfpm", "nfpm.yaml"),
		regexp.MustCompile(`(?m)^(version:\s*")[^"]+(".*)$`),
		`${1}`+version+`${2}`,
	)
}

// replaceExactlyOnce 只允许目标配置出现一次，格式漂移时立即让发布构建失败。
func replaceExactlyOnce(path string, pattern *regexp.Regexp, replacement string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取 %s: %w", path, err)
	}
	if matches := pattern.FindAllIndex(contents, -1); len(matches) != 1 {
		return fmt.Errorf("文件 %s 中期望匹配 1 处，实际匹配 %d 处", path, len(matches))
	}

	updated := pattern.ReplaceAll(contents, []byte(replacement))
	if err = os.WriteFile(path, updated, 0o644); err != nil {
		return fmt.Errorf("写入 %s: %w", path, err)
	}
	return nil
}
