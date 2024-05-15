.PHONY:	all build run

GOLANG_VERSION = 1.22.3

all:
	@echo 'for build :'
	@echo '   make build'
	@echo 'for development:'
	@echo '   make run'


run:
	cd compose && docker compose up --watch

build:
	cp -r src docker/rootfs/
	cd docker && docker build -t devops-test-task --build-arg GOLANG_VERSION=$(GOLANG_VERSION) .
