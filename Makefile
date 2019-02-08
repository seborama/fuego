deps:
	go get -d -t -u -v ./...

test: deps
	go test -cover -race -parallel 2 ./...

lint: deps

	echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./golangci-lint.sh
	echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gometalinter.sh
	echo "=~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~="
	./gomutesting.sh

goconvey:
	./goconvey.sh

