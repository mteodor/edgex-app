BUILD_DIR = build
CGO_ENABLED ?= 0
GOOS ?= linux

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ${BUILD_DIR}/edgex-app cmd/main.go
endef

all:  edgex-app

.PHONY: all 


edgex-app:
	$(call compile_service,$(@))

clean:
	rm -rf ${BUILD_DIR}

install:
	cp ${BUILD_DIR}/* $(GOBIN)

test:
	GOCACHE=off go test -v -race -tags test $(shell go list ./... | grep -v 'vendor\|cmd')



release:
	$(eval version = $(shell git describe --abbrev=0 --tags))
	git checkout $(version)


	

