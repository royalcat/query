package querymongo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const myVarName = "tgr"

func lookup(collection, idName, linkIdName, field string, children mongo.Pipeline) mongo.Pipeline {
	subPipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.M{
				"$expr": bson.M{
					"$in": bson.A{
						"$" + linkIdName,
						"$$" + myVarName,
					},
				},
			}},
		},
	}
	subPipeline = append(subPipeline, children...)

	return mongo.Pipeline{
		bson.D{{
			Key: "$set",
			Value: bson.M{
				idName: bson.M{
					"$ifNull": bson.A{
						"$" + idName,
						bson.A{},
					},
				},
			},
		}},
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from": collection,
				"let": bson.M{
					myVarName: "$" + idName,
				},
				"pipeline": subPipeline,
				"as":       field,
			}},
		},
	}
}

func lookupSingle(collection, idName, linkIdName, field string, children mongo.Pipeline) mongo.Pipeline {
	subPipeline := mongo.Pipeline{
		bson.D{{
			Key: "$match",
			Value: bson.M{
				"$expr": bson.M{
					"$eq": bson.A{
						"$" + linkIdName,
						"$$" + myVarName,
					},
				},
			}},
		},
	}
	subPipeline = append(subPipeline, children...)

	return mongo.Pipeline{
		bson.D{{
			Key: "$lookup",
			Value: bson.M{
				"from": collection,
				"let": bson.M{
					myVarName: "$" + idName,
				},
				"pipeline": subPipeline,
				"as":       field,
			},
		}},
		bson.D{{
			Key: "$set",
			Value: bson.M{
				field: bson.M{
					"$arrayElemAt": bson.A{"$" + field, 0},
				},
			},
		}},
	}
}
