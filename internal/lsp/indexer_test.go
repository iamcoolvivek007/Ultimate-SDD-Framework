package lsp

import (
	"testing"
)

var sampleGoCode = `
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

func (u *User) String() string {
	return fmt.Sprintf("%s (%d)", u.Name, u.Age)
}
`

func BenchmarkParseGoSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseGoSymbols(sampleGoCode, "main.go")
	}
}

var sampleJSCode = `
import { useState } from 'react';

export const Button = ({ onClick, children }) => {
  const [count, setCount] = useState(0);

  function handleClick() {
    setCount(count + 1);
    onClick();
  }

  return <button onClick={handleClick}>{children}</button>;
}

class Counter {
  constructor() {
    this.value = 0;
  }
}
`

func BenchmarkParseJSSymbols(b *testing.B) {
	for i := 0; i < b.N; i++ {
		parseJSSymbols(sampleJSCode, "button.js")
	}
}
