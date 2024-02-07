package query

import "slices"

type Query struct {
	Search     string
	Filter     Filter
	Sort       Sort
	Pagination Pagination
}

func (q Query) Copy() Query {
	return Query{
		Search:     q.Search,
		Filter:     slices.Clone(q.Filter),
		Sort:       q.Sort.Copy(),
		Pagination: q.Pagination,
	}
}

func (q Query) Fields() Fields {
	f := append(q.Filter.Fields(), q.Sort.Fields()...)
	f = SliceUnique(f)
	return f
}

type Pagination struct {
	Offset uint64
	Limit  uint64
}
