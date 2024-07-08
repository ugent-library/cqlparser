package cqlparser

import (
	"errors"
	"regexp"
	"strings"

	"github.com/wolfeidau/stringtokenizer"
)

var reSpace = regexp.MustCompile(`^\s+$`)

type token struct {
	Type       int
	Value      string
	Terminated bool
}

type lexer struct {
	tokens []*token
}

func newTokenFromString(val string) *token {
	var typ int
	if t, ok := cqlKeywords[strings.ToLower(val)]; ok {
		typ = t
	} else {
		typ = CQL_WORD
		if len(val) > 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
			val = strings.ReplaceAll(val, `\"`, `"`)
		}
	}
	return &token{
		Type:  typ,
		Value: val,
	}
}

func newLexerFromString(val string) (*lexer, error) {
	var strTokens []string
	t := stringtokenizer.NewStringTokenizer(strings.NewReader(val), "/<>=()\" ", true)
	for t.HasMoreTokens() {
		strTokens = append(strTokens, t.NextToken())
	}

	tokens := make([]*token, 0, len(strTokens))
	pos := 0
	for pos < len(strTokens) {
		strToken := strTokens[pos]
		var nextToken string
		if pos+1 < len(strTokens) {
			nextToken = strTokens[pos+1]
		}

		if strToken == "<" && nextToken == "=" {
			tokens = append(tokens, newTokenFromString("<="))
			pos++
		} else if strToken == "<" && nextToken == ">" {
			tokens = append(tokens, newTokenFromString("<>"))
			pos++
		} else if strToken == ">" && nextToken == "=" {
			tokens = append(tokens, newTokenFromString(">="))
			pos++
		} else if strToken == "\"" {
			escaping := false
			substr := strings.Builder{}
			substr.WriteString("\"")
			pos++
			for pos < len(strTokens) {
				subTokenStr := strTokens[pos]
				substr.WriteString(subTokenStr)
				if escaping {
					escaping = false
				} else if subTokenStr == "\"" {
					break
				} else if subTokenStr == "\\" {
					escaping = true
				}
				pos++
			}
			if pos >= len(strTokens) {
				return nil, errors.New("unterminated string")
			}

			token := newTokenFromString(substr.String())
			token.Terminated = true
			if len(tokens) > 0 {
				tokens[len(tokens)-1].Terminated = true
			}
			tokens = append(tokens, token)
		} else if reSpace.MatchString(strToken) {
			if len(tokens) > 0 {
				tokens[len(tokens)-1].Terminated = true
			}
		} else {
			tokens = append(tokens, newTokenFromString(strToken))
		}

		pos++
	}

	i := 0
	for i < len(tokens) {
		token := tokens[i]
		if token.Value == "\\" {
			var s string
			var replace bool = false
			if i > 0 {
				prevToken := tokens[i-1]
				if prevToken.Type == CQL_WORD && !prevToken.Terminated {
					s = prevToken.Value + s
					i--
					tokens = append(tokens[:i], tokens[i+1:]...)
					replace = true
				}
			}
			if token.Terminated && i < len(tokens) {
				nextToken := tokens[i+1]
				if nextToken.Type == CQL_WORD {
					s += nextToken.Value
					tokens = append(tokens[:i+1], tokens[i+2:]...)
					replace = true
				}
			}
			if replace {
				tokens[i] = newTokenFromString(s)
			}
		}
		i++
	}

	return &lexer{
		tokens: tokens,
	}, nil
}

func (l *lexer) HasMoreTokens() bool {
	return len(l.tokens) > 0
}

func (l *lexer) NextToken() *token {
	if len(l.tokens) > 0 {
		t := l.tokens[0]
		l.tokens = l.tokens[1:]
		return t
	}
	return nil
}
