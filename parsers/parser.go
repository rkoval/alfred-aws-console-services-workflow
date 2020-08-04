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
}

// NewParser returns a new instance of Parser.
func NewParser(reader io.Reader) *Parser {
	return &Parser{scanner: NewScanner(reader)}
}

func (p *Parser) scanNextWord(query *Query) (string, bool) {
	var token Token
	var lastReadLiteral string
	var hasTrailingWhitespace bool
	var lastHasTrailingWhitespace bool
	var i int
Loop:
	for token != WORD {
		if i >= 1000 {
			// prevent against accidental infinite loop
			log.Printf("infinite loop in parser.scanNextWord detected")
			break
		}
		i++
		lastHasTrailingWhitespace = hasTrailingWhitespace
		token, lastReadLiteral, hasTrailingWhitespace = p.scanner.Scan()
		// TODO add in region and other tokens
		switch token {
		case EOF:
			break Loop
		case WHITESPACE:
			continue
		case OPEN_ALL:
			query.HasOpenAll = true
		case WORD:
			return lastReadLiteral, hasTrailingWhitespace
		default:
			panic(fmt.Errorf("found %q, expected field", lastReadLiteral))
		}
	}
	return "", lastHasTrailingWhitespace
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
	query.HasTrailingWhitespace = hasTrailingWhitespace
	if lastReadLiteral == "" {
		return query
	} else if strings.HasPrefix(lastReadLiteral, defaultSearchAlias) {
		lastReadLiteral = lastReadLiteral[len(defaultSearchAlias):]
		query.HasDefaultSearchAlias = true
		query.RemainingQuery += lastReadLiteral
		if hasTrailingWhitespace {
			query.RemainingQuery += " "
		}
	} else {
		query.SubServiceId = lastReadLiteral
	}

	var i int
	for {
		if i >= 1000 {
			// prevent against accidental infinite loop
			log.Printf("infinite loop in parser.Parse detected")
			break
		}
		i++
		lastReadLiteral, hasTrailingWhitespace := p.scanNextWord(query)
		query.HasTrailingWhitespace = hasTrailingWhitespace
		if lastReadLiteral == "" {
			return query
		}
		query.RemainingQuery += lastReadLiteral
		if hasTrailingWhitespace {
			query.RemainingQuery += " "
		}
	}

	return query
}
