package main

import (
	"errors"
	"fmt"
	"strings"
)

// A rule is a URL fragment, URL regular expression, or content phrase
// that will be matched against a page in the process of determining its score.
type rule struct {
	t       ruleType
	content string
}

type ruleType int

const (
	defaultRule ruleType = iota
	urlMatch
	urlRegex
	hostRegex
	domainRegex
	pathRegex
	queryRegex
	contentPhrase
)

func (r rule) String() string {
	switch r.t {
	case urlMatch:
		return r.content
	case urlRegex, hostRegex, pathRegex, queryRegex:
		suffix := ""
		switch r.t {
		case hostRegex:
			suffix = "h"
		case pathRegex:
			suffix = "p"
		case queryRegex:
			suffix = "q"
		case domainRegex:
			suffix = "d"
		}
		return "/" + r.content + "/" + suffix
	case contentPhrase:
		return "<" + r.content + ">"
	}
	panic(fmt.Errorf("invalid rule type: %d", r.t))
}

// parseRule parses a rule from the beginning of s, returning the rule
// and any remaining unconsumed characters from s.
func parseRule(s string) (r rule, leftover string, err error) {
	s = strings.TrimLeft(s, " \t\r\n\f")
	if s == "" {
		return rule{}, "", errors.New("blank rule")
	}

	switch s[0] {
	case '/':
		r.t = urlRegex
		slash := strings.LastIndex(s, "/")
		if slash == 0 {
			return rule{}, s, errors.New("unmatched slash")
		}
		r.content = s[1:slash]
		s = s[slash+1:]
		if s != "" {
			switch s[0] {
			case 'h':
				r.t = hostRegex
				s = s[1:]
			case 'p':
				r.t = pathRegex
				s = s[1:]
			case 'q':
				r.t = queryRegex
				s = s[1:]
			case 'd':
				r.t = domainRegex
				s = s[1:]
			}
		}
	case '<':
		r.t = contentPhrase
		bracket := strings.LastIndex(s, ">")
		if bracket == -1 {
			return rule{}, s, errors.New("unmatched '<'")
		}
		r.content = wordString(s[1:bracket])
		s = s[bracket+1:]
	default:
		r.t = urlMatch
		space := strings.Index(s, " ")
		if space == -1 {
			r.content = strings.ToLower(s)
			s = ""
		} else {
			r.content = strings.ToLower(s[:space])
			s = s[space:]
		}
		if r.content == "default" {
			r.t = defaultRule
		}
	}

	return r, s, nil
}
