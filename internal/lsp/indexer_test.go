package lsp

import (
	"testing"
)

var goContent = `
package main

import (
	"fmt"
	"os"
)

type MyStruct struct {
	Field int
}

func (m *MyStruct) Method() {
	fmt.Println("hello")
}

func Function(a int) string {
	return "world"
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(goContent, "main.go")
	}
}

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(goContent, "main.go")
	// Expected: MyStruct (struct), Method (method), Function (function)
	// Actually, wait, let's check parseGoSymbols implementation.
	// It parses "func ..." and "type ...".
	// MyStruct is type struct -> yes.
	// Method is func -> yes.
	// Function is func -> yes.

	if len(symbols) != 3 {
		t.Errorf("expected 3 symbols, got %d", len(symbols))
		for _, s := range symbols {
			t.Logf("Found symbol: %s (%s)", s.Name, s.Kind)
		}
	}

	expected := map[string]string{
		"MyStruct": "struct",
		"Method":   "method",
		"Function": "function",
	}

	for _, s := range symbols {
		kind, ok := expected[s.Name]
		if !ok {
			t.Errorf("unexpected symbol %s", s.Name)
			continue
		}
		if s.Kind != kind {
			t.Errorf("expected kind %s for %s, got %s", kind, s.Name, s.Kind)
		}
	}
}
