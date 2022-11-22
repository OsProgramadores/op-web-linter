// This file is part of op-web-linter.
// See github.com/osprogramadores/op-web-linter for licensing and details.

package common

import (
	"fmt"
	"log"
	"net/http"
)

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
