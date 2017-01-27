#!/bin/bash
export CGO_ENABLED=0 
export GOOS=linux 

go get
go build -a --ldflags="-s" --installsuffix cgo -o go_million
./go_million