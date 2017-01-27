package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//Performance Profiling
	go func() {
		http.ListenAndServe(":8080", http.DefaultServeMux)
	}()

	// Delay for a bit so mysql is up
	fmt.Println("Sleeping for 3 seconds")
	time.Sleep(3000 * time.Millisecond)

	db, err := sql.Open("mysql", "root:mysql@tcp(aurora:3306)/close_event?parseTime=true")
	if err != nil && db != nil {
		fmt.Println("Couldn't connect to sql server with error:", err)
	} else {
		fmt.Println("Connected to DB")
	}

	fmt.Println("Executing Query...")
	rows, err := db.Query("SELECT ID, LIST_ID, SEGMENT_ID, LAST_ACTIVITY_DATE FROM event_activity")
	defer rows.Close()
	check(err)

	fmt.Println("Printing Results...")

	for rows.Next() {
		var id int
		var listId string
		var segmentId string
		var lastActivityDate time.Time
		if err := rows.Scan(&id, &listId, &segmentId, &lastActivityDate); err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("%d \t %s \t %s \t %v \n", id, listId, segmentId, lastActivityDate)
	}

	fmt.Println("Done, just listening now for profiling")
	select {}
}

func check(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}
