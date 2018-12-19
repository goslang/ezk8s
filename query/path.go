package query

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/oliveagle/jsonpath"
)

type Path struct {
	JsonPath string
	Target   interface{}
}

// Apply will iterate the JSON and look for the matching data. It will use
// Reflect to write the data to it's expected type for the user.
func (p *Path) Apply(data map[string]interface{}) (err error) {
	defer trapError(&err)

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

	// Slices must be handled specially.
	if targetV.Kind() == reflect.Slice {
		setAllInSlice(targetV, resV)
		return nil
	}

	targetV.Set(resV)
	return nil
}

// setAllInSlice iterates through sourceSlice and sets corresponding values in
// target. It panics if an error is encountered.
func setAllInSlice(target, sourceSlice reflect.Value) {
	// First, convert the sourceSlice to a real slice of interfaces. This
	// should work because the unknown values in JSON will always be parsed
	//as interfaces.
	values := sourceSlice.Interface().([]interface{})

	// Second, allocate enough memory in the target to fit all of of the
	// values.
	buf := reflect.MakeSlice(target.Type(), len(values), len(values))
	target.Set(buf)

	// Finally, iterate the values and set them at their correpsonding
	// locations in the target
	for i, v := range values {
		t := target.Index(i)
		t.Set(reflect.ValueOf(v))
	}
}

func trapError(targetErr *error) {
	err := recover()
	if err == nil {
		return
	}

	var ok bool
	*targetErr, ok = err.(error)
	if !ok {
		*targetErr = fmt.Errorf("%v", err)
	}
	return
}
