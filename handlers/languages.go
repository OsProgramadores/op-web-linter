// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/osprogramadores/op-web-linter/common"
)

// GetLangResponse contains the response to /languages.
type GetLangResponse struct {
	Languages []string `json:"Languages"` // JSON array with the list of languages.
}

// SupportedLangs holds the supported languages.
type SupportedLangs map[string]func(w http.ResponseWriter, r *http.Request, req LintRequest)

// LanguagesHandler defines the handler for /languages.
func LanguagesHandler(w http.ResponseWriter, r *http.Request, supported SupportedLangs) {
	log.Printf("LANGUAGES Request %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
	CORSHandler(w, r)
	if r.Method == "OPTIONS" {
		log.Printf("Got OPTIONS method. Returning.")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	langs := LanguagesList(supported)

	ret, err := json.Marshal(GetLangResponse{Languages: langs})
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}

// LanguagesList returns a string slice with all supported languages.
func LanguagesList(supported SupportedLangs) []string {
	var langs []string
	for lang, function := range supported {
		if function != nil {
			langs = append(langs, lang)
		}
	}
	sort.Strings(langs)
	return langs
}

// validLang returns true if the language is a supported language.
func validLang(lang string, supported SupportedLangs) bool {
	function, ok := supported[lang]
	if ok && function != nil {
		return true
	}
	return false
}
