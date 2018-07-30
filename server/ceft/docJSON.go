package ceft

import(
	"encoding/json"
)

// DocJSON is Single document struct
type DocJSON struct {
	Title string
	Text string
}

// DecodeDocJSON decode JSON text
func DecodeDocJSON(text string) DocJSON {
	var docj DocJSON
	json.Unmarshal([]byte(text), &docj)
	return docj
}
