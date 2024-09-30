package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const quoteUrl string = "http://api.quotable.io/random"

type QuoteResponse struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

type MonkeytypeQuote struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Length int    `json:"length"`
}

func generateQuoteTest(name string, lengthThreshold int) func() []segment {
	rand.Seed(time.Now().UnixNano())

	file, err := os.Open(name)
	if err != nil {
		die("Failed to open file: %v", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	_, err = decoder.Token()
	if err != nil {
		die("Failed to read JSON token: %v", err)
	}

	// First pass: Count quotes that meet the length threshold
	var matchingQuotes []int
	index := 0
	for decoder.More() {
		var quote MonkeytypeQuote
		err := decoder.Decode(&quote)
		if err != nil {
			die("Error decoding quote: %v", err)
		}

		// Collect the index of quotes that meet the length threshold
		if quote.Length >= lengthThreshold {
			matchingQuotes = append(matchingQuotes, index)
		}
		index++
	}

	// If no quotes meet the length threshold, handle it
	if len(matchingQuotes) == 0 {
		die("No quote found with length over %d.", lengthThreshold)
	}
	randomIndex := matchingQuotes[rand.Intn(len(matchingQuotes))]
	_, err = file.Seek(0, 0) // Reset file reader to the beginning
	if err != nil {
		die("Unable to reset file: %v", err)
	}
	decoder = json.NewDecoder(file)

	// Read the opening bracket again
	_, err = decoder.Token() // Skip the '[' again
	if err != nil {
		die("Improperly formatted quote file: %v", err)
	}

	// Iterate to the randomIndex
	index = 0
	for decoder.More() {
		var quote MonkeytypeQuote
		err := decoder.Decode(&quote)
		if err != nil {
			die("Error decoding quote: %v", err)
		}

		if index == randomIndex {
			return func() []segment {
				return []segment{{quote.Text, quote.Source}}
			}
		}
		index++
	}

	// If something goes wrong, return an empty quote
	die("Error: Could not retrieve random quote.")
	return nil
}

// getQuote returns a random quote and its author from the API
func getWebQuoteTest(length int) func() []segment {
	if length > 430 {
		length = 430
	}
	req := quoteUrl + "?minLength=" + strconv.Itoa(length)
	resp, err := http.Get(req)
	if err != nil {
		die("Failed to get quote: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		die("Failed to read quote: %v", err)
	}

	var quote QuoteResponse
	err = json.Unmarshal(body, &quote)
	if err != nil {
		die("Failed to parse quote: %v", err)
	}

	return func() []segment {
		return []segment{{quote.Content, quote.Author}}
	}
}
