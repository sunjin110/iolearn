package jsonutil

import (
	"encoding/json"
	"iolearn/pkg/common/chk"
)

// Marshal .
func Marshal(v interface{}) string {
	str, err := json.Marshal(v)
	chk.SE(err)
	return string(str)
}
