// Package lang defines all language specific components of op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package lang

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/osprogramadores/op-web-linter/common"
	"github.com/osprogramadores/op-web-linter/handlers"
)

// Regexp matching go build and go lint lines.
var goLineRegex = regexp.MustCompile("^([^:]+):([0-9]+):([0-9]+):[ ]*(.*)")

// LintGo lints programs written in Go.
func LintGo(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	tempdir, tempfile, err := handlers.SaveProgramToFile(req, "*.go")
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Golint.
	m, ok, err := runGolint(tempfile)
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		messages = append(messages, m...)
	}

	// Go Build.
	m, ok = runGoBuild(tempdir, tempfile)
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
	resp := handlers.LintResponse{
		Pass:          len(messages) == 0,
		ErrorMessages: messages,
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("JSON response: %v", string(jresp))
	w.Write(jresp)
	w.Write([]byte("\n"))
}

// runGolint runs golint on the source file and returns the output.
func runGolint(fname string) ([]string, bool, error) {
	// Golint to always exits with code 0 (no error). Any output
	// means the input program contains errors.
	out, err := Execute("golint", fname)
	if err != nil {
		return common.SlicePrefix(goErrorParse(out), "golint"), false, err
	}
	// No errors in the program.
	if len(out) == 0 {
		return []string{}, true, nil
	}
	return common.SlicePrefix(goErrorParse(out), "golint"), false, nil
}

// runGoBuild runs "go build" on the source file and returns the output.
func runGoBuild(dirname, fname string) ([]string, bool) {
	out, err := Execute("go", "build", "-o", dirname, fname)
	retcode := Exitcode(err)

	// No errors.
	if retcode == 0 {
		return []string{}, true
	}
	return common.SlicePrefix(goErrorParse(out), "go build"), false
}

// runGoFmt runs "go fmt -d" on the source file and indicates if any output
// exists (this means the user needs to run gofmt on their source).
func runGoFmt(fname string) ([]string, bool) {
	out, err := Execute("gofmt", "-d", fname)
	retcode := Exitcode(err)

	// No errors.
	if retcode != 0 || len(out) != 0 {
		ret := []string{"Gofmt detected differences. Please run \"gofmt\" to fix this"}
		return common.SlicePrefix(ret, "gofmt"), len(out) == 0
	}
	return []string{""}, true
}

// goErrorParse remove undesirable lines and formats the output from go build.
func goErrorParse(list []string) []string {
	var ret []string
	for _, v := range list {
		// Go builds adds lines starting with #
		if strings.HasPrefix(v, "#") {
			continue
		}
		// Go build and go lint prefix lines with filename:line:column. Remove
		// the filename since it's a temp file anyway.
		r := goLineRegex.FindStringSubmatch(v)

		// Unable to parse line. Include literally.
		if r == nil || len(r) < 5 {
			ret = append(ret, v)
			continue
		}

		ret = append(ret, fmt.Sprintf("Line %s Col %s: %s", r[2], r[3], r[4]))
	}
	return ret
}
