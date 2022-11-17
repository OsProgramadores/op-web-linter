// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

// URL path for static files.
const staticURLPath = "/static/"

// BuildVersion Holds the current git HEAD version number.
// This is filled in by the build process (make).
var BuildVersion string

//go:embed "templates/form.tmpl"
var tmpl string

func main() {
	var (
		port      = flag.Int("port", 10000, "Specify the TCP port to listen to")
		apiurl    = flag.String("url", "http://localhost:{port}", "Base URL for API requests")
		staticdir = flag.String("staticdir", "./static", "Directory where we serve static files")
	)
	flag.Parse()

	// Replace {port} with actual port.
	*apiurl = strings.ReplaceAll(*apiurl, "{port}", fmt.Sprintf("%d", port))

	// All information required to serve the form.
	fe := &frontend{
		LintPath:  *apiurl + "/lint",
		Languages: getLanguagesList(),
		StaticDir: *staticdir,
		Template:  template.Must(template.New("form").Parse(tmpl)),
	}

	u, err := url.Parse(*apiurl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	http.HandleFunc(u.Path+"/", fe.formHandler)                  // Serve form.
	http.HandleFunc(u.Path+"/getlanguages", getLanguagesHandler) // Send list of languages back to caller.
	http.HandleFunc(u.Path+"/lint", lintRequestHandler)          // Linter request.

	// Everything under staticURLPath is served as a regular file from rootdir.
	// This allows us to keep local javascript files and other accessory files.
	fs := http.FileServer(http.Dir(*staticdir))
	http.Handle(staticURLPath, http.StripPrefix(staticURLPath, fs))

	log.Printf("Listening on port %d", *port)
	log.Printf("URL for API requests: %s", *apiurl)
	log.Printf("Serving static files on %s", *apiurl+staticURLPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
