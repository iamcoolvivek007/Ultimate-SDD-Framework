package lsp

import (
	"testing"
)

var goCode = `
package main

import (
	"fmt"
	"os"
)

type MyStruct struct {
	Name string
}

func (m *MyStruct) Hello() {
	fmt.Println("Hello")
}

func main() {
	m := &MyStruct{Name: "World"}
	m.Hello()
}
`

var jsCode = `
import React from 'react';
const { useState } = require('react');

export class MyComponent extends React.Component {
  render() {
    return <div>Hello</div>;
  }
}

export const myFunc = () => {
  console.log('Hello');
}

function anotherFunc() {
  return true;
}
`

var pythonCode = `
import os
from sys import argv

class MyClass:
    def method(self):
        pass

def main():
    print("Hello")
`

var rustCode = `
use std::io;

struct MyStruct {
    field: i32,
}

impl MyStruct {
    fn new() -> Self {
        MyStruct { field: 0 }
    }
}

fn main() {
    println!("Hello");
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(goCode, "main.go")
	}
}

func BenchmarkParseGoImports(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoImports(goCode)
	}
}

func BenchmarkParseJSSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSSymbols(jsCode, "main.js")
	}
}

func BenchmarkParseJSImports(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSImports(jsCode)
	}
}

func BenchmarkParsePythonSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parsePythonSymbols(pythonCode, "main.py")
	}
}

func BenchmarkParseRustSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseRustSymbols(rustCode, "main.rs")
	}
}

func TestParseGoSymbols(t *testing.T) {
	symbols := parseGoSymbols(goCode, "main.go")
	if len(symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(symbols))
	}
	// Check for MyStruct
	found := false
	for _, s := range symbols {
		if s.Name == "MyStruct" && s.Kind == "struct" {
			found = true
			break
		}
	}
	if !found {
		t.Error("MyStruct not found")
	}
}

func TestParseJSSymbols(t *testing.T) {
	symbols := parseJSSymbols(jsCode, "main.js")
	// MyComponent (class), myFunc (const/func), anotherFunc (function)
	// Note: exportPattern might also capture MyComponent and myFunc, but they are already captured by class/func patterns.
	// The implementation appends symbols.

	// Let's check specifically for the names.
	names := map[string]bool{}
	for _, s := range symbols {
		names[s.Name] = true
	}

	if !names["MyComponent"] {
		t.Error("MyComponent not found")
	}
	if !names["myFunc"] {
		t.Error("myFunc not found")
	}
	if !names["anotherFunc"] {
		t.Error("anotherFunc not found")
	}
}

func TestParsePythonSymbols(t *testing.T) {
	symbols := parsePythonSymbols(pythonCode, "main.py")
	// MyClass, method, main
	if len(symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(symbols))
	}
}

func TestParseRustSymbols(t *testing.T) {
	symbols := parseRustSymbols(rustCode, "main.rs")
	// MyStruct, new (method), main
	// Note: `impl MyStruct` is not a symbol itself but a context for methods.
	// `struct MyStruct` -> symbol
	// `fn new` -> method
	// `fn main` -> function

	count := 0
	for _, s := range symbols {
		if s.Name == "MyStruct" { count++ }
		if s.Name == "new" { count++ }
		if s.Name == "main" { count++ }
	}
	if count != 3 {
		t.Errorf("Expected 3 specific symbols, got %d relevant ones", count)
	}
}
