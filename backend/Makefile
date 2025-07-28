# make debug
.PHONY:debug
debug:
	air -c .air.toml

.PHONY:debug-file-clear
debug-file-clear:
	rm -rf ./tmp/main

# using at .air.toml
debug-file: debug-file-clear
	swag init
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o tmp/main main.go


# make build
.PHONY:build-clear
build-clear:
	rm -rf ./dist/

.PHONY:build
build: build-clear
	mkdir -p dist
# 	GOOS=darwin GOARCH=arm64 go build -o dist/main-macos main.go
	GOOS=linux GOARCH=arm64 go build -o dist/main-linux main.go

