package query

import (
	"encoding/json"
	"io"
)

type Result interface {
	Error() error
	Decode(target interface{}) error
	Scan(paths ...Path) error
}

type errorResult func() error

func NewErrorResult(err error) errorResult {
	return func() error { return err }
}

func (er errorResult) Decode(_ interface{}) error {
	return er()
}

func (er errorResult) Error() error {
	return er()
}

func (er errorResult) Scan(_ ...Path) error {
	return er()
}

type decodeResult func(target interface{}) error

func NewDecodeResult(reader io.ReadCloser) decodeResult {
	return func(target interface{}) error {
		defer reader.Close()

		if target == nil {
			return nil
		}
		return json.NewDecoder(reader).Decode(target)
	}
}

func (dr decodeResult) Decode(target interface{}) error {
	return dr(target)
}

func (dr decodeResult) Error() error {
	return dr(nil)
}

func (dr decodeResult) Scan(paths ...Path) error {
	data := make(map[string]interface{})
	if err := dr.Decode(&data); err != nil {
		return err
	}

	for _, path := range paths {
		err := path.Apply(data)
		if err != nil {
			return err
		}
	}
	return nil
}
