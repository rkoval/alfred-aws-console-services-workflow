package parsers

import (
	"fmt"
	"log"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

// Parser represents a parser.
type Parser struct {
	rawQuery string
	scanner  *Scanner
}

// NewParser returns a new instance of Parser.
func NewParser(rawQuery string) *Parser {
	reader := strings.NewReader(rawQuery)
	return &Parser{
		rawQuery: rawQuery,
		scanner:  NewScanner(reader),
	}
}

func (p *Parser) scanIntoTokens() ([]Token, bool) {
	tokens := []Token{}
	var tokenType TokenType
	var lastReadLiteral string
	var hasTrailingWhitespace bool
	var lastHasTrailingWhitespace bool
	var i int
Loop:
	for tokenType != EOF {
		if i >= 1000 {
			// prevent against accidental infinite loop
			log.Printf("infinite loop in parser.scanNextWord detected")
			break
		}
		i++
		lastHasTrailingWhitespace = hasTrailingWhitespace
		tokenType, lastReadLiteral, hasTrailingWhitespace = p.scanner.Scan()
		switch tokenType {
		case EOF:
			break Loop
		case WHITESPACE:
			continue
		default:
			tokens = append(tokens, Token{
				Type:  tokenType,
				Value: lastReadLiteral,
			})
		}
	}
	return tokens, lastHasTrailingWhitespace
}

func (p *Parser) Parse(ymlPath string) (*Query, []awsworkflow.AwsService) {
	awsServices := ParseConsoleServicesYml(ymlPath)
	query := &Query{
		RawQuery: p.rawQuery,
	}

	tokens, hasTrailingWhitespace := p.scanIntoTokens()

	query.HasTrailingWhitespace = hasTrailingWhitespace
	// TODO make this use string builder?
	var remainingQuery string
	for i, token := range tokens {
		switch token.Type {
		case WORD:
			if remainingQuery != "" {
				remainingQuery += " "
			}
			remainingQuery += token.Value
			if query.Service == nil {
				awsService := getAwsServiceById(remainingQuery, awsServices)
				// if awsService != nil && (i < len(tokens)-1 && hasTrailingWhitespace) {
				if awsService != nil {
					query.Service = awsService
					remainingQuery = ""
				}
			} else if query.SubService == nil {
				awsService := getAwsServiceById(remainingQuery, query.Service.SubServices)
				// if awsService != nil && (i < len(tokens)-1 && hasTrailingWhitespace) {
				if awsService != nil {
					query.SubService = awsService
					remainingQuery = ""
				}
			}
		case OPEN_ALL:
			query.HasOpenAll = true
		case SEARCH_ALIAS:
			query.HasDefaultSearchAlias = true
			remainingQuery += token.Value
		case REGION_OVERRIDE:
			if query.ProfileQuery != nil {
				continue
			}
			for _, region := range awsconfig.AllAWSRegions {
				if token.Value == region.Name {
					query.regionOverride = &region
					break
				}
			}

			if query.regionOverride == nil || (i >= len(tokens)-1 && !hasTrailingWhitespace) {
				query.RegionQuery = &token.Value
			}
		case PROFILE_OVERRIDE:
			if query.RegionQuery != nil {
				continue
			}
			for _, profile := range awsconfig.GetAwsProfiles() {
				if token.Value == profile.Name {
					query.ProfileOverride = &profile
					break
				}
			}

			if query.ProfileOverride == nil || (i >= len(tokens)-1 && !hasTrailingWhitespace) {
				query.ProfileQuery = &token.Value
			}
		default:
			panic(fmt.Errorf("no handler for token: %#v", token))
		}
	}

	query.RemainingQuery = remainingQuery

	return query, awsServices
}

func getAwsServiceById(id string, awsServices []awsworkflow.AwsService) *awsworkflow.AwsService {
	var awsService *awsworkflow.AwsService
	for i := range awsServices {
		if awsServices[i].Id == id {
			awsService = &awsServices[i]
			break
		}
	}
	return awsService
}
