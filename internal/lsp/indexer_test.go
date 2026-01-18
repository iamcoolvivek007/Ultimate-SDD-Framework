package lsp

import (
	"testing"
)

var sampleGoContent = `
package main

import (
	"fmt"
	"os"
)

type MyStruct struct {
	Field1 string
	Field2 int
}

func (m *MyStruct) Method1() {
	fmt.Println("Method1")
}

func Function2(a int) string {
	return "result"
}

func main() {
	fmt.Println("Hello")
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(sampleGoContent, "main.go")
	}
}

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(sampleGoContent, "main.go")
	if len(symbols) != 4 { // MyStruct, Method1, Function2, main (Wait, main is func)
		// Check manual counting:
		// 1. MyStruct (struct)
		// 2. Method1 (method)
		// 3. Function2 (function)
		// 4. main (function)
		t.Errorf("Expected 4 symbols, got %d", len(symbols))
	}
}
