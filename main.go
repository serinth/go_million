package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	// Delay for a bit so mysql is up
	fmt.Println("Sleeping for 3 seconds")
	time.Sleep(3000 * time.Millisecond)

	db, err := sql.Open("mysql", "root:mysql/?parseTime=true")
	if err != nil && db != nil {
		fmt.Println("Couldn't connect to sql server with error:", err)
	} else {
		fmt.Println("Connected to DB")
	}

}

func getData() {

}
