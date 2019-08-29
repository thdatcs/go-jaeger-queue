GOPATH := $(PWD)
PATH := "$(PATH):$(GOPATH)/bin"
PKGS := go-jaeger-queue

.PHONY:

prepare:
	@go get -v github.com/golang/dep/cmd/dep
	@go get -v golang.org/x/lint/golint
	@go get -v github.com/golangci/golangci-lint/cmd/golangci-lint

dep:
	@bash $(GOPATH)/scripts/dep.sh $(PKGS)

lint:
	@bash $(GOPATH)/scripts/golangci-lint.sh $(PKGS)
	@bash $(GOPATH)/scripts/golint.sh $(PKGS)

run-consumer:
	@cd $(GOPATH)/src/go-jaeger-queue/cmd/consumer && go run main.go

run-producer:
	@cd $(GOPATH)/src/go-jaeger-queue/cmd/producer && go run main.go