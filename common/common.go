// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package common

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
)

// Maximum output line length (in characters).
const outputLineLength = 100

// HTTPError logs the error and returns the appropriate message & http code.
func HTTPError(w http.ResponseWriter, msg string, httpcode int) {
	m := fmt.Sprintf("Returned HTTP error %d: %v", httpcode, msg)
	log.Print(m)
	http.Error(w, msg, httpcode)
}

// SlicePrefix adds a prefix to every string line in the slice.
func SlicePrefix(slice []string, prefix string) []string {
	var ret []string

	// Print prefix on first line, indent on following ones.
	indent := strings.Repeat(" ", len(prefix)+2)

	for _, line := range slice {
		for i, subline := range wordwrap(line, outputLineLength) {
			// Add string on first line, indent on following ones.
			p := indent
			if i == 0 {
				p = "[" + prefix + "]"
			}
			ret = append(ret, fmt.Sprintf("%s %s", p, html.EscapeString(subline)))
		}
	}
	return ret
}

// wordwrap wraps the string to the number of columns given and returns a
// string slice with the new (wrapped) lines.
func wordwrap(s string, max int) []string {
	var (
		ret    []string
		totlen int
		word   string
		words  []string
	)

	for _, word = range strings.Split(s, " ") {
		if totlen+len(word) > max {
			ret = append(ret, strings.Join(words, " "))
			words = nil
			totlen = 0
		}
		words = append(words, word)
		// +1 == space at the end of the word.
		totlen += len(word) + 1
	}

	if totlen != 0 {
		ret = append(ret, strings.Join(words, " "))
	}
	return ret
}
