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

func NewUser(name string, age int) *User {
	return &User{Name: name, Age: age}
}

func (u *User) Greet() string {
	return fmt.Sprintf("Hello, %s!", u.Name)
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(sampleGoCode, "main.go")
	}
}

func BenchmarkParseGoImports(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoImports(sampleGoCode)
	}
}

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(sampleGoCode, "main.go")
	if len(symbols) == 0 {
		t.Fatal("Expected symbols, got none")
	}

	foundUser := false
	foundNewUser := false
	foundGreet := false

	for _, s := range symbols {
		if s.Name == "User" && s.Kind == "struct" {
			foundUser = true
		}
		if s.Name == "NewUser" && s.Kind == "function" {
			foundNewUser = true
		}
		if s.Name == "Greet" && s.Kind == "method" {
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
