package query

type SortOrder int8

const (
	ASC  SortOrder = 1
	DESC SortOrder = -1
)

type SortField struct {
	Key   string
	Order SortOrder
}

type Sort []SortField

func (s Sort) Fields() Fields {
	keys := make([]string, 0, len(s))
	for _, v := range s {
		keys = append(keys, v.Key)
	}
	return SliceUnique(keys)
}

func (s Sort) Copy() Sort {
	c := make(Sort, len(s))
	copy(c, s)
	return c
}

func (s Sort) Get(key string) (SortOrder, bool) {
	for _, v := range s {
		if v.Key == key {
			return v.Order, true
		}
	}
	return 0, false
}

func (s *Sort) Set(key string, order SortOrder) {
	if s == nil {
		m := make(Sort, 0)
		s = &m
	}

	for i := range *s {
		if (*s)[i].Key == key {
			(*s)[i].Order = order
			return
		}
	}
	*s = append(*s, SortField{Key: key, Order: order})
}
