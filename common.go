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
func httpError(w http.ResponseWriter, err error, httpcode int) {
	m := fmt.Sprintf("Returned HTTP error %d: %v", httpcode, err)
	log.Print(m)
	http.Error(w, m, httpcode)
	return
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
