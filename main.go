package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// SupportedLangs contains the supported linter languages.
var SupportedLangs = []string{
	"c", "cpp", "csharp", "java", "javascript", "go", "php", "python",
}

// LintRequest contains a request to lint a source program.
type LintRequest struct {
	Text []byte `json:"text"` // Text of the program.
	Lang string `json:"lang"` // Language (must be in SupportedLangs)
}

// LintNotice contains a linter "notice" (error, warning, note, etc)
type LintNotice struct {
	Line    int    `json:"line"`    // Line of the notice
	Column  int    `json:"column"`  // Column of the notice (not all linters support this).
	Type    string `json:"type"`    // Type: "error", "warning", "note", "other"
	Message string `json:"message"` // The linter message.
}

// LintResponse contains a response to a lint request.
type LintResponse struct {
	Pass         boolean      // Pass or not?
	ErrorMessage string       // Used to send global linter failures back (usually blank).
	Notices      []LintNotice // Linter messages
}

func linterResults(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //Parse url parameters passed, then parse the response packet for the POST body (request body)
	// attention: If you do not call ParseForm method, the following data can not be obtained form
	fmt.Println(r.Form) // print information on server side.
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Erros do linter serão exibidos aqui") // write data to response
}

func linterForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func main() {
	http.HandleFunc("/op-web-linter", linterForm)           // linter input route
	http.HandleFunc("/op-web-linter/result", linterResults) // linter results route

	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
