## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-22 - [Regex Compilation in Hot Loops]
**Learning:** Moving `regexp.MustCompile` from inside function calls (hot loops) to package-level variables reduced execution time by ~77% (4.45x faster) and memory allocations by ~90% in the LSP indexer.
**Action:** Always define `regexp.MustCompile` as package-level variables or in `var` blocks, never inside functions that are called frequently.
