// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"log"
	"net/http"
	"text/template"
)

// frontend holds the parameters passed to the frontend form.
type frontend struct {
	LintPath  string
	Languages []string
	StaticDir string
	Template  *template.Template
}

func (x *frontend) formHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving root request from: %v", r.RemoteAddr)

	if err := x.Template.Execute(w, x); err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
