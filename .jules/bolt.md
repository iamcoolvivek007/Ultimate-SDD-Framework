## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-03-24 - [Regex Recompilation in Hot Paths]
**Learning:** `regexp.MustCompile` inside function calls (like `parseGoSymbols`) causes recompilation on every invocation, leading to massive allocation churn and CPU usage.
**Action:** Always move `regexp.MustCompile` to package-level variables or `var` blocks, especially for functions called in loops or file traversals.
