package query

import (
	"encoding/json"
)

type Result struct {
	Data map[string]interface{}

	decoder json.Decoder
}

func NewResult() *Result {
	return &Result{
		Data: make(map[string]interface{}),
	}
}

func (r *Result) Scan(paths ...Path) error {
	for _, path := range paths {
		err := path.Apply(r.Data)
		if err != nil {
			return err
		}
	}
	return nil
}
