CURDIR = $(shell pwd)
GOPATH = "$(CURDIR)/.gopath"
VERSION = $(shell grep -e 'const VERSION' rivescript.go | head -n 1 | cut -d '"' -f 2)
UNAME = $(shell uname)
all: build

# `make setup` to set up git submodules
setup:
	git submodule init
	git submodule update

# `make run` to run the rivescript cmd
run: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run cmd/rivescript/main.go eg/brain

# `make debug` to run the rivescript cmd in debug mode
debug: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run cmd/rivescript/main.go -debug eg/brain

# `make fmt` to run gofmt
fmt:
	gofmt -w .

# `make test` to run unit tests
test: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go test ./src

# `make build` to build the binary
build: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 \
		go build -i -o bin/rivescript cmd/rivescript/main.go

# `make build.win32` to build a Windows binary
build.win32: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 GOOS=windows GOARCH=386 \
		go build -v -i -o bin/rivescript.exe cmd/rivescript/main.go

# `make dist` to create a binary shippable distribution
dist: build dist.common
	cp bin/rivescript dist/rivescript/
	cd dist; tar -czvf ../rivescript-$(VERSION)-$(UNAME).tar.gz rivescript

# `make dist.win32` to cross compile and distribute for Windows.
dist.win32: build.win32 dist.common
	echo -e "@echo off\nrivescript eg/brain" > dist/rivescript/example.bat
	cp bin/rivescript.exe dist/rivescript/
	cd dist; zip -r ../rivescript-$(VERSION)-win32.zip rivescript

dist.common:
	# Reset the dist directory and copy relevant files.
	rm -rf dist; mkdir -p dist/rivescript
	cp -r README.md LICENSE eg dist/rivescript/

# Sets up the gopath / build environment
gopath:
	mkdir -p .gopath/src/github.com/aichaos bin
	ln -sf "$(CURDIR)" .gopath/src/github.com/aichaos/

# Cleans everything up.
clean:
	rm -rf .gopath bin dist
