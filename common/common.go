// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package common

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	execTimeout = time.Duration(5 * time.Second)
)

// Execute runs the program specified by name with the command-line specified
// in slice args. Returns the error code and a string slice containing all
// non-blank lines in the program's combined output.
func Execute(name string, args ...string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	out, err := cmd.CombinedOutput()

	var ret []string
	for _, line := range strings.Split(string(out), "\n") {
		if line != "" {
			ret = append(ret, line)
		}
	}
	return ret, err
}

// Exitcode fetches the numeric return code from the return of exec.Run.
// There's no portable way of retrieving the exit code. This function returns
// 255 if there is an error in the code and we are in a platform that does not
// have syscall.WaitStatus.
func Exitcode(err error) int {
	if err == nil {
		return 0
	}
	retcode := 255
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			retcode = status.ExitStatus()
		}
	}
	return retcode
}

// HttpError logs the error and returns the appropriate message & http code.
func HttpError(w http.ResponseWriter, msg string, httpcode int) {
	m := fmt.Sprintf("Returned HTTP error %d: %v", httpcode, msg)
	log.Print(m)
	http.Error(w, msg, httpcode)
}

// SlicePrefix adds a prefix to every string line in the slice.
func SlicePrefix(slice []string, prefix string) []string {
	var ret []string
	for _, line := range slice {
		ret = append(ret, fmt.Sprintf("[%s] %s", prefix, line))
	}
	return ret
}
