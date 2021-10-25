SHELL := /bin/bash -o pipefail
VERSION := 1.0.0
SERVICE_NAME ?= $(shell git remote -v | head -n1 | awk '{print $$2}' | sed 's/.*\///' | sed 's/\.git//')
VERSION_DATE ?= $(shell date +"%y%m%d%H%M")
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_COMMIT_SHORT = $(shell a=$(GIT_COMMIT); echo $${a:0:7})
GIT_BRANCH ?= $(shell git symbolic-ref HEAD | sed -e 's,.*/\(.*\),\1,')
ifeq ($(GIT_BRANCH),master)
    BRANCH :=
else ifeq ($(GIT_BRANCH),)
	BRANCH := _unknown
else
	BRANCH := _$(GIT_BRANCH)
endif
FULL_VERSION ?= $(VERSION).$(VERSION_DATE)-$(GIT_COMMIT_SHORT)$(BRANCH)
PKG_TO_TEST := $(shell go list ./... | grep -v docs | grep -v test | tr '\n' ',' | sed 's/,$$//')
DOCKER_IMAGE := viveportsocial/$(SERVICE_NAME):$(FULL_VERSION)

echo:
	@echo SERVICE_NAME: $(SERVICE_NAME)
	@echo VERSION: $(VERSION)
	@echo VERSION_DATE: $(VERSION_DATE)
	@echo GIT_COMMIT_SHORT: $(GIT_COMMIT_SHORT)
	@echo GIT_BRANCH: $(GIT_BRANCH)
	@echo FULL_VERSION: $(FULL_VERSION)
	@echo DOCKER_IMAGE: $(DOCKER_IMAGE)

set-version:
	@echo echo "##vso[task.setvariable variable=buildVersion;isOutput=true]$(FULL_VERSION)"

install:
	@go mod download
	@go get github.com/swaggo/swag/cmd/swag@v1.6.9
	@go get -u github.com/jstemmer/go-junit-report
	@go get -u github.com/t-yuki/gocover-cobertura

fmt:
	@goimports -w .

build:
	@swag init --parseDependency --parseInternal
	@make fmt
	@go build -v -o $(SERVICE_NAME)

run:
	@make build
	@./$(SERVICE_NAME)

clean:
	@rm -rf $(SERVICE_NAME)
	@rm -rf coverage
	@rm -rf docs
	@rm -rf gen
	@go clean -i ./...

test:
	@rm -rf coverage
	@mkdir coverage
	@go test -v ./... -coverpkg=$(PKG_TO_TEST) -coverprofile=coverage/coverage.out -covermode=count 2>&1 | tee coverage/test.log
	@cat coverage/test.log | go-junit-report > coverage/test_result.xml
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@gocover-cobertura < coverage/coverage.out > coverage/coverage.xml

docker:
	@docker build --no-cache --build-arg "FULL_VERSION=${FULL_VERSION}" -t $(DOCKER_IMAGE) .

help:
	@echo "make install: install all dependency"
	@echo "make fmt: format the coding style of this project"
	@echo "make build: build this project to binary file"
	@echo "make clean: remove binary file and dir when building/testing"
