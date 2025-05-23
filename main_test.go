package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWordGenerator_ThreeWords(t *testing.T) {
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	generator := NewWordGenerator(words)
	
	result := generator.ThreeWords()
	
	parts := strings.Split(result, "-")
	if len(parts) != 3 {
		t.Errorf("Expected 3 words separated by hyphens, got: %s", result)
	}
	
	for _, part := range parts {
		found := false
		for _, word := range words {
			if part == word {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Word '%s' not found in the provided word list", part)
		}
	}
}

func TestDetermineListenAddress(t *testing.T) {
	// Test default case
	os.Unsetenv("PORT")
	addr, err := determineListenAddress()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if addr != ":8080" {
		t.Errorf("Expected address :8080, got %s", addr)
	}
	
	// Test with PORT set
	os.Setenv("PORT", "9000")
	addr, err = determineListenAddress()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if addr != ":9000" {
		t.Errorf("Expected address :9000, got %s", addr)
	}
}

func TestLoadWords(t *testing.T) {
	// Skip this test if the word file doesn't exist locally
	if _, err := os.Stat("many_words.txt"); os.IsNotExist(err) {
		t.Skip("Skipping test as many_words.txt doesn't exist")
	}
	
	words, err := loadWords("many_words.txt")
	if err != nil {
		t.Fatalf("Failed to load words: %v", err)
	}
	
	if len(words) == 0 {
		t.Error("Expected words to be loaded, but got empty list")
	}
}

func TestHTTPHandler(t *testing.T) {
	// Create test words
	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
	generator := NewWordGenerator(words)
	
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	
	// Create the handler function
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, generator.ThreeWords())
	})
	
	// Our handler fulfills http.Handler, so we can call ServeHTTP method directly
	handler.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check that response contains three words separated by hyphens
	response := rr.Body.String()
	parts := strings.Split(response, "-")
	if len(parts) != 3 {
		t.Errorf("Expected 3 words separated by hyphens, got: %s", response)
	}
}