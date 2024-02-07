package querymongo

import (
	"slices"
	"strconv"
	"strings"

	"github.com/royalcat/query"
)

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
