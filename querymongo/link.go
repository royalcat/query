package querymongo

import (
	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func FullMongoAgg(l query.ModelLink) mongo.Pipeline {
	p := mongo.Pipeline{}
	for f, m := range l.Resolvers {
		linkIdName := f.LinkIdName
		if linkIdName == "" || linkIdName == "id" {
			linkIdName = "_id"
		}
		idName := f.IdName
		if idName == "" || idName == "id" {
			idName = "_id"
		}

		switch f.Type {
		case query.Array:
			p = append(p, lookup(
				m.Collection, idName, linkIdName, f.ResolvedName,
				FullMongoAgg(m),
			)...)
		case query.Single:
			p = append(p, lookupSingle(
				m.Collection, idName, linkIdName, f.ResolvedName,
				FullMongoAgg(m),
			)...)
		case query.SingleLast:
			childrenPipeline := mongo.Pipeline{
				bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}},
				bson.D{{Key: "$limit", Value: 1}},
			}
			childrenPipeline = append(childrenPipeline, FullMongoAgg(m)...)
			p = append(p, lookupSingle(
				m.Collection, idName, linkIdName, f.ResolvedName,
				childrenPipeline,
			)...)
		}
	}
	return p
}
