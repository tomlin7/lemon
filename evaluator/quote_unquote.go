package evaluator

import (
	"lemon/ast"
	"lemon/object"
)

func quote(node ast.Node) object.Object {
	return &object.Quote{Node: node}
}
