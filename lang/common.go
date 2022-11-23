// Package lang defines all language specific components of op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package lang

import (
	"log"
	"net/url"
	"os"

	"github.com/osprogramadores/op-web-linter/handlers"
)

// saveProgramToFile saves the program in req.text into a temporary
// file and returns the name of the temporary directory and file.
// The template parameter specifies how the filename will appear.
// Use "*.foo" to have a temporary filename with extension foo.
// Callers must use defer os.Removeall(tempdir) in their functions.
func saveProgramToFile(req handlers.LintRequest, template string) (string, string, error) {
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
