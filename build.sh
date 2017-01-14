#!/bin/bash

export CGO_ENABLED=0 
export GOOS=linux 

mkdir -p $PWD/mysql_data

sed -i s@#_VOLUME_MAPPING_#@$PWD:/var/lib/mysql@g docker-compose.yml 

go build -a --ldflags="-s" --installsuffix cgo -o go_million

docker-compose up
