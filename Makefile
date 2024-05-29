BINARY=bin/monitoring-system.out

CMD_DIR=./src/cmd
CONFIG_DIR=./src/config
DOMAIN_DIR=./src/domain
INTERNAL_DIR=./src/internal
PKG_DIR=./src/pkg

GO=go
GOFMT=gofmt

CONFIG_FILE=config/config.go

PKGS=$(shell $(GO) list ./... | grep -v /vendor/)

all: build

build: fmt vet
	$(GO) build -o $(BINARY) $(CMD_DIR)/main.go

run: build
	$(BINARY)

fmt:
	$(GOFMT) -w $(CMD_DIR) $(CONFIG_DIR) $(DOMAIN_DIR) $(INTERNAL_DIR) $(PKG_DIR)

vet:
	$(GO) vet $(PKGS)

test:
	$(GO) test -v $(PKGS)

clean:
	$(GO) clean
	rm -f $(BINARY)

deps:
	$(GO) get -u ./...

.PHONY: all build run fmt vet test clean deps env