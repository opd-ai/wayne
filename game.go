//go:build windows || darwin || android || ios

package wayne

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// ebitenToWayneKeyMap maps ebiten.Key constants to wayne.Key constants.
var ebitenToWayneKeyMap = map[ebiten.Key]Key{
	// Special keys
	ebiten.KeyEscape:     KeyEscape,
	ebiten.KeyEnter:      KeyReturn,
	ebiten.KeyTab:        KeyTab,
	ebiten.KeyBackspace:  KeyBackspace,
	ebiten.KeyDelete:     KeyDelete,
	ebiten.KeyArrowLeft:  KeyLeft,
	ebiten.KeyArrowUp:    KeyUp,
	ebiten.KeyArrowRight: KeyRight,
	ebiten.KeyArrowDown:  KeyDown,
	ebiten.KeyHome:       KeyHome,
	ebiten.KeyEnd:        KeyEnd,
	ebiten.KeyPageUp:     KeyPageUp,
	ebiten.KeyPageDown:   KeyPageDown,
	ebiten.KeySpace:      KeySpace,

	// Modifier keys
	ebiten.KeyShiftLeft:    KeyShiftL,
	ebiten.KeyShiftRight:   KeyShiftR,
	ebiten.KeyControlLeft:  KeyControlL,
	ebiten.KeyControlRight: KeyControlR,
	ebiten.KeyAlt:          KeyAltL,
	ebiten.KeyAltRight:     KeyAltR,
	ebiten.KeyMetaLeft:     KeySuperL,
	ebiten.KeyMetaRight:    KeySuperR,

	// Letter keys (a-z → 0x0061-0x007A)
	ebiten.KeyA: Key('a'),
	ebiten.KeyB: Key('b'),
	ebiten.KeyC: Key('c'),
	ebiten.KeyD: Key('d'),
	ebiten.KeyE: Key('e'),
	ebiten.KeyF: Key('f'),
	ebiten.KeyG: Key('g'),
	ebiten.KeyH: Key('h'),
	ebiten.KeyI: Key('i'),
	ebiten.KeyJ: Key('j'),
	ebiten.KeyK: Key('k'),
	ebiten.KeyL: Key('l'),
	ebiten.KeyM: Key('m'),
	ebiten.KeyN: Key('n'),
	ebiten.KeyO: Key('o'),
	ebiten.KeyP: Key('p'),
	ebiten.KeyQ: Key('q'),
	ebiten.KeyR: Key('r'),
	ebiten.KeyS: Key('s'),
	ebiten.KeyT: Key('t'),
	ebiten.KeyU: Key('u'),
	ebiten.KeyV: Key('v'),
	ebiten.KeyW: Key('w'),
	ebiten.KeyX: Key('x'),
	ebiten.KeyY: Key('y'),
	ebiten.KeyZ: Key('z'),

	// Digit keys (0-9 → 0x0030-0x0039)
	ebiten.KeyDigit0: Key('0'),
	ebiten.KeyDigit1: Key('1'),
	ebiten.KeyDigit2: Key('2'),
	ebiten.KeyDigit3: Key('3'),
	ebiten.KeyDigit4: Key('4'),
	ebiten.KeyDigit5: Key('5'),
	ebiten.KeyDigit6: Key('6'),
	ebiten.KeyDigit7: Key('7'),
	ebiten.KeyDigit8: Key('8'),
	ebiten.KeyDigit9: Key('9'),
}

// ebitenGame implements the ebiten.Game interface for wayne's main loop.
type ebitenGame struct {
	app *App
}

// Update processes one tick of the game loop: handles input and dispatches events.
func (g *ebitenGame) Update() error {
	a := g.app

	// Check for quit signal.
	if a.shouldQuit() {
		return ebiten.Termination
	}

	// Drain notify channel.
	a.drainNotify()

	// Process input and dispatch events.
	g.processInput()

	return nil
}

// Draw renders the widget tree to the screen.
func (g *ebitenGame) Draw(screen *ebiten.Image) {
	a := g.app

	root := a.primaryRoot()
	if root == nil {
		return
	}

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	// Resolve layout first.
	childW, childH := computeChildPixelSize(root, sw, sh)
	resolveTree(root, 0, 0, childW, childH)

	// Draw to the screen.
	canvas := newEbitenCanvas(screen)
	root.Draw(canvas)
}

// Layout returns the logical screen dimensions.
func (g *ebitenGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	w, h := g.app.dimensions()
	if w <= 0 {
		w = outsideWidth
	}
	if h <= 0 {
		h = outsideHeight
	}
	return w, h
}

// processInput reads Ebitengine input state each tick and dispatches events to the app.
func (g *ebitenGame) processInput() {
	g.processMouseInput()
	g.processKeyboardInput()
	g.processTouchInput()
}

