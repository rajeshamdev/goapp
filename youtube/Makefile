GO = go
#GOBUILD = $(GO) build -mod vendor
GOBUILD = $(GO) build
GOTEST  = $(GO) test
GOCLEAN = $(GO) clean
APP = server
APPLINUX = server-linux

DEBUGFLAGS = -race -gcflags="-m -l"
DEBUGBUILD = $(GO) build -mod vendor $(DEBUGFLAGS)

.PHONY: all build test clean

all: test build

build:
	$(GOBUILD) -o $(APP) main.go

linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(APPLINUX) main.go

debug:
	$(DEBUGBUILD) -o $(APP) main.go

test:
	$(GOTEST) -v ./...

clean:
	rm $(APP)
