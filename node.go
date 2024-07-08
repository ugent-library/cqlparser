package cqlparser

import "strings"

type Node interface {
	String() string
}

type TermNode struct {
	Index    string `json:"index"`
	Relation string `json:"relation"`
	Value    string `json:"value"`
}

func (t *TermNode) String() string {
	qVal := t.Value
	if strings.Contains(qVal, " ") {
		qVal = "\"" + qVal + "\""
	} else if qVal == "" {
		qVal = "\"\""
	}
	return t.Index + " " + t.Relation + " " + qVal
}

type BooleanNode struct {
	Op    string `json:"op"`
	Left  Node   `json:"left"`
	Right Node   `json:"right"`
}

func (n *BooleanNode) String() string {
	return "(" + n.Left.String() + ") " + n.Op + " (" + n.Right.String() + ")"
}
