GO := go
GOFLAGS=-mod=vendor

GOARCH := amd64
GOOS := linux

GOLINT_VERSION:=1.24.0

BUILD_FLAGS = GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) GOFLAGS=$(GOFLAGS)
APPLICATION_NAME := anodot-kube-events

AWS_REGION := us-east-1
AWS_ECR := 932213950603.dkr.ecr.$(AWS_REGION).amazonaws.com
DOCKER_IMAGE_NAME := $(AWS_ECR)/$(APPLICATION_NAME)

ENVIRONMENT := dev

VERSION := $(shell grep 'VERSION' pkg/version/version.go | awk '{ print $$4 }' | tr -d '"')
PREVIOUS_VERSION := $(shell git show HEAD:pkg/version/version.go | grep 'VERSION' | awk '{ print $$4 }' | tr -d '"' )
GIT_COMMIT := $(shell git describe --dirty --always)

all: clean format vet test build build-container
publish-container: clean format vet test build build-container push-container
lint: check-formatting errorcheck vet
test-all: test build

clean:
	@rm -rf $(APPLICATION_NAME)
	docker rmi -f `docker images $(DOCKER_IMAGE_NAME):$(VERSION) -a -q` || true

check-formatting:
	./utils/check_if_formatted.sh

format:
	@$(GO) fmt ./...

vet:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $$(go env GOPATH)/bin v$(GOLINT_VERSION)
	$(BUILD_FLAGS) $$(go env GOPATH)/bin/golangci-lint run

errorcheck: install-errcheck
	$$(go env GOPATH)/bin/errcheck ./pkg/...

install-errcheck:
	which errcheck || GO111MODULE=off go get -u github.com/kisielk/errcheck

build:
	@echo ">> building binaries with version $(VERSION)"
	$(BUILD_FLAGS) $(GO) build -ldflags "-s -w -X github.com/anodot/github.com/anodot/kube-events/pkg/version.REVISION=$(GIT_COMMIT)" -o $(APPLICATION_NAME)

build-container: build
	docker build -t $(APPLICATION_NAME):$(VERSION) .
	docker tag $(APPLICATION_NAME):$(VERSION) $(DOCKER_IMAGE_NAME):$(VERSION)
	@echo ">> created docker image $(DOCKER_IMAGE_NAME):$(VERSION)"

test:
	GOFLAGS=$(GOFLAGS) $(GO) test -v -race -coverprofile=coverage.txt -covermode=atomic -timeout 10s ./pkg/...

push-container:
	aws ecr get-login-password --region $(AWS_REGION) | docker login --username AWS --password-stdin $(AWS_ECR)
	docker push $(DOCKER_IMAGE_NAME):$(VERSION)

dockerhub-login:
	docker login -u $(DOCKER_USERNAME) -p $(DOCKER_PASSWORD)

version-set:
	echo "Version $(VERSION) set in code, deployment, chart"

vendor-update:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor

deploy:
	helm upgrade $(APPLICATION_NAME) ./helm --install --values=helm/values-$(ENVIRONMENT).yaml --set base-chart.image.tag=$(VERSION) -n $(APPLICATION_NAME) --debug