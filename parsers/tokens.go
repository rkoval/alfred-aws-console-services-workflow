package parsers

// Token represents a lexical token.
type TokenType int

const (
	ILLEGAL TokenType = iota

	EOF
	WHITESPACE
	WORD

	// special tokens
	OPEN_ALL
	SEARCH_ALIAS
	REGION_OVERRIDE
	PROFILE_OVERRIDE
)

type Token struct {
	Type  TokenType
	Value string
}
