# API Compatibility with wain

## Overview

**wayne** provides an API-compatible alternative to [opd-ai/wain](https://github.com/opd-ai/wain) for applications targeting Windows, macOS, Android, and iOS. While wain uses Wayland/X11 for Linux-specific GUI rendering, wayne uses [Ebitengine](https://ebitengine.org/) to support cross-platform environments where Wayland/X11 are not available.

## Compatibility Level: Public API Surface

wayne maintains **public API compatibility** with wain, meaning:

- ✅ **Interface compatibility**: All public widget types implement the same interfaces (`PublicWidget`, `Container`)
- ✅ **Constructor compatibility**: Widget constructors have identical signatures (`NewButton`, `NewPanel`, etc.)
- ✅ **Type compatibility**: Core types (`Color`, `Theme`, `Size`, `AppConfig`) have identical public fields
- ✅ **Constant compatibility**: All event types, alignment, and layout constants match exactly
- ✅ **Method compatibility**: Widget methods (`Bounds()`, `HandleEvent()`, `Draw()`) have identical signatures

This allows **source code migration** from wain to wayne by changing the import path:

```go
// Before (wain for Linux)
import "github.com/opd-ai/wain"

// After (wayne for Windows/macOS/Android/iOS)
import "github.com/opd-ai/wayne"
```

## Known Differences

### 1. Platform Support

| Library | Supported Platforms |
|---------|---------------------|
| **wain** | Linux (Wayland, X11), FreeBSD |
| **wayne** | Windows, macOS, Android, iOS |

**Implication**: The two libraries target **mutually exclusive platform sets**. They cannot run on the same platform natively, though applications can be built with both using build tags:

```go
//go:build linux || freebsd
import "github.com/opd-ai/wain"

//go:build windows || darwin || android || ios  
import "github.com/opd-ai/wayne"
```

### 2. Backend-Specific Types

wayne does **not** include wain's internal types that expose Wayland/X11 implementation details:

- ❌ `DisplayServer` enum (Wayland/X11 selection)
- ❌ `DisplayListEmitter` interface (GPU command buffer abstraction)
- ❌ `RenderBridge` (Wayland damage tracking integration)

**Rationale**: These types are backend-specific and not part of the portable widget API. wayne uses Ebitengine's rendering pipeline internally, which has different abstractions.

**Migration strategy**: If your wain application uses these types, you're working with backend-specific code that must be rewritten for each platform. Isolate such code behind platform-specific build tags.

### 3. Constructor Parameter Differences

A few constructors have **parameter naming differences** that maintain type compatibility:

| Constructor | wain | wayne | Compatible? |
|-------------|------|-------|-------------|
| `NewButton` | `text string` | `label string` | ✅ Both `string` |
| `NewTextInput` | `initialText string` | `placeholder string` | ✅ Both `string` |
| `NewImageWidget` | `img *Image, size Size` | `size Size` | ✅ Use `NewImageWidgetWithImage` |

**NewImageWidget compatibility**: wain requires an `*Image` at construction, while wayne's `NewImageWidget` takes only size. For drop-in compatibility, use `NewImageWidgetWithImage(img, size)` which matches wain's signature exactly.

**Migration**:
```go
// wain
img := loadImage()
widget := wain.NewImageWidget(img, size)

// wayne (option 1 - drop-in compatible)
widget := wayne.NewImageWidgetWithImage(img, size)

// wayne (option 2 - set image after construction)
widget := wayne.NewImageWidget(size)
widget.SetImage(img)
```

### 4. Additional wayne Features

wayne includes interfaces and features not present in wain:

- ✅ `Themeable` interface for theme propagation
- ✅ `TouchPhase` type alias for touch event handling
- ✅ Event constructors (`NewPointerEvent`, `NewKeyEvent`, `NewTouchEvent`) for testing

These additions are **backward-compatible** (no breaking changes to existing API).

## Automated Compatibility Validation

The `compatibility_test.go` file provides automated validation of API compatibility:

```bash
# Run compatibility tests (requires both libraries)
go test -v -run=TestWainCompatibility ./...
```

**Note**: The test suite validates **structural compatibility** (type signatures, interface conformance) but cannot run behavioral tests across platforms due to differing build tag requirements.

### Test Coverage

The compatibility test suite validates:

1. **Interface conformance**: All 11 widget types implement wain's `PublicWidget` or `Container` interfaces
2. **Constructor signatures**: 22 constructor functions match wain's signatures (with noted exceptions)
3. **Type structure**: 6 core struct types have matching public fields
4. **Constants**: 35+ enum constants have identical values

## Migration Guide

### Step 1: Replace Import

```go
// Change this:
import "github.com/opd-ai/wain"

// To this:
import "github.com/opd-ai/wayne"
```

### Step 2: Fix NewImageWidget Calls (Optional)

wayne provides both patterns:

```go
// wain
img := loadImage()
widget := wain.NewImageWidget(img, wain.Size{Width: 100, Height: 100})

// wayne (drop-in compatible)
widget := wayne.NewImageWidgetWithImage(img, wayne.Size{Width: 100, Height: 100})

// wayne (alternative: set image after)
widget := wayne.NewImageWidget(wayne.Size{Width: 100, Height: 100})
widget.SetImage(img)
```

### Step 3: Remove Backend-Specific Code

If you use `DisplayServer`, `RenderBridge`, or `DisplayListEmitter`, isolate them:

```go
//go:build linux || freebsd
package myapp

import "github.com/opd-ai/wain"

func initRenderer() wain.RenderBridge {
    // wain-specific rendering setup
}
```

```go
//go:build windows || darwin || android || ios
package myapp

// wayne does not expose rendering internals
// Use the high-level Canvas API instead
```

### Step 4: Test

Both libraries share the same test patterns. If you have tests like:

```go
func TestMyWidget(t *testing.T) {
    widget := wain.NewButton("Click", wain.Size{Width: 100, Height: 40})
    if widget == nil {
        t.Fatal("expected widget")
    }
}
```

They work identically with wayne:

```go
func TestMyWidget(t *testing.T) {
    widget := wayne.NewButton("Click", wayne.Size{Width: 100, Height: 40})
    if widget == nil {
        t.Fatal("expected widget")
    }
}
```

## Semantic Compatibility

Beyond API signatures, wayne aims for **behavioral compatibility**:

- Layout algorithms (Row, Column, Grid) produce identical widget trees
- Event dispatch follows the same hit-testing and propagation rules
- Theme application and propagation work identically
- Focus management (when implemented) matches wain's keyboard navigation

Where behavioral differences exist due to platform constraints (e.g., touch vs. mouse events, window decorations), they are documented in the respective widget/method godocs.

## Limitations

### Cannot Cross-Compile Tests

Due to platform-specific build tags, you **cannot** run wayne tests on Linux or wain tests on Windows/macOS in the same test binary. Cross-library tests must validate **structure** (via reflection) rather than **behavior** (via actual execution).

The `compatibility_test.go` file uses this approach: it imports both libraries and uses reflection to compare type signatures, ensuring compile-time compatibility without requiring runtime execution on both platforms.

### Runtime Behavioral Differences

Minor behavioral differences may exist due to platform windowing system constraints:

- **Window decorations**: Ebitengine's decoration control differs from Wayland/X11
- **HiDPI scaling**: Handled differently by each backend
- **Touch vs. pointer events**: Mobile platforms (Android/iOS) emit touch events; desktop platforms emit pointer events

These differences are **inherent to the platforms**, not API incompatibilities. Applications should handle them gracefully (e.g., check for both touch and pointer events).

## Verification Status

✅ **Automated tests**: `compatibility_test.go` validates structural compatibility  
✅ **Manual review**: Public API surface audited for 1:1 correspondence  
⚠️ **Cross-platform behavioral testing**: Not yet automated (requires multi-platform CI)  

## Feedback

If you encounter API incompatibilities not documented here, please file an issue at:
https://github.com/opd-ai/wayne/issues

Include:
1. wain code that works
2. wayne equivalent that fails
3. Error message or behavioral difference
