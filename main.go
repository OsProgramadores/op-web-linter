// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var (
		port = flag.Int("port", 10000, "Specify the TCP port to listen to")
	)

	http.HandleFunc("/getlanguages", getLanguagesHandler) // Send list of languages back to caller.
	http.HandleFunc("/lint", lintRequestHandler)          // Linter request.

	log.Printf("Listening on port %d", *port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
