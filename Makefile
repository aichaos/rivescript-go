VERSION=$(shell grep -e 'const Version' rivescript.go | head -n 1 | cut -d '"' -f 2)
BUILD=$(shell git describe --always)
CURDIR=$(shell pwd)
GOPATH="$(CURDIR)/.gopath"

# Inject the build version (commit hash) into the executable.
LDFLAGS := -ldflags "-X main.Build=$(BUILD)"

# `make setup` to set up git submodules
.PHONY: setup
setup:
	git submodule init
	git submodule update

# `make build` to build the binary
.PHONY: build
build: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 \
		go build $(LDFLAGS) -i -o bin/rivescript cmd/rivescript/main.go

# `make run` to run the rivescript cmd
.PHONY: run
run: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run $(LDFLAGS) cmd/rivescript/main.go eg/brain

# `make debug` to run the rivescript cmd in debug mode
.PHONY: debug
debug: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go run $(LDFLAGS) cmd/rivescript/main.go -debug eg/brain

# `make fmt` to run gofmt
.PHONY: fmt
fmt:
	gofmt -w .

# `make test` to run unit tests
.PHONY: test
test: gopath
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 go test

# `make clean` cleans up everything
.PHONY: clean
clean:
	rm -rf .gopath bin dist

################################################################################
## Below are commands for shipping distributable binaries for each platfomr.  ##
################################################################################

PLATFORMS := linux/amd64 linux/386 darwin/amd64
WIN32     := windows/amd64 windows/386
release: $(PLATFORMS) $(WIN32)
.PHONY: release $(PLATFORMS)

# Handy variables to pull OS and arch from $PLATFORMS.
temp = $(subst /, ,$@)
os   = $(word 1, $(temp))
arch = $(word 2, $(temp))

$(PLATFORMS): gopath
	mkdir -p dist/rivescript-$(VERSION)-$(os)-$(arch)
	cp -r README.md LICENSE Changes.md eg dist/rivescript-$(VERSION)-$(os)-$(arch)/
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 GOOS=$(os) GOARCH=$(arch) \
		go build $(LDFLAGS) -v -i -o bin/rivescript cmd/rivescript/main.go
	cp bin/rivescript* dist/rivescript-$(VERSION)-$(os)-$(arch)/
	cd dist; tar -czvf ../rivescript-$(VERSION)-$(os)-$(arch).tar.gz rivescript-$(VERSION)-$(os)-$(arch)

$(WIN32): gopath
	mkdir -p dist/rivescript-$(VERSION)-$(os)-$(arch)
	cp -r README.md LICENSE Changes.md eg dist/rivescript-$(VERSION)-$(os)-$(arch)/
	GOPATH=$(GOPATH) GO15VENDOREXPERIMENT=1 GOOS=$(os) GOARCH=$(arch) \
		go build $(LDFLAGS) -v -i -o bin/rivescript.exe cmd/rivescript/main.go
	cp bin/rivescript.exe dist/rivescript-$(VERSION)-$(os)-$(arch)/
	echo -e "@echo off\nrivescript eg/brain" > dist/rivescript-$(VERSION)-$(os)-$(arch)/example.bat
	cd dist; zip -r ../rivescript-$(VERSION)-$(os)-$(arch).zip rivescript-$(VERSION)-$(os)-$(arch)

# Sets up the gopath / build environment
gopath:
	mkdir -p .gopath/src/github.com/aichaos bin
	ln -sf "$(CURDIR)" .gopath/src/github.com/aichaos/
	ln -sf "$(CURDIR)/vendor/gopkg.in" .gopath/src/
