// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// BuildVersion Holds the current git HEAD version number.
// This is filled in by the build process (make).
var BuildVersion string

//go:embed "frontend/form.tmpl"
var tmpl string

func main() {
	var (
		port = flag.Int("port", 10000, "Specify the TCP port to listen to")
	)

	fe := &frontend{
		LintPath:  fmt.Sprintf("http://localhost:%d/lint", *port),
		Languages: getLanguagesList(),
		Template:  template.Must(template.New("form").Parse(tmpl)),
	}

	http.HandleFunc("/", fe.formHandler)                  // Serve form.
	http.HandleFunc("/getlanguages", getLanguagesHandler) // Send list of languages back to caller.
	http.HandleFunc("/lint", lintRequestHandler)          // Linter request.

	log.Printf("Listening on port %d", *port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
