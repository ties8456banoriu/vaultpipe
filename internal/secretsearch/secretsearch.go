// Package secretsearch provides full-text and pattern-based search
// over a secrets map, returning matching key-value pairs.
package secretsearch

import (
	"errors"
	"regexp"
	"strings"
)

// Mode controls how the query is interpreted.
type Mode string

const (
	ModeExact  Mode = "exact"
	ModePrefix Mode = "prefix"
	ModeRegex  Mode = "regex"
)

// Result holds a single search match.
type Result struct {
	Key   string
	Value string
}

// Searcher searches secrets by key or value.
type Searcher struct {
	mode      Mode
	searchKey bool
	searchVal bool
}

// New returns a Searcher configured with the given mode.
// searchKey and searchVal control which fields are matched against.
func New(mode Mode, searchKey, searchVal bool) (*Searcher, error) {
	switch mode {
	case ModeExact, ModePrefix, ModeRegex:
	default:
		return nil, errors.New("secretsearch: unsupported mode: " + string(mode))
	}
	if !searchKey && !searchVal {
		return nil, errors.New("secretsearch: at least one of searchKey or searchVal must be true")
	}
	return &Searcher{mode: mode, searchKey: searchKey, searchVal: searchVal}, nil
}

// Search returns all secrets whose key or value matches the query.
func (s *Searcher) Search(secrets map[string]string, query string) ([]Result, error) {
	if len(secrets) == 0 {
		return nil, errors.New("secretsearch: secrets map is empty")
	}
	if query == "" {
		return nil, errors.New("secretsearch: query must not be empty")
	}

	var re *regexp.Regexp
	if s.mode == ModeRegex {
		var err error
		re, err = regexp.Compile(query)
		if err != nil {
			return nil, errors.New("secretsearch: invalid regex: " + err.Error())
		}
	}

	var results []Result
	for k, v := range secrets {
		if s.matches(k, v, query, re) {
			results = append(results, Result{Key: k, Value: v})
		}
	}
	return results, nil
}

func (s *Searcher) matches(key, value, query string, re *regexp.Regexp) bool {
	check := func(target string) bool {
		switch s.mode {
		case ModeExact:
			return target == query
		case ModePrefix:
			return strings.HasPrefix(target, query)
		case ModeRegex:
			return re.MatchString(target)
		}
		return false
	}
	if s.searchKey && check(key) {
		return true
	}
	if s.searchVal && check(value) {
		return true
	}
	return false
}
