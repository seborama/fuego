#!/bin/bash

golangci-lint run ./... --enable-all --disable=dupl

