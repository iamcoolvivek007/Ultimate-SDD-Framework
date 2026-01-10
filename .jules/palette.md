## 2025-05-23 - Context-Aware Shortcuts & Truncation
**Learning:** Hardcoded keyboard hints in the main UI create visual noise and inconsistency. Unifying them into a context-aware footer significantly cleans up the interface and improves discoverability.
**Action:** Always check for repeated UI patterns (like "[Key] Action") and centralize them in a dedicated area (footer/header) rather than scattering them in content.

**Learning:** Reading files in `View()` functions is a major anti-pattern in TUI frameworks as it blocks the render loop.
**Action:** Move all data loading to `Init` or `Update` phases and store in the model, keeping `View` pure and fast.
