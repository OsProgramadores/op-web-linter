// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// SupportedLangs contains the supported linter languages.
var SupportedLangs = []string{
	"c", "cpp", "csharp", "java", "javascript", "go", "php", "python",
}

// GetLangResponse contains the response to /getlanguages.
type GetLangResponse struct {
	Languages []string `json:"Languages"` // JSON array with the list of languages.
}

func getLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Returning list of languages to %v", r.RemoteAddr)
	ret, err := json.Marshal(GetLangResponse{Languages: SupportedLangs})
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}
