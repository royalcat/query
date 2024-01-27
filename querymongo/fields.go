package querymongo

import (
	"slices"
	"strconv"
	"strings"

	"github.com/royalcat/query"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Project(fields query.Fields, res query.ModelLink) (mongo.Pipeline, error) {
	agg := mongo.Pipeline{}
	q := fields
	prjFields := []string{}
	for k, v := range res.Resolvers {
		var g query.Fields
		g, q = q.GroupFilter(k.ResolvedName)
		if len(g) == 0 {
			continue
		}

		linkIdName := k.LinkIdName
		if linkIdName == "" || linkIdName == "id" {
			linkIdName = "_id"
		}

		switch k.Type {
		case query.Array:
			g = cleanFilterFirstPart(g)
			childAgg, err := Project(g, v)
			if err != nil {
				return nil, err
			}
			prjFields = append(prjFields, k.ResolvedName)
			agg = append(agg, lookup(
				v.Collection, k.IdName, linkIdName, k.ResolvedName,
				childAgg,
			)...)
		case query.Single:
			childAgg, err := Project(g, v)
			if err != nil {
				return nil, err
			}
			prjFields = append(prjFields, k.ResolvedName)
			agg = append(agg, lookupSingle(
				v.Collection, k.IdName, linkIdName, k.ResolvedName,
				childAgg,
			)...)
		case query.SingleLast:
			childAgg, err := Project(g, v)
			if err != nil {
				return nil, err
			}
			prjFields = append(prjFields, k.ResolvedName)

			childrenPipeline := mongo.Pipeline{
				bson.D{{Key: "$sort", Value: bson.M{"_id": -1}}},
				bson.D{{Key: "$limit", Value: 1}},
			}
			childrenPipeline = append(childrenPipeline, childAgg...)
			lk := lookupSingle(
				v.Collection, k.IdName, linkIdName, k.ResolvedName,
				childrenPipeline,
			)

			agg = append(agg, lk...)
		}
	}
	prjFields = append(prjFields, q...)
	prjFields = query.SliceUnique(prjFields)

	if len(prjFields) > 0 {
		project := bson.D{}
		for _, k := range prjFields {
			k = clearKeyForMongo(k)
			k = cleanProjectPath(k)
			project = append(project, bson.E{
				Key: k, Value: 1,
			})
		}
		agg = append(agg, bson.D{{Key: "$project", Value: project}})
	}

	return agg, nil
}

func cleanFirstPart(path string) string {
	parts := strings.Split(path, ".")
	if _, err := strconv.ParseUint(parts[0], 10, 64); err == nil {
		return strings.Join(parts[1:], ".")
	}
	return path
}

func cleanProjectPath(s string) string {
	parts := strings.Split(s, ".")
	parts = slices.DeleteFunc(parts, func(p string) bool {
		if _, err := strconv.ParseUint(p, 10, 64); err == nil {
			return true
		}
		return false
	})
	return strings.Join(parts, ".")
}

func cleanFilterFirstPart(f query.Fields) query.Fields {
	for k, v := range f {
		f[k] = cleanFirstPart(v)
	}
	return f
}

func clearKeyForMongo(k string) string {
	if k == "id" {
		k = "_id"
	} else {
		k = strings.ReplaceAll(k, ".id", "._id")
	}
	return k
}
