SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BUILD_PATH=build
BINARY=terraform-provider-zabbix
BIN=bin
TARGETS=darwin linux windows
VERSION=0.0.4
RELEASE_DIR=release
RELEASE=terraform-provider-zabbix_$(VERSION).tar.gz

build=GOOS=$(1) GOARCH=$(2) go build -o ${BUILD_PATH}/$(BINARY)/$(1)_$(2)/${BINARY}`if [ "$(1)" = "windows" ]; then echo ".exe"; fi` main.go

.DEFAULT_GOAL: $(BINARY)

$(BIN)/$(BINARY): $(SOURCES)
	go build -o $(BIN)/${BINARY} main.go

$(TARGETS): $(SOURCES)
	$(call build,$@,amd64,$(ext))
	$(call build,$@,386,$(ext))

release: $(TARGETS)
	mkdir -p $(RELEASE_DIR)
	tar -czvf $(RELEASE_DIR)/$(RELEASE) -C $(BUILD_PATH) $(BINARY)

install:
	go install ./...

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@if [ -d ${BUILD_PATH} ]; then rm -r ${BUILD_PATH} ; fi
	@if [ -d ${RELEASE_DIR} ]; then rm -r $(RELEASE_DIR) ; fi

.PHONY: clean install build release $(TARGETS)