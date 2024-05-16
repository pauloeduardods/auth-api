BINARY=bin/monitoring-system-server.out

API_DIR=./api
CMD_DIR=./cmd
CONFIG_DIR=./config
INTERNAL_DIR=./internal
PKG_DIR=./pkg

GO=go
GOFMT=gofmt

CONFIG_FILE=config/config.go

PKGS=$(shell $(GO) list ./... | grep -v /vendor/)

include .env

all: build

build: fmt vet
	$(GO) build -o $(BINARY) $(CMD_DIR)/main.go

run: build
	APP_ENV=$(APP_ENV) HOST=$(HOST) PORT=$(PORT) COGNITO_CLIENT_ID=$(COGNITO_CLIENT_ID) COGNITO_USER_POOL_ID=$(COGNITO_USER_POOL_ID) REGION=$(REGION) $(BINARY)

fmt:
	$(GOFMT) -w $(CMD_DIR) $(CONFIG_DIR) $(INTERNAL_DIR) $(PKG_DIR) $(API_DIR)

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