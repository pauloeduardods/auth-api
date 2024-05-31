BINARY=bin/auth-api.out

API_DIR=./src/api
CMD_DIR=./src/cmd
CONFIG_DIR=./src/config
INTERNAL_DIR=./src/internal
PKG_DIR=./src/pkg
DB_DIR=./db

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
	$(GOFMT) -w $(API_DIR) $(CMD_DIR) $(CONFIG_DIR) $(INTERNAL_DIR) $(PKG_DIR)

vet:
	$(GO) vet $(PKGS)

test:
	$(GO) test -v $(PKGS)

clean:
	$(GO) clean
	rm -f $(BINARY)

deps:
	$(GO) get -u ./...

database:
	docker-compose -f docker-compose.database.yml up  -d 

database-down:
	docker-compose -f docker-compose.database.yml down

db-create-tables:
	docker-compose -f docker-compose.database.yml exec postgres psql -U myuser -d mydatabase -a -f /scripts/schema.sql

.PHONY: all build run fmt vet test clean deps database database-down db-create-tables
