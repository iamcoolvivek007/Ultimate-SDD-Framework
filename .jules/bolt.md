## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-22 - [Regex Compilation in Hot Paths]
**Learning:** Moving `regexp.MustCompile` from inside parsing functions (called per-file) to package-level variables resulted in ~4x speedup and ~90% reduction in allocations for the LSP indexer.
**Action:** Always define `regexp.MustCompile` as package-level variables for any regex used in repeated operations like file parsing.
