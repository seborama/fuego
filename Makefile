MACHINE := $(shell uname -m)
ifneq ($(MACHINE),aarch64)
	GORACE := -race
endif

.PHONY: deps
deps:
	go mod tidy && go mod download

.PHONY: test
test: deps
	go test -timeout 5s -cover $(GORACE) -parallel 100

.PHONY: generate
generate: deps
	go build -o bin/maptoXXX ./generate/maptoXXX.go
	go generate

.PHONY: lint
lint: deps
	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh

# .PHONY: mutations
# mutations: deps
# 	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
# 	./gomutesting.sh

