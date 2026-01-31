package lsp

import "testing"

var goContent = `package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

type MyStruct struct {
	Field string
}

func (m *MyStruct) Method() {
}
`

var jsContent = `
function myFunction() {
	return true;
}

class MyClass {
	constructor() {}
}

export const myConst = () => {};
`

var pyContent = `
def my_func():
    pass

class MyClass:
    def method(self):
        pass
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(goContent, "main.go")
	}
}

func BenchmarkParseJSSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSSymbols(jsContent, "test.js")
	}
}

func BenchmarkParsePythonSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePythonSymbols(pyContent, "test.py")
	}
}
