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

func generateQuoteTest(name string) func() []segment {
	var quotes []segment

	if b := readResource("quotes", name); b == nil {
		die("%s does not appear to be a valid quote file. See '-list quotes' for a list of builtin quotes.", name)
	} else {
		err := json.Unmarshal(b, &quotes)
		if err != nil {
			die("Improperly formatted quote file: %v", err)
		}
	}

	return func() []segment {
		idx := rand.Int() % len(quotes)
		return quotes[idx : idx+1]
	}
}

// getQuote returns a random quote and its author from the API
func getQuoteTest(length int) func() []segment {
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
