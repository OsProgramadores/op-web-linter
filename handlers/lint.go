// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/osprogramadores/op-web-linter/common"
)

// LintRequest contains a request to lint a source program.
type LintRequest struct {
	Text string `json:"text"` // Text of the program.
	Lang string `json:"lang"` // Language (must be in SupportedLangs)
}

// LintResponse contains a response to a lint request.
type LintResponse struct {
	Pass          bool     // Pass or not?
	ErrorMessages []string // Used to send global linter failures back (usually blank).
}

// LintRequestHandler handles /lint. The entire JSON request needs
// to be posted as field "request" in the form.
func LintRequestHandler(w http.ResponseWriter, r *http.Request, supported SupportedLangs) {
	var req LintRequest

	log.Printf("LINT Request %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	// Always set CORS.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Set CORS header on preflight request (method == OPTIONS)
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Method", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusNoContent)
		log.Printf("OPTIONS Response headers:")
		for h := range w.Header() {
			for _, v := range w.Header().Values(h) {
				log.Printf("%s: %s", h, v)
			}
		}
		return
	}

	// Only POST request.
	if r.Method != "POST" {
		common.HttpError(w, "Only POST requested accepted", http.StatusMethodNotAllowed)
		return
	}

	// Content-type must be application/json.
	if !strings.Contains(r.Header.Get("content-type"), "application/json") {
		common.HttpError(w, "Incorrect content-type. Expected: application/json", http.StatusUnsupportedMediaType)
		return
	}

	d := json.NewDecoder(r.Body)
	d.Decode(&req)
	log.Printf("Received form data: %+v", req)

	// Program text must not be null.
	if len(req.Text) == 0 {
		common.HttpError(w, "Program text cannot be empty", http.StatusBadRequest)
		return
	}

	// Validate as JSON.
	jreq, err := json.Marshal(req)
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Parsed JSON: %v\n", string(jreq))

	// Test valid languages.
	if !validLang(req.Lang, supported) {
		common.HttpError(w, "Invalid Language", http.StatusBadRequest)
		return
	}

	// Call the appropriate linter.
	supported[req.Lang](w, r, req)
}
