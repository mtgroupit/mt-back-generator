GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=./bin/generator

TARGET_DIRECTORY=~/Desktop
CONFIG=./config.yaml

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/gen

gen:
	sh gen.sh $(TARGET_DIRECTORY) $(CONFIG)
	