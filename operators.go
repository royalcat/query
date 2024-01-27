package query

type Operator string

const (
	OperatorEmpty          Operator = ""
	OperatorEqual          Operator = "eq"
	OperatorIn             Operator = "in"
	OperatorNotEqual       Operator = "ne"
	OperatorGreater        Operator = "gt"
	OperatorGreaterOrEqual Operator = "gte"
	OperatorLess           Operator = "lt"
	OperatorLessOrEqual    Operator = "lte"
	OperatorSubString      Operator = "substr"
)
