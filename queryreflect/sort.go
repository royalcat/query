package queryreflect

import (
	"reflect"
	"slices"

	"github.com/royalcat/query"
)

func ApplySort[D any](s query.Sort, in []D) ([]D, error) {
	cmps := generateReflectSort[D](s)

	for _, cmp := range cmps {
		slices.SortFunc(in, cmp)
	}

	return in, nil
}

type compare[D any] func(v1, v2 D) int

func generateReflectSort[D any](s query.Sort) []compare[D] {
	out := make([]compare[D], 0, len(s))
	for _, f := range s {

		if f.Order == query.ASC {
			out = append(out, func(v1, v2 D) int {
				vs1, _ := getValueByPath(reflect.ValueOf(v1), f.Key)
				vs2, _ := getValueByPath(reflect.ValueOf(v2), f.Key)
				if reflectCompare(query.OperatorGreater, vs1[0], vs2[0]) {
					return 1
				} else if reflectCompare(query.OperatorLess, vs1[0], vs2[0]) {
					return -1
				} else {
					return 0
				}
			})
		} else {
			out = append(out, func(v1, v2 D) int {
				vs1, _ := getValueByPath(reflect.ValueOf(v1), f.Key)
				vs2, _ := getValueByPath(reflect.ValueOf(v2), f.Key)
				if reflectCompare(query.OperatorGreater, vs1[0], vs2[0]) {
					return -1
				} else if reflectCompare(query.OperatorLess, vs1[0], vs2[0]) {
					return 1
				} else {
					return 0
				}
			})
		}
	}

	return out
}
