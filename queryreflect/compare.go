package queryreflect

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/royalcat/query"
)

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
		case query.OperatorEqual, query.OperatorDefault:
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
		case query.OperatorEqual, query.OperatorDefault:
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
		case query.OperatorEqual, query.OperatorDefault:
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
		key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if key == name {
			return val.Field(i), true
		}
	}

	return reflect.Value{}, false
}
