# Wayne – Project Completion Plan

This document lists all known issues and the concrete fixes required to bring
the `wayne` library to a shippable state.  Issues are grouped by theme and
ordered roughly from "highest impact / easiest to fix" to "lower impact /
requires more design work".  Each item references the source file and the
review-comment that surfaced it.

---

## 1. Dead code / Wiring gaps

### 1a. Wire `Window.dispatcher` (app.go line 270)
**Problem:** `App.dispatchEvent` tries to use `win.dispatcher`, but
`Window.dispatcher` is never initialised, so hit-testing and focus management
are dead code.  
**Fix:** In `App.NewWindow`, create an `EventDispatcher` for the new window
(e.g. `win.dispatcher = NewEventDispatcher()`).  When the root widget is set
via `Window.SetRoot`, call `win.dispatcher.SetWidgetRoot(...)` so hit-tests
resolve correctly.

### 1b. - [x] `App.SetTheme` has no effect (app.go line 108)
**Problem:** `a.theme` is stored but never read; all widgets fall back to
`DefaultDark()` independently.  
**Fix:** Thread `App.theme` into `resolveTree` / `Draw` calls (e.g. pass a
`*Theme` argument down the call chain, or expose a `CurrentTheme()` method and
have each widget call it when `widget.theme == nil`).

### 1c. - [x] `SetAlign` on `Panel` is a silent no-op (layout.go line 235)
**Problem:** `Panel.align` is stored but never consulted in `resolveChildren`.  
**Fix:** Read `p.align` in `resolveChildren` and adjust the cross-axis position
of each child for `AlignCenter`, `AlignEnd`, and `AlignStretch`.

### 1d. - [x] `SetTheme` only propagates to `*Panel` children (layout.go line 256)
**Problem:** `Button`, `Label`, `TextInput`, etc. will not inherit a theme set
on their parent `Panel`.  
**Fix:** Define a `Themeable` interface (`SetTheme(Theme)`) and use a type
assertion loop over all `PublicWidget` children rather than a concrete type
check.

### 1e. - [x] `EventDispatcher` uses the internal `Widget` interface; no concrete widget implements it (dispatcher.go line 19)
**Problem:** Hit-testing is written against `Widget` (float64 bounds) but all
concrete widgets implement `PublicWidget` (int pixel bounds), so the dispatcher
can never resolve a hit.  
**Fix (option A – preferred):** Change `EventDispatcher.widgetRoot` to
`PublicWidget` and re-implement `hitTest` using `BasePublicWidget.Position()` +
`Bounds()`.  
**Fix (option B):** Make concrete widgets also satisfy the `Widget` interface
by adding the float64-based `Contains/Children/HandlePointer/HandleKey/HandleTouch/SetFocused/IsFocused` methods (possibly as thin wrappers).

---

## 2. Incorrect event routing (events dispatched without bounds-checking)

### 2a. `Panel.HandleEvent` broadcasts to all children (layout.go line 181)
**Problem:** Every pointer / touch event reaches every child widget regardless
of where the cursor is.  Buttons and text inputs can be activated by clicks
anywhere on the panel.  
**Fix:** Before forwarding a `PointerEvent` or `TouchEvent`, check whether
`pe.X()/pe.Y()` falls within the child's bounds (`child.BasePublicWidget.Position()` +
`Bounds()`).  Only send the event to the topmost containing child; once a child
consumes the event (`HandleEvent` returns `true`), stop iterating.

### 2b. ✅ `Button.HandleEvent` has no bounds check (concretewidgets.go line 97)
**Problem:** Any press/release dispatched to the button fires the click, even
if the pointer was elsewhere.  
**Fix:** At the start of press/release handling, verify that
`pe.X()/pe.Y()` is within `b.Position()` + `b.Bounds()`.

### 2c. ✅ `TextInput` can steal focus without a bounds check (concretewidgets.go line 299)
**Problem:** `PointerButtonPress` always sets `focused = true`.  
**Fix:** Only focus the input when `pe.X()/pe.Y()` is within the widget's
bounds.  Conversely, unfocus when a press occurs outside the bounds.

### 2d. ✅ `ScrollView.HandleEvent` scrolls regardless of pointer position (concretewidgets.go line 470)
**Problem:** All scroll views scroll on every `PointerScroll` event.  
**Fix:** Check that the event coordinates fall within `sv.Position()` +
`sv.Bounds()` before adjusting `sv.scrollY`.  Only handle the relevant axis.

---

## 3. Game-loop / input translation bugs

### 3a. `|| true` makes mouse-move condition unconditional (game.go line 76)
**Problem:** The guard `if mx != lastMX || my != lastMY || true` always fires,
defeating the delta optimisation.  
**Fix:** Remove `|| true`.

