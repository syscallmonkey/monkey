files = $(shell find . -iname "*.go")

build: $(files)
	BIN=monkey ./build/build.sh

clean:
	rm -rf bin
