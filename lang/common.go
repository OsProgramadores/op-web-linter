// Package lang defines all language specific components of op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package lang

import (
	"os"
)

// saveProgramToFile saves the program in into a temporary file and returns the
// name of the temporary directory and file.  The template parameter specifies
// how the filename will appear.  Use "*.foo" to have a temporary filename with
// extension foo.  Callers must use defer os.Removeall(tempdir) in their
// functions.
func saveProgramToFile(data string, template string) (string, string, error) {
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

	if _, err = tempfd.Write([]byte(data)); err != nil {
		os.RemoveAll(tempdir)
		return "", "", err
	}
	return tempdir, tempfd.Name(), nil
}
