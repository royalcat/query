package queryreflect

import "github.com/royalcat/query"

func ApplyQuery[D any](q query.Query, in []D) ([]D, error) {
	var err error
	if len(q.Filter) > 0 {
		in, err = ApplyFilter(q.Filter, in)
		if err != nil {
			return nil, err
		}
	}
	if len(q.Sort) > 0 {
		in, err = ApplySort(q.Sort, in)
		if err != nil {
			return nil, err
		}
	}

	if len(in) <= int(q.Pagination.Offset) {
		return []D{}, nil
	} else if len(in) < int(q.Pagination.Offset+q.Pagination.Limit) || q.Pagination.Limit == 0 {
		return in[q.Pagination.Offset:], nil
	} else {
		return in[q.Pagination.Offset : q.Pagination.Offset+q.Pagination.Limit], nil
	}
}

type PageGetter[D any] func(q query.Query) ([]D, error)

func ApplyQueryWithNext[D any](q query.Query, getPage PageGetter[D]) (out []D, err error) {
	pageQuery := q.Copy()
	for {
		page, err := getPage(pageQuery)
		if err != nil {
			return nil, err
		}

		page, err = ApplyFilter(q.Filter, page)
		if err != nil {
			return nil, err
		}

		out = append(out, page...)

		if len(page) != int(q.Pagination.Limit) {
			break
		}

		if len(out) >= int(q.Pagination.Offset+q.Pagination.Limit) {
			break
		}
	}

	out, err = ApplyQuery(q, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
