package cqlparser

import (
	"errors"
	"strings"
)

func Parse(query string) (Node, error) {
	parser, err := newParser(query)
	if err != nil {
		return nil, err
	}
	return parser.parse()
}

type parser struct {
	query    string
	lexer    *lexer
	curToken *token
}

func newParser(query string) (*parser, error) {
	lexer, err := newLexerFromString(query)
	if err != nil {
		return nil, err
	}
	return &parser{
		query: query,
		lexer: lexer,
	}, nil
}

func (p *parser) nextToken() *token {
	p.curToken = p.lexer.NextToken()
	return p.curToken
}

func (p *parser) parse() (Node, error) {
	p.nextToken()
	root, err := p.parseQuery("srw.ServerChoice", CQL_SCR)
	if err != nil {
		return nil, err
	}
	if t := p.nextToken(); t != nil && t.Type != CQL_EOF {
		return nil, errors.New("junk after end: " + t.Value)
	}
	return root, nil
}

func (p *parser) parseQuery(qualifier string, relation int) (Node, error) {
	term, err := p.parseTerm(qualifier, relation)
	if err != nil {
		return nil, err
	}

	for p.curToken != nil && p.curToken.Type != CQL_EOF && p.curToken.Type != CQL_RPAREN {
		switch p.curToken.Type {
		case CQL_AND:
			p.nextToken()
			term2, err := p.parseTerm(qualifier, relation)
			if err != nil {
				return nil, err
			}
			term = &BooleanNode{
				Op:    "and",
				Left:  term,
				Right: term2,
			}
		case CQL_OR:
			p.nextToken()
			term2, err := p.parseTerm(qualifier, relation)
			if err != nil {
				return nil, err
			}
			term = &BooleanNode{
				Op:    "or",
				Left:  term,
				Right: term2,
			}
		case CQL_NOT:
			p.nextToken()
			term2, err := p.parseTerm(qualifier, relation)
			if err != nil {
				return nil, err
			}
			term = &BooleanNode{
				Op:    "not",
				Left:  term,
				Right: term2,
			}
		default:
			return nil, errors.New("expected boolean got " + p.curToken.Value)
		}
	}
	return term, nil
}

func (p *parser) parseTerm(qualifier string, relation int) (Node, error) {
	var word string

	for p.curToken != nil {
		if p.curToken.Type == CQL_LPAREN {
			p.nextToken()
			expr, err := p.parseQuery(qualifier, relation)
			if err != nil {
				return nil, err
			}
			if p.curToken.Type != CQL_RPAREN {
				return nil, errors.New("missing )")
			}
			p.nextToken()
			return expr, nil
		}

		word = p.curToken.Value
		p.nextToken()

		if !p.isBaseRelation(p.curToken) {
			break
		}

		qualifier = word
		relation = p.curToken.Type
		p.nextToken()
	}

	if word == "" {
		return nil, errors.New("missing term")
	}

	return &TermNode{
		Index:    qualifier,
		Relation: relationToStr(relation),
		Value:    word,
	}, nil
}

func (p *parser) isBaseRelation(t *token) bool {
	if t == nil {
		return false
	}
	// unknown first class relation
	if t.Type == CQL_WORD && !strings.Contains(t.Value, ".") {
		return false
	}
	return p.isProxRelation(t) || t.Type == CQL_ALL || t.Type == CQL_ANY || t.Type == CQL_EXACT ||
		t.Type == CQL_SCR || t.Type == CQL_WORD || t.Type == CQL_WITHIN
}

func (p *parser) isProxRelation(t *token) bool {
	if t == nil {
		return false
	}
	return t.Type == CQL_LT || t.Type == CQL_GT || t.Type == CQL_EQ ||
		t.Type == CQL_LE || t.Type == CQL_GE || t.Type == CQL_NE
}
