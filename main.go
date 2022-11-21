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

	"github.com/osprogramadores/op-web-linter/handlers"
	"github.com/osprogramadores/op-web-linter/lang"
)

// URL path for static files.
const staticURLPath = "/static/"

// BuildVersion Holds the current git HEAD version number.
// This is filled in by the build process (make).
var BuildVersion string

// supported contains the supported linter languages.
var supported = handlers.SupportedLangs{
	"c":          nil,
	"cpp":        nil,
	"csharp":     nil,
	"java":       nil,
	"javascript": nil,
	"go":         lang.LintGo,
	"php":        nil,
	"python":     nil,
}

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
	*apiurl = strings.ReplaceAll(*apiurl, "{port}", fmt.Sprintf("%d", *port))

	// All information required to serve the form.
	fe := &handlers.Frontend{
		LintPath:  *apiurl + "/lint",
		Languages: handlers.GetLanguagesList(supported),
		StaticDir: *staticdir,
		Template:  template.Must(template.New("form").Parse(tmpl)),
	}

	u, err := url.Parse(*apiurl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	staticpath := u.Path + staticURLPath

	// Main HTML form for interactive access.
	http.HandleFunc(u.Path+"/", fe.FormHandler)

	// Send list of languages back to caller.
	http.HandleFunc(u.Path+"/getlanguages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetLanguagesHandler(w, r, supported)
	})

	// Lint request.
	http.HandleFunc(u.Path+"/lint", func(w http.ResponseWriter, r *http.Request) {
		handlers.LintRequestHandler(w, r, supported)
	})

	// Everything under staticURLPath is served as a regular file from rootdir.
	// This allows us to keep local javascript files and other accessory files.
	fs := http.FileServer(http.Dir(*staticdir))
	http.Handle(staticpath, http.StripPrefix(staticpath, fs))

	log.Printf("Listening on port %d", *port)
	log.Printf("URL for API requests: %s", *apiurl)
	log.Printf("Serving static files on path: %s", staticpath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
