package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func produceMap(file string) (map[string]string, error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	reader := csv.NewReader(fd)
	qta := map[string]string{}
	for {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return qta, err
		}
		qta[row[0]] = row[1]
	}
}

func printResult(c, q int) {
	fmt.Println("Correct answers:", c)
	fmt.Println("Total answers:", q)
}

func main() {
	fileNamePtr := flag.String("csv", "problems.csv", "a csv file containing questions and answers")
	timePtr := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	qta, err := produceMap(*fileNamePtr)
	if err != nil {
		fmt.Println("Failed: ", err)
		return
	}
	qTotal := 0
	cTotal := 0
	fmt.Println(*timePtr)

	c1 := make(chan bool, 1)

	go func() {
		fmt.Println("Welcome to the timed quiz")
		for q, a := range qta {
			fmt.Println(q)
			var answer string
			fmt.Scanln(&answer)
			if answer == a {
				fmt.Println("Correct!")
				cTotal++
			} else {
				fmt.Println("Incorrect")
			}
			qTotal++
		}
		c1 <- true
	}()

	select {
	case <-c1:
		printResult(cTotal, qTotal)
	case <-time.After(time.Duration(*timePtr) * time.Second):
		fmt.Println("Time ran out")
		printResult(cTotal, qTotal)
	}
}
