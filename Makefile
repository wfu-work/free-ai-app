SHELL := /bin/sh

APP_NAME ?= free-ai-app
APP_DISPLAY_NAME ?= FreeAi
PLATFORM ?= $(shell go env GOOS 2>/dev/null || uname -s | tr '[:upper:]' '[:lower:]')
ARCH ?= $(shell go env GOARCH 2>/dev/null || uname -m)
DEV ?= false
FORMAT ?= nsis
INSTALL_SCOPE ?= machine
TAG ?= $(APP_NAME):latest
PORT ?= 8080
VERSION ?= $(shell awk '/^[[:space:]]+version:/ { gsub(/["'\'']/, "", $$2); print $$2; exit }' build/config.yml 2>/dev/null)
DMG_NAME ?= $(APP_DISPLAY_NAME)-$(VERSION)-$(ARCH).dmg
DMG_STAGING ?= bin/dmg-staging

TASK ?= $(shell if command -v task >/dev/null 2>&1; then printf 'task'; elif command -v wails3 >/dev/null 2>&1; then printf 'wails3 task'; fi)
HOST_OS := $(shell go env GOOS 2>/dev/null || uname -s | tr '[:upper:]' '[:lower:]')

TASK_VARS = ARCH=$(ARCH)
TASK_PLATFORM_VARS = $(TASK_VARS) FORMAT=$(FORMAT) INSTALL_SCOPE=$(INSTALL_SCOPE)

.DEFAULT_GOAL := help

.PHONY: help check-tools check-darwin tidy bindings icons dev dev-dist run build build-dev package package-universal dmg \
	build-platform build-platform-dev package-platform build-server run-server build-docker run-docker \
	setup-docker clean

help:
	@printf '%s\n' 'FreeAi Make targets'
	@printf '%s\n' ''
	@printf '%s\n' 'Development:'
	@printf '%s\n' '  make dev                         Build and run the desktop app for the host platform'
	@printf '%s\n' '  make dev-dist                    Same as dev, using embedded frontend assets'
	@printf '%s\n' '  make run                         Run the already-built host binary/app bundle'
	@printf '%s\n' ''
	@printf '%s\n' 'Build and package:'
	@printf '%s\n' '  make build                       Build for the host platform'
	@printf '%s\n' '  make build-dev                   Build a debug binary for the host platform'
	@printf '%s\n' '  make package                     Package for the host platform'
	@printf '%s\n' '  make dmg                         Build macOS .dmg installer image'
	@printf '%s\n' '  make package-universal           Package universal macOS app (darwin only)'
	@printf '%s\n' ''
	@printf '%s\n' 'Cross-platform:'
	@printf '%s\n' '  make build-platform PLATFORM=darwin ARCH=arm64'
	@printf '%s\n' '  make build-platform PLATFORM=windows ARCH=amd64'
	@printf '%s\n' '  make build-platform PLATFORM=linux ARCH=amd64'
	@printf '%s\n' '  make package-platform PLATFORM=windows FORMAT=nsis'
	@printf '%s\n' '  make package-platform PLATFORM=windows FORMAT=msix'
	@printf '%s\n' '  make package-platform PLATFORM=linux'
	@printf '%s\n' ''
	@printf '%s\n' 'Server and Docker:'
	@printf '%s\n' '  make build-server                Build no-GUI HTTP server binary'
	@printf '%s\n' '  make run-server                  Build and run server binary'
	@printf '%s\n' '  make build-docker TAG=free-ai:dev'
	@printf '%s\n' '  make run-docker PORT=8080 TAG=free-ai:dev'
	@printf '%s\n' ''
	@printf '%s\n' 'Maintenance:'
	@printf '%s\n' '  make tidy                        Run go mod tidy through Task'
	@printf '%s\n' '  make bindings                    Generate Wails bindings'
	@printf '%s\n' '  make icons                       Generate app icons'
	@printf '%s\n' '  make setup-docker                Build the cross-compilation Docker image'
	@printf '%s\n' '  make clean                       Remove generated binaries and app bundles'

check-tools:
	@[ -n "$(TASK)" ] || { printf '%s\n' 'Missing task runner. Install Wails with: go install github.com/wailsapp/wails/v3/cmd/wails3@latest'; exit 1; }
	@command -v wails3 >/dev/null 2>&1 || { printf '%s\n' 'Missing wails3. Install it with: go install github.com/wailsapp/wails/v3/cmd/wails3@latest'; exit 1; }

check-darwin:
	@[ "$(HOST_OS)" = "darwin" ] || { printf '%s\n' 'DMG packaging must run on macOS because it uses hdiutil and codesign.'; exit 1; }
	@command -v hdiutil >/dev/null 2>&1 || { printf '%s\n' 'Missing hdiutil. It should be available on macOS.'; exit 1; }

tidy: check-tools
	$(TASK) common:go:mod:tidy

bindings: check-tools
	$(TASK) common:generate:bindings

icons: check-tools
	$(TASK) common:generate:icons

dev: check-tools
	$(TASK) dev

dev-dist: check-tools
	$(TASK) dev:dist

run: check-tools
	$(TASK) run

build: check-tools
	$(TASK) build $(TASK_VARS)

build-dev: check-tools
	$(TASK) build $(TASK_VARS) DEV=true

package: check-tools
	$(TASK) package $(TASK_PLATFORM_VARS)

dmg: check-tools check-darwin
	$(TASK) darwin:package $(TASK_VARS)
	rm -rf "$(DMG_STAGING)"
	mkdir -p "$(DMG_STAGING)"
	cp -R "bin/$(APP_DISPLAY_NAME).app" "$(DMG_STAGING)/"
	ln -s /Applications "$(DMG_STAGING)/Applications"
	hdiutil create -volname "$(APP_DISPLAY_NAME)" -srcfolder "$(DMG_STAGING)" -ov -format UDZO "bin/$(DMG_NAME)"
	@printf '%s\n' "Created bin/$(DMG_NAME)"

package-universal: check-tools
	$(TASK) darwin:package:universal

build-platform: check-tools
	$(TASK) $(PLATFORM):build $(TASK_VARS)

build-platform-dev: check-tools
	$(TASK) $(PLATFORM):build $(TASK_VARS) DEV=true

package-platform: check-tools
	$(TASK) $(PLATFORM):package $(TASK_PLATFORM_VARS)

build-server: check-tools
	$(TASK) build:server

run-server: check-tools
	$(TASK) run:server

build-docker: check-tools
	$(TASK) build:docker TAG=$(TAG)

run-docker: check-tools
	$(TASK) run:docker TAG=$(TAG) PORT=$(PORT)

setup-docker: check-tools
	$(TASK) setup:docker

clean:
	rm -rf bin/free-ai-app bin/free-ai-app.exe bin/free-ai-app-server bin/free-ai-app-server.exe bin/FreeAi.app bin/free-ai-app.app bin/free-ai-app.dev.app bin/dmg-staging bin/*.dmg
