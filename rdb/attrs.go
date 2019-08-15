package rdb

import (
	"fmt"
	"strings"
)

// See https://github.com/c-bata/goptuna/issues/34
// for the reason why we need following code.

func encodeToOptunaInternalAttr(xr string) string {
	return fmt.Sprintf("\"%s\"", strings.Replace(xr, "\"", "\\\"", -1))
}

func decodeFromOptunaInternalAttr(j string) string {
	l := len(j)
	if l < 2 {
		return j
	}
	return strings.Replace(j[1:l-1], "\\\"", "\"", -1)
}