// processMouseInput handles mouse movement, button presses/releases, and scroll wheel events.
func (g *ebitenGame) processMouseInput() {
	a := g.app
	mx, my := ebiten.CursorPosition()
	mxf, myf := float64(mx), float64(my)

	// Mouse move
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) ||
		ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) ||
		true { // always emit move events
		prev := a.lastMousePos()
		if prev[0] != mx || prev[1] != my {
			a.setLastMousePos(mx, my)
			a.dispatchEvent(NewPointerEvent(PointerMove, mxf, myf, 0, ScrollAxisVertical, 0))
		}
	}

	// Mouse button press/release
	for _, pair := range []struct {
		ebBtn    ebiten.MouseButton
		wayneBtn PointerButton
	}{
		{ebiten.MouseButtonLeft, PointerButtonLeft},
		{ebiten.MouseButtonRight, PointerButtonRight},
		{ebiten.MouseButtonMiddle, PointerButtonMiddle},
	} {
		if inpututil.IsMouseButtonJustPressed(pair.ebBtn) {
			a.dispatchEvent(NewPointerEvent(PointerButtonPress, mxf, myf, pair.wayneBtn, ScrollAxisVertical, 0))
		}
		if inpututil.IsMouseButtonJustReleased(pair.ebBtn) {
			a.dispatchEvent(NewPointerEvent(PointerButtonRelease, mxf, myf, pair.wayneBtn, ScrollAxisVertical, 0))
		}
	}

	// Scroll wheel
	wx, wy := ebiten.Wheel()
	if wx != 0 {
		a.dispatchEvent(NewPointerEvent(PointerScroll, mxf, myf, 0, ScrollAxisHorizontal, wx))
	}
	if wy != 0 {
		a.dispatchEvent(NewPointerEvent(PointerScroll, mxf, myf, 0, ScrollAxisVertical, wy))
	}
}

// processKeyboardInput handles keyboard key presses, releases, and text input.
func (g *ebitenGame) processKeyboardInput() {
	a := g.app

	// Collect text input (printable characters).
	var inputChars []rune
	inputChars = ebiten.AppendInputChars(inputChars)

	// Build a map of rune → already handled, to associate with key events.
	justPressedKeys := inpututil.AppendJustPressedKeys(nil)
	for _, ebKey := range justPressedKeys {
		wKey := ebitenKeyToWayne(ebKey)
		mods := currentModifiers()
		var r rune
		if len(inputChars) > 0 {
			r = inputChars[0]
			inputChars = inputChars[1:]
		}
		a.dispatchEvent(NewKeyEvent(KeyPress, wKey, mods, r))
	}

	// Any remaining input chars (typed without a direct key mapping).
	for _, r := range inputChars {
		a.dispatchEvent(NewKeyEvent(KeyPress, 0, currentModifiers(), r))
	}

	justReleasedKeys := inpututil.AppendJustReleasedKeys(nil)
	for _, ebKey := range justReleasedKeys {
		wKey := ebitenKeyToWayne(ebKey)
		mods := currentModifiers()
		a.dispatchEvent(NewKeyEvent(KeyRelease, wKey, mods, 0))
	}
}

// processTouchInput handles touch events for mobile and touch-enabled devices.
func (g *ebitenGame) processTouchInput() {
	a := g.app

	justPressedTouches := inpututil.AppendJustPressedTouchIDs(nil)
	for _, tid := range justPressedTouches {
		tx, ty := ebiten.TouchPosition(tid)
		a.dispatchEvent(NewTouchEvent(TouchDown, int32(tid), float64(tx), float64(ty)))
	}

	var activeTouches []ebiten.TouchID
	activeTouches = ebiten.AppendTouchIDs(activeTouches)
	for _, tid := range activeTouches {
		tx, ty := ebiten.TouchPosition(tid)
		a.dispatchEvent(NewTouchEvent(TouchMotion, int32(tid), float64(tx), float64(ty)))
	}

	justReleasedTouches := inpututil.AppendJustReleasedTouchIDs(nil)
	for _, tid := range justReleasedTouches {
		a.dispatchEvent(NewTouchEvent(TouchUp, int32(tid), 0, 0))
	}
}

// currentModifiers reads the current keyboard modifier state.
func currentModifiers() Modifier {
	var mods Modifier
	if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
		mods |= ModShift
	}
	if ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight) {
		mods |= ModControl
	}
	if ebiten.IsKeyPressed(ebiten.KeyAlt) || ebiten.IsKeyPressed(ebiten.KeyAltRight) {
		mods |= ModAlt
	}
	if ebiten.IsKeyPressed(ebiten.KeyMetaLeft) || ebiten.IsKeyPressed(ebiten.KeyMetaRight) {
		mods |= ModSuper
	}
	return mods
}

// ebitenKeyToWayne maps an ebiten.Key to a wayne Key (X11 keysym-compatible).
func ebitenKeyToWayne(k ebiten.Key) Key {
	if wayneKey, ok := ebitenToWayneKeyMap[k]; ok {
		return wayneKey
	}
	return Key(k)
}
