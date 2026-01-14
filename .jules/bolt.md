## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-22 - [Regex Compilation in Hot Loops]
**Learning:** Compiling regex patterns with `regexp.MustCompile` inside hot loops (like file parsing functions) causes massive performance degradation due to repeated compilation overhead.
**Action:** Always move `regexp.MustCompile` calls to package-level variables or `init()` blocks to ensure they are compiled only once. Benchmarks showed a 60% speedup and 85% memory reduction.
