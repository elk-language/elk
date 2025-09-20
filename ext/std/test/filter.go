package test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
)

var Filters []Filter

func RegisterFilter(filter Filter) {
	Filters = append(Filters, filter)
}

func SuiteMatchesFilters(suite *Suite) SuiteMatch {
	if suite.FullMatch {
		return SUITE_MATCH_FULL
	}

	var result SuiteMatch

	for _, filter := range Filters {
		suiteMatch := filter.SuiteMatches(suite)
		switch suiteMatch {
		case SUITE_MATCH_FALSE:
			return suiteMatch
		case SUITE_MATCH_FULL:
			if result == SUITE_MATCH_FALSE {
				result = SUITE_MATCH_FULL
			}
		case SUITE_MATCH_TRUE:
			result = SUITE_MATCH_TRUE
		}
	}

	if result == SUITE_MATCH_FALSE {
		return SUITE_MATCH_TRUE
	}
	return result
}

func CaseMatchesFilters(testCase *Case) bool {
	if testCase.FullMatch() {
		return true
	}

	for _, filter := range Filters {
		if !filter.CaseMatches(testCase) {
			return false
		}
	}

	return true
}

type SuiteMatch uint8

const (
	SUITE_MATCH_FALSE SuiteMatch = iota
	SUITE_MATCH_TRUE
	SUITE_MATCH_FULL
)

// Filter checks if a test case/suite matches
// and discards the test if it does not match
type Filter interface {
	CaseMatches(test *Case) bool
	SuiteMatches(suite *Suite) SuiteMatch
}

// PathFilter matches tests on filepaths and lines
type PathFilter struct {
	pattern string
	line    int
}

func NewPathFilter(path string) (*PathFilter, error) {
	var err error

	pathLine := strings.SplitN(path, ":", 2)
	line := -1
	pattern := pathLine[0]
	if len(pathLine) > 1 {
		line, err = strconv.Atoi(pathLine[1])
		if err != nil {
			return nil, fmt.Errorf("invalid path filter line: %w", err)
		}
	}

	if !doublestar.ValidatePattern(pattern) {
		return nil, fmt.Errorf("invalid path filter: %s", path)
	}

	return &PathFilter{
		pattern: pattern,
		line:    line,
	}, nil
}

func (p *PathFilter) LocationMatches(loc *position.Location) bool {
	pathMatches := doublestar.MatchUnvalidated(p.pattern, loc.FilePath)
	if !pathMatches {
		return false
	}

	if p.line < 0 {
		return true
	}
	return p.line >= loc.StartPos.Line && p.line <= loc.EndPos.Line
}

func (p *PathFilter) CaseMatches(test *Case) bool {
	return p.LocationMatches(test.Location())
}

func (p *PathFilter) SuiteMatches(suite *Suite) SuiteMatch {
	loc := suite.Location
	if loc == nil {
		return SUITE_MATCH_TRUE
	}

	pathMatches := doublestar.MatchUnvalidated(p.pattern, loc.FilePath)
	if !pathMatches {
		return SUITE_MATCH_FALSE
	}

	if p.line < 0 {
		return SUITE_MATCH_TRUE
	}
	if p.line == loc.StartPos.Line {
		return SUITE_MATCH_FULL
	}
	if p.line >= loc.StartPos.Line && p.line <= loc.EndPos.Line {
		return SUITE_MATCH_TRUE
	}
	return SUITE_MATCH_FALSE
}

// RegexFilter matches tests based on a regular expression
type RegexFilter struct {
	Regex *value.Regex
}

func NewRegexFilter(pattern string) (*RegexFilter, error) {
	regex, err := value.CompileRegex(pattern, bitfield.BitField8{})
	if err != nil {
		return nil, fmt.Errorf("invalid regex filter: %w", err)
	}

	return &RegexFilter{
		Regex: regex,
	}, nil
}

func (r *RegexFilter) CaseMatches(test *Case) bool {
	return r.Regex.MatchesString(test.FullNameWithSeparator())
}

func (r *RegexFilter) SuiteMatches(suite *Suite) SuiteMatch {
	return SUITE_MATCH_TRUE
}
