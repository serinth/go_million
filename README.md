This spike is to measure how Go performs when publishing a million messages serialized with protobuf.

# In-Progress #
- Create script to import data to mysql
- Update main.go to actually process the data

# QuickStart #

```bash
mkdir mysql_data
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

This generates data.tsv with 1 million records to be used to import into mysql with the following columns:

ID (auto generated sequential int, 0 indexed)
LIST_ID (fixed to L1)
SEGMENT_ID (fixed to S1)
LAST_ACTIVITY_DATE (Auto generated randomized for -d option)

# Connect to MySQL using CLI #

```bash
docker run --link gomillion_aurora_1:mysql --volume $PWD/data:/var/lib/mysql-files -it --rm mysql mysql -hgomillion_aurora_1 -uroot -pmysql
```
The Mysql container runs with --secure-file-priv flag so files must be put into `/var/lib/mysql-files` to use commands like `LOAD DATA INFILE`.
We can tell which directory by running the following command once the CLI has loaded:
```mysql
SHOW VARIABLES LIKE "secure_file_priv";
```


```mysql
create database close_event;
use close_event;
```

```mysql
CREATE TABLE event_activity
  (ID INT NOT NULL PRIMARY KEY,
  LIST_ID VARCHAR(10) NOT NULL,
  SEGMENT_ID VARCHAR(10) NOT NULL,
  LAST_ACTIVITY_DATE DATETIME
  );
LOAD DATA LOCAL INFILE '/var/lib/mysql-files/data.tsv' INTO TABLE event_activity;
```

## Available Options ##
-n Number of rows. Default is 10

-f output filename. Default is `data.tsv`

-d days back to generate last activity date. Default is 10

# Performance Profiling #

- `--alloc-space` for number of megabytes that have been allocated
- `--inuse-space` for number of megabytes still in use

```bash
go tool pprof --alloc_space go_million http://localhost:8080/debug/pprof/heap
```

Then once in pprof we can use:
`topk` where k is a number of top memory hogs or just omit it for everything
`list` to show which lines of code used the most