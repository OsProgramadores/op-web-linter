// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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
	var (
		port = flag.Int("port", 10000, "Specify the TCP port to listen to")
	)

	http.HandleFunc("/getlanguages", getLanguagesHandler) // Send list of languages back to caller.

	log.Printf("Listening on port %d", *port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
