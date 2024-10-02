package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strconv"
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
	var monkeytypequotes []MonkeytypeQuote

	if b := readResource("quotes", name); b == nil {
		die("%s does not appear to be a valid quote file. See '-list quotes' for a list of builtin quotes.", name)
	} else {
		err := json.Unmarshal(b, &monkeytypequotes)
		if err != nil {
			die("Improperly formatted quote file: %v", err)
		}
	}

	// Get quotes that meet the length threshold
	var quotes []segment
	for _, quote := range monkeytypequotes {
		// Collect the index of quotes that meet the length threshold
		if quote.Length >= lengthThreshold {
			quotes = append(quotes, segment{quote.Text, quote.Source})
		}
	}

	// If no quotes meet the length threshold, handle it
	if len(quotes) == 0 {
		die("No quote found with length over %d.", lengthThreshold)
	}

	return func() []segment {
		idx := rand.Int() % len(quotes)
		return quotes[idx : idx+1]
	}
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
