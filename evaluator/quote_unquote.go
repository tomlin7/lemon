package evaluator

import (
	"fmt"
	"lemon/ast"
	"lemon/object"
	"lemon/token"
)

func quote(node ast.Node, env *object.Environment) object.Object {
	node = evalUnquoteCalls(node, env)
	return &object.Quote{Node: node}
}

func evalUnquoteCalls(quoted ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(quoted, func(node ast.Node) ast.Node {
		if !isUnquoteCall(node) {
			return node
		}

		call, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}

		if len(call.Arguments) != 1 {
			return node
		}

		unquoted := Eval(call.Arguments[0], env)
		return convertObjectToASTNode(unquoted)
	})
}

func convertObjectToASTNode(obj object.Object) ast.Node {
	switch obj := obj.(type) {
	case *object.Integer:
		t := token.Token{
			Type:    token.INT,
			Literal: fmt.Sprintf("%d", obj.Value),
		}
		return &ast.IntegerLiteral{Token: t, Value: obj.Value}

	case *object.Boolean:
		var t token.Token
		if obj.Value {
			t = token.Token{Type: token.TRUE, Literal: "true"}
		} else {
			t = token.Token{Type: token.FALSE, Literal: "false"}
		}
		return &ast.Boolean{Token: t, Value: obj.Value}

	case *object.Null:
		// TODO: hack, find better way to represent null
		t := token.Token{Type: token.INT, Literal: "null"}
		return &ast.IntegerLiteral{Token: t, Value: 0}

	case *object.String:
		t := token.Token{
			Type:    token.STRING,
			Literal: obj.Value,
		}
		return &ast.StringLiteral{Token: t, Value: obj.Value}

	case *object.Array:
		elements := []ast.Expression{}
		for _, el := range obj.Elements {
			elements = append(elements, convertObjectToASTNode(el).(ast.Expression))
		}
		return &ast.ArrayLiteral{Token: token.Token{Type: token.LBRACKET, Literal: "["}, Elements: elements}

	case *object.Map:
		pairs := make(map[ast.Expression]ast.Expression)
		for _, value := range obj.Pairs {
			key := convertObjectToASTNode(value.Key).(ast.Expression)
			val := convertObjectToASTNode(value.Value).(ast.Expression)
			pairs[key] = val
		}
		return &ast.MapLiteral{Token: token.Token{Type: token.LBRACE, Literal: "{"}, Pairs: pairs}

	case *object.Function:
		return &ast.Identifier{
			Token: token.Token{Type: token.IDENT},
			Value: "function", // replace with the name of the function
		}

	case *object.Builtin:
		return &ast.Identifier{
			Token: token.Token{Type: token.IDENT},
			Value: obj.Value,
		}
	case *object.Quote:
		return obj.Node

	default:
		return nil
	}
}

func isUnquoteCall(node ast.Node) bool {
	callExpression, ok := node.(*ast.CallExpression)
	if !ok {
		return false
	}
	return callExpression.Function.TokenLiteral() == "unquote"
}
