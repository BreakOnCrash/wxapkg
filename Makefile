BINARY_NAME=wxapkgx
BINARY_ARM=wxapkgx-arm64

build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -trimpath -o bin/${BINARY_NAME} cmd/wxapkgx/main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -trimpath -o bin/${BINARY_ARM} cmd/wxapkgx/main.go

clean:
	go clean
	rm bin/${BINARY_NAME} & rm bin/${BINARY_ARM}