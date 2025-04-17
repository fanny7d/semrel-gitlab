.PHONY: build test clean lint

# 构建参数
BINARY_NAME=semrel-gitlab
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# 默认目标
all: build

# 构建
build:
	go build $(LDFLAGS) -o $(BINARY_NAME)

# 测试
test:
	go test -v ./...

# 清理
clean:
	go clean
	rm -f $(BINARY_NAME)

# 代码检查
lint:
	golangci-lint run

# 安装依赖
deps:
	go mod tidy
	go mod vendor

# 运行
run: build
	./$(BINARY_NAME)

# 安装
install: build
	cp $(BINARY_NAME) /usr/local/bin/

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build    - 构建项目"
	@echo "  make test     - 运行测试"
	@echo "  make clean    - 清理构建文件"
	@echo "  make lint     - 运行代码检查"
	@echo "  make deps     - 更新依赖"
	@echo "  make run      - 构建并运行"
	@echo "  make install  - 安装到系统"
	@echo "  make help     - 显示帮助信息"

docker-build:
	docker build --build-arg VERSION=${VERSION} -t go-semrel-gitlab .

docker-run:
	docker run -it --rm go-semrel-gitlab

# 添加自动补全支持
install-completion:
	cp completions/go-semrel-gitlab.bash /etc/bash_completion.d/
	cp completions/go-semrel-gitlab.zsh /usr/local/share/zsh/site-functions/ 