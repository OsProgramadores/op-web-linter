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

// API paths.
const (
	staticURLPath    = "/static"
	lintURLPath      = "/lint"
	languagesURLPath = "/languages"
)

// BuildVersion Holds the current git HEAD version number.
// This is filled in by the build process (make).
var BuildVersion string

// supported contains the supported linter languages.
var supported = handlers.SupportedLangs{
	"c":          {Display: "C", LintFn: lang.LintC},
	"cpp":        {Display: "C++", LintFn: lang.LintCPP},
	"golang":     {Display: "Go", LintFn: lang.LintGo},
	"java":       {Display: "Java  (reformat only)", LintFn: lang.LintJava},
	"javascript": {Display: "Javascript (lint only)", LintFn: lang.LintJavascript},
	"python":     {Display: "Python  (lint only)", LintFn: lang.LintPython},
}

//go:embed "templates/form.tmpl"
var tmpl string

func main() {
	var (
		port      = flag.Int("port", 10000, "Specify the TCP port to listen to")
		apiurl    = flag.String("url", "http://localhost:{port}", "Base URL for API requests (no slash at the end)")
		staticdir = flag.String("staticdir", "./static", "Directory where we serve static files")
	)
	flag.Parse()

	// Replace {port} with actual port.
	*apiurl = strings.ReplaceAll(*apiurl, "{port}", fmt.Sprintf("%d", *port))

	u, err := url.Parse(*apiurl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}
	// All information required to serve the form. All paths end in slash.
	fe := &handlers.Frontend{
		RootPath:       u.Path + "/",
		LintPath:       u.Path + lintURLPath + "/",
		StaticPath:     u.Path + staticURLPath + "/",
		LanguagesPath:  u.Path + languagesURLPath + "/",
		StaticDir:      *staticdir,
		SupportedLangs: supported,
		Template:       template.Must(template.New("form").Parse(tmpl)),
	}

	// Send list of languages back to caller.
	http.HandleFunc(u.Path+"/languages", func(w http.ResponseWriter, r *http.Request) {
		handlers.LanguagesHandler(w, r, supported)
	})

	// Lint request.
	http.HandleFunc(fe.LintPath, func(w http.ResponseWriter, r *http.Request) {
		handlers.LintRequestHandler(w, r, supported)
	})

	// Everything under staticURLPath is served as a regular file from rootdir.
	// This allows us to keep local javascript files and other accessory files.
	fs := http.FileServer(http.Dir(*staticdir))
	http.Handle(fe.StaticPath, http.StripPrefix(fe.StaticPath, fs))

	// Main HTML form for interactive access. This is also the "catch-all" URL
	// for anything not matched in the more specific handlers above. The
	// function will emit a 404 if the path is anything other than "/".
	http.HandleFunc(fe.RootPath, fe.FormHandler)

	log.Printf("Started op-web-linte, version %s", BuildVersion)
	log.Printf("Listening on port %d", *port)
	log.Printf("URL for API requests: %s", *apiurl)
	log.Printf("Serving static files on path: %s", fe.StaticPath)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
