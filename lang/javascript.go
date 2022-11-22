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
	tempdir, tempfile, err := handlers.SaveProgramToFile(req, "*.js")
	if err != nil {
		common.HttpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempdir)

	var messages []string

	// eslint.
	messages, ok, err := runEslint(tempfile)

	// Create response, convert to JSON and return.
	resp := handlers.LintResponse{
		Pass:          ok,
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

// runEslint runs eslint on the source file and returns the output.
func runEslint(fname string) ([]string, bool, error) {
	cmd := `npx eslint -c /tmp/build/src/op-web-linter/config/eslintrc.yaml ` + fname
	c := strings.Split(cmd, " ")
	out, err := Execute(c[0], c[1:]...)
	retcode := Exitcode(err)

	// No errors.
	if retcode == 0 {
		return []string{}, true, nil
	}
	return common.SlicePrefix(JavascriptErrorParse(out, fname), "eslint"), false, err
}

// JavascriptErrorParse remove undesirable messages from the eslint output.
func JavascriptErrorParse(list []string, tempfile string) []string {
	var ret []string
	for _, v := range list {
		// eslint adds a line with the filename.
		if strings.HasPrefix(v, tempfile) {
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
