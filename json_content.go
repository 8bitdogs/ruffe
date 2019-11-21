package ruffe

import (
	"encoding/json"
	"io"
)

type jsonContent struct{}

func (c jsonContent) Unmarshal(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (c jsonContent) ContentType() string {
	return "application/json"
}

func (c jsonContent) Marshal(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}
