default: build

# 构建编译
build:
	go build .

# 生成对应的文档
generatedoc:
	rm -rf website/*
	tfplugindocs generate --rendered-website-dir website --examples-dir examples

# 跨平台编译
build-all:
	rm -rf bin/windows_arm64/*
	rm -rf bin/windows_amd64/*
	rm -rf bin/darwin_arm64/*
	rm -rf bin/darwin_amd64/*
	rm -rf bin/linux_arm64/*
	rm -rf bin/linux_amd64/*
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o bin/windows_arm64/ .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/windows_amd64/ .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64/ .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64/ .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64/ .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64/ .
