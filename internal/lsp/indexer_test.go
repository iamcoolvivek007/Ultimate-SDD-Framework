package lsp

import (
	"testing"
)

func TestParseGoSymbols(t *testing.T) {
	content := `
package main

type MyStruct struct {
	Field int
}

func (m *MyStruct) Method() {}

func Function(a int) string {
	return ""
}

type MyInterface interface {
	Do()
}
`
	symbols := parseGoSymbols(content, "main.go")

	expected := []Symbol{
		{Name: "MyStruct", Kind: "struct", File: "main.go", Line: 4},
		{Name: "Method", Kind: "method", File: "main.go", Line: 8, Signature: "func (m *MyStruct) Method() {}", Parent: "MyStruct"},
		{Name: "Function", Kind: "function", File: "main.go", Line: 10, Signature: "func Function(a int) string"},
		{Name: "MyInterface", Kind: "interface", File: "main.go", Line: 14},
	}

	if len(symbols) != len(expected) {
		t.Fatalf("Expected %d symbols, got %d", len(expected), len(symbols))
	}

	for i, sym := range symbols {
		exp := expected[i]
		if sym.Name != exp.Name {
			t.Errorf("Symbol %d: expected name %s, got %s", i, exp.Name, sym.Name)
		}
		if sym.Kind != exp.Kind {
			t.Errorf("Symbol %d: expected kind %s, got %s", i, exp.Kind, sym.Kind)
		}
		if sym.Line != exp.Line {
			t.Errorf("Symbol %d: expected line %d, got %d", i, exp.Line, sym.Line)
		}
		if sym.Parent != exp.Parent {
			t.Errorf("Symbol %d: expected parent %s, got %s", i, exp.Parent, sym.Parent)
		}
	}
}

func BenchmarkParseGoSymbols(b *testing.B) {
	content := `
package main

import "fmt"

type User struct {
	ID   int
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (u *User) GetName() string {
	return u.Name
}

type Repository interface {
	Save(u *User) error
}
`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseGoSymbols(content, "benchmark.go")
	}
}
