package cqlparser

import "testing"

func BenchmarkParse(b *testing.B) {
	q := `t=a or t=b and t=c`
	for i := 0; i < b.N; i++ {
		_, err := Parse(q)
		if err != nil {
			b.Errorf("unable to create parser: %s", err)
		}
	}
}

func TestParse(t *testing.T) {
	tests := map[string]string{
		"first_name=Nicolas":                      "first_name = Nicolas",
		"first_name=Nicolas and last_name=Franck": "(first_name = Nicolas) and (last_name = Franck)",
		"a=b and b=c and c=d":                     "((a = b) and (b = c)) and (c = d)",
		"a=b and (b=c and c=d)":                   "(a = b) and ((b = c) and (c = d))",
		"dna":                                     "srw.ServerChoice scr dna",
		"title any \"a b c\"":                     "title any \"a b c\"",
		"year >= 2003 and year <= 2005":           "(year >= 2003) and (year <= 2005)",
	}
	for got, expected := range tests {
		node, err := Parse(got)
		if err != nil {
			t.Fatalf("parse error: %s", err)
			continue
		}
		if node.String() != expected {
			t.Fatalf("%s <> %s", node.String(), expected)
		}
	}
}
