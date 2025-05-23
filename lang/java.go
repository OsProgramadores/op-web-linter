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

	"github.com/osprogramadores/op-web-linter/common"
	"github.com/osprogramadores/op-web-linter/handlers"
)

// LintJava lints programs written in Java. For now, only reformats code with google-java-format.
func LintJava(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	// Save program text in request to file.
	tempdir, tempfile, err := saveRequestToFile(req.Text, "*.java")
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// Reformat source code with google-java-format.
	reformatted, err := Execute("/usr/lib/jvm/java-17-openjdk/bin/java", "-jar", "/home/op/google-java-format-1.24.0-all-deps.jar", tempfile)
	if err != nil {
		messages = append(messages, fmt.Sprintf("Reformat failed: %v", err))
	}

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:            err == nil,
		ErrorMessages:   messages,
		Reformatted:     reformatted != req.Text && err == nil,
		ReformattedText: reformatted,
	}
	jresp, err := json.Marshal(resp)
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("JSON response:\n%s", prettyJSONString(jresp))
	w.Write(jresp)
	w.Write([]byte("\n"))
}
