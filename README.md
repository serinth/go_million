This spike is to measure how Go performs when publishing a million messages serialized with protobuf.

# In-Progress #
- Create script to import data to mysql
- Update main.go to actually process the data

# QuickStart #

```bash
chmod +x build.sh
./build.sh
```

This will update the `docker-compose` file with a volume mapping for mysql to the current working dir.
It will then statically compile the go_million and run the docker-compose for mysql and go_million.

# Generate Data #

```bash
cd data
go build .

./data -n 100000
```

This generates data.tsv with 1 million records to be used to import into mysql

## Available Options ##
-n Number of rows. Default is 10
-f output filename. Default is `data.tsv`
-d days back to generate last activity date. Default is 10
