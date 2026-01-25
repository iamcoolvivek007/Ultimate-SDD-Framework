package lsp

import (
	"testing"
)

var goContent = `
package main

import "fmt"

type MyStruct struct {
	Field string
}

type MyInterface interface {
	Method()
}

func NewMyStruct() *MyStruct {
	return &MyStruct{}
}

func (m *MyStruct) Method() {
	fmt.Println("Method called")
}

func Helper() {
	// do nothing
}
`

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(goContent, "test.go")

	// Expecting: MyStruct (struct), MyInterface (interface), NewMyStruct (function), Method (method), Helper (function)
	expectedCount := 5
	if len(symbols) != expectedCount {
		t.Errorf("Expected %d symbols, got %d", expectedCount, len(symbols))
		for _, s := range symbols {
			t.Logf("Found symbol: %s (%s)", s.Name, s.Kind)
		}
	}

	expected := map[string]string{
		"MyStruct":    "struct",
		"MyInterface": "interface",
		"NewMyStruct": "function",
		"Method":      "method",
		"Helper":      "function",
	}

	found := make(map[string]string)
	for _, s := range symbols {
		found[s.Name] = s.Kind
	}

	for name, kind := range expected {
		if k, ok := found[name]; !ok {
			t.Errorf("Missing symbol: %s", name)
		} else if k != kind {
			t.Errorf("Symbol %s has wrong kind: expected %s, got %s", name, kind, k)
		}
	}
}

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(goContent, "test.go")
	}
}
