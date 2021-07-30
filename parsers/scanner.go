package parsers

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
)

var eof = rune(0)

// Scanner represents a lexical scanner.
type Scanner struct {
	reader *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.reader.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (TokenType, string, bool) {
	ch := s.read()

	if ch == eof {
		return EOF, "", false
	}

	s.unread()
	if isWhitespace(ch) {
		token, literal := s.scanWhitespace()
		return token, literal, true
	}
	token, literal, hasTrailingWhitespace := s.scanWord()
	return token, literal, hasTrailingWhitespace
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok TokenType, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var i int
	for {
		if i >= 1000 {
			// prevent against accidental infinite loop
			log.Println("infinite loop in scanner.scanWord detected")
			break
		}
		i++
		ch := s.read()
		if ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WHITESPACE, buf.String()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

// scanWord consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanWord() (TokenType, string, bool) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	var i int
	var hasTrailingWhitespace bool
	for {
		if i >= 1000 {
			// prevent against accidental infinite loop
			log.Println("infinite loop in scanner.scanWord detected")
			break
		}
		i++
		ch := s.read()
		if ch == eof {
			break
		} else if isWhitespace(ch) {
			hasTrailingWhitespace = true
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	stringBuf := buf.String()

	if strings.HasPrefix(stringBuf, aliases.Search) {
		return SEARCH_ALIAS, stringBuf[len(aliases.Search):], hasTrailingWhitespace
	}

	if strings.HasPrefix(stringBuf, aliases.OverrideAwsRegion) {
		return REGION_OVERRIDE, stringBuf[len(aliases.OverrideAwsRegion):], hasTrailingWhitespace
	}

	switch stringBuf {
	case "OPEN_ALL":
		return OPEN_ALL, stringBuf, hasTrailingWhitespace
	}

	return WORD, stringBuf, hasTrailingWhitespace
}
