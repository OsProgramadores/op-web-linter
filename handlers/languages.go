// Package handlers contains http handler code for op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"

	"github.com/osprogramadores/op-web-linter/common"
)

// GetLangResponse contains the response to /getlanguages.
type GetLangResponse struct {
	Languages []string `json:"Languages"` // JSON array with the list of languages.
}

// SupportedLangs holds the supported languages.
type SupportedLangs map[string]func(w http.ResponseWriter, r *http.Request, req LintRequest)

// GetLanguagesHandler defines the handler for /getlanguages.
func GetLanguagesHandler(w http.ResponseWriter, r *http.Request, supported SupportedLangs) {
	log.Printf("Returning list of languages to %v", r.RemoteAddr)
	langs := GetLanguagesList(supported)

	ret, err := json.Marshal(GetLangResponse{Languages: langs})
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}

// SaveProgramToFile saves the program in req.text into a temporary
// file and returns the name of the temporary directory and file.
// The template parameter specifies how the filename will appear.
// Use "*.foo" to have a temporary filename with extension foo.
// Callers must use defer os.Removeall(tempdir) in their functions.
func SaveProgramToFile(req LintRequest, template string) (string, string, error) {
	tempdir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", "", err
	}
	tempfd, err := os.CreateTemp(tempdir, template)
	if err != nil {
		os.RemoveAll(tempdir)
		return "", "", err
	}
	defer tempfd.Close()

	// Save program text in request to file.
	program, err := url.QueryUnescape(req.Text)
	if err != nil {
		return "", "", err
	}

	log.Printf("Decoded program: %s\n", program)
	if _, err = tempfd.Write([]byte(program)); err != nil {
		os.RemoveAll(tempdir)
		return "", "", err
	}
	return tempdir, tempfd.Name(), nil
}

// GetLanguagesList returns a string slice with all supported languages.
func GetLanguagesList(supported SupportedLangs) []string {
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
