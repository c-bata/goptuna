package rdb

import "encoding/json"

// See https://github.com/c-bata/goptuna/issues/34
// for the reason why we need following code.

type attrJSONRepresentation struct {
	value string `json:"value"`
}

func encodeToAttrJSON(j string) (string, error) {
	var r attrJSONRepresentation
	err := json.Unmarshal([]byte(j), &r)
	if err != nil {
		return "", err
	}
	return r.value, nil
}

func decodeAttrFromJSON(xr string) (string, error) {
	j, err := json.Marshal(&attrJSONRepresentation{
		value: xr,
	})
	if err != nil {
		return "", err
	}
	return string(j), nil
}
