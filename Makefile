MACHINE := $(shell uname -m)
ifneq ($(MACHINE),aarch64)
	GORACE := -race
endif

.PHONY: deps
deps:
	go get -d -t -v ./...

.PHONY: test
test: deps
	go test -timeout 5s -cover $(GORACE) -parallel 100

.PHONY: lint
lint: deps
	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh

.PHONY: mutations
mutations: deps
	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gomutesting.sh
