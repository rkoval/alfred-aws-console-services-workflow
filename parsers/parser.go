package parsers

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	scanner *Scanner
	buffer  struct {
		token           Token  // last read token
		lastReadLiteral string // last read literal
	}
}

// NewParser returns a new instance of Parser.
func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: NewScanner(reader)}
}

func (p *Parser) scanNextWord(query *Query) (string, bool) {
	var token Token
	var lastToken Token
	var lastReadLiteral string
	var i int
Loop:
	for token != WORD {
		if i >= 100 {
			// prevent against accidental infinite loop
			log.Printf("infinite loop in parser.scanNextWord detected")
			break
		}
		i++
		lastToken = token
		token, lastReadLiteral = p.scanner.Scan()
		// TODO add in region and other tokens
		switch token {
		case EOF:
			break Loop
		case WHITESPACE:
			continue
		case OPEN_ALL:
			query.HasOpenAll = true
		case WORD:
			return lastReadLiteral, lastToken == WHITESPACE
		default:
			panic(fmt.Errorf("found %q, expected field", lastReadLiteral))
		}
	}
	return "", lastToken == WHITESPACE
}

func (p *Parser) Parse() *Query {
	query := &Query{}

	lastReadLiteral, hasTrailingWhitespace := p.scanNextWord(query)
	if lastReadLiteral == "" {
		query.HasTrailingWhitespace = hasTrailingWhitespace
		return query
	}
	query.ServiceId = lastReadLiteral

	defaultSearchAlias := os.Getenv("ALFRED_AWS_CONSOLE_SERVICES_WORKFLOW_SEARCH_ALIAS")
	if defaultSearchAlias == "" {
		defaultSearchAlias = ","
	}
	lastReadLiteral, hasTrailingWhitespace = p.scanNextWord(query)
	if lastReadLiteral == "" {
		query.HasTrailingWhitespace = hasTrailingWhitespace
		return query
	} else if strings.HasPrefix(lastReadLiteral, defaultSearchAlias) {
		lastReadLiteral = lastReadLiteral[len(defaultSearchAlias):]
		query.SearchTerms = append(query.SearchTerms, lastReadLiteral)
	} else {
		query.SubServiceId = lastReadLiteral
	}

	var i int
	for {
		if i >= 100 {
			// prevent against accidental infinite loop
			log.Printf("infinite loop in parser.Parse detected")
			break
		}
		i++
		lastReadLiteral, hasTrailingWhitespace := p.scanNextWord(query)
		if lastReadLiteral == "" {
			query.HasTrailingWhitespace = hasTrailingWhitespace
			return query
		}
		query.SearchTerms = append(query.SearchTerms, lastReadLiteral)
	}

	return query
}
