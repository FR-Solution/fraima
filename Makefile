#change for new project
project = fraimactl
#change for new release
release = v1.0.0

tag = $(DOCKER_USER)/$(project):$(release)
pwd = $(shell pwd)
module = $(shell go list -m)

build-bin:
	go install ./...

build-and-push:
	docker build -t $(tag) --build-arg VERSION=$(release) --build-arg PROJECT=$(project) -f Dockerfile .
	docker image push $(tag)
	echo $(tag)

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --fix -v

release:
	sh hack/release.sh
