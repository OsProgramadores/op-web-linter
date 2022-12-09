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
	BasePath       string         // The base path for the server (default = "/").
	LintPath       string         // API Linter path (/lint).
	SupportedLangs SupportedLangs // Supported Languages.
	StaticDir      string         // Directory for static files.
	StaticPath     string         // Path for static files (/static).
	Template       *template.Template
}

// FormHandler serves the main form to the user.
func (x *Frontend) FormHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("FORM Request %s %s %s\n", common.RealRemoteAddress(r), r.Method, r.URL)

	// Allow in iframes.
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")

	// This is the catch-all URL since Handle("/",...) will drop anything not
	// caught by other patterns here. If path is not our BasePath ("/" in a
	// system that does not use a reverse proxy), then we're dealing with an
	// unhandled path, so just return 404.
	if r.URL.Path != x.BasePath {
		http.NotFound(w, r)
		return
	}

	if err := x.Template.Execute(w, x); err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
