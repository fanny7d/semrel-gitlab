.PHONY: build test clean install

VERSION ?= DEV

build:
	go build -ldflags="-X main.version=${VERSION}" -o go-semrel-gitlab

test:
	go test -v ./...

clean:
	rm -f go-semrel-gitlab

install: build
	cp go-semrel-gitlab /usr/local/bin/

docker-build:
	docker build --build-arg VERSION=${VERSION} -t go-semrel-gitlab .

docker-run:
	docker run -it --rm go-semrel-gitlab

# 添加自动补全支持
install-completion:
	cp completions/go-semrel-gitlab.bash /etc/bash_completion.d/
	cp completions/go-semrel-gitlab.zsh /usr/local/share/zsh/site-functions/ 