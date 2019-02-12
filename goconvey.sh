#!/bin/bash

tput sc ; echo -n -e "Go getting go-mutesting..." ; tput rc
go get github.com/smartystreets/goconvey

tput sc ; echo -n -e "Running go-mutesting...   " ; tput rc
pgrep -f "goconvey.*6020" || \
    ${GOPATH}/bin/goconvey -depth=20 -timeout=10s -excludedDirs=.git,.vscode,.idea -packages=2 -cover -poll=5000ms -port=6020 1>/dev/null 2>&1 &
open http://localhost:6020

