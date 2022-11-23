// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import "net/http"

// CORSHandler sends the necessary CORS headers.
func CORSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Method", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "*")
}
