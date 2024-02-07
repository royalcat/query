package query

import (
	"reflect"
)

type Filter []FieldFilter

type FieldFilter struct {
	Field string
	Op    Operator
	Value any
}

func (q Filter) Fields() Fields {
	fields := []string{}
	for _, f := range q {
		fields = append(fields, f.Field)
	}
	return fields
}

func (q Filter) Operators() map[string]Operator {
	ops := map[string]Operator{}
	for _, f := range q {
		// name, operator := ParseKey(k)
		ops[f.Field] = f.Op
	}
	return ops
}

// func ParseKey(key FilterKey) (name string, operator Operator) {
// 	opStart := strings.Index(key, "{")

// 	if opStart == -1 {
// 		return key, operator

// 	}
// 	name = key[:opStart]
// 	operator = Operator(key[opStart : len(key)-1])
// 	return name, operator
// }

func IsNumber(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}
