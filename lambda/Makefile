
# build lambda func

GO = go
GOBUILD = $(GO) build
GOTEST  = $(GO) test
GOCLEAN = $(GO) clean
APP = bootstrap
SOURCES = main.go
ZIP_FILE = bowbowLambdaFunc.zip

DEBUGFLAGS = -race -gcflags="-m -l"

.PHONY: all build test clean

all: build test

build:
	@echo "Building production ready bowbow app lambda code..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(APP) $(SOURCES)
	@echo "Creating ZIP package for AWS Lambda"
	zip $(ZIP_FILE) $(APP)

debug:
	@echo "Building the debug mode bowbow app lambda code..."
	GOOS=linux GOARCH=amd64 $(DEBUGBUILD) -o $(APP) $(SOURCES)

test:
	$(GOTEST) -v ./...

clean:
	rm $(APP) $(ZIP_FILE)
