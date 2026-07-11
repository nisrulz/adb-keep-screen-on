.PHONY: help build install uninstall clean test vet build-all run start stop status

BINARY := adb-keep-screen-on
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build                 Build the binary ($(BINARY))"
	@echo "  install               Build and copy to ~/go/bin"
	@echo "  uninstall             Remove from ~/go/bin"
	@echo "  clean                 Remove the binary"
	@echo "  test                  Run all tests"
	@echo "  vet                   Run go vet (static analysis)"
	@echo "  build-all             Cross-compile for all platforms (via build.sh)"
	@echo "  run                   Build and run in foreground"
	@echo "  start                 Start the daemon"
	@echo "  stop                  Stop the daemon"
	@echo "  status                Show daemon status"

build:
	@go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./src/
	@echo "  ✓ Built $(BINARY) ($(VERSION))"

install: build
	@mkdir -p ~/go/bin
	@cp $(BINARY) ~/go/bin/$(BINARY)
	@echo "  ✓ Installed to ~/go/bin/$(BINARY)"

test:
	@go test -v ./...
	@echo "  ✓ All tests passed"

vet:
	@go vet ./...
	@echo "  ✓ go vet passed"

build-all:
	@./build.sh

run: build
	@./$(BINARY) --foreground

start: build
	@./$(BINARY)

stop:
	@./$(BINARY) stop

status: build
	@./$(BINARY) status

clean:
	@rm -f $(BINARY)
	@echo "  ✓ Removed $(BINARY)"

uninstall:
	@rm -f ~/go/bin/$(BINARY)
	@echo "  ✓ Removed $(BINARY) from ~/go/bin"
