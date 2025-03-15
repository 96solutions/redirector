TARGET=redirector_$(REVISION)
INSTALL_DIR=./bin/
REVISION=$(shell sh -c "git config --global --add safe.directory $(CURRENT_DIR)" && git rev-parse --short HEAD | awk '{print $$1}')
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(patsubst %/,%,$(dir $(MKFILE_PATH)))

.PHONY: mocks
mocks:
	@mockgen -package=mocks -destination=mocks/mock_clicks_repository.go -source=domain/repository/clicks_repository.go ClicksRepository
	@mockgen -package=mocks -destination=mocks/mock_click_handler.go -source=domain/interactor/click_handler.go ClickHandlerInterface
	@mockgen -package=mocks -destination=mocks/mock_redirect_interactor.go -source=domain/interactor/redirect_interactor.go RedirectInteractor
	@mockgen -package=mocks -destination=mocks/mock_tracking_links_repository.go -source=domain/repository/tracking_links_repository.go TrackingLinksRepositoryInterface
	@mockgen -package=mocks -destination=mocks/mock_ip_address_parser.go -source=domain/service/ip_address_parser.go IPAddressParserInterface
	@mockgen -package=mocks -destination=mocks/mock_user_agent_parser.go -source=domain/service/user_agent_parser.go UserAgentParser

lint:
	golangci-lint --exclude-use-default=false --out-format tab run ./...

all: clean build

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -o $(TARGET) .

.PHONY: build
build:
	@echo ">>> Current commit hash $(REVISION)"
	@echo ">>> go build -o $(TARGET)"
	@go mod vendor && make compile

clean:
	rm -rf $(TARGET)

rmlink:
	rm -f $(INSTALL_DIR)redirector

mklink:
	ln -sf $(CURRENT_DIR)/bin/$(TARGET) $(CURRENT_DIR)/bin/redirector

install:
	make all && mkdir -p $(INSTALL_DIR) && mv $(TARGET) $(INSTALL_DIR) && make rmlink && make mklink

compiledaemon:
	@make clean
	@make compile
	@mkdir -p $(INSTALL_DIR)
	@mv $(TARGET) $(INSTALL_DIR)
	@make rmlink
	@make mklink

#
#mkmigration:
#	@echo ">>> Current commit hash $(REVISION)"

#install:
# 	curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$os-$arch.tar.gz | tar xvz

.PHONY: compile
compilefull:
#	@ls -la
	@echo ">>> GIT fix"
	@git config --global --add safe.directory $(CURRENT_DIR)
	@echo "rm vendor"
	@rm -rf vendor
	@echo ">>> Current commit hash $(REVISION)"
	@echo ">>> go build -o $(TARGET)"
	@make clean
	@go mod vendor && CGO_ENABLED=0 GOOS=linux go build -a -o $(TARGET) .
	@mkdir -p $(INSTALL_DIR)
	@mv $(TARGET) $(INSTALL_DIR)
	@make rmlink && make mklink

#TODO: use builder user instead of root
.PHONY: dockercompile
dockercompile:
	@docker run --name redirector_builder --rm --interactive --tty --volume $(CURRENT_DIR):/usr/src/app -w /usr/src/app -u $(id -u):$(id -g) golang:1.23 make compilefull