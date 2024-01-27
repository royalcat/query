package queryreflect

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/royalcat/query"
)

func ApplyQuery[D any](q query.Query, in []D) ([]D, error) {
	var err error
	if len(q.Filter) > 0 {
		in, err = ApplyFilter(q.Filter, in)
		if err != nil {
			return nil, err
		}
	}
	if len(q.Sort) > 0 {
		in, err = ApplySort(q.Sort, in)
		if err != nil {
			return nil, err
		}
	}

	if len(in) <= int(q.Pagination.Offset) {
		return []D{}, nil
	} else if len(in) < int(q.Pagination.Offset+q.Pagination.Limit) || q.Pagination.Limit == 0 {
		return in[q.Pagination.Offset:], nil
	} else {
		return in[q.Pagination.Offset : q.Pagination.Offset+q.Pagination.Limit], nil
	}
}

func ApplySort[D any](s query.Sort, in []D) ([]D, error) {
	cmps := generateReflectSort[D](s)

	for _, cmp := range cmps {
		slices.SortFunc(in, cmp)
	}

	return in, nil
}

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

type pageGetter[D any] func(q query.Query) ([]D, error)

func ApplyQueryWithNext[D any](q query.Query, getPage pageGetter[D]) (out []D, err error) {
	pageQuery := q.Copy()
	for {
		page, err := getPage(pageQuery)
		if err != nil {
			return nil, err
		}

		page, err = ApplyFilter(q.Filter, page)
		if err != nil {
			return nil, err
		}

		out = append(out, page...)

		if len(page) != int(q.Pagination.Limit) {
			break
		}

		if len(out) >= int(q.Pagination.Offset+q.Pagination.Limit) {
			break
		}
	}

	out, err = ApplyQuery(q, out)
	if err != nil {
		return nil, err
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

func reflectCompare(o query.Operator, v1, v2 reflect.Value) bool {
	t := v1.Type()
	switch t.Kind() {
	case reflect.Bool:
		if v1.Type().Kind() != v2.Type().Kind() {
			return false
		}
		switch o {
		case query.OperatorEqual:
			return v1.Bool() == v2.Bool()
		case query.OperatorNotEqual:
			return v1.Bool() != v2.Bool()
		}
	case reflect.Float32, reflect.Float64:
		if v1.Type().Kind() != v2.Type().Kind() {
			return false
		}
		switch o {
		case query.OperatorEqual:
			return v1.Float() == v2.Float()
		case query.OperatorNotEqual:
			return v1.Float() != v2.Float()
		case query.OperatorGreater:
			return v1.Float() > v2.Float()
		case query.OperatorLess:
			return v1.Float() < v2.Float()
		case query.OperatorGreaterOrEqual:
			return v1.Float() >= v2.Float()
		case query.OperatorLessOrEqual:
			return v1.Float() <= v2.Float()
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if v1.Type().Kind() != v2.Type().Kind() {
			return false
		}
		switch o {
		case query.OperatorEqual:
			return v1.Int() == v2.Int()
		case query.OperatorNotEqual:
			return v1.Int() != v2.Int()
		case query.OperatorGreater:
			return v1.Int() > v2.Int()
		case query.OperatorLess:
			return v1.Int() < v2.Int()
		case query.OperatorGreaterOrEqual:
			return v1.Int() >= v2.Int()
		case query.OperatorLessOrEqual:
			return v1.Int() <= v2.Int()
		case query.OperatorSubString:
			return strings.Contains(strconv.Itoa(int(v1.Int())), strconv.Itoa(int(v1.Int())))
		}
	case reflect.String:
		if v1.Type().Kind() != v2.Type().Kind() {
			return false
		}
		switch o {
		case query.OperatorEqual:
			return v1.String() == v2.String()
		case query.OperatorNotEqual:
			return v1.String() != v2.String()
		case query.OperatorGreater:
			return v1.String() > v2.String()
		case query.OperatorLess:
			return v1.String() < v2.String()
		case query.OperatorGreaterOrEqual:
			return v1.String() >= v2.String()
		case query.OperatorLessOrEqual:
			return v1.String() <= v2.String()
		case query.OperatorSubString:
			return strings.Contains(v1.String(), v2.String())
		}

	case reflect.Struct:
		if v1.Type().Kind() != v2.Type().Kind() {
			return false
		}

		switch t {
		case reflect.TypeOf(time.Time{}):
			return reflectCompare(o,
				reflect.ValueOf(v1.Interface().(time.Time).Unix()),
				reflect.ValueOf(v2.Interface().(time.Time).Unix()),
			)
			// case reflect.TypeOf(types.Date{}):
			// 	return reflectCompare(o,
			// 		reflect.ValueOf(v1.Interface().(types.Date).Unix()),
			// 		reflect.ValueOf(v2.Interface().(types.Date).Unix()),
			// 	)
			// case reflect.TypeOf(uuid.UUID{}):
			// 	return o.reflectCompare(
			// 		reflect.ValueOf(v1.Interface().(uuid.UUID).String()),
			// 		reflect.ValueOf(v2.Interface().(uuid.UUID).String()),
			// 	)
			// }
		}
	case reflect.Array, reflect.Slice:
		if v1.Type().Elem().Kind() != v2.Type().Kind() {
			return false
		}
		for i := 0; i < v1.Len(); i++ {
			if reflectCompare(o, v1.Index(i), v2) {
				return true
			}
		}
	}

	return false
}

func getValueByPath(modelValue reflect.Value, path string) ([]reflect.Value, error) {
	parts := strings.Split(path, ".")

	t := modelValue

	for i := 0; i < len(parts); {
		switch t.Kind() {
		case reflect.Struct:
			f, found := getFieldValueByJsonTag(t, parts[i])
			if !found {
				return nil, fmt.Errorf("invalid path part: %s", parts[i])
			}
			t = f
			i++
		case reflect.Slice, reflect.Array:
			if idx, err := strconv.Atoi(parts[i]); err == nil {
				if idx >= t.Len() {
					return []reflect.Value{}, nil
				}
				t = t.Index(idx)
				i++
			} else {
				out := []reflect.Value{}
				for idx := 0; idx < t.Len(); idx++ {
					vals, err := getValueByPath(t.Index(idx), strings.Join(parts[1:], "."))
					if err != nil {
						return nil, err
					}
					out = append(out, vals...)
				}
				return out, nil
			}

		case reflect.Pointer:
			t = t.Elem()
		default:
			return nil, fmt.Errorf("invalid path part: %s", parts[i])
		}

	}

	return []reflect.Value{t}, nil
}

func getFieldValueByJsonTag(val reflect.Value, name string) (reflect.Value, bool) {
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if key == name {
			return val.Field(i), true
		}
	}

	return reflect.Value{}, false
}
