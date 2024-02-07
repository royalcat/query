package query

type Operator string

const (
	OperatorDefault        Operator = ""
	OperatorEqual          Operator = "eq"
	OperatorIn             Operator = "in"
	OperatorNotEqual       Operator = "ne"
	OperatorGreater        Operator = "gt"
	OperatorGreaterOrEqual Operator = "gte"
	OperatorLess           Operator = "lt"
	OperatorLessOrEqual    Operator = "lte"
	OperatorSubString      Operator = "substr"
)

func isOperator(op Operator) bool {
	return op == OperatorDefault ||
		op == OperatorEqual || op == OperatorIn || op == OperatorNotEqual ||
		op == OperatorGreater || op == OperatorGreaterOrEqual ||
		op == OperatorLess || op == OperatorLessOrEqual || op == OperatorSubString
}
