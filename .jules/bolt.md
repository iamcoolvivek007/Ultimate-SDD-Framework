## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-22 - [Regex Compilation in Loops]
**Learning:** Compiling regexes inside hot loops (or recursive functions) is a major bottleneck. Moving `regexp.MustCompile` to package-level variables improved parsing speed by ~5-6x and reduced allocations by >90%.
**Action:** Always define `regexp.MustCompile` at package level for patterns used repeatedly.
