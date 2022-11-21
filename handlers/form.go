// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"log"
	"net/http"
	"text/template"

	"github.com/osprogramadores/op-web-linter/common"
)

// Frontend holds the parameters passed to the frontend form.
type Frontend struct {
	LintPath  string
	Languages []string
	StaticDir string
	Template  *template.Template
}

// FormHandler serves the main form to the user.
func (x *Frontend) FormHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving root request from: %v", r.RemoteAddr)

	log.Printf("FORM Request %s %s %s\n", r.RemoteAddr, r.Method, r.URL)

	if err := x.Template.Execute(w, x); err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}