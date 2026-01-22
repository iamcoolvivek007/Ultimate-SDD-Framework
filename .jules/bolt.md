## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-22 - [Regex Compilation Overhead]
**Learning:** Defining `regexp.MustCompile` inside frequently called functions (like `parseGoSymbols`) causes massive overhead (recompiling per call). Moving to package-level vars reduced execution time by ~81% (27µs -> 5µs).
**Action:** strictly enforce package-level regex compilation for all parsers and hot paths.
