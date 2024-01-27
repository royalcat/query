package querymongo

import (
	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func ToMongoAggIds(q query.Query, ml query.ModelLink) (mongo.Pipeline, error) {
	agg := mongo.Pipeline{
		bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}},
	}

	prjAgg, err := q.Fields().toMongoProject(ml)
	if err != nil {
		return nil, err
	}
	agg = append(agg, prjAgg...)

	m, err := q.Filter.ToMongo(ml.FullModelType)
	if err != nil {
		return nil, err
	}
	if len(m) > 0 {
		agg = append(agg, bson.D{{Key: "$match", Value: m}})
	}

	s := q.Sort.ToMongoSort()
	if len(s) > 0 {
		sort := append(
			mtoD(s),
			bson.E{Key: "_id", Value: -1},
		)
		agg = append(agg, bson.D{{Key: "$sort", Value: sort}})
	}

	if q.Pagination.Offset != 0 {
		agg = append(agg, bson.D{{Key: "$skip", Value: q.Pagination.Offset}})
	}

	if q.Pagination.Limit != 0 {
		agg = append(agg, bson.D{{Key: "$limit", Value: q.Pagination.Limit}})
	}

	// if settings.Config.Debug {
	// 	aggData, err := bson.MarshalExtJSON(bson.M{"agg": agg}, true, true)
	// 	if err != nil {
	// 		logrus.Error(err)
	// 	} else {
	// 		logrus.Debug("Generated agg data: ", string(aggData))
	// 	}
	// }

	return agg, nil
}

func mtoD(m bson.M) bson.D {
	d := bson.D{}
	for k, v := range m {
		d = append(d, bson.E{Key: k, Value: v})
	}
	return d
}