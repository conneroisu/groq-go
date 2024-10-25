package test

import (
	"encoding/json"
	"io"
)

func encode(v any) []byte {
	res, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return res
}
func decode[T any](r io.Reader) T {
	// make a new instance of the type
	v := new(T)
	// read the JSON data from the request body
	bod, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bod, &v); err != nil {
		panic(err)
	}
	return *v
}
