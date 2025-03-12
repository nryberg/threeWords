package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var wordSpace []string

func threeWords(wordSpace []string, max int) string {
	first := rand.Intn(len(wordSpace))
	second := rand.Intn(len(wordSpace))
	third := rand.Intn(len(wordSpace))

	return wordSpace[first] + "-" + wordSpace[second] + "-" + wordSpace[third]
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080", nil
		//return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
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
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}
	//NewRand(NewSource(time.Now().UTC().UnixNano()))
	rand.Seed(time.Now().UTC().UnixNano())

	http.HandleFunc("/", handler)
	
	serverAddress := "http://localhost" + addr
	if addr != ":8080" {
		// If using Heroku or other platform with custom PORT
		serverAddress = "Server started on port" + addr
	}
	fmt.Println("Server started! Navigate to " + serverAddress + " to see three random words.")
	
	log.Fatal(http.ListenAndServe(addr, nil))
}
