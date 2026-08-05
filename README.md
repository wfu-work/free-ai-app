<div align="center">
  <img src="build/appicon.png" width="128" alt="FreeAI Logo">
  <h1>FreeAI Desktop</h1>
  <p>本地官方 AI 账号池与 OpenAI 兼容网关桌面应用。</p>
</div>

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v3-DF0000?logo=wails&logoColor=white)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux-blue)
![License](https://img.shields.io/badge/License-MIT-green)

FreeAi 是一个基于 Go 和 Wails 构建的本地 AI 网关桌面应用。它将 WebView 界面、本地后端、API 代理中间件、系统托盘能力以及跨平台构建资源打包到同一个桌面应用中。

这个项目适合在个人电脑上运行私有、轻量的 AI 代理管理工具：可以管理网关访问，让应用常驻系统托盘，暴露 OpenAI 兼容的 `/v1` 路由，并为不同平台构建桌面安装包。

> 当前仓库包含桌面外壳、后端启动器、构建脚本，以及位于 `frontend/dist` 下的已构建前端资源。如果需要修改 UI 源码，请先同步或恢复对应的前端源码项目，再重新构建前端产物。

## 功能特性

- 通过本地 `/v1` 路由提供 OpenAI 兼容的网关入口。
- 内置本地 API 代理中间件，用于转发来自 WebView 的后端请求。
- 基于 Wails v3 和 Go 的桌面应用外壳。
- macOS、Windows 和 Linux 关闭主窗口后继续在托盘运行。
- 单实例运行，重复启动时自动唤醒已有窗口。
- 系统托盘提供账号、模型、API 密钥和用量分析快速入口。
- 支持登录时自动启动、打开数据目录和日志目录。
- 使用 `embed.FS` 嵌入前端资源。
- 提供 macOS、Windows、Linux、iOS 和 Android 模板的跨平台构建任务。
- Makefile 封装开发、打包、Docker、服务端模式和 macOS DMG 输出等常用命令。
- 可选的无界面服务端构建，适合仅 HTTP 部署场景。

## 截图

截图暂未提交。如果发布正式版本，建议补充：

- 主控制台
- 网关访问指引
- 账号和密钥管理
- macOS 托盘菜单

## 技术栈

- Go 1.26+
- Wails v3
- Angular 构建后的前端资源
- Make
- 通过 `wails3 task` 执行的 Taskfile 兼容任务
- Docker，可选，用于交叉编译和服务端镜像
- 根据目标平台按需安装 Xcode、NSIS/MSIX、Linux 打包工具、Android/iOS 工具链

## 仓库结构

```text
.
├── backend/              # 后端宿主、目录解析、桌面服务和 API 代理
├── build/                # Wails 构建资源和各平台任务文件
│   ├── android/          # Android 模板资源
│   ├── darwin/           # macOS 应用包、签名和打包资源
│   ├── docker/           # 服务端和交叉编译 Dockerfile
│   ├── ios/              # iOS 模板资源
│   ├── linux/            # Linux AppImage/deb/rpm 打包资源
│   └── windows/          # Windows manifest、NSIS 和 MSIX 打包资源
├── frontend/dist/        # 从 free-ai-web 同步的浏览器构建产物
├── tools/                # 本地辅助工具
├── main.go               # Wails 桌面端入口
├── Makefile              # 开发和发布常用命令
├── Taskfile.yml          # Wails/Task 任务定义
└── go.mod                # Go 模块元数据
```

## 环境要求

必需：

- Go 1.26 或更高版本
- Wails v3 CLI
- Make

安装 Wails：

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

可选：

- `task`：如果你更习惯使用独立 Task CLI 可以安装；Makefile 会自动回退到 `wails3 task`，因此不是必需项。
- Docker：用于交叉编译和服务端镜像。
- Xcode Command Line Tools：用于 macOS 构建。
- NSIS 或 MSIX 工具：用于 Windows 安装包。
- Linux 打包工具：用于 AppImage/deb/rpm 输出。

## 快速开始

桌面、后端和前端仓库应位于同一个父目录：

```text
free-ai/
├── free-ai-app/
├── free-ai-go/
└── free-ai-web/
```

先构建前端，再运行桌面应用：

```bash
cd free-ai-web
pnpm install --frozen-lockfile
pnpm run build

cd free-ai-app
make dev
```

桌面应用会启动本地后端，并通过 Wails 内部 loopback 服务加载内嵌界面资源：

```text
Wails 内部地址（桌面窗口直接访问，不需要手动打开）
```

如果你想直接使用 Wails 任务命令：

```bash
wails3 task dev
```

## 常用命令

查看所有 Make 目标：

```bash
make help
```

开发：

```bash
make dev          # 构建并运行桌面应用
make dev-dist     # 使用内嵌前端资源运行开发模式
make run          # 运行已经构建好的宿主应用
```

构建与打包：

```bash
make build        # 为当前平台构建二进制文件
make build-dev    # 为当前平台构建调试版本
make package      # 打包当前平台应用
make dmg          # macOS：创建 bin/FreeAi-<version>-<arch>.dmg
```

服务端模式：

```bash
make build-server
make run-server
```

Docker：

```bash
make build-docker TAG=free-ai-app:dev
make run-docker TAG=free-ai-app:dev PORT=8080
```

维护：

```bash
make tidy
make bindings
make icons
make clean
```

## 跨平台构建

为指定目标平台构建：

```bash
make build-platform PLATFORM=darwin ARCH=arm64
make build-platform PLATFORM=windows ARCH=amd64
make build-platform PLATFORM=linux ARCH=amd64
```

为指定目标平台打包：

```bash
make package-platform PLATFORM=windows FORMAT=nsis
make package-platform PLATFORM=windows FORMAT=msix
make package-platform PLATFORM=linux
```

准备基于 Docker 的交叉编译环境：

```bash
make setup-docker
```

平台说明：

- macOS 打包需要在 macOS 上运行，以完成临时签名和 DMG 创建。
- Windows 打包可能需要 NSIS、MSIX 工具以及 WebView2 bootstrapper 生成能力。
- Linux 打包可通过 Wails 生成的任务输出 AppImage、deb、rpm 和 Arch 包。
- 仓库中包含 iOS 和 Android 模板，但移动端构建需要对应平台 SDK。

## macOS 分发

创建标准 macOS DMG：

```bash
make dmg
```

DMG 包含：

- `FreeAi.app`
- 用于拖拽安装的 `Applications` 快捷入口

输出文件会写入 `bin/`，例如：

```text
bin/FreeAi-0.0.1-arm64.dmg
```

如需发布签名和公证，请先在 `build/darwin/Taskfile.yml` 中配置签名变量，再使用 Wails 的签名任务。

## GitHub 自动构建与发布

仓库通过 [`.github/workflows/release.yml`](.github/workflows/release.yml) 构建并发布桌面安装包，覆盖：

| 平台 | 架构 | 发布产物 |
| --- | --- | --- |
| macOS | Intel、Apple Silicon | DMG |
| Windows | x64、ARM64 | NSIS 用户级安装程序 |
| Linux | x64、ARM64 | AppImage、deb、rpm、Arch Linux 包 |

正式发布使用符合 SemVer 的 `v` 前缀标签：

```bash
git tag -a v1.0.0 -m "FreeAI v1.0.0"
git push origin v1.0.0
```

标签推送后，GitHub Actions 会自动完成以下步骤：

1. 校验版本号并构建固定提交的 `free-ai-web`。
2. 将前端资源同步进 Wails 应用，更新各平台安装包内的版本元数据。
3. 并行构建 macOS、Windows 和 Linux 的 x64/ARM64 安装包。
4. 生成 `SHA256SUMS.txt`，创建 GitHub Release 并上传所有制品。
5. 使用 [`.github/RELEASE_NOTES.md`](.github/RELEASE_NOTES.md) 写入安装、升级和校验说明，再由 GitHub 追加本版本变更记录。

在 Actions 页面手动执行工作流只进行测试构建，不会创建 Release。测试构建可指定其他前端分支、标签或提交；产物会保留 14 天。

为保证标签可以重复构建，正式发布默认使用工作流中 `FRONTEND_REF` 固定的前端提交。发布新版 UI 前应先更新该提交 SHA。若 `free-ai-web` 是私有仓库，需要在 `free-ai-app` 仓库中增加具有只读权限的 Actions Secret：`FREE_AI_WEB_TOKEN`。

当前自动发布包尚未接入 Apple Developer ID、公证和 Windows 商业代码签名，系统首次运行时可能显示安全确认。发布说明会明确提示这一状态；获得证书后应将签名步骤接入发布工作流。

## 运行行为

- 在 macOS、Windows 和 Linux 上关闭窗口会隐藏应用，而不是退出应用。
- 应用会继续保留在系统托盘中。
- 托盘菜单支持快速导航、自动启动、打开数据/日志目录和退出。
- 第二次启动 FreeAI 会激活已经运行的主窗口。
- 只有明确执行退出操作时，应用进程才会结束。

## API 代理

Wails 入口会启动本地后端，并为内嵌前端资源配置 API 中间件：

- 前端页面由 Wails 本地资源服务直接返回。
- `/api` 请求会转发到本地后端 `127.0.0.1:8787`。
- 后端自己的静态站点不再作为桌面应用入口使用。
- 代理错误会返回 `503 Service Unavailable`，并在日志中限频输出。

相关文件：

- `main.go`
- `assets.go`
- `backend/host.go`
- `backend/proxy.go`

## 配置

项目元数据位于：

```text
build/config.yml
```

它控制产品名称、应用标识、版本、版权文本、开发模式、文件关联等信息。

修改生成类构建元数据后，请刷新 Wails 资源：

```bash
wails3 task common:update:build-assets
```

注意：该命令可能会覆盖生成的平台资源。提交前请认真检查 diff。

## 前端资源

应用会嵌入：

```text
frontend/dist
```

当前仓库不重复保存前端源码。修改 UI 后，在相邻的 `free-ai-web` 工程执行：

```bash
pnpm run build
```

桌面构建任务会校验 `dist/freeai-web/browser/index.html`，并完整同步到 `free-ai-app/frontend/dist`；不需要手动复制文件：

```bash
wails3 task common:sync:frontend
```

直接修改 `frontend/dist` 中的压缩文件虽然可行，但不利于长期维护，不建议作为常规方式。

桌面端应用图标以 `free-ai-web/src/assets/logo.svg` 为品牌源，彩色版本用于应用与安装包，单色版本用于 macOS 菜单栏模板图标。

## 参与贡献

欢迎提交贡献。一个好的 Pull Request 应该聚焦、可复现，并且便于评审。

建议流程：

1. Fork 本仓库。
2. 创建功能分支：`git checkout -b feat/your-feature`。
3. 完成代码修改。
4. 运行相关检查。
5. 使用清晰的提交信息提交代码。
6. 创建 Pull Request，并补充背景说明；如果涉及 UI 修改，请附上截图；如果涉及构建，请说明构建注意事项。

提交 PR 前建议运行：

```bash
go test ./...
make build
```

如果修改了平台打包逻辑，也请运行对应的打包目标，例如：

```bash
make dmg
make package-platform PLATFORM=windows FORMAT=nsis
```

## 已知说明

- `go test ./...` 或 `go build ./...` 可能会包含不适合直接在宿主机上构建的平台模板包。构建桌面应用本身时，请使用 `make build` 或 `go build .`。
- 交叉编译依赖平台 SDK 和 Docker 可用性。
- 部分生成资源会被有意提交到仓库中，因为应用需要嵌入已构建的前端资源。

## 许可证

FreeAi 基于 [MIT License](./LICENSE) 发布。
