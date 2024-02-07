package querymongo

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Find[Model any](q query.Query) (bson.D, *options.FindOptions, error) {
	d, err := Filter[Model](q.Filter)
	if err != nil {
		return nil, nil, err
	}

	opts := options.Find().
		SetSkip(int64(q.Pagination.Offset)).
		SetLimit(int64(q.Pagination.Limit)).
		SetSort(Sort(q.Sort))

	return d, opts, nil
}

func Filter[Model any](q query.Filter) (bson.D, error) {
	mongoFilter := bson.D{}
	for _, filter := range q {
		e, err := mongoOperator(
			filter.Op, filter.Field, filter.Value,
			reflect.TypeOf((*Model)(nil)).Elem(),
		)
		if err != nil {
			return nil, fmt.Errorf("query parsing error: %w", err)
		}
		mongoFilter = append(mongoFilter, e)
	}

	return mongoFilter, nil
}

func mongoOperator(q query.Operator, name string, value any, v reflect.Type) (bson.E, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t, err := query.GetTypeByPath(v, name)
	if err != nil {
		return bson.E{}, err
	}

	// val, err := query.GetValueForType(t, value)
	// if err != nil {
	// 	return bson.E{}, fmt.Errorf("cant get value for type %s, error: %s", t.Kind().String(), err.Error())
	// }

	// TODO
	// if _, ok := val.(uuid.UUID); ok && strings.HasSuffix(name, "id") {
	// 	name = strings.TrimSuffix(name, "id") + "_id"
	// }

	e := bson.E{
		Key:   name,
		Value: value,
	}

	switch q {
	case query.OperatorEqual:
		e.Value = bson.M{"$eq": value}
	case query.OperatorIn:
		values, err := interfacesSlice(v)
		if err != nil {
			return e, err
		}
		e = inFilter(name, values)
	case query.OperatorNotEqual:
		e.Value = bson.M{"$ne": value}
	case query.OperatorGreater:
		e.Value = bson.M{"$gt": value}
	case query.OperatorGreaterOrEqual:
		e.Value = bson.M{"$gte": value}
	case query.OperatorLess:
		e.Value = bson.M{"$lt": value}
	case query.OperatorLessOrEqual:
		e.Value = bson.M{"$lte": value}
	case query.OperatorSubString:
		if query.IsNumber(t) {
			return bson.E{ // special case for substring in numbers
				Key: "$where",
				Value: fmt.Sprintf(`function() {
					return String(this.%s).includes('%d');
					}`, name, value),
			}, nil
		} else {
			pattern, _ := value.(string)

			e.Value = primitive.Regex{Pattern: pattern, Options: "i"}
		}
	case query.OperatorDefault:
		e.Value = value
	}
	return e, nil
}

func interfacesSlice(v any) ([]any, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice {
		var out []any
		for i := 0; i < rv.Len(); i++ {
			out = append(out, rv.Index(i).Interface())
		}
		return out, nil
	}
	return nil, fmt.Errorf("type %s is not a slice", rv.Type().String())
}

func inFilter(name string, values []any) bson.E {
	if slices.Contains(values, nil) {
		values = slices.DeleteFunc(values, func(v any) bool {
			return v == nil
		})

		return bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{name: nil},
				bson.M{name: bson.M{"$in": values}},
			},
		}
	} else {
		return bson.E{
			Key:   name,
			Value: bson.M{"$in": values},
		}
	}

}
