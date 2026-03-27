# Implementation Gaps — 2026-03-27

This document identifies gaps between wayne's stated goals and its current implementation, with actionable steps to close each gap.

---

## Gap 1: HiDPI Scale Factor Not Applied

- **Stated Goal**: "Theme.Scale is the HiDPI scale factor (1.0 = standard, 2.0 = retina)" (theme.go:55)
- **Current State**: `Theme.Scale` field exists and is initialized to 1.0 in all theme presets (DefaultDark, DefaultLight, HighContrast), but no rendering or layout code references it. Setting `Theme{Scale: 2.0}` has no visible effect.
- **Impact**: On HiDPI displays (Retina Macs, 4K Windows monitors), UI elements render at physical pixel size rather than logical pixel size, making widgets appear too small and text difficult to read. This is a significant accessibility and usability issue for the stated target platforms.
- **Closing the Gap**:
  1. In `resolveTree()` (layout.go:61), after computing pixel dimensions, multiply by `theme.Scale` from the current rendering context
  2. In `Canvas.DrawText()` (render.go:95), multiply the effective font size by scale
  3. In widget `Draw()` methods, apply scale to padding and border calculations
  4. Expose `Canvas.Scale()` method or pass scale through draw context
  5. Update `doc.go` Quick Start example to demonstrate scale usage
  - **Validation**: Set `Theme{Scale: 2.0}` and verify widgets render at 2× size visually; add unit test comparing resolved bounds at scale=1.0 vs scale=2.0

---

## Gap 2: Android CI Tests Build Only

- **Stated Goal**: "Android - Supported but requires manual testing (CI in progress)" (README.md:12)
- **Current State**: `.github/workflows/test-android.yml` exists and runs `gomobile bind -target=android`, verifying the package compiles for Android. However, no tests execute on an Android emulator. The workflow comment acknowledges "Full test execution requires running emulator which adds significant CI time."
- **Impact**: Android-specific bugs (touch input mapping, screen orientation, lifecycle events) could regress without detection. Contributors cannot verify their changes work on Android without significant local setup. The "CI in progress" claim remains partially fulfilled.
- **Closing the Gap**:
  1. Add `reactivecircus/android-emulator-runner@v2` step to test-android.yml
  2. Configure emulator with API 30+ and arm64-v8a system image
  3. After emulator boots, run `gomobile test -target=android ./...`
  4. Cache emulator snapshots to reduce CI time on subsequent runs
  5. Update README.md testing matrix: Android CI → ✅ Automated (tests run on emulator)
  - **Validation**: CI workflow shows test output with pass/fail counts; GitHub badge reflects test status

---

## Gap 3: iOS CI Tests Build Only

- **Stated Goal**: "iOS - Supported but requires manual testing (CI in progress)" (README.md:13)
- **Current State**: `.github/workflows/test-ios.yml` exists and runs `gomobile bind -target=ios`, producing an XCFramework. No tests execute on iOS Simulator. The workflow is analogous to Android CI limitations.
- **Impact**: iOS-specific bugs (Metal rendering differences, touch event timing, safe area insets) could regress without detection. The "CI in progress" claim remains partially fulfilled.
- **Closing the Gap**:
  1. Add `xcrun simctl boot "iPhone 15 Pro"` step after gomobile init
  2. Run `gomobile test -target=ios/simulator ./...`
  3. Use `macos-14` runner (M1/M2 available) for faster simulator performance
  4. Add simulator device availability check and fallback selection
  5. Update README.md testing matrix: iOS CI → ✅ Automated (tests run on simulator)
  - **Validation**: CI workflow shows test output with pass/fail counts; GitHub badge reflects test status

---

## Gap 4: Linux Build Tags Contradict Documentation

- **Stated Goal**: "Wayne explicitly does NOT support Linux or BSD." (doc.go:70-71), "Linux is not a supported platform." (TESTING.md:92-93)
- **Current State**: All source files include `linux` in their build tags: `//go:build windows || darwin || android || ios || linux`. This means `go build` succeeds on Linux, contrary to documentation.
- **Impact**: Developers on Linux may expect the library to work, only to encounter runtime issues (Ebitengine works on Linux but wayne claims not to target it). The mixed signals create confusion about platform support.
- **Closing the Gap**:
  1. Remove `|| linux` from all source file build tags (app.go, concretewidgets.go, dispatcher.go, doc.go, event.go, game.go, layout.go, publicwidget.go, render.go, resource.go, theme.go, widget.go, color.go)
  2. Verify `GOOS=linux go list ./...` returns "no Go files" after change
  3. Add explicit note in README.md: "For Linux, use github.com/opd-ai/wain instead"
  - **Validation**: `GOOS=linux go build ./...` fails with build constraint error; `GOOS=windows go build ./...` succeeds

