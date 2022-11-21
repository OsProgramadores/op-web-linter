// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// httpError logs the error and returns the appropriate message & http code.
func httpError(w http.ResponseWriter, msg string, httpcode int) {
	m := fmt.Sprintf("Returned HTTP error %d: %v", httpcode, msg)
	log.Print(m)
	http.Error(w, msg, httpcode)
}
