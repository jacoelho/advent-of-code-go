# disable default rules
.SUFFIXES:
MAKEFLAGS+=-r -R
GOBIN = $(shell go env GOPATH)/bin
DATE  = $(shell date +%Y%m%d%H%M%S)

default: test

.PHONY: test
test:
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

.PHONY: tmpl-%
tmpl-%:
	@year=$(shell echo $* | cut -d- -f1); \
	day=$(shell echo $* | cut -d- -f2); \
	go run ./cmd/template/template.go -year $$year -day $$day

.PHONY: input-fetch-%
input-fetch-%:
	@year=$(shell echo $* | cut -d- -f1); \
	day=$(shell echo $* | cut -d- -f2); \
	go run ./cmd/input-fetch/main.go -year $$year -day $$day

test-timings-runner: cmd/test-timings/main.go
	go build -o test-timings-runner ./cmd/test-timings

.PHONY: test-timings-%
test-timings-%: test-timings-runner
	@year=$(shell echo $* | cut -d- -f1); \
	./test-timings-runner -year $$year

.PHONY: test-timings
test-timings: test-timings-runner
	./test-timings-runner

.PHONY: test-%
test-%:
	go test -race -shuffle=on -timeout=2m -v ./internal/aoc$*/...