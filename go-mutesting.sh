#!/bin/bash

tput sc ; echo -n -e "Go getting go-mutesting..." ; tput rc
go get -t -v github.com/zimmski/go-mutesting/...

tput sc ; echo -n -e "Running go-mutesting...   " ; tput rc
go-mutesting --test-recursive ${GOPATH}/src/github.com/seborama/fuego/... | grep -Ev "^PASS"

