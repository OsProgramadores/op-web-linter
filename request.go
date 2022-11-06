// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// LintNotice contains a linter "notice" (error, warning, note, etc)
type LintNotice struct {
	Line    int    `json:"line"`    // Line of the notice
	Column  int    `json:"column"`  // Column of the notice (not all linters support this).
	Type    string `json:"type"`    // Type: "error", "warning", "note", "other"
	Message string `json:"message"` // The linter message.
}

// LintRequest contains a request to lint a source program.
type LintRequest struct {
	Text string `json:"text"` // Text of the program.
	Lang string `json:"lang"` // Language (must be in SupportedLangs)
}

// LintResponse contains a response to a lint request.
type LintResponse struct {
	Pass         bool         // Pass or not?
	ErrorMessage string       // Used to send global linter failures back (usually blank).
	Notices      []LintNotice // Linter messages
}

// lintRequestHandler handles /lint. The entire JSON request needs
// to be posted as field "request" in the form.
func lintRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req LintRequest

	log.Printf("Serving request from: %v", r.RemoteAddr)

	// Only POST request.
	if r.Method != "POST" {
		httpError(w, fmt.Errorf("Only POST requested accepted"), http.StatusMethodNotAllowed)
		return
	}

	// Content-type must be application/json.
	if !strings.Contains(r.Header.Get("content-type"), "application/json") {
		httpError(w, fmt.Errorf("Incorrect content-type. Expected: application/json"), http.StatusUnsupportedMediaType)
		return
	}

	d := json.NewDecoder(r.Body)
	d.Decode(&req)
	log.Printf("Received form data: %+v", req)

	// Program text must not be null.
	if len(req.Text) == 0 {
		httpError(w, fmt.Errorf("Program text cannot be empty"), http.StatusBadRequest)
		return
	}

	// Validate as JSON.
	jreq, err := json.Marshal(req)
	if err != nil {
		httpError(w, fmt.Errorf("Invalid json: %v", err), http.StatusBadRequest)
		return
	}
	log.Printf("Parsed JSON: %v\n", string(jreq))

	// TODO: Test valid languages.
	// TODO: Run actual linters, parse and return results.

	return
}
