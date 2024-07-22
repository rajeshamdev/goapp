GO = go
#GOBUILD = $(GO) build -mod vendor
GOBUILD = $(GO) build
GOTEST  = $(GO) test
GOCLEAN = $(GO) clean
APP = server

DEBUGFLAGS = -race -gcflags="-m -l"

.PHONY: all build test clean

all: test build

build:
	$(GOBUILD) -o $(APP) main.go

test:
	$(GOTEST) -v ./...

clean:
	rm $(APP)
