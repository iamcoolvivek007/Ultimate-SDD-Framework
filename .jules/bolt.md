## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2026-01-19 - [Regex Compilation in Hot Paths]
**Learning:** `regexp.MustCompile` is expensive and should not be called inside frequently executed functions or loops. Compiling it once as a package-level variable reduced execution time by ~2.7x and allocations by ~8x.
**Action:** Always define `regexp.MustCompile` patterns as package-level variables or `var` blocks, especially for patterns used in parsing or loops.
