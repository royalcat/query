package queryreflect

import (
	"reflect"

	"github.com/royalcat/query"
)

func ApplyFilter[D any](f query.Filter, in []D) ([]D, error) {
	cond, err := generateReflectFilter[D](f)
	if err != nil {
		return nil, err
	}
	out := []D{}
	for _, v := range in {
		if cond(v) {
			out = append(out, v)
		}
	}
	return out, nil
}

type condition[D any] func(v D) bool

func generateReflectFilter[D any](f query.Filter) (condition[D], error) {
	conditions := []condition[D]{}
	var m D
	t := reflect.TypeOf(m)

	for key, value := range f {
		name, operator := query.ParseKey(key)
		t2, err := query.GetTypeByPath(t, name)
		if err != nil {
			return nil, err
		}
		val2, err := query.GetValueForType(t2, value)
		if err != nil {
			return nil, err
		}

		f := func(data D) bool {
			vs1, _ := getValueByPath(reflect.ValueOf(data), name)
			for _, v1 := range vs1 {
				if reflectCompare(operator, v1, reflect.ValueOf(val2)) {
					return true
				}
			}
			return false
		}

		conditions = append(conditions, f)
	}

	return func(v D) bool {
		for _, c := range conditions {
			if !c(v) {
				return false
			}
		}
		return true
	}, nil
}
