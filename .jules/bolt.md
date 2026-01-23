## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2026-01-23 - [Optimizing Regex Compilation in Hot Paths]
**Learning:** `regexp.MustCompile` is expensive. Calling it inside high-frequency functions (like file parsers) causes massive allocs and CPU burn. Moving to package-level vars yielded 4.8x speedup (28µs -> 5.7µs) and 12x fewer allocs.
**Action:** Always define `regexp.MustCompile` at package level (global vars) or use `sync.Once` for lazy loading, never inside loops or hot functions.
