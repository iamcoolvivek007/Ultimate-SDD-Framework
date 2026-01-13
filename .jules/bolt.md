## 2024-03-24 - [Optimizing File Walk in LSP]
**Learning:** `filepath.Walk` calls `lstat` for every file, which is slow for large codebases. `filepath.WalkDir` (Go 1.16+) avoids this by using directory entry type info.
**Action:** Use `filepath.WalkDir` for file traversals, especially when filtering by directory or file name before reading metadata.

## 2024-05-23 - [Regex Compilation in Loops]
**Learning:** `regexp.MustCompile` inside function calls (especially loops) is a major bottleneck. Moving them to package-level variables improved indexing speed by ~6x.
**Action:** Always pre-compile regexes at package level or using `sync.Once` if they are used repeatedly.
