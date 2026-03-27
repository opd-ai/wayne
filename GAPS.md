# Implementation Gaps — 2026-03-26

This document identifies gaps between wayne's stated goals and its current implementation, with actionable steps to close each gap.

---

## Gap 1: Android CI Automation

- **Stated Goal**: "Android - Supported but requires manual testing (CI in progress)" (README.md:12)
- **Current State**: Build tags include `android`; `app_android_test.go` exists with smoke tests. No automated CI workflow runs tests on Android emulator/device.
- **Impact**: Android platform regressions may go undetected. Contributors cannot verify their changes work on Android without manual setup. The "CI in progress" claim has no corresponding workflow file.
- **Closing the Gap**:
  1. Create `.github/workflows/test-android.yml` using `macos-latest` runner with Android emulator
  2. Install gomobile via `go install golang.org/x/mobile/cmd/gomobile@latest`
  3. Configure Android SDK/NDK using `setup-android-sdk` action
  4. Run `gomobile test -target=android/arm64 ./...`
  5. Update README.md testing matrix to show ✅ for Android CI once passing
  - **Validation**: CI badge in README shows green for Android workflow

---

## Gap 2: iOS CI Automation

- **Stated Goal**: "iOS - Supported but requires manual testing (CI in progress)" (README.md:13)
- **Current State**: Build tags include `ios`; `app_ios_test.go` exists with smoke tests. No automated CI workflow runs tests on iOS Simulator.
- **Impact**: iOS platform regressions may go undetected. The "CI in progress" claim has no corresponding workflow file.
- **Closing the Gap**:
  1. Create `.github/workflows/test-ios.yml` using `macos-latest` runner
  2. Install gomobile and initialize iOS toolchain
  3. Run `gomobile test -target=ios/simulator ./...` (requires Xcode)
  4. Update README.md testing matrix to show ✅ for iOS CI once passing
  - **Validation**: CI badge in README shows green for iOS workflow

---

## Gap 3: Custom Font Loading

- **Stated Goal**: "LoadFont loads a font from the specified path at the given size. Supported formats: TrueType (.ttf), OpenType (.otf)" (resource.go:240-241)
- **Current State**: `LoadFont()` accepts path and size parameters but ignores the path entirely. It always returns the embedded `basicfont.Face7x13` regardless of what font file is specified. The comment says "currently not implemented, falls back to embedded font."
- **Impact**: Applications that attempt to use custom fonts will silently get the default font. This is particularly problematic for localization (non-Latin scripts) and branding (custom typefaces).
- **Closing the Gap**:
  1. Add dependency on `golang.org/x/image/font/opentype`
  2. In `LoadFont()`, read the file at `path` using `os.ReadFile()`
  3. Parse with `opentype.Parse(fontData)` and create face with `opentype.NewFace()`
  4. Store the requested size in `Font.size` for scaling
  5. Update documentation to clarify supported font formats
  - **Validation**: Test loads NotoSans.ttf and renders text; verify different from basicfont output

---

## Gap 4: Grid Layout Column Behavior

- **Stated Goal**: "Grid is a fixed-column grid container" (layout.go:382)
- **Current State**: `Grid` embeds `*Panel` and stores a `columns` field, but it does not override `resolveChildren()`. Children are arranged in a single flow direction (vertical or horizontal) rather than in a grid pattern with the specified number of columns.
- **Impact**: Users who create a `Grid` with 3 columns expect children to wrap to the next row after every 3 items. Currently, all children stack vertically regardless of the `columns` setting.
- **Closing the Gap**:
  1. Add `resolveChildren()` method to `Grid` struct
  2. Compute cell width as `contentW / g.columns`
  3. Place child `i` at column `i % g.columns`, row `i / g.columns`
  4. Account for `padding` and `gap` between cells
  - **Validation**: Create grid with 3 columns and 9 children; verify 3×3 arrangement

---

## Gap 5: Stack Layout Z-Ordering

- **Stated Goal**: "Stack is a layering container that places children on top of each other" (layout.go:372)
- **Current State**: `Stack` embeds `*Panel` without overriding `resolveChildren()`. Children are arranged vertically (FlowColumn) instead of being positioned at the same coordinates for true z-order layering.
- **Impact**: Users expecting overlay/layering behavior (e.g., modal dialogs, tooltips, background images with foreground content) get vertical stacking instead.
- **Closing the Gap**:
  1. Add `resolveChildren()` method to `Stack` struct
  2. Position all children at the same (parentX + padding, parentY + padding) coordinates
  3. Give all children the same (contentW, contentH) dimensions
  4. Later children render on top due to Draw() order
  - **Validation**: Create stack with background panel and foreground button; verify button overlays panel

---

## Gap 6: Example Applications

- **Stated Goal**: The README.md Quick Start section shows code but there are no runnable example applications in the repository.
- **Current State**: The `doc.go` package comment contains example code snippets, but no `examples/` directory or standalone example programs exist.
- **Impact**: New users must piece together documentation to create their first application. Lack of examples slows adoption and increases onboarding friction.
- **Closing the Gap**:
  1. Create `examples/hello/main.go` — minimal window with button and label
  2. Create `examples/form/main.go` — TextInput, Labels, Button with validation
  3. Create `examples/scrollview/main.go` — large list with scroll behavior
  4. Add build instructions in `examples/README.md`
  - **Validation**: Each example compiles and runs on Windows and macOS with `go run .`

---

## Gap 7: HiDPI/Scale Factor Support

- **Stated Goal**: Theme.Scale exists "HiDPI scale factor (1.0 = standard, 2.0 = retina)" (theme.go:55)
- **Current State**: `Theme.Scale` is defined but never applied anywhere in rendering. Widget dimensions, font sizes, and padding are not multiplied by the scale factor.
- **Impact**: On HiDPI displays (Retina Macs, 4K Windows), UI elements appear too small. Users must manually adjust widget sizes rather than setting a scale factor.
- **Closing the Gap**:
  1. In `resolveTree()`, multiply pixel dimensions by `theme.Scale`
  2. In `Canvas.DrawText()`, multiply font size by scale
  3. In widget Draw() methods, apply scale to padding and border calculations
  4. Document that Scale affects all pixel-based values
  - **Validation**: Set `Theme{Scale: 2.0}` and verify widgets render at 2× size

---

## Summary

| Gap | Severity | Effort | Priority |
|-----|----------|--------|----------|
| Android CI | Medium | Medium | P1 |
| iOS CI | Medium | Medium | P1 |
| Custom Font Loading | High | Medium | P2 |
| Grid Layout | Medium | Low | P2 |
| Stack Layout | Medium | Low | P2 |
| Example Applications | Low | Low | P3 |
| HiDPI Support | Low | Medium | P3 |

**Recommendation**: Address P1 items (mobile CI) first to fulfill the "CI in progress" claim in the README. Then address P2 items to fix functional gaps in documented features.
