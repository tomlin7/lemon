package evaluator

import (
	"lemon/lexer"
	"lemon/object"
	"lemon/parser"
	"testing"
)

func TestAllBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len([1, 2, 3])`, 3},
		{`len("hello")`, 5},
		{`first([1, 2, 3])`, 1},
		{`first("hello")`, "h"},
		{`last([1, 2, 3])`, 3},
		{`last("hello")`, "o"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest("hello")`, "ello"},
		{`push([1, 2], 3)`, []int{1, 2, 3}},
		{`pop([1, 2, 3])`, 3},
		{`clone([1, 2, 3])`, []int{1, 2, 3}},
		{`keys({"a": 1})`, []string{"a"}}, // order of keys is not guaranteed
		{`values({"a": 1, "b": 2})`, []int{1, 2}},
		{`merge([1, 2], [3, 4])`, []int{1, 2, 3, 4}},
		{`merge("hello", " world")`, "hello world"},
		{`merge({"a": 1}, {"b": 2})`, map[string]int{"a": 1, "b": 2}},
		{`merged([1, 2], [3, 4])`, []int{1, 2, 3, 4}},
		{`merged("hello", " world")`, "hello world"},
		{`merged({"a": 1}, {"b": 2})`, map[string]int{"a": 1, "b": 2}},
		{`sort([3, 1, 2])`, []int{1, 2, 3}},
		{`sort(["c", "a", "b"])`, []string{"a", "b", "c"}},
		{`sorted([3, 1, 2])`, []int{1, 2, 3}},
		{`sorted(["c", "a", "b"])`, []string{"a", "b", "c"}},
		{`print("hello")`, nil},
		{`println("hello")`, nil},
		{`int("123")`, 123},
		{`str(123)`, "123"},
		{`bool(1)`, true},
		{`bool(0)`, false},
	}

	for _, tt := range tests {
		evaluated := testBuiltinEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testForIntegerObject(t, evaluated, int64(expected))
		case string:
			testForStringObject(t, evaluated, expected)
		case []int:
			testForIntegerArrayObject(t, evaluated, expected)
		case []string:
			testForStringArrayObject(t, evaluated, expected)
		case bool:
			testForBooleanObject(t, evaluated, expected)
		case nil:
			if evaluated != NULL {
				t.Errorf("object is not NULL. got=%T (%+v)", evaluated, evaluated)
			}
		}
	}
}

func testBuiltinEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testForIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testForStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func testForIntegerArrayObject(t *testing.T, obj object.Object, expected []int) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}
	if len(result.Elements) != len(expected) {
		t.Errorf("array has wrong length. got=%d, want=%d", len(result.Elements), len(expected))
		return false
	}
	for i, expectedElem := range expected {
		if !testForIntegerObject(t, result.Elements[i], int64(expectedElem)) {
			return false
		}
	}
	return true
}

func testForStringArrayObject(t *testing.T, obj object.Object, expected []string) bool {
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("object is not Array. got=%T (%+v)", obj, obj)
		return false
	}
	if len(result.Elements) != len(expected) {
		t.Errorf("array has wrong length. got=%d, want=%d", len(result.Elements), len(expected))
		return false
	}
	for i, expectedElem := range expected {
		if !testForStringObject(t, result.Elements[i], expectedElem) {
			return false
		}
	}
	return true
}

func testForBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}
