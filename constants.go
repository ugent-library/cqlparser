package cqlparser

const (
	CQL_LT = iota
	CQL_GT
	CQL_EQ
	CQL_LE
	CQL_GE
	CQL_NE
	CQL_AND
	CQL_OR
	CQL_NOT
	CQL_ANY
	CQL_ALL
	CQL_EXACT
	CQL_WITHIN
	CQL_SCR
	CQL_WORD
	CQL_LPAREN
	CQL_RPAREN
	CQL_EOF
)

var cqlKeywords = map[string]int{
	"<":      CQL_LT,
	">":      CQL_GT,
	"=":      CQL_EQ,
	"<=":     CQL_LE,
	">=":     CQL_GE,
	"<>":     CQL_NE,
	"and":    CQL_AND,
	"or":     CQL_OR,
	"not":    CQL_NOT,
	"any":    CQL_ANY,
	"within": CQL_WITHIN,
	"all":    CQL_ALL,
	"exact":  CQL_EXACT,
	"(":      CQL_LPAREN,
	")":      CQL_RPAREN,
	"":       CQL_EOF,
	"scr":    CQL_SCR,
}

var cqlRelations = map[string]int{
	"<":      CQL_LT,
	">":      CQL_GT,
	"=":      CQL_EQ,
	"<=":     CQL_LE,
	">=":     CQL_GE,
	"<>":     CQL_NE,
	"any":    CQL_ANY,
	"within": CQL_WITHIN,
	"all":    CQL_ALL,
	"exact":  CQL_EXACT,
	"scr":    CQL_SCR,
}

func relationToStr(rel int) string {
	for k, v := range cqlRelations {
		if v == rel {
			return k
		}
	}
	return ""
}
