package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func threeWords(wordSpace []string, max int) string {
	first := rand.Intn(len(wordSpace))
	second := rand.Intn(len(wordSpace))
	third := rand.Intn(len(wordSpace))

	return wordSpace[first] + "-" + wordSpace[second] + "-" + wordSpace[third]
}

func handler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("many_words.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	words := strings.Fields(lines[0])
	max := len(words)

	result := threeWords(words, max)
	fmt.Fprintf(w, result)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
