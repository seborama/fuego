deps:
	go get -d -t -v ./...

test: deps
	go test -cover -race -parallel 5 ./...

lint: deps

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh ||:

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gometalinter.sh ||:

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gomutesting.sh ||:

goconvey:
	./goconvey.sh

