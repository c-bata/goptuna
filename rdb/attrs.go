package rdb

import (
	"encoding/base64"
	"fmt"
)

// See https://github.com/c-bata/goptuna/issues/34
// for the reason why we need following code.

// Caution "_number" in trial_system_attributes must not be encoded.

func encodeAttrValue(xr string) string {
	return fmt.Sprintf("\"%s\"",
		base64.StdEncoding.EncodeToString([]byte(xr)))
}

func decodeAttrValue(j string) (string, error) {
	l := len(j)
	if l < 2 {
		return j, nil
	}
	encoded, err := base64.StdEncoding.DecodeString(j[1 : l-1])
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
