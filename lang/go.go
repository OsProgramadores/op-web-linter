// Package lang defines all language speicfic components of op-web-linter.
//
// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.
package lang

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	original, err := url.QueryUnescape(req.Text)
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Decoded program: %s\n", original)

	// Save program text in request to file.
	tempdir, tempfile, err := saveProgramToFile(original, "*.go")
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Attempt to reformat source with gofmt (+simplify).
	// Indicate formatting failure if necessary.
	reformatted, err := Execute("gofmt", "-s", tempfile)

	if err != nil {
		messages = append(messages, fmt.Sprintf("Reformat failed: %v", err))
	} else {
		// Rewrite reformatted program to tempfile.
		if err := os.WriteFile(tempfile, []byte(reformatted), 0644); err != nil {
			common.HttpError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

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

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:            len(messages) == 0,
		ErrorMessages:   messages,
		Reformatted:     reformatted != original,
		ReformattedText: reformatted,
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
	o, err := Execute("golint", fname)
	out := strings.Split(o, "\n")

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
	o, err := Execute("go", "build", "-o", dirname, fname)
	out := strings.Split(o, "\n")
	retcode := Exitcode(err)

	// No errors.
	if retcode == 0 {
		return []string{}, true
	}
	return common.SlicePrefix(goErrorParse(out), "go build"), false
}

// goErrorParse remove undesirable lines and formats the output from go build.
func goErrorParse(list []string) []string {
	var ret []string
	for _, v := range list {
		// Go builds adds lines starting with #
		if strings.HasPrefix(v, "#") {
			continue
		}
		// Remove blank lines.
		if strings.TrimSpace(v) == "" {
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
