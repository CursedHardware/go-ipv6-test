package ipv6test

import (
	"bytes"
	"encoding/json"
	"io"
)

func unmarshalJSONP(r io.Reader, v any) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	start := bytes.IndexByte(body, '(')
	end := bytes.LastIndexByte(body, ')')
	if start != -1 && end != -1 {
		body = body[start+1 : end]
	}
	return json.Unmarshal(body, v)
}
