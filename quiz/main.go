package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func cleanString(str string) string {
	str = strings.TrimSpace(strings.ToLower(str))
	return str
}

func main() {

	var csvFile string
	var quizTime int
	var shuffle bool

	flag.StringVar(&csvFile, "csv", "problems.csv", "CSV file path")
	flag.IntVar(&quizTime, "time", 30, "Quiz time in seconds")
	flag.BoolVar(&shuffle, "shuffle", true, "Shuffle the questions")
	flag.Parse()

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(file)
	var questions [][2]string
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		questions = append(questions, [2]string{row[0], row[1]})
	}

	if shuffle {
		for i := range questions {
			j := rand.Intn(i + 1)
			questions[i], questions[j] = questions[j], questions[i]
		}
	}

	// questions
	score := 0
	go func() {
		time.Sleep(time.Duration(quizTime) * time.Second)
		fmt.Print("\n\nQuiz time ran out!!\n")
		fmt.Printf("\nScore: %d\nTotal: %d\n\n", score, len(questions))
		os.Exit(0)
	}()

	for _, question := range questions {
		fmt.Printf("%v: ", question[0])
		var answer string

		fmt.Scan(&answer)

		if cleanString(answer) == cleanString(question[1]) {
			score++
		}
	}
	fmt.Printf("\nScore: %d\nTotal: %d\n\n", score, len(questions))
}
