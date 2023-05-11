install:
	go install github.com/go-swagger/go-swagger/cmd/swagger@latest
	go install github.com/wailsapp/wails/v2/cmd/wails@latest

generate:
	go generate ./...

test:
	go test ./...

fmt:
	gofumpt -w .
	go mod tidy
	golangci-lint run

build-front-standalone:
	cd web-frontend && npm run build:standalone && cd ..

build-macos:
	mkdir -p api/server
	mkdir -p "build/wallet-plugin/"
	make generate
	make build-front-standalone
	CGO_LDFLAGS="-framework UniformTypeIdentifiers"  go build -tags desktop,production -ldflags "-w -s" -o build/wallet-plugin/wallet-plugin main.go

run-standalone:
	STANDALONE=1 ./build/wallet-plugin/wallet-plugin

wails-dev:
	STANDALONE=1 wails dev

install-plugin:
	make build
	mkdir -p ~/.config/thyra/plugins/wallet-plugin
	cp build/wallet-plugin/wallet-plugin ~/.config/thyra/plugins/wallet-plugin
