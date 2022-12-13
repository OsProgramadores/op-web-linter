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

	"github.com/osprogramadores/op-web-linter/common"
	"github.com/osprogramadores/op-web-linter/handlers"
	"github.com/osprogramadores/op-web-linter/lang"
)

// API paths.
const (
	lintURLPath      = "/lint"
	languagesURLPath = "/languages"
	pingURLPath      = "/ping"
	staticURLPath    = "/static"
	tmplURLPath      = "/t"
	formTmplFile     = "form.html"
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
		tmpldir   = flag.String("templates", "./t", "Directory where we serve templates")
	)
	flag.Parse()

	// Replace {port} with actual port.
	*apiurl = strings.ReplaceAll(*apiurl, "{port}", fmt.Sprintf("%d", *port))

	u, err := url.Parse(*apiurl)
	if err != nil {
		log.Fatalf("Error parsing URL: %v", err)
	}

	log.Printf("Started op-web-linter, version %s", BuildVersion)
	log.Printf("Listening on port %d", *port)
	log.Printf("URL for API requests: %s", *apiurl)

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

	// Pre-parse templates and register handlers.
	if err := handlers.TmplSetup(*tmpldir, u.Path+tmplURLPath+"/", fe); err != nil {
		log.Fatalf("Error setting up template handlers: %v", err)
	}

	// Everything under staticURLPath is served as a regular file from rootdir.
	// This allows us to keep local javascript files and other accessory files.
	fs := http.FileServer(http.Dir(*staticdir))
	http.Handle(fe.StaticPath, http.StripPrefix(fe.StaticPath, fs))

	// This is a simple /ping handler that just returns "pong" and does not
	// log anything. Useful for health probers.
	http.HandleFunc(u.Path+pingURLPath+"/", func(w http.ResponseWriter, r *http.Request) {
		// Only GET requests.
		if r.Method != "GET" {
			common.HTTPError(w, "Only POST requested accepted", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintln(w, "pong")
	})

	// Main HTML form for interactive access. This is also the "catch-all" URL
	// for anything not matched in the more specific handlers above. The
	// function will emit a 404 if the path is anything other than "/".
	// If everything is OK, it emits a 302 to the form path (served as a template).
	http.HandleFunc(fe.RootPath, func(w http.ResponseWriter, r *http.Request) {
		log.Printf("FORM Request %s %s %s\n", common.RealRemoteAddress(r), r.Method, r.URL)

		if r.URL.Path != fe.RootPath {
			http.NotFound(w, r)
			return
		}
		u := u.Path + tmplURLPath + "/" + formTmplFile
		http.Redirect(w, r, u, http.StatusTemporaryRedirect)
	})

	log.Printf("Serving static files on path: %s", fe.StaticPath)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
