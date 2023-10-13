GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" \                                
                -trimpath -o bin/wxapkgx-arm64 cmd/wxapkgx/main.go
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" \                                
                -trimpath -o bin/wxapkgx-amd64 cmd/wxapkgx/main.go