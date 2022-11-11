// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

const (
	execTimeout = time.Duration(5 * time.Second)
)

// SupportedLangs contains the supported linter languages.
var SupportedLangs = map[string]func(w http.ResponseWriter, r *http.Request, req LintRequest){
	"c":          nil,
	"cpp":        nil,
	"csharp":     nil,
	"java":       nil,
	"javascript": nil,
	"go":         lintGo,
	"php":        nil,
	"python":     nil,
}

// GetLangResponse contains the response to /getlanguages.
type GetLangResponse struct {
	Languages []string `json:"Languages"` // JSON array with the list of languages.
}

func getLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Returning list of languages to %v", r.RemoteAddr)

	var langs []string
	for lang, function := range SupportedLangs {
		if function != nil {
			langs = append(langs, lang)
		}
	}
	sort.Strings(langs)

	ret, err := json.Marshal(GetLangResponse{Languages: langs})
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(ret))
}

// validLang returns true if the language is a supported language.
func validLang(lang string) bool {
	function, ok := SupportedLangs[lang]
	if ok && function != nil {
		return true
	}
	return false
}

// saveProgramToFile saves the program in req.text into a temporary
// file and returns the name of the temporary directory and file.
// Callers must use defer os.Removeall(tempdir) in their functions.
func saveProgramToFile(req LintRequest) (string, string, error) {
	tempdir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", "", err
	}
	tempfd, err := os.CreateTemp(tempdir, "*.go")
	if err != nil {
		os.RemoveAll(tempdir)
		return "", "", err
	}
	defer tempfd.Close()

	// Save program text in request to file.
	program, err := base64.StdEncoding.DecodeString(req.Text)
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
