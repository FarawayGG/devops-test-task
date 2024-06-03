BUILD_SRC=.
REV=$(shell git rev-parse HEAD)

IMAGE_NAME="app"
IMAGE_TAG=$(shell git rev-parse --short HEAD)
IMAGE=${IMAGE_NAME}:${IMAGE_TAG}

default: help

.PHONY: fmt
fmt:
	@gofmt -w -s ./

clean:
	@rm -rf ${BUILD_DIR}

all: build

.PHONY: build
build: fmt
	DOCKER_BUILDKIT=1 docker build --no-cache --ssh default -f Dockerfile -t ${IMAGE} .
	docker tag ${IMAGE} ${IMAGE_NAME}:latest

.PHONY: run
run: run
	docker-compose -f devops/docker-compose.yml up -d

.PHONY: stop
stop: stop
	docker-compose -f devops/docker-compose.yml stop

.PHONY: restart
restart: restart
	docker-compose -f devops/docker-compose.yml restart
