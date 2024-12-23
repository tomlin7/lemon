package evaluator

import "lemon/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			return NULL
		},
	},
}
