## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-03-24 - [Optimizing Regexp Compilation in Loops]
**Learning:** Compiling regexes inside functions (e.g., `regexp.MustCompile`) causes severe performance degradation when those functions are called repeatedly (e.g., during indexing). Moving them to package-level `var` blocks yielded a ~3.5x speedup (7.6µs vs 26.6µs per op).
**Action:** Always define `regexp.MustCompile` at the package level or in a singleton, never inside a frequently called function.
