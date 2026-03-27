# AUDIT — 2026-03-27

## Project Goals

**wayne** is an API-compatible cross-platform GUI library for Windows, macOS, Android, and iOS using Ebitengine as the rendering backend. It provides source-level migration from [opd-ai/wain](https://github.com/opd-ai/wain) for non-Linux platforms.

### Stated Claims (from README.md, doc.go, COMPATIBILITY.md)

1. **Cross-platform support**: Windows, macOS, Android, iOS
2. **API compatibility with wain**: Interface, constructor, type, constant, and method compatibility
3. **Source-level migration**: Change import path from wain to wayne
4. **Widget system**: Button, Label, TextInput, Panel, Row, Column, Grid, Stack, ScrollView, ImageWidget, Spacer
5. **Theme system**: App-wide themes, per-widget StyleOverride, dark/light/high-contrast presets
6. **Event dispatch**: Pointer, keyboard, touch, window events with hit-testing and focus management
7. **Layout system**: Percentage-based sizing, flow direction (row/column), padding, gap, alignment
8. **Focus/Tab navigation**: Automatic focus chain from widget tree, Tab/Shift+Tab navigation
9. **Custom font loading**: LoadFont supports TrueType (.ttf) and OpenType (.otf)
10. **Image loading**: LoadImage supports PNG, JPEG, GIF
11. **CI coverage**: Windows and macOS automated; Android/iOS CI "in progress"

---

## Goal-Achievement Summary

| Goal | Status | Evidence |
|------|--------|----------|
| Windows support | ✅ Achieved | `.github/workflows/test.yml`:10-31; build tags in all files |
| macOS support | ✅ Achieved | `.github/workflows/test.yml`:33-54; build tags in all files |
| Android support | ⚠️ Partial | `.github/workflows/test-android.yml` runs gomobile bind only; no emulator tests |
| iOS support | ⚠️ Partial | `.github/workflows/test-ios.yml` runs gomobile bind only; no simulator tests |
| API compatibility with wain | ✅ Achieved | `api_alignment_test.go`:20-68 validates 22+ types, 20+ constructors |
| Widget system | ✅ Achieved | `concretewidgets.go`:54-797; all 11 widget types implemented |
| Theme system | ✅ Achieved | `theme.go`:26-177; DefaultDark, DefaultLight, HighContrast, StyleOverride |
| Event dispatch | ✅ Achieved | `dispatcher.go`:8-354; hit-testing, focus management, Tab navigation |
| Layout system | ✅ Achieved | `layout.go`:59-510; resolveTree, Row, Column, Grid, Stack, alignment |
| Focus/Tab navigation | ✅ Achieved | `dispatcher.go`:140-151; FocusManager with FocusNext/FocusPrev |
| Custom font loading | ✅ Achieved | `resource.go`:100-145; uses golang.org/x/image/font/opentype |
| Image loading | ✅ Achieved | `resource.go`:153-222; PNG, JPEG, GIF via standard library decoders |
| HiDPI/Scale support | ❌ Missing | `theme.go`:55 defines Scale but it is never applied in rendering |
| Example applications | ✅ Achieved | `examples/hello/`, `examples/form/`, `examples/scrollview/` |

**Overall: 11/14 goals fully achieved, 2/14 partial (mobile CI), 1/14 missing (HiDPI)**

---

## Findings

### CRITICAL

*No critical findings.*

### HIGH

- [ ] **HiDPI Scale Factor Not Applied** — `theme.go:55`, `render.go`, `layout.go` — Theme.Scale (documented as "HiDPI scale factor") is defined but never used. Widget dimensions, font sizes, and rendering are not scaled by this value. On HiDPI displays (Retina Macs, 4K Windows), UI elements will appear too small. — **Remediation:** In `resolveTree()` (layout.go:61), multiply final pixel dimensions by `theme.Scale`. In `DrawText()` (render.go:95), scale the font size. Validation: `go test -v ./... -run TestTheme` after adding scale multiplication.

- [ ] **Mobile CI Verifies Build Only, Not Tests** — `.github/workflows/test-android.yml:50`, `.github/workflows/test-ios.yml:39` — README claims Android/iOS "CI in progress" but workflows only run `gomobile bind`, not `gomobile test`. Platform-specific regressions could go undetected. — **Remediation:** Add Android emulator step using `reactivecircus/android-emulator-runner@v2` and run `gomobile test -target=android ./...`. For iOS, boot simulator with `xcrun simctl boot` and run `gomobile test -target=ios/simulator ./...`. Validation: CI badges show test results, not just build success.

### MEDIUM

- [ ] **Build Tags Include Linux But README Excludes It** — `app.go:1`, `doc.go:1`, etc. — All source files have `//go:build windows || darwin || android || ios || linux` but README states "Wayne does NOT support Linux". Including Linux in build tags contradicts documentation and may cause confusion. — **Remediation:** Remove `|| linux` from all build tags to match documented platform support. Validation: `GOOS=linux go build ./...` should fail with "no Go files".

- [ ] **Event Constants Potentially Unreferenced** — `event.go:146-169` — Many Key constants (KeyShiftL, KeyControlL, KeyAltL, etc.) are exported but appear unreferenced in internal code. While these are valid public API for consumers, 21 functions are flagged as dead code (7.5% of total). — **Remediation:** Document these as intentionally exported API surface via godoc comments, or verify usage in consumer code. Validation: Review dead code list from `go-stats-generator` and add `// Exported for consumer use` comments where appropriate.

- [ ] **Widget.go Defines Separate Widget Interface** — `widget.go:5-27` — Defines `Widget` interface with `HandlePointer`, `HandleKey`, `HandleTouch` methods, but main widget types implement `PublicWidget` interface with single `HandleEvent(Event)` method. Two competing interfaces for the same purpose. — **Remediation:** Either deprecate `Widget` interface in widget.go or document its specific use case. Consider consolidating to single interface pattern. Validation: `grep -r "Widget interface" *.go` shows only one primary interface.

### LOW

- [ ] **Naming Convention Stuttering** — `app.go:157`, `app.go:178` — Methods `applyWindowSizeLimits` and `applyWindowDecorations` contain redundant "Window" prefix (receiver is already `*App` managing windows). — **Remediation:** Rename to `applySizeLimits` and `applyDecorations`. Validation: `go vet ./...` passes.

- [ ] **Package Name Directory Mismatch** — Go package is `wayne` but repository directory is `wayne` (lowercase). This is actually correct convention but flagged by linter due to directory structure. — **Remediation:** No action required; this is a false positive from static analysis.

- [ ] **Minor Code Duplication** — `examples/form/main.go:52-59,66-73` — 8 lines duplicated for creating label+input rows. — **Remediation:** Extract helper function `createInputRow(label, placeholder string) (*Row, *TextInput)`. Validation: `go-stats-generator` duplication ratio drops below 0.1%.

---

## Metrics Snapshot

| Metric | Value | Assessment |
|--------|-------|------------|
| Lines of Code | 1,747 | Appropriate for scope |
| Total Functions/Methods | 281 | Well-organized |
| Avg Function Length | 7.0 lines | Excellent (<10 is good) |
| Functions >50 lines | 2 (0.7%) | Acceptable |
| Avg Cyclomatic Complexity | 2.8 | Good (<5 is ideal) |
| High Complexity (>10) | 1 function | Acceptable (Draw method) |
| Documentation Coverage | 92.4% | Strong (>80% target) |
| Packages | 2 (wayne, main) | Minimal, appropriate |
| Duplication Ratio | 0.18% | Excellent (<5% target) |
| Dead Code | 21 functions (7.5%) | Review needed |
| Magic Numbers | 207 | Most are color RGB values (acceptable) |

### Top Complex Functions (require no action — complexity justified)

| Function | File | Lines | Complexity | Justification |
|----------|------|-------|------------|---------------|
| Draw (Button) | concretewidgets.go | 44 | 11.9 | Multiple state combinations (enabled/disabled/hovered/pressed/focused) |
| LoadFont | resource.go | 44 | 9.6 | File I/O with error handling and font parsing |
| Draw (TextInput) | concretewidgets.go | 43 | 9.6 | Background, border, text, placeholder, cursor rendering |
| hitTest | dispatcher.go | 23 | 9.3 | Recursive widget tree traversal with bounds checking |

---

## Dependency Status

| Dependency | Version | Status | Notes |
|------------|---------|--------|-------|
| github.com/hajimehoshi/ebiten/v2 | v2.9.9 | ✅ Current | Active development; uses text/v2 (correct) |
| golang.org/x/image | v0.36.0 | ✅ Current | Standard x/image; basicfont for fallback |

**No known vulnerabilities or deprecations affecting this codebase.**

---

## Validation Commands

```bash
# Build verification (all platforms)
GOOS=windows go build ./...
GOOS=darwin go build ./...

# Test suite
CGO_ENABLED=1 go test -race -v ./...

# Static analysis
go vet ./...

# Metrics (requires go-stats-generator)
go-stats-generator analyze . --skip-tests

# API alignment
go test -v -run TestAPIAlignment ./...
```

---

## Recommendations

### Priority 1: Fix HiDPI Support (HIGH)
Theme.Scale is documented but non-functional. This affects all HiDPI displays. Implement scale multiplication in layout resolution and text rendering.

### Priority 2: Enhance Mobile CI (HIGH)
Current workflows verify builds but not behavior. Add emulator/simulator-based test runs to catch platform-specific regressions.

### Priority 3: Resolve Build Tag Inconsistency (MEDIUM)
Remove Linux from build tags to match documented platform exclusion, preventing confusion for developers.

### Priority 4: Clean Up Dead Code Flags (MEDIUM)
Review 21 flagged dead code items. Document intentionally exported API surface or remove truly unused code.
