#!/bin/bash

export CGO_ENABLED=0 
export GOOS=linux 

go build -a --ldflags="-s" --installsuffix cgo -o go_million

