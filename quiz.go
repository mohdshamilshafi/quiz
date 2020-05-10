package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "data/problems.csv", "CSV file with \"problem\",\"answer\" tuples")
	timeLimit := flag.Int("limit", 2, "Time limit per question")
	flag.Parse()

	file, err := os.Open(*csvFilename)

	if err != nil {
		exit(fmt.Sprintln("Failed to open %s", *csvFilename))
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		exit(fmt.Sprintln("Failed to parse the file %s", *csvFilename))
	}

	correct := 0

	for i, p := range parseProblems(lines) {
		fmt.Printf("Question #%d: %s\n", i+1, p.q)
		fmt.Print("Answer: ")
		timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
		answerCh := make(chan string)
		go func() {
			var ans string
			fmt.Scanf("%s\n", &ans)
			answerCh <- ans
		}()
		select {
		case <-timer.C:
			fmt.Println("Time expired!")
		case ans := <-answerCh:
			if ans == p.a {
				correct++
			}
			continue
		}
	}
	fmt.Printf("You got %d out of %d question correct! Congrats!\n", correct, len(lines))
}

func parseProblems(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
