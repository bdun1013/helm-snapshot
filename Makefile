PLUGIN_NAME := unittest

.PHONY: build
build: build_linux build_mac build_windows

build_windows: export GOARCH=amd64
build_windows: export GO111MODULE=on
build_windows: export GOPROXY=https://goproxy.io
build_windows:
	@GOOS=windows go build -v --ldflags="-w -X main.version=$(VERSION)" \
		-o bin/windows/amd64/unittest main.go  # windows

link_windows:
	@cp bin/windows/amd64/unittest ./bin/unittest

build_linux: export GOARCH=amd64
build_linux: export CGO_ENABLED=0
build_linux: export GO111MODULE=on
build_linux: export GOPROXY=https://goproxy.io
build_linux:
	@GOOS=linux go build -v --ldflags="-w -X main.version=$(VERSION)" \
		-o bin/linux/amd64/unittest main.go  # linux

link_linux:
	@cp bin/linux/amd64/unittest ./bin/unittest

build_mac: export GOARCH=amd64
build_mac: export CGO_ENABLED=0
build_mac: export GO111MODULE=on
build_mac: export GOPROXY=https://goproxy.io
build_mac:
	@GOOS=darwin go build -v --ldflags="-w -X main.version=$(VERSION)" \
		-o bin/darwin/amd64/unittest main.go # mac osx
	@cp bin/darwin/amd64/unittest ./bin/unittest # For use w make install

link_mac:
	@cp bin/darwin/amd64/unittest ./bin/unittest

.PHONY: clean
clean:
	@git status --ignored --short | grep '^!! ' | sed 's/!! //' | xargs rm -rf

.PHONY: covhtml
covhtml:
	@go tool cover -html=.cover/cover.out

.PHONY: tree
tree:
	@tree -I vendor

.PHONY: install
install:
	HELM_PUSH_PLUGIN_NO_INSTALL_HOOK=1 helm plugin install $(shell pwd)

.PHONY: remove
remove:
	helm plugin remove $(PLUGIN_NAME)