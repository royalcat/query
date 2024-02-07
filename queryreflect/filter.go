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
		res, err := cond(v)
		if err != nil {
			return nil, err
		}
		if res {
			out = append(out, v)
		}
	}
	return out, nil
}

type conditionErr[D any] func(v D) (bool, error)

func generateReflectFilter[D any](f query.Filter) (conditionErr[D], error) {
	conditions := []conditionErr[D]{}

	for _, filter := range f {
		f := func(data D) (bool, error) {
			vs1, err := getValueByPath(reflect.ValueOf(data), filter.Field)
			if err != nil {
				return false, err
			}
			for _, v1 := range vs1 {
				if reflectCompare(filter.Op, v1, reflect.ValueOf(filter.Value)) {
					return true, nil
				}
			}
			return false, nil
		}

		conditions = append(conditions, f)
	}

	return func(v D) (bool, error) {
		for _, c := range conditions {
			res, err := c(v)
			if err != nil {
				return false, err
			}
			if !res {
				return false, nil
			}
		}
		return true, nil
	}, nil
}
