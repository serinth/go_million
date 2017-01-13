package main

import (
  "fmt"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "time"
)

func main() {

  // Delay for a bit so mysql is up

  fmt.Printf("Sleeping for 3 seconds");

  time.Sleep(3000 * time.Millisecond)

  db, err := sql.Open("mysql", "root:mysql/")
  if err != nil && db != nil {
    fmt.Printf("Couldn't connect to sql server with error:", err)
  } else {
    fmt.Printf("Connected to DB")
  }


}

