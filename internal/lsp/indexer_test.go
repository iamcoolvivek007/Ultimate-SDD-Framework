package lsp

import (
	"testing"
)

const sampleGoCode = `
package main

import (
	"fmt"
	"strings"
)

type User struct {
	Name string
	Age  int
}

func NewUser(name string) *User {
	return &User{Name: name}
}

func (u *User) Greet() string {
	return fmt.Sprintf("Hello, %s", u.Name)
}
`

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(sampleGoCode, "main.go")

	expectedCount := 3 // User struct, NewUser func, Greet method
	if len(symbols) != expectedCount {
		t.Errorf("Expected %d symbols, got %d", expectedCount, len(symbols))
	}

	foundUser := false
	foundNewUser := false
	foundGreet := false

	for _, sym := range symbols {
		if sym.Name == "User" && sym.Kind == "struct" {
			foundUser = true
		}
		if sym.Name == "NewUser" && sym.Kind == "function" {
			foundNewUser = true
		}
		if sym.Name == "Greet" && sym.Kind == "method" && sym.Parent == "User" {
			foundGreet = true
		}
	}

	if !foundUser {
		t.Error("Did not find User struct")
	}
	if !foundNewUser {
		t.Error("Did not find NewUser function")
	}
	if !foundGreet {
		t.Error("Did not find Greet method")
	}
}

func BenchmarkParseGoSymbols(b *testing.B) {
	// Pre-load content to avoid IO overhead in benchmark loop, though here it's a constant string
	content := sampleGoCode
	file := "main.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parseGoSymbols(content, file)
	}
}
