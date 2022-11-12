// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package main

import (
	"encoding/json"
	"net/http"
	"os"
)

// runGolint runs golint on the source file and returns the output.
func runGolint(fname string) ([]string, bool, error) {
	// Golint to always exits with code 0 (no error). Any output
	// means the input program contains errors.
	out, err := execute("golint", fname)
	if err != nil {
		return slicePrefix(out, "golint"), false, err
	}
	// No errors in the program.
	if len(out) == 0 {
		return []string{}, true, nil
	}
	return slicePrefix(out, "golint"), false, nil
}

// runGoBuild runs "go build" on the source file and returns the output.
func runGoBuild(fname string) ([]string, bool) {
	out, err := execute("go", "build", fname)
	retcode := exitcode(err)

	// No errors.
	if retcode == 0 {
		return []string{}, true
	}
	return slicePrefix(out, "go build"), false
}

// runGoFmt runs "go fmt -d" on the source file and indicates if any output
// exists (this means the user needs to run gofmt on their source).
func runGoFmt(fname string) ([]string, bool) {
	out, err := execute("gofmt", "-d", fname)
	retcode := exitcode(err)

	// No errors.
	if retcode != 0 || len(out) != 0 {
		ret := []string{"Gofmt detected differences. Please run \"gofmt\" to fix this"}
		return slicePrefix(ret, "gofmt"), len(out) == 0
	}
	return []string{""}, true
}

// lintGo is a test linter for a fake "test" language.
func lintGo(w http.ResponseWriter, r *http.Request, req LintRequest) {
	tempdir, tempfile, err := saveProgramToFile(req)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Golint.
	m, ok, err := runGolint(tempfile)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	if !ok {
		messages = append(messages, m...)
	}

	// Go Build.
	m, ok = runGoBuild(tempfile)
	if !ok {
		messages = append(messages, m...)
	}

	// Go fmt.
	m, ok = runGoFmt(tempfile)
	if !ok {
		messages = append(messages, m...)
	}

	// Return an empty JSON array if no messages
	if len(messages) == 0 {
		messages = []string{}
	}

	// Create response, convert to JSON and return.
	resp := LintResponse{
		Pass:          len(messages) == 0,
		ErrorMessages: messages,
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		httpError(w, err, http.StatusInternalServerError)
		return
	}
	w.Write(jresp)
	w.Write([]byte("\n"))
}