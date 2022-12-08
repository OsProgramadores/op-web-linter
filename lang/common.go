// Package lang defines all language specific components of op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package lang

import (
	"log"
	"net/url"
	"os"
)

// saveRequestToFile unescapes the passed request data and saves it into a
// temporary file, returning its directory and name. The template parameter
// specifies how the filename will appear.  Use "*.foo" to have a temporary
// filename with extension foo.  Callers must use defer os.Removeall(tempdir)
// in their functions.
func saveRequestToFile(data string, template string) (string, string, error) {
	unescaped, err := url.QueryUnescape(data)
	if err != nil {
		return "", "", err
	}
	log.Printf("Decoded Request: %s\n", unescaped)

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

	if _, err = tempfd.Write([]byte(unescaped)); err != nil {
		os.RemoveAll(tempdir)
		return "", "", err
	}
	return tempdir, tempfd.Name(), nil
}
