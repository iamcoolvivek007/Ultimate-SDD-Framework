package lsp

import (
	"testing"
)

var goContent = `
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

func (u *User) String() string {
	return fmt.Sprintf("User: %s", u.Name)
}

func main() {
	u := NewUser("Alice")
	fmt.Println(u)
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(goContent, "main.go")
	}
}

func BenchmarkParseGoImports(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoImports(goContent)
	}
}

var jsContent = `
const React = require('react');

export const Button = ({ onClick, children }) => {
	return <button onClick={onClick}>{children}</button>;
}

export default function App() {
	return <Button onClick={() => console.log('click')}>Click me</Button>;
}

class ErrorBoundary extends React.Component {
	render() {
		return this.props.children;
	}
}
`

func BenchmarkParseJSSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSSymbols(jsContent, "App.js")
	}
}

func BenchmarkParseJSImports(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSImports(jsContent)
	}
}

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(goContent, "main.go")

	expected := []struct{
		Name string
		Kind string
	}{
		{"User", "struct"},
		{"NewUser", "function"},
		{"String", "method"},
		{"main", "function"},
	}

	if len(symbols) != len(expected) {
		t.Errorf("Expected %d symbols, got %d", len(expected), len(symbols))
	}

	for i, sym := range symbols {
		if i < len(expected) {
			if sym.Name != expected[i].Name || sym.Kind != expected[i].Kind {
				t.Errorf("Symbol %d: expected %s (%s), got %s (%s)", i, expected[i].Name, expected[i].Kind, sym.Name, sym.Kind)
			}
		}
	}
}

func TestParseGoImports(t *testing.T) {
	imports := parseGoImports(goContent)
	expected := []string{"fmt", "strings"}

	if len(imports) != len(expected) {
		t.Errorf("Expected %d imports, got %d", len(expected), len(imports))
	}

	for i, imp := range imports {
		if i < len(expected) {
			if imp != expected[i] {
				t.Errorf("Import %d: expected %s, got %s", i, expected[i], imp)
			}
		}
	}
}

func TestParseJSSymbols(t *testing.T) {
	symbols := parseJSSymbols(jsContent, "App.js")

	expected := []struct{
		Name string
		Kind string
	}{
		{"Button", "function"},
		{"App", "function"},
		{"ErrorBoundary", "class"},
	}

	if len(symbols) != len(expected) {
		t.Errorf("Expected %d symbols, got %d", len(expected), len(symbols))
	}

	for i, sym := range symbols {
		if i < len(expected) {
			if sym.Name != expected[i].Name || sym.Kind != expected[i].Kind {
				t.Errorf("Symbol %d: expected %s (%s), got %s (%s)", i, expected[i].Name, expected[i].Kind, sym.Name, sym.Kind)
			}
		}
	}
}

func TestParseJSImports(t *testing.T) {
	imports := parseJSImports(jsContent)
	expected := []string{"react"}

	if len(imports) != len(expected) {
		t.Errorf("Expected %d imports, got %d", len(expected), len(imports))
	}

	for i, imp := range imports {
		if i < len(expected) {
			if imp != expected[i] {
				t.Errorf("Import %d: expected %s, got %s", i, expected[i], imp)
			}
		}
	}
}
