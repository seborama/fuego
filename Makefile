deps:
	go get -d -t -v ./...

test: deps
	go test -timeout 5s -cover -race -parallel 100

lint: deps

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh ||:

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gometalinter.sh ||:

	@echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gomutesting.sh ||:

goconvey:
	./goconvey.sh

