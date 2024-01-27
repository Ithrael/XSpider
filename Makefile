.PHONY: all build run gotool clean help

Binary="xspider"

Version="v0.0.1"


all: gotool build

build:
	go build -o ${Binary}

build_all: gotool build_win_amd64 build_win_arm64 build_linux_amd64 build_linux_arm64 build_darwin_amd64 build_darwin_arm64

build_win_amd64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./artifacts/${Binary}-${Version}-windows-amd64.exe

build_win_arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./artifacts/${Binary}-${Version}-windows-arm64.exe

build_linux_amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./artifacts/${Binary}-${Version}-linux-amd64

build_linux_arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./artifacts/${Binary}-${Version}-linux-arm64

build_darwin_amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./artifacts/${Binary}-${Version}-darwin-amd64

build_darwin_arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./artifacts/${Binary}-${Version}-darwin-arm64

run:
	@go run ./

test:
	go test -v ./...

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${Binary} ] ; then rm ${Binary} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make test - 直接运行单元测试代码"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"
	@echo "make build_all - 编译所有平台的版本"