---

## Gap 5: Widget Interface Duality

- **Stated Goal**: "PublicWidget is the stable public interface for all UI widgets" (publicwidget.go:5-8)
- **Current State**: Two widget interfaces exist:
  - `PublicWidget` (publicwidget.go:9-24) with `HandleEvent(Event) bool`
  - `Widget` (widget.go:5-27) with `HandlePointer(*PointerEvent)`, `HandleKey(*KeyEvent)`, `HandleTouch(*TouchEvent)`
  
  All concrete widgets implement `PublicWidget`. The `Widget` interface in widget.go appears to be an earlier design or alternative API that coexists without clear purpose distinction.
- **Impact**: API consumers may be confused about which interface to implement for custom widgets. The `BaseWidget` type (widget.go:29-109) provides the `Widget` interface pattern, while `BasePublicWidget` (publicwidget.go:100-200) provides the `PublicWidget` pattern.
- **Closing the Gap**:
  Option A (Deprecate Widget):
  1. Add deprecation notice to `Widget` interface in widget.go
  2. Update `BaseWidget` documentation to recommend `BasePublicWidget` instead
  3. Remove or hide `Widget` interface in v2.0
  
  Option B (Differentiate Purpose):
  1. Document `Widget` as low-level internal interface for event granularity
  2. Document `PublicWidget` as high-level public interface for consumers
  3. Add godoc examples showing when to use each
  - **Validation**: `go doc wayne.Widget` and `go doc wayne.PublicWidget` show clear differentiation or deprecation notice

---

## Gap 6: Dead Code Export Documentation

- **Stated Goal**: Comprehensive API surface for wain compatibility (COMPATIBILITY.md)
- **Current State**: go-stats-generator reports 21 unreferenced exported functions/constants (7.5% of codebase). Many are event-related constants like `KeyShiftL`, `KeyControlL`, `ModAlt`, `PointerButtonRight`, `ScrollAxisHorizontal`, etc. These appear unused internally but are exported for consumer applications.
- **Impact**: Static analysis tools flag these as dead code, creating noise in code quality reports. Without documentation, maintainers may incorrectly remove "unused" exports that are intentionally part of the public API.
- **Closing the Gap**:
  1. Audit all 21 flagged items to categorize as:
     - Intentional API export (document with comment)
     - Truly unused (remove or unexport)
  2. For intentional exports, add godoc comment: `// KeyShiftL is exported for consumer applications handling keyboard events.`
  3. Consider grouping related constants under documented sections in event.go
  - **Validation**: `go-stats-generator` dead code count drops to <5%; all remaining exports have documenting comments

---

## Summary

| Gap | Severity | Effort | Priority |
|-----|----------|--------|----------|
| HiDPI Scale Not Applied | High | Medium | P1 |
| Android CI Tests Build Only | Medium | Medium | P2 |
| iOS CI Tests Build Only | Medium | Medium | P2 |
| Linux Build Tags Contradiction | Medium | Low | P2 |
| Widget Interface Duality | Low | Low | P3 |
| Dead Code Export Documentation | Low | Low | P3 |

**Recommendation**: Address P1 (HiDPI) first as it affects user-facing functionality on target platforms. Then address P2 items to fulfill CI claims and resolve documentation inconsistencies. P3 items are code hygiene improvements.

---

## Closed Gaps (Since Previous Audit)

The following gaps from the previous GAPS.md have been addressed:

1. ✅ **Custom Font Loading** — Now implemented via `golang.org/x/image/font/opentype` in resource.go:100-145
2. ✅ **Grid Layout Column Behavior** — Grid.resolveChildren() now implements proper row/column placement (layout.go:460-510)
3. ✅ **Stack Layout Z-Ordering** — Stack.resolveChildren() places children at same position for overlay behavior (layout.go:396-432)
4. ✅ **Example Applications** — examples/hello, examples/form, examples/scrollview created
5. ✅ **Android CI Workflow** — test-android.yml created (though only verifies build)
6. ✅ **iOS CI Workflow** — test-ios.yml created (though only verifies build)
