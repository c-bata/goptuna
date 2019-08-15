package rdb

import "testing"

func TestEncodeAttr(t *testing.T) {
	var tests = []struct {
		str string
	}{
		{str: "foo bar"},
		{str: "\"foo bar"},
		{str: "f\"oo\" bar"},
		{str: "fo\"o\" b\"ar"},
		{str: "\"\" b\"aro\""},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			x := encodeAttrValue(tt.str)
			y, err := decodeAttrValue(x)
			if err != nil {
				t.Errorf("error should be nil, but got %s", err)
				return
			}
			if tt.str != y {
				t.Errorf("%s != %s", tt.str, y)
			}
		})
	}
}
