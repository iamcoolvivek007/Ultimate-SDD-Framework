package lsp

import (
	"testing"
)

func TestParseGoSymbols(t *testing.T) {
	content := `
package main

import "fmt"

func main() {
	fmt.Println("Hello")
}

func (s *Service) Process(data string) error {
	return nil
}

type Service struct {
	Name string
}
`
	file := "main.go"

	symbols := parseGoSymbols(content, file)

	if len(symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(symbols))
	}

	foundMain := false
	foundProcess := false
	foundService := false

	for _, s := range symbols {
		if s.Name == "main" && s.Kind == "function" {
			foundMain = true
		}
		if s.Name == "Process" && s.Kind == "method" && s.Parent == "Service" {
			foundProcess = true
		}
		if s.Name == "Service" && s.Kind == "struct" {
			foundService = true
		}
	}

	if !foundMain {
		t.Error("Did not find main function")
	}
	if !foundProcess {
		t.Error("Did not find Process method")
	}
	if !foundService {
		t.Error("Did not find Service struct")
	}
}

func TestParseJSSymbols(t *testing.T) {
	content := `
import fs from 'fs';

export const config = () => {
	port: 3000
};

export function startServer() {
	console.log('Starting...');
}

class User {
	constructor(name) {
		this.name = name;
	}
}
`
	file := "index.js"

	symbols := parseJSSymbols(content, file)

	// Expected symbols: config (const function), startServer (function), User (class)

	if len(symbols) != 3 {
		t.Errorf("Expected 3 symbols, got %d", len(symbols))
	}

	foundConfig := false
	foundStartServer := false
	foundUser := false

	for _, s := range symbols {
		if s.Name == "config" {
			foundConfig = true
		}
		if s.Name == "startServer" && s.Kind == "function" {
			foundStartServer = true
		}
		if s.Name == "User" && s.Kind == "class" {
			foundUser = true
		}
	}

	if !foundConfig {
		t.Error("Did not find config const")
	}
	if !foundStartServer {
		t.Error("Did not find startServer function")
	}
	if !foundUser {
		t.Error("Did not find User class")
	}
}
