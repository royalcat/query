package query

import (
	"strings"
)

type Fields []string

func (f Fields) groupFilter(prefix string) (Fields, Fields) {
	g := Fields{}
	r := Fields{}
	for _, k := range f {
		if kCut, ok := strings.CutPrefix(k, prefix+"."); ok {
			g = append(g, kCut)
		} else {
			r = append(r, k)
		}
	}
	return g, r
}
