# Implementation Plan: HiDPI Scale Factor Support

## Project Context
- **What it does**: API-compatible cross-platform GUI library for Windows, macOS, Android, and iOS using Ebitengine as rendering backend, providing source-level migration from opd-ai/wain for non-Linux platforms.
- **Current goal**: Implement functional HiDPI scaling (Theme.Scale)
- **Estimated Scope**: Medium (5–15 items above threshold)

## Goal-Achievement Status

| Stated Goal | Current Status | This Plan Addresses |
|-------------|---------------|---------------------|
| HiDPI/Scale support | ❌ Missing | **Yes** |
| Android CI tests (not just build) | ⚠️ Partial | No |
| iOS CI tests (not just build) | ⚠️ Partial | No |
| Linux build tags contradict docs | ⚠️ Inconsistent | Yes |
| Widget interface duality cleanup | ⚠️ Confusing | Deferred |
| API compatibility with wain | ✅ Achieved | No |
| Windows support | ✅ Achieved | No |
| macOS support | ✅ Achieved | No |
| Theme system | ✅ Achieved | No |
| Event dispatch | ✅ Achieved | No |
| Layout system | ✅ Achieved | No |
| Documentation coverage | ✅ Achieved (92.4%) | No |

## Metrics Summary

- **Complexity hotspots on goal-critical paths**: 1 function above threshold (Draw in Button: 11.9)
- **Duplication ratio**: 0.18% (excellent — below 5% target)
- **Doc coverage**: 92.4% overall (functions: 92.6%, types: 84.6%)
- **Package coupling**: Low (cohesion 5.17, coupling 3 for main package)
- **Anti-patterns detected**: 2 panics in library code (layout.go:22,25 — validated as intentional API constraints), 23 unused receivers (low severity)

### Functions Requiring Scale-Aware Modifications

| Function | File | Complexity | Role in HiDPI |
|----------|------|------------|---------------|
| `resolveTree` | layout.go | 4.0 | Layout resolution — must apply scale to final pixel dimensions |
| `DrawText` | render.go | 5.6 | Text rendering — must scale font size |
| `Draw` (Button) | concretewidgets.go | 11.9 | Widget rendering — must scale padding/borders |
| `Draw` (TextInput) | concretewidgets.go | 9.6 | Widget rendering — must scale padding/borders |
| `Draw` (Label) | concretewidgets.go | 4.2 | Widget rendering — must scale padding |
| `Draw` (Panel) | layout.go | 3.8 | Container rendering — must scale gaps/padding |

---

## Implementation Steps

### Step 1: Add Scale-Aware Layout Resolution

- **Deliverable**: Modify `resolveTree()` in `layout.go` to multiply final pixel dimensions by `theme.Scale`
- **Dependencies**: None (foundational change)
- **Goal Impact**: Enables all widgets to render at scaled dimensions
- **Files**: `layout.go` (lines 49–95)
- **Acceptance**: Setting `Theme{Scale: 2.0}` doubles resolved widget bounds in unit tests
- **Validation**: 
  ```bash
  GOOS=windows go test -v -run TestScaleLayout ./...
  ```

### Step 2: Add Scale-Aware Text Rendering

- **Deliverable**: Modify `Canvas.DrawText()` in `render.go` to multiply effective font size by scale
- **Dependencies**: Step 1 (layout must provide scaled bounds context)
- **Goal Impact**: Text becomes readable on HiDPI displays
- **Files**: `render.go` (lines 95–135)
- **Acceptance**: Text at scale=2.0 renders at 2× point size; visual test confirms readability
- **Validation**: 
  ```bash
  GOOS=windows go test -v -run TestScaleText ./...
  ```

### Step 3: Add Scale to Canvas Context

- **Deliverable**: Add `Scale() float64` method to `Canvas` interface in `publicwidget.go`, implement in `canvasImpl`
- **Dependencies**: None (interface extension)
- **Goal Impact**: Widgets can query current scale for padding/border calculations
- **Files**: `publicwidget.go` (Canvas interface ~line 45), `render.go` (canvasImpl)
- **Acceptance**: `Canvas.Scale()` returns current theme scale; interface documented
- **Validation**: 
  ```bash
  GOOS=windows go test -v -run TestCanvasScale ./...
  go-stats-generator analyze . --skip-tests --format json --sections documentation | grep -A5 '"coverage"'
  ```

