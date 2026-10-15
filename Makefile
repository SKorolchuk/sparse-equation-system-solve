# -------------------------------------------------------------------------
# DIRECTORY PARAMS (Standard Go Project Layout adjusted for static WASM web)
# -------------------------------------------------------------------------
WEB_APP_DIR       := ./web/app
WASM_SRC          := $(WEB_APP_DIR)/wasm
WASM_OUT          := $(WEB_APP_DIR)/public/solver.wasm
INTERNAL_DIR      := ./internal/solver

# -------------------------------------------------------------------------
# SYSTEM UTILS CONTEXT
# -------------------------------------------------------------------------
.PHONY: all init sync build build-wasm build-frontend test bench lint help

## all: Default task - runs building process (WASM + Frontend)
all: init build

## build: Compiles the Go solver library into WASM and bundles React application
build: build-wasm build-frontend
	@echo "[SUCCESS] React web bundle generated successfully."

## init: First-time project setup - restores Go modules and npm packages
init:
	@echo "[INIT] Downloading Go module dependencies..."
	go mod download
	@echo "[INIT] Restoring frontend node modules..."
	cd $(WEB_APP_DIR) && npm install
	@echo "[SUCCESS] Project dependencies are successfully installed."

## sync: Update go.mod file and sync module dependecies
sync:
	go mod tidy

## build-wasm: Cross-compiles Go package to JS/WASM target architecture and copy wasm-exec.js bootstrap script from Go SDK
build-wasm:
	@echo "[BUILD] Compiling Go library into WebAssembly target..."
	GOOS=js GOARCH=wasm go build -o $(WASM_OUT) $(WASM_SRC)
	@echo "[BUILD] Injecting GOROOT/wasm_exec.js bridge layer..."
	@go run ./scripts/copy_wasm_bridge.go
	@echo "[SUCCESS] solver.wasm and bridge layers successfully placed in React static environment."

## build-frontend: Runs package installation and triggers Vite build
build-frontend:
	@echo "[BUILD] Restoring node modules..."
	cd $(WEB_APP_DIR) && npm install
	@echo "[BUILD] Build React SPA assets via Vite bundler..."
	cd $(WEB_APP_DIR) && npm run build
	@echo "[SUCCESS] Production build ready."

## test: Invokes validation suites for the core library
test:
	@echo "[TEST] Triggering test suite validation..."
	go test -v -cover $(INTERNAL_DIR)/...

## bench: Executes CPU performance benchmarks for runtime estimation
bench:
	@echo "[BENCHMARK] Auditing runtime operational efficiency and allocations..."
	go test -bench=. -benchmem $(INTERNAL_DIR)/...

## lint: Audits Go code quality indicators using standard static checking tools
lint:
	@echo "[LINT] Invoking static code quality scanners..."
	go vet ./...

## help: Shows descriptive details for all accessible targets inside this workspace automation system
help:
	@echo "Available automation targets inside this project:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
