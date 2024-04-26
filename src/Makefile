VERSION=dev
BINARY_NAME=$(VERSION)_ssh2ws

amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)_amd64 main.go

amd64-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_amd64_static main.go

arm32:
	GOOS=linux GOARCH=arm go build -o bin/$(BINARY_NAME)_arm32 main.go

arm32-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_arm32_static main.go

arm64:
	GOOS=linux GOARCH=arm64 go build -o bin/$(BINARY_NAME)_arm64 main.go

arm64-static:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_arm64_static main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME).exe main.go

windows-static:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_static.exe main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)_darwin main.go

darwin-static:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_darwin_static main.go

darwinarm64:
	GOOS=darwin GOARCH=arm64 go build -o bin/$(BINARY_NAME)_darwinarm64 main.go

darwinarm64-static:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -installsuffix cgo -o bin/$(BINARY_NAME)_darwinarm64_static main.go

all: amd64 amd64-static arm32 arm32-static arm64 arm64-static windows windows-static darwin darwin-static darwinarm64 darwinarm64-static