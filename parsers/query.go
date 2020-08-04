package parsers

type Query struct {
	ServiceId             string
	SubServiceId          string
	HasTrailingWhitespace bool
	HasOpenAll            bool
	HasDefaultSearchAlias bool
	RemainingQuery        string
}

func (q *Query) IsEmpty() bool {
	return q.ServiceId == "" && q.SubServiceId == "" && q.RemainingQuery == ""
}
