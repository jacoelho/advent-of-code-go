# disable default rules
.SUFFIXES:
MAKEFLAGS+=-r -R
GOBIN = $(shell go env GOPATH)/bin
DATE  = $(shell date +%Y%m%d%H%M%S)

default: test

.PHONY: test
test: staticcheck
	go test -race -shuffle=on -v ./...

.PHONY: fmt
fmt:
	go fmt ./...


.PHONY: ci-tidy
ci-tidy:
	go mod tidy
	git status --porcelain go.mod go.sum || { echo "Please run 'go mod tidy'."; exit 1; }

$(GOBIN)/staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest

$(GOBIN)/gcassert:
	go install github.com/jordanlewis/gcassert/cmd/gcassert@latest

.PHONY: staticcheck
staticcheck: $(GOBIN)/staticcheck
	$(GOBIN)/staticcheck ./...