name ?= monkey
version ?= v0.0.1
bin ?= monkey
pkg ?= github.com/syscallmonkey/monkey
tag = $(name):$(version)
goos ?= linux
namespace ?= ghcr.io/syscallmonkey/
files = $(shell find . -iname "*.go")


bin/$(bin): $(files)
	GOOS=${goos} PKG=${pkg} ARCH=amd64 VERSION=${version} BIN=${bin} ./build/build.sh

build:
	docker build -t $(tag) .

tag:
	docker tag $(tag) $(namespace)$(tag)

push:
	docker push $(namespace)$(tag)

run:
	docker run --rm -ti --entrypoint /bin/bash --cap-add SYS_PTRACE -t $(tag)

clean:
	rm -rf bin

generate:
	python3 build/generate.py > pkg/syscall/syscalls_linux.go

.PHONY: clean build tag push run generate
