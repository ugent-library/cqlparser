SRU CQL query parser
--------------------

Parses a [CQL](https://www.loc.gov/standards/sru/cql/spec.html) query into a tree of nodes that you can use to further translate to another query (e.g. Lucene, ElasticSearch, Solr).

# Usage

## `Parse(query string) (Node, error)`

Returns interface Node that can be either a TermNode or a BooleanNode.
Use type switches to detect the right type.

Example:

```
package main

import (
	"fmt"
	"github.com/ugent-library/cqlparser"
)

func visit(n cqlparser.Node, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Print(" ")
	}
	switch v := n.(type) {
	case *cqlparser.BooleanNode:
		fmt.Printf("%s %s\n", "BOOLEAN", v.Op)
		visit(v.Left, indent+1)
		visit(v.Right, indent+1)
	case *cqlparser.TermNode:
		fmt.Printf("%s\n", v.String())
	}
}

func main() {
	q := "(year >= 2003 and year <= 2005) and (title=Plato)"
	node, _ := cqlparser.Parse(q)
	visit(node, 0)
}
```

# Support

* term code. But no relation modifiers
* boolean node

# References

This code is heavily based on the perl module [CQL::Parser](https://github.com/bricas/cql-parser) the work of [Brian Cassidy](https://github.com/bricas), whose code I have tried to convert partially to golang.