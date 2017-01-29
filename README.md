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

# Spike Findings #

Context: Events greater than now minus 7 days ago (Those still in 7 day segment)
Rows: 1 million

## SQL Query ##

```
mysql> select count(*) from event_activity where LAST_ACTIVITY_DATE > DATE_SUB(NOW(), INTERVAL 7 DAY);
+----------+
| count(*) |
+----------+
|   499614 |
+----------+
1 row in set (0.21 sec)
```

## Processing in Go ##

```
app_1     | Executing Query...
app_1     | Printing Results...
app_1     | Number of subscribers who's last activity date is greater than 7 days ago:  499614
app_1     | Processing Time including SQL connection:  1.579584004s
```

### Memory Usage ###

**Areas of code consuming the most memory:**
```
ROUTINE ======================== main.main in main.go
      66MB   261.01MB (flat, cum)   100% of Total
         .          .     37:	check(err)
         .          .     38:	fmt.Println("Printing Results...")
         .          .     39:
         .          .     40:	// Actually process the rows and see which ones don't fall in our last activity date
         .          .     41:	var cnt = 0
         .   183.51MB     42:	for rows.Next() {
       8MB        8MB     43:		var id int
      15MB       15MB     44:		var listId string
      12MB       12MB     45:		var segmentId string
      31MB       31MB     46:		var lastActivityDate time.Time
         .    11.50MB     47:		if err := rows.Scan(&id, &listId, &segmentId, &lastActivityDate); err != nil {
         .          .     48:			log.Fatal(err)
         .          .     49:		}
         .          .     50:
         .          .     51:		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
         .          .     52:
ROUTINE ======================== runtime.goexit in /usr/local/go/src/runtime/asm_amd64.s
         0   261.01MB (flat, cum)   100% of Total
```
**Top memory consumption areas:**

```
261.01MB of 261.01MB total (  100%)
      flat  flat%   sum%        cum   cum%
  183.51MB 70.31% 70.31%   183.51MB 70.31%  github.com/go-sql-driver/mysql.(*textRows).readRow
      66MB 25.29% 95.59%   261.01MB   100%  main.main
       7MB  2.68% 98.28%    11.50MB  4.41%  database/sql.convertAssign
    4.50MB  1.72%   100%     4.50MB  1.72%  database/sql.asString
         0     0%   100%   183.51MB 70.31%  database/sql.(*Rows).Next
         0     0%   100%    11.50MB  4.41%  database/sql.(*Rows).Scan
         0     0%   100%   183.51MB 70.31%  github.com/go-sql-driver/mysql.(*textRows).Next
         0     0%   100%   261.01MB   100%  runtime.goexit
         0     0%   100%   261.01MB   100%  runtime.main
```


