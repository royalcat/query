package query

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Filter map[string]string

type QueryUnmarshaler interface {
	QueryUnmarshal(v string) (any, error)
}

func (q Filter) Fields() Fields {
	fields := []string{}
	for k := range q {
		name, _ := ParseKey(k)
		fields = append(fields, name)
	}
	return fields
}

func (q Filter) Operators() map[string]Operator {
	ops := map[string]Operator{}
	for k := range q {
		name, operator := ParseKey(k)
		ops[name] = operator
	}
	return ops
}

func ParseKey(key string) (name string, operator Operator) {
	parts := strings.SplitN(key, "{", 2)
	name = parts[0]
	if len(parts) == 2 {
		operator = Operator(strings.TrimSuffix(parts[1], "}"))
	} else {
		operator = OperatorEqual
	}

	return name, operator
}

func GetValueForType(t reflect.Type, s string) (any, error) {
	switch t.Kind() {
	case reflect.Slice, reflect.Array, reflect.Pointer:
		return GetValueForType(t.Elem(), s)
	case reflect.String:
		return s, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.Atoi(s)
	case reflect.Float32, reflect.Float64:
		return strconv.ParseFloat(s, 64)
	case reflect.Bool:
		switch s {
		case "true", "True":
			return true, nil
		case "false", "False":
			return false, nil
		}
		return nil, fmt.Errorf("unknow bool value: %s", s)
	case reflect.Struct:
		if val, ok := reflect.New(t).Interface().(QueryUnmarshaler); ok {
			val, err := val.QueryUnmarshal(s)
			if err != nil {
				return nil, fmt.Errorf("QueryUnmarshal error: %s", err.Error())
			}
			return val, err
		}

		switch t {
		case reflect.TypeOf(time.Time{}):
			ts, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return nil, fmt.Errorf("cant parse as timestamp: %s", err.Error())
			}
			return ts, nil
		}
		val, ok := reflect.New(t).Interface().(json.Unmarshaler)
		if ok {
			err := val.UnmarshalJSON([]byte("\"" + s + "\""))
			if err != nil {
				return nil, fmt.Errorf("json.Unmarshal error: %s", err.Error())
			}
			return val, nil
		}
	}
	return nil, fmt.Errorf("unsupported type: %s", t.String())
}

func IsNumber(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

func GetTypeByPath(modelType reflect.Type, path string) (reflect.Type, error) {
	parts := strings.Split(path, ".")

	t := modelType

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

	return t, nil
}

func getFieldByJsonTag(typ reflect.Type, name string) (reflect.StructField, bool) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		key := strings.Split(field.Tag.Get("json"), ",")[0]
		if key == name {
			return field, true
		}
	}

	return reflect.StructField{}, false
}
