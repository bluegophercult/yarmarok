package service

// Direction represents the sorting direction for ordering query results.
type Direction int

const (
	ASC  Direction = iota + 1 // Ascending order
	DESC                      // Descending order
)

// Operator represents the comparison operator used for filtering query results.
type Operator int

const (
	LT  Operator = iota + 1 // Less than
	LTE                     // Less than or equal to
	EQ                      // Equal to
	GT                      // Greater than
	GTE                     // Greater than or equal to
	NE                      // Not equal to
	IN                      // In a specified list
	NI                      // Not in a specified list
	CN                      // Contains specified value
	CA                      // Contains any of the specified values
)

// Query represents a query that can be executed against a storage service.
type Query struct {
	Filter  *Filter
	OrderBy *OrderBy
	Limit   int
	Offset  int
}

// Filter represents a single filter criterion for a query.
type Filter struct {
	Field    string
	Operator Operator
	Value    any
	Next     *Filter
}

// OrderBy represents a single ordering criterion for a query.
type OrderBy struct {
	Field     string
	Direction Direction
}

// WithFilter adds a filter criterion to the query.
func (q Query) WithFilter(field string, op Operator, val any) Query {
	filter := &Filter{Field: field, Operator: op, Value: val}

	head := q.Filter
	if head == nil {
		q.Filter = filter
		return q
	}

	last := head
	for last.Next != nil {
		last = last.Next
	}

	last.Next = filter

	return q
}

// WithOrderBy sets the ordering criterion for the query.
func (q Query) WithOrderBy(field string, dir Direction) Query {
	q.OrderBy = &OrderBy{Field: field, Direction: dir}
	return q
}

// WithLimit sets the maximum number of results to return for the query.
func (q Query) WithLimit(lim int) Query {
	q.Limit = lim
	return q
}

// WithOffset sets the number of results to skip before returning results for the query.
func (q Query) WithOffset(off int) Query {
	q.Offset = off
	return q
}
