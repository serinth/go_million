package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Pallinder/go-randomdata"
)

var outputFile = flag.String("f", "data.tsv", "output filename")
var numRepetitions = flag.Int("n", 3, "number of rows to generate")
var daysBack = flag.Int("d", 10, "number of days to randomize for last activity day")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	flag.Parse()
	pwd, _ := os.Getwd()
	now := time.Now()

	fmt.Println("file output:", *outputFile)
	fmt.Println("number of reps:", *numRepetitions)
	fmt.Println("days back to randomize:", *daysBack)
	fmt.Println("working dir:", pwd)

	f, err := os.Create(pwd + "/" + *outputFile)
	check(err)
	defer f.Close()

	for i := 0; i < *numRepetitions; i++ {
		r := now.AddDate(0, 0, randomdata.Number(*daysBack)*-1)
		//Only one list L1 and one segment S1
		_, err := f.WriteString(strconv.Itoa(i) + "\t" + "L1\t" + "S1\t" + r.Format("2006-01-02 15:04:05") + "\n")
		check(err)
	}
	f.Sync()

}
