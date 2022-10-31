// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"fmt"
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
	Text []byte `json:"text"` // Text of the program.
	Lang string `json:"lang"` // Language (must be in SupportedLangs)
}

// LintResponse contains a response to a lint request.
type LintResponse struct {
	Pass         bool         // Pass or not?
	ErrorMessage string       // Used to send global linter failures back (usually blank).
	Notices      []LintNotice // Linter messages
}

func lintRequestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Erros do linter serão exibidos aqui") // write data to response
}