### 3b. `PointerEnter` / `PointerLeave` are never emitted (game.go line 82)
**Problem:** Button hover state relies on `PointerEnter`/`PointerLeave`; neither
is ever dispatched, so hover visuals never activate.  
**Fix:** Track the widget that was "under the cursor" on the previous frame.
When the hovered widget changes, emit `PointerLeave` for the old widget and
`PointerEnter` for the new one before emitting `PointerMove`.  A simple
last-hovered pointer is enough for now.

### 3c. `TouchUp` events carry (0, 0) coordinates (game.go line 157)
**Problem:** `ebiten.AppendJustReleasedTouchIDs` gives no position; the code
emits `TouchEvent{x:0, y:0}`, breaking hit-testing on release.  
**Fix:** Keep a `map[ebiten.TouchID][2]float64` that records the last-known
position for every active touch, populated during `TouchBegan`/`TouchMoved`.
Use that map to supply coordinates in the `TouchEnded` handler.

---

## 4. Resource & lifecycle correctness

### 4a. `LoadFont`/`LoadImage` after `Close()` doesn't return `ErrNotRunning` (app.go line 198)
**Problem:** `Close()` calls `resources.cleanup()` but does not nil out
`a.resources`, so subsequent calls succeed against a cleaned-up manager.  
**Fix:** Set `a.resources = nil` at the end of `Close()` (or introduce an
explicit `closed atomic.Bool` and check it before resource operations).

---

## 5. Rendering gaps

### 5a. `DrawText` ignores `Font.size` (render.go line 107)
**Problem:** When a `*Font` is passed, the text face is used as-is regardless
of the size stored in the `Font`.  
**Fix (Ebitengine v2 / text/v2):** Use `text.NewGoXFaceWithOptions` or a
`text.Face` wrapper that respects the requested size.  At minimum, document
that size scaling is not yet implemented if a quick fix is not feasible.

### 5b. `ScrollView.Draw` does not clip children to the viewport (concretewidgets.go line 506)
**Problem:** Child widgets can render outside the scroll region.  
**Fix:** Render children to an offscreen `*ebiten.Image` the size of the
viewport and `DrawImage` only the visible slice onto the main canvas.
Alternatively, add `SetClip(x, y, w, h)` / `ClearClip()` to the `Canvas`
interface and implement it with an Ebitengine stencil.

---

## 6. API / configuration correctness

### 6a. `WindowConfig` fields `MinWidth`, `MaxWidth`, `Fullscreen`, `Decorations` are silently ignored (app.go line 167)
**Problem:** Callers who set these fields get no effect and no error.  
**Fix:** In `App.NewWindow` (or in `App.Run`), apply supported fields via
`ebiten.SetWindowSizeLimits`, `ebiten.SetFullscreen`, and
`ebiten.SetWindowDecorated`.  For unsupported combinations, log a warning
(if `a.verbose`) or document the limitation.

### 6b. `go.mod` pins `go 1.24.0` + `toolchain go1.24.13` (go.mod)
**Problem:** Forces all users and CI to that exact toolchain.  
**Fix:** Reduce to the minimum version actually required (Ebitengine v2 targets
Go 1.21+); change the `go` directive to `go 1.21` and remove the `toolchain`
line (or use `toolchain go1.21.0`).

---

## 7. Testing

- Add unit tests for layout resolution (`resolveTree`) covering `FlowRow`,
  `FlowColumn`, `Grid`, padding, and gap.
- Add unit tests for event routing: ensure `Panel.HandleEvent` only forwards to
  children whose bounds contain the event coordinates.
- Add tests for `Color`, `Theme`, and `StyleOverride` merging.
- Run `GOOS=windows go vet ./...` as a CI gate (already passes; keep it green).

---

## 8. Order of execution (suggested)

1. - [x] Fix `|| true` in game.go (trivial, highest signal-to-noise).
2. - [x] Fix bounds checks in `Button`, `TextInput`, `ScrollView` event handlers.
3. - [x] Wire `Window.dispatcher` and adapt `EventDispatcher` to `PublicWidget`.
4. - [x] Fix `Panel.HandleEvent` hit-testing.
5. - [x] Emit `PointerEnter`/`PointerLeave`; track last-touch positions for `TouchUp`.
6. - [x] Fix `App.SetTheme` propagation and `SetAlign` layout.
7. - [x] Fix `SetTheme` child propagation in `Panel`.
8. Fix `DrawText` font size; add `ScrollView` clipping.
9. Apply `WindowConfig` extra fields to Ebitengine.
10. Fix `LoadFont`/`LoadImage` after `Close()`.
11. Downgrade `go.mod` to `go 1.21`.
12. Add tests (layout, event routing, color/theme).
