package query

import (
	"errors"
	"reflect"

	"github.com/oliveagle/jsonpath"
)

type Path struct {
	JsonPath string
	Target   interface{}
}

func (p *Path) Apply(data map[string]interface{}) error {
	pat, err := jsonpath.Compile(p.JsonPath)
	if err != nil {
		return err
	}

	res, err := pat.Lookup(data)
	if err != nil {
		return err
	}

	resV := reflect.ValueOf(res)
	targetV := reflect.ValueOf(p.Target).Elem()

	if !targetV.CanAddr() || !targetV.CanSet() {
		return errors.New("Cannot set value for path " + p.JsonPath)
	}

	if targetV.Kind() != resV.Kind() {
		return errors.New("Cannot set value because its kind is incorrect")
	}

	targetV.Set(resV)
	return nil
}
