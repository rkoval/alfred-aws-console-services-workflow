package parsers

// Token represents a lexical token.
type Token int

const (
	ILLEGAL Token = iota

	EOF
	WHITESPACE
	WORD

	// special tokens
	OPEN_ALL
)
