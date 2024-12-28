package evaluator

import (
	"fmt"
	"lemon/object"
	"sort"
	"strconv"
)

var builtins = map[string]*object.Builtin{
	// Iterables/Sequences

	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return arg.Elements[0]
			case *object.String:
				return &object.String{Value: string(arg.Value[0])}
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				return arg.Elements[length-1]
			case *object.String:
				length := len(arg.Value)
				return &object.String{Value: string(arg.Value[length-1])}
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				if length > 0 {
					newElements := make([]object.Object, length-1)
					copy(newElements, arg.Elements[1:length])
					return &object.Array{Elements: newElements}
				}
				return NULL
			case *object.String:
				length := len(arg.Value)
				if length > 0 {
					return &object.String{Value: string(arg.Value[1:length])}
				}
				return NULL
			default:
				return newError("argument to `last` not supported, got %s", args[0].Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			arr.Elements = append(arr.Elements, args[1])
			return arr
		},
	},
	"pop": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `pop` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				popped := arr.Elements[length-1]
				arr.Elements = arr.Elements[:length-1]
				return popped
			}

			return NULL
		},
	},
	"clone": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				newElements := make([]object.Object, len(arg.Elements))
				copy(newElements, arg.Elements)
				return &object.Array{Elements: newElements}
			case *object.String:
				return &object.String{Value: arg.Value}
			case *object.Map:
				newPairs := make(map[object.HashKey]object.MapPair)
				for key, value := range arg.Pairs {
					newPairs[key] = value
				}
				return &object.Map{Pairs: newPairs}
			default:
				return newError("argument to `clone` not supported, got %s", args[0].Type())
			}
		},
	},
	"keys": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.MAP_OBJ {
				return newError("argument to `keys` must be MAP, got %s",
					args[0].Type())
			}

			m := args[0].(*object.Map)
			keys := make([]object.Object, len(m.Pairs))
			i := 0
			for _, v := range m.Pairs {
				keys[i] = v.Key
				i++
			}

			return &object.Array{Elements: keys}
		},
	},
	"values": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.MAP_OBJ {
				return newError("argument to `values` must be MAP, got %s",
					args[0].Type())
			}

			m := args[0].(*object.Map)
			values := make([]object.Object, len(m.Pairs))
			i := 0
			// order of keys is not guaranteed
			for _, v := range m.Pairs {
				values[i] = v.Value
				i++
			}

			return &object.Array{Elements: values}
		},
	},
	"merge": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1+", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				for _, a := range args[1:] {
					if a.Type() != object.ARRAY_OBJ {
						return newError("all arguments to `merge` must be of same type, got %s",
							a.Type())
					}
					arr := a.(*object.Array)
					arg.Elements = append(arg.Elements, arr.Elements...)
				}
				return arg
			case *object.String:
				for _, a := range args[1:] {
					if a.Type() != object.STRING_OBJ {
						return newError("all arguments to `merge` must be of same type, got %s",
							a.Type())
					}
					str := a.(*object.String)
					arg.Value += str.Value
				}
				return arg
			case *object.Map:
				for _, a := range args[1:] {
					if a.Type() != object.MAP_OBJ {
						return newError("all arguments to `merge` must be of same type, got %s",
							a.Type())
					}
					m := a.(*object.Map)
					for key, value := range m.Pairs {
						arg.Pairs[key] = value
					}
				}
				return arg
			default:
				return newError("argument to `merge` not supported, got %s", args[0].Type())
			}
		},
	},
	"merged": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) < 1 {
				return newError("wrong number of arguments. got=%d, want=1+", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				newElements := make([]object.Object, len(arg.Elements))
				copy(newElements, arg.Elements)
				for _, a := range args[1:] {
					if a.Type() != object.ARRAY_OBJ {
						return newError("all arguments to `merged` must be of same type, got %s",
							a.Type())
					}
					arr := a.(*object.Array)
					newElements = append(newElements, arr.Elements...)
				}
				return &object.Array{Elements: newElements}
			case *object.String:
				newValue := arg.Value
				for _, a := range args[1:] {
					if a.Type() != object.STRING_OBJ {
						return newError("all arguments to `merged` must be of same type, got %s",
							a.Type())
					}
					str := a.(*object.String)
					newValue += str.Value
				}
				return &object.String{Value: newValue}
			case *object.Map:
				newPairs := make(map[object.HashKey]object.MapPair)
				for key, value := range arg.Pairs {
					newPairs[key] = value
				}
				for _, a := range args[1:] {
					if a.Type() != object.MAP_OBJ {
						return newError("all arguments to `merged` must be of same type, got %s",
							a.Type())
					}
					m := a.(*object.Map)
					for key, value := range m.Pairs {
						newPairs[key] = value
					}
				}
				return &object.Map{Pairs: newPairs}

			default:
				return newError("argument to `merged` not supported, got %s", args[0].Type())
			}
		},
	},
	"sort": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `sort` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) == 0 {
				return arr
			}

			switch arr.Elements[0].(type) {
			case *object.Integer:
				ints := make([]int, len(arr.Elements))
				for i, e := range arr.Elements {
					if intObj, ok := e.(*object.Integer); ok {
						ints[i] = int(intObj.Value)
					} else {
						return newError("all elements in array must be INTEGER, got %s", e.Type())
					}
				}
				sort.Ints(ints)
				for i, e := range ints {
					arr.Elements[i] = &object.Integer{Value: int64(e)}
				}
			case *object.String:
				strs := make([]string, len(arr.Elements))
				for i, e := range arr.Elements {
					if strObj, ok := e.(*object.String); ok {
						strs[i] = strObj.Value
					} else {
						return newError("all elements in array must be STRING, got %s", e.Type())
					}
				}
				sort.Strings(strs)
				for i, e := range strs {
					arr.Elements[i] = &object.String{Value: e}
				}
			default:
				return newError("argument to `sort` not supported, got %s", arr.Elements[0].Type())
			}

			return arr
		},
	},
	"sorted": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `sorted` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) == 0 {
				return arr
			}

			newElements := make([]object.Object, len(arr.Elements))
			copy(newElements, arr.Elements)

			switch newElements[0].(type) {
			case *object.Integer:
				ints := make([]int, len(newElements))
				for i, e := range newElements {
					if intObj, ok := e.(*object.Integer); ok {
						ints[i] = int(intObj.Value)
					} else {
						return newError("all elements in array must be INTEGER, got %s", e.Type())
					}
				}
				sort.Ints(ints)
				for i, e := range ints {
					newElements[i] = &object.Integer{Value: int64(e)}
				}
			case *object.String:
				strs := make([]string, len(newElements))
				for i, e := range newElements {
					if strObj, ok := e.(*object.String); ok {
						strs[i] = strObj.Value
					} else {
						return newError("all elements in array must be STRING, got %s", e.Type())
					}
				}
				sort.Strings(strs)
				for i, e := range strs {
					newElements[i] = &object.String{Value: e}
				}
			default:
				return newError("argument to `sorted` not supported, got %s", newElements[0].Type())
			}

			return &object.Array{Elements: newElements}
		},
	},

	// I/O functions

	"print": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Printf("%s ", arg.Inspect())
			}

			return NULL
		},
	},
	"println": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Printf("%s ", arg.Inspect())
			}
			fmt.Println()

			return NULL
		},
	},
	"input": {
		Fn: func(args ...object.Object) object.Object {
			var input string
			fmt.Scanln(&input)
			return &object.String{Value: input}
		},
	},

	// Type conversion functions
	"int": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				value, err := strconv.ParseInt(arg.Value, 0, 64)
				if err != nil {
					return newError("could not parse %q as integer", arg.Value)
				}
				return &object.Integer{Value: value}
			case *object.Integer:
				return arg
			case *object.Boolean:
				if arg.Value {
					return &object.Integer{Value: 1}
				} else {
					return &object.Integer{Value: 0}
				}
			default:
				return newError("argument to `int` not supported, got %s", args[0].Type())
			}
		},
	},
	"str": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			return &object.String{Value: args[0].Inspect()}
		},
	},
	"bool": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				if arg.Value == 0 {
					return FALSE
				}
				return TRUE
			case *object.Boolean:
				return arg
			case *object.String:
				if arg.Value == "" {
					return FALSE
				}
				return TRUE
			case *object.Array:
				if len(arg.Elements) == 0 {
					return FALSE
				}
				return TRUE
			case *object.Map:
				if len(arg.Pairs) == 0 {
					return FALSE
				}
				return TRUE
			default:
				if arg == NULL {
					return FALSE
				}
				return TRUE
			}
		},
	},
}
