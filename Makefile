BUILD_DIR = build
CGO_ENABLED ?= 0
GOOS ?= linux

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -ldflags "-s -w" -o ${BUILD_DIR}/edgex-app cmd/main.go
endef

all:  edgex-app

.PHONY: all docker


edgex-app:
	$(call compile_service,$(@))

clean:
	rm -rf ${BUILD_DIR}

install:
	cp ${BUILD_DIR}/* $(GOBIN)

docker:
	docker build --no-cache --build-arg edgex-app --tag=edgex-app -f docker/Dockerfile .


test:
	GOCACHE=off go test -v -race -tags test $(shell go list ./... | grep -v 'vendor\|cmd')


	

