package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// WordGenerator handles generation of random three-word combinations
type WordGenerator struct {
	words []string
	mu    sync.Mutex // Protects rand
	rand  *rand.Rand
}

// NewWordGenerator creates a new word generator with the given words
func NewWordGenerator(words []string) *WordGenerator {
	source := rand.NewSource(time.Now().UTC().UnixNano())
	return &WordGenerator{
		words: words,
		rand:  rand.New(source),
	}
}

// ThreeWords generates a string of three random words separated by hyphens
func (wg *WordGenerator) ThreeWords() string {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	first := wg.rand.Intn(len(wg.words))
	second := wg.rand.Intn(len(wg.words))
	third := wg.rand.Intn(len(wg.words))

	var sb strings.Builder
	sb.WriteString(wg.words[first])
	sb.WriteString("-")
	sb.WriteString(wg.words[second])
	sb.WriteString("-")
	sb.WriteString(wg.words[third])

	return sb.String()
}

// loadWords loads words from a file
func loadWords(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("word list is empty")
	}

	return strings.Fields(lines[0]), nil
}

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return ":80", nil // Default to port 80
	}
	return ":" + port, nil
}

func main() {
	// Determine listen address
	listenAddr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	// Load words file
	wordFilePath := os.Getenv("WORD_FILE_PATH")
	if wordFilePath == "" {
		wordFilePath = "many_words.txt" // Default
	}

	words, err := loadWords(wordFilePath)
	if err != nil {
		log.Fatalf("Failed to load words: %v", err)
	}

	// Create word generator
	wordGen := NewWordGenerator(words)

	// Create server and routes
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, wordGen.ThreeWords())
	})

	server := &http.Server{
		Addr:    listenAddr,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		serverURL := "http://localhost" + listenAddr
		if listenAddr != ":8080" {
			serverURL = "Server started on port" + listenAddr
		}
		fmt.Println("Server started! Navigate to " + serverURL + " to see three random words.")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Shutdown gracefully with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("\nShutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
