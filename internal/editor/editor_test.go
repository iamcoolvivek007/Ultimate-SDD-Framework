package editor

import (
	"testing"
)

func BenchmarkParseCodeBlocks(b *testing.B) {
	editor := NewEditor(".")
	response := `
Here is some code:

` + "```go:main.go" + `
package main
import "fmt"
func main() {
	fmt.Println("Hello")
}
` + "```" + `

And another file:
` + "```python script.py" + `
print("Hello World")
` + "```" + `
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		editor.ParseCodeBlocks(response)
	}
}

func TestParseCodeBlocks(t *testing.T) {
	editor := NewEditor(".")
	response := `
` + "```go:main.go" + `
package main
` + "```" + `
`
	blocks := editor.ParseCodeBlocks(response)
	if len(blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(blocks))
	}
	if blocks[0].Language != "go" {
		t.Errorf("Expected language go, got %s", blocks[0].Language)
	}
	if blocks[0].Filename != "main.go" {
		t.Errorf("Expected filename main.go, got %s", blocks[0].Filename)
	}
}
