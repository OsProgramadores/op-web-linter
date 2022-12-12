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

// Regexp matching eslint lines.
var eslintLineRegex = regexp.MustCompile("^[ \t]*([0-9]+):([0-9]+)[ ]*(.*)")

// LintJavascript lints programs written in Javascript.
func LintJavascript(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	tempdir, tempfile, err := saveRequestToFile(req.Text, "*.js")
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	// eslint.
	o, err := Execute("npx", "eslint", "-c", "/tmp/build/src/op-web-linter/config/eslintrc.json", tempfile)
	out := strings.Split(o, "\n")

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:          err == nil,
		ErrorMessages: common.SlicePrefix(JavascriptFilterOutput(out, tempfile), "eslint"),
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

// JavascriptFilterOutput remove undesirable messages from the eslint output.
func JavascriptFilterOutput(list []string, tempfile string) []string {
	var ret []string
	for _, v := range list {
		// eslint adds a line with the filename.
		if strings.HasPrefix(v, tempfile) {
			continue
		}
		// Remove blank lines.
		if strings.TrimSpace(v) == "" {
			continue
		}
		// Parse line:column message error lines.
		r := eslintLineRegex.FindStringSubmatch(v)

		// Unable to parse line. Include literally.
		if r == nil || len(r) < 4 {
			ret = append(ret, v)
			continue
		}
		ret = append(ret, fmt.Sprintf("Line %s Col %s: %s", r[1], r[2], r[3]))
	}
	return ret
}
