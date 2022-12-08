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
	"strings"

	"github.com/osprogramadores/op-web-linter/common"
	"github.com/osprogramadores/op-web-linter/handlers"
)

// LintC lints programs written in C using clang-format and clang-tidy.
func LintC(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	var clangChecks = []string{
		"readability*",
		"clang-analyzer-*",
		"concurrency-*",
		"cppcoreguidelines-*",
		"google-*",
		"-readability-identifier-length",
		"-readability-magic-numbers",
		"-cppcoreguidelines-avoid-magic-numbers",
	}

	tempdir, tempfile, err := saveRequestToFile(req.Text, "*.c")
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Reformat source code using clang-format. In case of errors, we move ahead
	// with the old code and attempt linting anyway.
	reformatted, err := Execute("clang-format", "--assume-filename=c",
		"--style={BasedOnStyle: google, IndentWidth: 4}", tempfile)
	if err != nil {
		messages = append(messages, fmt.Sprintf("Error reformatting C code: %v", err))
		messages = append(messages, strings.Split(reformatted, "\n")...)
	} else {
		// Rewrite reformatted program to tempfile.
		if err := os.WriteFile(tempfile, []byte(reformatted), 0644); err != nil {
			common.HTTPError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	reformatErr := err

	// clang-tidy returns an error code (1) on errors, but nothing on warnings.
	// We want to indicate every situation, so we ignore it here and look for
	// the output. Blank output means no errors.
	out, err := Execute("clang-tidy", "--checks="+strings.Join(clangChecks, ","), tempfile, "--")
	// Use cppFilterOutput since it's basically a clang-tidy output beautifier.
	lines := cppFilterOutput(strings.Split(out, "\n"), tempfile)
	messages = append(messages, common.SlicePrefix(lines, "clang-tidy")...)

	// Pass if no messages from the reformatter or linter.
	pass := len(messages) == 0

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:            pass,
		ErrorMessages:   messages,
		Reformatted:     reformatted != req.Text && reformatErr == nil,
		ReformattedText: reformatted,
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("JSON response: %v", string(jresp))
	w.Write(jresp)
	w.Write([]byte("\n"))
}
