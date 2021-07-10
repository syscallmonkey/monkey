name ?= syscall-monkey
version ?= v0.0.1
bin ?= monkey
pkg ?= "github.com/seeker89/syscall-monkey"
tag = $(name):$(version)
goos ?= linux
namespace ?= "seeker89/"
files = $(shell find . -iname "*.go")


bin/$(bin): $(files)
	GOOS=${goos} PKG=${pkg} ARCH=amd64 VERSION=${version} BIN=${bin} ./build/build.sh

build:
	docker build -t $(tag) .

tag:
	docker tag $(tag) $(namespace)$(tag)

push:
	docker push $(namespace)$(tag)

clean:
	rm -rf bin


.PHONY: clean build tag push
