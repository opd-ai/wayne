//go:build windows || darwin || android || ios

package wayne

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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
	a := g.app

	// --- Mouse ---
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

	// --- Keyboard ---
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

	// --- Touch ---
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
	switch k {
	case ebiten.KeyEscape:
		return KeyEscape
	case ebiten.KeyEnter:
		return KeyReturn
	case ebiten.KeyTab:
		return KeyTab
	case ebiten.KeyBackspace:
		return KeyBackspace
	case ebiten.KeyDelete:
		return KeyDelete
	case ebiten.KeyArrowLeft:
		return KeyLeft
	case ebiten.KeyArrowUp:
		return KeyUp
	case ebiten.KeyArrowRight:
		return KeyRight
	case ebiten.KeyArrowDown:
		return KeyDown
	case ebiten.KeyHome:
		return KeyHome
	case ebiten.KeyEnd:
		return KeyEnd
	case ebiten.KeyPageUp:
		return KeyPageUp
	case ebiten.KeyPageDown:
		return KeyPageDown
	case ebiten.KeySpace:
		return KeySpace
	case ebiten.KeyShiftLeft:
		return KeyShiftL
	case ebiten.KeyShiftRight:
		return KeyShiftR
	case ebiten.KeyControlLeft:
		return KeyControlL
	case ebiten.KeyControlRight:
		return KeyControlR
	case ebiten.KeyAlt:
		return KeyAltL
	case ebiten.KeyAltRight:
		return KeyAltR
	case ebiten.KeyMetaLeft:
		return KeySuperL
	case ebiten.KeyMetaRight:
		return KeySuperR

	// Letter keys (a-z → 0x0061-0x007A).
	case ebiten.KeyA:
		return Key('a')
	case ebiten.KeyB:
		return Key('b')
	case ebiten.KeyC:
		return Key('c')
	case ebiten.KeyD:
		return Key('d')
	case ebiten.KeyE:
		return Key('e')
	case ebiten.KeyF:
		return Key('f')
	case ebiten.KeyG:
		return Key('g')
	case ebiten.KeyH:
		return Key('h')
	case ebiten.KeyI:
		return Key('i')
	case ebiten.KeyJ:
		return Key('j')
	case ebiten.KeyK:
		return Key('k')
	case ebiten.KeyL:
		return Key('l')
	case ebiten.KeyM:
		return Key('m')
	case ebiten.KeyN:
		return Key('n')
	case ebiten.KeyO:
		return Key('o')
	case ebiten.KeyP:
		return Key('p')
	case ebiten.KeyQ:
		return Key('q')
	case ebiten.KeyR:
		return Key('r')
	case ebiten.KeyS:
		return Key('s')
	case ebiten.KeyT:
		return Key('t')
	case ebiten.KeyU:
		return Key('u')
	case ebiten.KeyV:
		return Key('v')
	case ebiten.KeyW:
		return Key('w')
	case ebiten.KeyX:
		return Key('x')
	case ebiten.KeyY:
		return Key('y')
	case ebiten.KeyZ:
		return Key('z')

	// Digit keys (0-9 → 0x0030-0x0039).
	case ebiten.KeyDigit0:
		return Key('0')
	case ebiten.KeyDigit1:
		return Key('1')
	case ebiten.KeyDigit2:
		return Key('2')
	case ebiten.KeyDigit3:
		return Key('3')
	case ebiten.KeyDigit4:
		return Key('4')
	case ebiten.KeyDigit5:
		return Key('5')
	case ebiten.KeyDigit6:
		return Key('6')
	case ebiten.KeyDigit7:
		return Key('7')
	case ebiten.KeyDigit8:
		return Key('8')
	case ebiten.KeyDigit9:
		return Key('9')

	default:
		return Key(k)
	}
}
