## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-23 - [Regex Compilation in Hot Paths]
**Learning:** Moving `regexp.MustCompile` from function scope (called per file) to package-level variables reduced execution time by ~80% and memory allocations by ~90% in LSP indexing.
**Action:** Always define `regexp.MustCompile` as package-level variables or `var` blocks, never inside functions called frequently.
