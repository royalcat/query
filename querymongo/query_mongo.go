package querymongo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Find(q query.Query, m reflect.Type) (bson.D, *options.FindOptions, error) {
	d, err := q.Filter.ToMongo(m)
	if err != nil {
		return nil, nil, err
	}

	opts := options.Find().
		SetSkip(int64(q.Pagination.Offset)).
		SetLimit(int64(q.Pagination.Limit)).
		SetSort(q.Sort.ToMongoSort())

	return d, opts, nil
}

func Filter(q query.Filter, model reflect.Type) (bson.D, error) {
	mongoFilter := bson.D{}
	for k, v := range q {
		name, operator := query.ParseKey(k)
		e, err := operator.toMongo(name, v, model)
		if err != nil {
			return nil, fmt.Errorf("query parsing error: %w", err)
		}
		mongoFilter = append(mongoFilter, e)
	}

	return mongoFilter, nil
}

func toMongo(q query.Operator, name string, value string, v reflect.Type) (bson.E, error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t, err := query.GetTypeByPath(v, name)
	if err != nil {
		return bson.E{}, err
	}

	val, err := query.GetValueForType(t, value)
	if err != nil {
		return bson.E{}, fmt.Errorf("cant get value for type %s, error: %s", t.Kind().String(), err.Error())
	}

	// TODO
	// if _, ok := val.(uuid.UUID); ok && strings.HasSuffix(name, "id") {
	// 	name = strings.TrimSuffix(name, "id") + "_id"
	// }

	e := bson.E{
		Key:   name,
		Value: val,
	}

	switch q {
	case query.OperatorEqual:
		e.Value = bson.M{"$eq": val}
	case query.OperatorIn:
		e = inFilter(name, value)
	case query.OperatorNotEqual:
		e.Value = bson.M{"$ne": val}
	case query.OperatorGreater:
		e.Value = bson.M{"$gt": val}
	case query.OperatorGreaterOrEqual:
		e.Value = bson.M{"$gte": val}
	case query.OperatorLess:
		e.Value = bson.M{"$lt": val}
	case query.OperatorLessOrEqual:
		e.Value = bson.M{"$lte": val}
	case query.OperatorSubString:
		if query.IsNumber(t) {
			return bson.E{ // special case for substring in numbers
				Key: "$where",
				Value: fmt.Sprintf(`function() {
					return String(this.%s).includes('%d');
					}`, name, val),
			}, nil
		} else {
			pattern, _ := val.(string)

			e.Value = primitive.Regex{Pattern: pattern, Options: "i"}
		}
	case query.OperatorEmpty:
		e.Value = val
	}
	return e, nil
}

func inFilter(name, value string) bson.E {
	if strings.HasPrefix(value, ",") {
		vals := strings.Split(value, ",")
		return bson.E{
			Key: "$or",
			Value: bson.A{
				bson.M{name: nil},
				bson.M{name: bson.M{"$in": vals}},
			},
		}
	} else {
		return bson.E{
			Key:   name,
			Value: bson.M{"$in": strings.Split(value, ",")},
		}
	}

}
