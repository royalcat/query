package querymongo

import (
	"fmt"

	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
)

func Sort(s query.Sort) bson.M {
	m := bson.M{}

	for _, f := range s {
		k := clearKeyForMongo(f.Key)
		switch f.Order {
		case query.ASC:
			m[k] = 1
		case query.DESC:
			m[k] = -1
		default:
			panic(fmt.Errorf("unknown sort order: %d", f.Order))
		}
	}
	return m
}
