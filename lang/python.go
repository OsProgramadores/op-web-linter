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
// Sample line: /tmp/smartcd.py:187:0: W0311: Bad indentation. Found 4 spaces, expected 8 (bad-indentation)
var pylintLineRegex = regexp.MustCompile("^[^:]+:([0-9]+):([0-9]+):[ ]*(.*)")

// LintPython lints programs written in Python (v3).
func LintPython(w http.ResponseWriter, r *http.Request, req handlers.LintRequest) {
	tempdir, tempfile, err := saveRequestToFile(req.Text, "*.py")
	if err != nil {
		common.HTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	// pylint.
	homedir := os.Getenv("HOME")
	out, err := Execute("pylint", "--rcfile="+homedir+"/op-web-linter/config/pylint3.rc", tempfile)

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:          err == nil,
		ErrorMessages: common.SlicePrefix(PythonFilterOutput(out, tempfile), "pylint"),
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

// PythonFilterOutput remove undesirable messages from the pylint output.
// pylint3 is very verbose. Limit output to the lines starting with our
// filename.
func PythonFilterOutput(output string, tempfile string) []string {
	var ret []string
	for _, v := range strings.Split(output, "\n") {
		if !strings.HasPrefix(v, tempfile) {
			continue
		}
		// Remove blank lines.
		if strings.TrimSpace(v) == "" {
			continue
		}
		// Parse line:column message error lines.
		r := pylintLineRegex.FindStringSubmatch(v)

		// Unable to parse line, Include literally (this should not happen).
		if len(r) < 4 {
			ret = append(ret, v)
			continue
		}
		ret = append(ret, fmt.Sprintf("Line %s Col %s: %s", r[1], r[2], r[3]))
	}
	return ret
}
