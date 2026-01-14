package lsp

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func createTestProject(t *testing.T) string {
	t.Helper()
	tmpDir, err := os.MkdirTemp("", "lsp-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	files := map[string]string{
		"main.go": `package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello")
}

type User struct {
	Name string
}

func (u *User) Greet() string {
	return "Hello " + u.Name
}
`,
		"script.js": `
import fs from 'fs';

function calculate(a, b) {
	return a + b;
}

class Calculator {
	add(a, b) {
		return a + b;
	}
}

export const version = "1.0";
`,
		"script.py": `
import os
from sys import argv

def process_data(data):
	pass

class Processor:
	def __init__(self):
		pass

	def run(self):
		pass
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	return tmpDir
}

func TestIndexer(t *testing.T) {
	tmpDir := createTestProject(t)
	defer os.RemoveAll(tmpDir)

	indexer := NewIndexer(tmpDir)
	if err := indexer.Index(); err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	// Test Go symbols
	symbols := indexer.GetFileSymbols("main.go")
	foundMain := false
	foundUser := false
	foundGreet := false

	for _, sym := range symbols {
		if sym.Name == "main" && sym.Kind == "function" {
			foundMain = true
		}
		if sym.Name == "User" && sym.Kind == "struct" {
			foundUser = true
		}
		if sym.Name == "Greet" && sym.Kind == "method" {
			foundGreet = true
		}
	}

	if !foundMain {
		t.Error("Did not find main function in main.go")
	}
	if !foundUser {
		t.Error("Did not find User struct in main.go")
	}
	if !foundGreet {
		t.Error("Did not find Greet method in main.go")
	}

	// Test JS symbols
	symbols = indexer.GetFileSymbols("script.js")
	foundCalc := false
	foundClass := false

	for _, sym := range symbols {
		if sym.Name == "calculate" && sym.Kind == "function" {
			foundCalc = true
		}
		if sym.Name == "Calculator" && sym.Kind == "class" {
			foundClass = true
		}
	}

	if !foundCalc {
		t.Error("Did not find calculate function in script.js")
	}
	if !foundClass {
		t.Error("Did not find Calculator class in script.js")
	}

	// Test Python symbols
	symbols = indexer.GetFileSymbols("script.py")
	foundProcess := false
	foundProcessor := false

	for _, sym := range symbols {
		if sym.Name == "process_data" && sym.Kind == "function" {
			foundProcess = true
		}
		if sym.Name == "Processor" && sym.Kind == "class" {
			foundProcessor = true
		}
	}

	if !foundProcess {
		t.Error("Did not find process_data function in script.py")
	}
	if !foundProcessor {
		t.Error("Did not find Processor class in script.py")
	}

	// Test search
	results := indexer.Search("Greet")
	if len(results) == 0 {
		t.Error("Search failed to find Greet")
	}
}

func BenchmarkIndexer(b *testing.B) {
	// Create a temporary directory for benchmarking
	bDir, err := os.MkdirTemp("", "lsp-bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(bDir)

	content := `package main
func f1() {}
func f2() {}
type S1 struct {}
`
	for i := 0; i < 100; i++ {
		// Use fmt.Sprintf for simpler zero padding
		fname := filepath.Join(bDir, fmt.Sprintf("file%03d.go", i))
		if err := os.WriteFile(fname, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		indexer := NewIndexer(bDir)
		if err := indexer.Index(); err != nil {
			b.Fatal(err)
		}
	}
}