### Step 4: Apply Scale in Widget Draw Methods

- **Deliverable**: Update `Draw()` methods in `concretewidgets.go` to use `canvas.Scale()` for padding/border calculations
- **Dependencies**: Step 3 (Canvas.Scale() must exist)
- **Goal Impact**: Widget visual elements (borders, padding, focus rings) scale proportionally
- **Files**: `concretewidgets.go` (Button.Draw ~line 143, TextInput.Draw ~line 280, Label.Draw ~line 363)
- **Acceptance**: Visual inspection shows scaled borders/padding at scale=2.0
- **Validation**: 
  ```bash
  GOOS=windows go test -v -run TestScaleWidgetDraw ./...
  ```

### Step 5: Update Quick Start Example

- **Deliverable**: Add HiDPI scale demonstration to `doc.go` Quick Start section
- **Dependencies**: Steps 1–4 (feature must be functional)
- **Goal Impact**: Developers learn how to enable HiDPI support
- **Files**: `doc.go` (Quick Start example ~line 25)
- **Acceptance**: Example compiles; godoc shows scale usage
- **Validation**: 
  ```bash
  GOOS=windows go build ./...
  go doc wayne | grep -i scale
  ```

### Step 6: Add Scale Unit Tests

- **Deliverable**: Create `scale_test.go` with tests for layout, text, and widget scaling
- **Dependencies**: Steps 1–4 (implementations must exist)
- **Goal Impact**: Ensures scale feature doesn't regress
- **Files**: `scale_test.go` (new file)
- **Acceptance**: Tests verify bounds at scale=1.0 vs scale=2.0 differ by 2×; tests pass with `-race`
- **Validation**: 
  ```bash
  GOOS=windows CGO_ENABLED=1 go test -v -race -run TestScale ./...
  ```

---

## Deferred Work (Not In This Plan)

These items from GAPS.md/ROADMAP.md are valid but lower priority than HiDPI:

| Item | Reason for Deferral |
|------|---------------------|
| Android/iOS emulator CI | Medium priority; requires significant CI infrastructure; current build verification is adequate |
| Remove Linux from build tags | Medium priority; cosmetic inconsistency; no user-facing bug |
| Widget interface duality | Low priority; code hygiene; no functional impact |
| Dead code documentation | Low priority; affects static analysis noise only |
| Godoc examples | Low priority; 92% coverage is already strong |

---

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Existing tests break from scale changes | Medium | Medium | Add scale=1.0 default; existing tests should be unaffected |
| Performance regression from scale calculations | Low | Low | Scale is a single float multiply; negligible overhead |
| Layout rounding errors at non-integer scales | Medium | Low | Document that non-integer scales may cause minor pixel alignment issues |

---

## Success Criteria

1. `Theme{Scale: 2.0}` visually doubles all UI elements
2. All existing tests pass with `-race` on Windows/macOS
3. `go vet ./...` reports no new issues
4. Documentation coverage remains ≥90%
5. `go-stats-generator` metrics remain stable or improve

---

## Validation Commands

```bash
# Full test suite (on Windows/macOS)
CGO_ENABLED=1 go test -race -v ./...

# Static analysis
go vet ./...

# Metrics verification
go-stats-generator analyze . --skip-tests --format json --sections documentation | grep -A5 '"coverage"'

# Build verification
GOOS=windows go build ./...
GOOS=darwin go build ./...
```

---

## Appendix: Metrics Baseline (2026-03-27)

| Metric | Value |
|--------|-------|
| Total Lines of Code | 1,747 |
| Total Functions | 50 |
| Total Methods | 231 |
| Average Function Complexity | 2.8 |
| High Complexity (>10) | 1 function (Button.Draw: 11.9) |
| Documentation Coverage | 92.4% |
| Duplication Ratio | 0.18% |
| Panic Calls in Library | 2 (intentional — validateSize constraints) |
| Unused Receivers | 23 (low severity — interface conformance) |

Generated by: `go-stats-generator analyze . --skip-tests --format json --sections functions,duplication,documentation,packages,patterns`
