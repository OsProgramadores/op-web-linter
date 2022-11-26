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

	"github.com/osprogramadores/op-web-linter/common"
	"github.com/osprogramadores/op-web-linter/handlers"
)

// LintC lints programs written in C. For now, only reformats code with indent.
func LintC(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	original, err := url.QueryUnescape(req.Text)
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Decoded program: %s\n", original)

	// Save program text in request to file.
	tempdir, tempfile, err := saveProgramToFile(original, "*.c")
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Reformat source code with indent.
	reformatted, err := Execute("indent", "-st", "--no-tabs", "--tab-size4", "--indent-level4",
		"--braces-on-if-line", "--cuddle-else", "--braces-on-func-def-line", "--braces-on-struct-decl-line",
		"--cuddle-do-while", "--no-space-after-function-call-names", "--no-space-after-parentheses",
		"--dont-break-procedure-type", tempfile)
	if err != nil {
		messages = append(messages, fmt.Sprintf("Reformat failed: %v", err))
	}

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:            err == nil,
		ErrorMessages:   messages,
		Reformatted:     reformatted != original && err == nil,
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
