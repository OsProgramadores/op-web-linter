// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

// httpError logs the error and returns the appropriate message & http code.
func httpError(w http.ResponseWriter, msg string, httpcode int) {
	m := fmt.Sprintf("Returned HTTP error %d: %v", httpcode, msg)
	log.Print(m)
	http.Error(w, msg, httpcode)
}

// exitcode fetches the numeric return code from the return of exec.Run.
// There's no portable way of retrieving the exit code. This function returns
// 255 if there is an error in the code and we are in a platform that does not
// have syscall.WaitStatus.
func exitcode(err error) int {
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

// slicePrefix adds a prefix to every string line in the slice.
func slicePrefix(slice []string, prefix string) []string {
	var ret []string
	for _, line := range slice {
		ret = append(ret, fmt.Sprintf("[%s] %s", prefix, line))
	}
	return ret
}
