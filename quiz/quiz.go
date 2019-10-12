package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	FILEPATH = "data/hard.csv"
	QUIZTIME = 20 * time.Second
)

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err)
	}
	return file
}

func extractQuestionAnswers (file *os.File) [][]string {
	var output [][]string
	reader := csv.NewReader(file)
	for {
		if line, err := reader.Read(); err == nil {
			output = append(output, line)
		} else {
			return output
		}
	}
}

func makeQuiz (questionAnswers [][]string, input *bufio.Reader, correctCountCh chan int, quizFinishedCh chan int) {
	for _, questionAnswer := range questionAnswers {
		fmt.Print(questionAnswer[0] + " > ")
		if answer, err := input.ReadString('\n'); err == nil {
			if strings.TrimRight(answer, "\n") == questionAnswer[1] {
				correctCountCh <- 1
			}
		}
	}
	quizFinishedCh <- 1
}

func main() {
	file := openFile(FILEPATH)
	questionAnswers := extractQuestionAnswers(file)

	input := bufio.NewReader(os.Stdin)

	quizTimeoutEvent := time.After(QUIZTIME)
	correctCountCh, quizFinishedCh := make(chan int), make(chan int)
	go makeQuiz(questionAnswers, input, correctCountCh, quizFinishedCh)

	var correctAnswerCount int
	for {
		select {
		case <-quizTimeoutEvent:
			fmt.Println("\nQuiz timeout.")
			fmt.Println(correctAnswerCount, "out of", len(questionAnswers))
			return
		case <-quizFinishedCh:
			fmt.Print("\nQuiz finished.")
			fmt.Println(correctAnswerCount, "out of", len(questionAnswers))
			return
		case <-correctCountCh:
			correctAnswerCount++
		}
	}
}
