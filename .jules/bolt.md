## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2025-05-18 - [Optimizing Regex Compilation]
**Learning:** Compiling regexes inside functions (`regexp.MustCompile`) caused massive allocation churn (129 allocs/op) and slow execution (~30us/op) in the LSP indexer. Moving them to package-level `var` blocks reduced time by ~78% and allocations by ~90%.
**Action:** Always define `regexp.MustCompile` at package level or in a `var` block for any regex used in a loop or frequently called function.
