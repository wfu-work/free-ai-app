# FreeAi

![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v3-DF0000?logo=wails&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Platform](https://img.shields.io/badge/Platform-macOS%20%7C%20Windows%20%7C%20Linux%20%7C%20Mobile-blue)

FreeAi 是一个基于 Wails v3 与 Go 构建的跨平台 AI 桌面应用。项目将本地 WebView、Go 后端服务和前端静态资源打包到同一个应用中，适合用于构建私有化、轻量级、可扩展的 AI 客户端。

> 当前仓库以应用壳、后端启动器、构建脚本和已构建前端资源为主。前端源码如需二次开发，请确保同步对应的 `frontend` 工程源码。

## 特性

- 跨平台桌面应用：基于 Wails v3，支持 macOS、Windows、Linux 构建。
- 本地后端服务：应用启动时初始化 Go 后端能力，并通过 `/api` 路由转发请求。
- 前端资源内嵌：使用 `embed.FS` 打包 `frontend/dist/freeai-web/browser` 静态资源。
- 开发与发布脚本：通过 Taskfile 统一管理开发、构建、打包、Docker 和移动端任务。
- 移动端构建模板：仓库内包含 iOS 与 Android 的 Wails 构建资产和任务脚本。
- 服务模式：支持无 GUI 的 server build，便于部署为纯 HTTP 服务。

## 技术栈

- Go 1.26+
- Wails v3
- Task
- Docker，可选，用于交叉编译或 server 镜像
- Android / iOS 构建工具，可选，用于移动端构建

## 项目结构

```text
.
├── backend/              # 后端启动与 API 代理中间件
├── build/                # Wails 跨平台构建资产与 Taskfile
│   ├── android/          # Android 构建模板
│   ├── darwin/           # macOS 构建与签名配置
│   ├── docker/           # server / cross compile Dockerfile
│   ├── ios/              # iOS 构建模板
│   ├── linux/            # Linux 打包配置
│   └── windows/          # Windows 打包配置
├── data/                 # 本地运行数据
├── frontend/dist/        # 已构建的前端静态资源
├── logback/              # 运行日志
├── main.go               # Wails 应用入口
├── Taskfile.yml          # 统一任务入口
└── go.mod                # Go 模块依赖
```

## 快速开始

### 环境要求

请先安装：

- Go 1.26 或更高版本
- Wails v3 CLI
- Task
- Node.js / npm，仅在需要重新构建前端时使用

安装 Wails CLI：

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

安装 Task：

```bash
go install github.com/go-task/task/v3/cmd/task@latest
```

### 拉取项目

```bash
git clone https://github.com/<your-org>/free-ai-app.git
cd free-ai-app
```

### 开发模式运行

```bash
task dev
```

或者直接使用 Wails：

```bash
wails3 dev
```

应用默认在本机启动服务，并通过 `127.0.0.1:8787` 加载窗口内容。

## 构建

### 构建当前平台应用

```bash
task build
```

构建产物默认输出到 `bin/` 目录。

### 打包当前平台应用

```bash
task package
```

macOS 会生成 `.app` bundle；Windows、Linux 会根据对应平台 Taskfile 生成安装包或可执行文件。

### Server 模式

如果只需要运行后端 HTTP 服务，不需要桌面 GUI：

```bash
task build:server
task run:server
```

### Docker 镜像

```bash
task build:docker
task run:docker
```

默认会构建 server 模式镜像，可通过 `TAG` 和 `PORT` 覆盖参数：

```bash
task run:docker TAG=free-ai-app:dev PORT=8080
```

## 跨平台构建

本项目的 `build/` 目录包含 Wails 生成的多平台构建任务：

```bash
task darwin:build
task windows:build
task linux:build
task android:build
task ios:build
```

部分平台需要额外工具链：

- macOS / iOS：Xcode、Command Line Tools、可用的模拟器或签名配置。
- Windows：WebView2、NSIS / MSIX 相关工具链。
- Linux：C 编译器、AppImage / nfpm 相关依赖。
- Android：Android Studio、SDK、NDK、Gradle。

跨平台编译可使用 Docker：

```bash
task setup:docker
task linux:build
task windows:build
```

## API 代理

应用入口会启动后端初始化逻辑，并为 Wails 静态资源服务挂载 API 中间件：

- 非 `/api` 请求交给前端静态资源处理。
- `/api` 请求会转发到本地后端服务。
- 代理异常时返回 `503 Service Unavailable`，并限制重复日志输出频率。

相关代码位于：

- `main.go`
- `backend/backend.go`

## 开发说明

常用命令：

```bash
task dev              # 开发模式
task build            # 构建当前平台
task package          # 打包当前平台
task build:server     # 构建 server 模式
task run:server       # 运行 server 模式
task build:docker     # 构建 Docker 镜像
task run:docker       # 运行 Docker 镜像
```

更新构建资产：

```bash
wails3 task common:update:build-assets
```

生成绑定：

```bash
wails3 task common:generate:bindings
```

## 配置

项目元信息位于 `build/config.yml`，包括：

- 应用名称
- 产品标识
- 版本号
- 版权信息
- 开发模式监听规则
- 文件关联配置

修改 `build/config.yml` 中的应用信息后，建议同步更新构建资产：

```bash
wails3 task common:update:build-assets
```

## 贡献

欢迎提交 Issue 和 Pull Request。

建议流程：

1. Fork 本仓库。
2. 创建功能分支：`git checkout -b feat/your-feature`。
3. 提交修改：`git commit -m "feat: add your feature"`。
4. 推送分支：`git push origin feat/your-feature`。
5. 创建 Pull Request。

提交前建议运行：

```bash
go test ./...
task build
```

## License

本项目基于 [MIT License](./LICENSE) 开源。
