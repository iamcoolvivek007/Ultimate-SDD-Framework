package lsp

import (
	"testing"
)

func TestParseGoSymbols(t *testing.T) {
	content := `
package main

import "fmt"

type User struct {
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (u *User) Greet() string {
	return "Hello " + u.Name
}
`
	symbols := parseGoSymbols(content, "test.go")

	expected := map[string]string{
		"User":    "struct",
		"NewUser": "function",
		"Greet":   "method",
	}

	if len(symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(symbols))
	}

	for _, sym := range symbols {
		kind, ok := expected[sym.Name]
		if !ok {
			t.Errorf("Unexpected symbol: %s", sym.Name)
			continue
		}
		if sym.Kind != kind {
			t.Errorf("Expected symbol %s to be %s, got %s", sym.Name, kind, sym.Kind)
		}
	}
}

func BenchmarkParseGoSymbols(b *testing.B) {
	content := `
package main

import "fmt"

type User struct {
	Name string
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (u *User) Greet() string {
	return "Hello " + u.Name
}
`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseGoSymbols(content, "test.go")
	}
}
