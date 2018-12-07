package query

type Result struct {
	mapping map[string]interface{}

	Data map[string]interface{}
}

func NewResult() *Result {
	return &Result{
		mapping: make(map[string]interface{}),
		Data:    make(map[string]interface{}),
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
