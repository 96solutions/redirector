TARGET=redirective_$(REVISION)
INSTALL_DIR=./bin/
REVISION=$(shell sh -c "git rev-parse --short HEAD" | awk '{print $$1}')
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

lint:
	golangci-lint --exclude-use-default=false --out-format tab run ./...

#all: clean build
#
#.PHONY: build
#build:
#	@echo ">>> Current commit hash $(REVISION)"
#	@echo ">>> go build -o $(TARGET)"
#	@go mod vendor && CGO_ENABLED=0 GOOS=linux go build -a -o $(TARGET) .

#.PHONY: dockerbuild
#dockerbuild:
#	@docker run --name redirector_builder_c --rm --interactive --tty --volume $(CURRENT_DIR):/src redirector_builder make build
#
#clean:
#	rm -rf $(TARGET)
#
#rmlink:
#	rm -f $(INSTALL_DIR)redirector
#
#mklink:
#	ln -sf $(CURRENT_DIR)/bin/$(TARGET) $(CURRENT_DIR)/bin/redirector
#
#install:
#	make all && mkdir -p $(INSTALL_DIR) && mv $(TARGET) $(INSTALL_DIR) && make rmlink && make mklink
#
#mkmigration:
#	@echo ">>> Current commit hash $(REVISION)"

install:
 	curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz
