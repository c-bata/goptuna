package rdb

import (
	"encoding/json"
)

// See https://github.com/c-bata/goptuna/issues/34
// for the reason why we need following code.

type GoptunaAttr struct {
	Library string `json:"library"`
	Param   string `json:"param"`
}

func encodeAttrValue(xr string) string {
	jsonBytes, _ := json.Marshal(&GoptunaAttr{
		Library: "Goptuna",
		Param:   xr,
	})
	return string(jsonBytes)
}

func decodeAttrValue(j string) string {
	var attr GoptunaAttr
	err := json.Unmarshal([]byte(j), &attr)
	// Return an empty string if couldn't parse attr.
	if err != nil || attr.Library != "Goptuna" {
		return ""
	}
	return attr.Param
}
