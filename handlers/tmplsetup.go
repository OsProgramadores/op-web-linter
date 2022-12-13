// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/osprogramadores/op-web-linter/common"
)

// FormData holds the parameters passed to the
type FormData struct {
	RootPath       string         // The base path for the server (default = "/").
	LanguagesPath  string         // Path for API languages calls.
	LintPath       string         // Path for API linter calls.
	SupportedLangs SupportedLangs // Supported Languages.
	StaticDir      string         // Directory for static files.
	StaticPath     string         // Path for static files (/static).
	Template       *template.Template
}

// TmplSetup parses all templates under dir and sets up handlers under path for
// each of the files it finds.
func TmplSetup(dir, path string, tmpldata *FormData) error {
	log.Printf("Setting up templates from: %s", dir)

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		fname := file.Name()
		fpath := filepath.Join(dir, fname)

		if !file.Type().IsRegular() {
			log.Printf("Ignoring %s (not a plain file)", fpath)
			continue
		}

		tmpl, err := template.ParseFiles(fpath)
		if err != nil {
			return err
		}

		// Create handlers for each file.
		urlpath := path + fname
		http.HandleFunc(urlpath, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Serving template %s", urlpath)
			if err := tmpl.Execute(w, tmpldata); err != nil {
				log.Printf("Error serving template: %v", err)
				common.HTTPError(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
		log.Printf("Registered template handler at: %s", urlpath)
	}
	return nil
}
