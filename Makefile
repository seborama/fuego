MACHINE := $(shell uname -m)
ifneq ($(MACHINE),aarch64)
	GORACE := -race
endif

deps:
	go get -d -t -v ./...

test: deps
	go test -timeout 5s -cover $(GORACE) -parallel 100

lint: deps

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh ||:

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gomutesting.sh ||:

goconvey:
	./goconvey.sh

