package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func genericType[Model any]() reflect.Type {
	return reflect.TypeOf((*Model)(nil)).Elem()
}

type Unmarshaler interface {
	QueryUnmarshal(v string) (any, error)
}

func ParseStringFilter[Model any](values map[string]string) (Filter, error) {
	t := genericType[Model]()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	f := Filter{}

	for k, v := range values {
		name, op, err := parseMapKey(k)
		if err != nil {
			return nil, err
		}

		t, err := GetTypeByPath(t, name)
		if err != nil {
			return f, err
		}

		if op == OperatorIn {
			vals := strings.Split(v, ",")
			filterValue := reflect.MakeSlice(reflect.SliceOf(t), 0, len(vals))

			for _, v := range vals {
				val, err := parseStringForType(t, v)
				if err != nil {
					return f, fmt.Errorf("cant get value for type %s, error: %s", t.Kind().String(), err.Error())
				}
				filterValue = reflect.Append(filterValue, reflect.ValueOf(val).Convert(t))
			}
			f = append(f, FieldFilter{
				Field: name,
				Op:    op,
				Value: filterValue.Interface(),
			})
			continue
		}

		val, err := parseStringForType(t, v)
		if err != nil {
			return f, fmt.Errorf("cant get value for type %s, error: %s", t.Kind().String(), err.Error())
		}
		f = append(f, FieldFilter{
			Field: name,
			Op:    op,
			Value: val,
		})

	}

	return f, nil
}

func parseMapKey(key string) (name string, operator Operator, err error) {
	opStart := strings.Index(key, "{")

	if opStart == -1 {
		return key, "", nil
	}
	name = key[:opStart]
	operator = Operator(key[opStart+1 : len(key)-1])
	if !isOperator(operator) {
		return "", "", fmt.Errorf("unknow operator: %s", operator)
	}
	return name, operator, nil
}

var pathTypesCache syncmap[reflect.Type, *syncmap[string, reflect.Type]]

func GetTypeByPath(t reflect.Type, path string) (reflect.Type, error) {
	if typeCache, ok := pathTypesCache.Load(t); ok && typeCache != nil {
		if cached, ok := typeCache.Load(path); ok {
			return cached, nil
		}
	}

	parts := strings.Split(path, ".")

	for i := 0; i < len(parts); {
		switch t.Kind() {
		case reflect.Struct:
			f, found := getFieldByJsonTag(t, parts[i])
			if !found {
				return nil, fmt.Errorf("invalid path part: %s", parts[i])
			}
			t = f.Type
			i++
		case reflect.Slice, reflect.Array:
			t = t.Elem()
			if _, err := strconv.Atoi(parts[i]); err == nil {
				i++
			}
		case reflect.Pointer:
			t = t.Elem()
		default:
			return nil, fmt.Errorf("invalid path part: %s", parts[i])
		}
	}

	typeCache, _ := pathTypesCache.LoadOrStore(t, &syncmap[string, reflect.Type]{}) //nolint:exhaustruct
	typeCache.Store(path, t)

	return t, nil
}

func getFieldByJsonTag(typ reflect.Type, name string) (field reflect.StructField, ok bool) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		key, _ := strings.CutPrefix(field.Tag.Get("json"), ",")
		if key == name {
			return field, true
		}
	}

	return field, false
}

// func GetValueForType(t reflect.Type, v any) (any, error) {
// 	vt := reflect.ValueOf(v)

// 	switch vt.Kind() {
// 	case reflect.Interface, reflect.Pointer:
// 		return GetValueForType(t, vt.Elem().Interface())
// 	case reflect.String:
// 		return parseStringForType(t, v.(string))
// 	}
// 	if t.Kind() == vt.Kind() {
// 		return v, nil
// 	}

// 	return nil, fmt.Errorf("unknow kind of type: %v", vt.String())
// }

func parseStringForType(t reflect.Type, v string) (any, error) {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return parseStringForType(t.Elem(), v)
	case reflect.String:
		return v, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.Atoi(v)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(v, 64)
	case reflect.Bool:
		switch v {
		case "true", "True":
			return true, nil
		case "false", "False":
			return false, nil
		}
		return nil, fmt.Errorf("unknow bool value: %s", v)
	case reflect.Struct:
		if val, ok := reflect.New(t).Interface().(Unmarshaler); ok {
			val, err := val.QueryUnmarshal(v)
			if err != nil {
				return nil, fmt.Errorf("QueryUnmarshal error: %s", err.Error())
			}
			return val, err
		}

		switch t {
		case reflect.TypeOf(time.Time{}):
			ts, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return nil, fmt.Errorf("cant parse as timestamp: %s", err.Error())
			}
			return ts, nil
		}
		val, ok := reflect.New(t).Interface().(json.Unmarshaler)
		if ok {
			err := val.UnmarshalJSON([]byte("\"" + v + "\""))
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal error: %s", err.Error())
			}
			return val, nil
		}
	}
	return nil, fmt.Errorf("unsupported type: %s", t.String())
}
