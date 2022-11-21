// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const (
	execTimeout = time.Duration(5 * time.Second)
)

// GetLangResponse contains the response to /getlanguages.
type GetLangResponse struct {
	Languages []string `json:"Languages"` // JSON array with the list of languages.
}

// getLanguagesList returns a string slice with all supported languages.
func getLanguagesList() []string {
	var langs []string
	for lang, function := range SupportedLangs {
		if function != nil {
			langs = append(langs, lang)
		}
	}
	sort.Strings(langs)
	return langs
}

func getLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Returning list of languages to %v", r.RemoteAddr)
	langs := getLanguagesList()

	ret, err := json.Marshal(GetLangResponse{Languages: langs})
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
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

// execute runs the program specified by name with the command-line specified
// in slice args. Returns the error code and a string slice containing all
// non-blank lines in the program's combined output.
func execute(name string, args ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()

	var ret []string
	for _, line := range strings.Split(string(out), "\n") {
		if line != "" {
			ret = append(ret, line)
		}
	}
	return ret, err
}
