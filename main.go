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

// URL path for static files.
const staticURLPath = "/static/"

// BuildVersion Holds the current git HEAD version number.
// This is filled in by the build process (make).
var BuildVersion string

//go:embed "frontend/form.tmpl"
var tmpl string

func main() {
	var (
		port      = flag.Int("port", 10000, "Specify the TCP port to listen to")
		host      = flag.String("host", "localhost", "Hostname (used for API requests)")
		staticdir = flag.String("staticdir", "./static", "Directory where we serve static files")
	)
	flag.Parse()

	// All information required to serve the form.
	fe := &frontend{
		LintPath:  fmt.Sprintf("http://%s:%d/lint", *host, *port),
		Languages: getLanguagesList(),
		StaticDir: *staticdir,
		Template:  template.Must(template.New("form").Parse(tmpl)),
	}

	http.HandleFunc("/", fe.formHandler)                  // Serve form.
	http.HandleFunc("/getlanguages", getLanguagesHandler) // Send list of languages back to caller.
	http.HandleFunc("/lint", lintRequestHandler)          // Linter request.

	// Everything under staticURLPath is served as a regular file from rootdir.
	// This allows us to keep local javascript files and other accessory files.
	fs := http.FileServer(http.Dir(*staticdir))
	http.Handle(staticURLPath, http.StripPrefix(staticURLPath, fs))

	url := fmt.Sprintf("http://%s:%d", *host, *port)
	log.Printf("Listening on %s", url)
	log.Printf("Serving static files on %s", url+staticURLPath)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
