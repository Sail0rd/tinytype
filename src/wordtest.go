package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const wordUrl string = "https://random-word-api.herokuapp.com/word"

func generateWordTest(name string, n int, g int) func() []segment {
	var b []byte

	if b = readResource("words", name); b == nil {
		die("%s does not appear to be a valid word list. See '-list words' for a list of builtin word lists.", name)
	}

	words := regexp.MustCompile("\\s+").Split(string(b), -1)

	return func() []segment {
		segments := make([]segment, g)
		for i := 0; i < g; i++ {
			segments[i] = segment{randomText(n, words), ""}
		}

		return segments
	}
}

// getWords returns a string of multiple lowercased random words from the API
func getWordTest(nb int) func() []segment {
	req := wordUrl + "?number=" + strconv.Itoa(nb)
	resp, err := http.Get(req)
	if err != nil {
		die("Failed to get words: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		die("Failed to read words: %v", err)
	}

	var words []string
	err = json.Unmarshal(body, &words)
	if err != nil {
		die("Failed to parse words: %v", err)
	}
	return func() []segment {
		return []segment{{strings.Join(words, " "), ""}}
	}
}
