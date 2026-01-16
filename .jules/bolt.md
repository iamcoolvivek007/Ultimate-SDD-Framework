## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2025-05-18 - [Regex Compilation in Loops]
**Learning:** Calling `regexp.MustCompile` inside a function that is called frequently (e.g., for every file in a loop) causes severe performance degradation due to repeated compilation.
**Action:** Always define `regexp.MustCompile` calls as package-level variables or inside a `var` block for constant patterns.
