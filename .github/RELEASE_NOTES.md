## FreeAI Desktop {{VERSION}}

FreeAI 是面向本地单机使用的官方 AI 账号池与 OpenAI 兼容网关。本版本由 GitHub Actions 从固定源码版本自动构建。

### 下载选择

| 操作系统 | 架构 | 推荐文件 |
| --- | --- | --- |
| macOS 12 或更高版本 | Apple Silicon | `FreeAi-*-macos-arm64.dmg` |
| macOS 12 或更高版本 | Intel | `FreeAi-*-macos-amd64.dmg` |
| Windows 10/11 | x64 | `FreeAi-*-windows-amd64-setup.exe` |
| Windows 10/11 | ARM64 | `FreeAi-*-windows-arm64-setup.exe` |
| Linux | x64 | `FreeAi-*-linux-amd64.AppImage`、`.deb`、`.rpm` 或 `.pkg.tar.zst` |
| Linux | ARM64 | `FreeAi-*-linux-arm64.AppImage`、`.deb`、`.rpm` 或 `.pkg.tar.zst` |

### 安装说明

- macOS：打开 DMG，将 `FreeAi.app` 拖入“应用程序”。当前自动构建采用临时签名、尚未进行 Apple 公证；若系统阻止首次运行，请在“系统设置 → 隐私与安全性”中确认打开。
- Windows：运行对应架构的安装程序。当前安装包尚未配置商业代码签名，SmartScreen 可能要求手动确认。
- Linux：Ubuntu/Debian 推荐 `.deb`，Fedora/RHEL 推荐 `.rpm`，Arch Linux 推荐 `.pkg.tar.zst`；其他发行版可使用 AppImage。

### 升级与数据

- 升级安装不会主动删除账号、API 密钥、用量记录或本地配置。
- 升级前仍建议在系统的备份功能中导出一份账号池备份。
- 桌面程序会在本机启动 FreeAI 后端，不依赖 Redis，默认不向第三方中转供应商发送请求。

### 完整性校验

所有安装包的 SHA-256 摘要都记录在 `SHA256SUMS.txt`。Linux/macOS 可执行：

```bash
sha256sum -c SHA256SUMS.txt
```

macOS 也可以使用 `shasum -a 256 <文件名>` 单独核对。

### 构建信息

- 桌面端标签：`{{VERSION}}`
- 前端提交：`{{FRONTEND_SHA}}`
- 发布内容下方的变更列表由 GitHub 根据本标签与上一版本之间的合并记录自动生成。
