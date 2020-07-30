package parsers

type Query struct {
	ServiceId    string
	SubServiceId string
	HasTrailingWhitespace bool
	HasOpenAll   bool
	SearchTerms  []string
}

func (q *Query) ShouldUseDefaultSearch() bool {
	return q.SubServiceId == "" && q.SearchTerms != nil && len(q.SearchTerms) > 0
}
