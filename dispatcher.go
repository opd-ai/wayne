//go:build windows || darwin || android || ios || linux

package wayne

import "sync"

// EventDispatcher manages event routing from platform sources to widget handlers.
type EventDispatcher struct {
	mu sync.RWMutex

	pointerHandlers []func(*PointerEvent)
	keyHandlers     []func(*KeyEvent)
	touchHandlers   []func(*TouchEvent)
	windowHandlers  []func(*WindowEvent)
	customHandlers  []func(*CustomEvent)

	focusManager *FocusManager
	widgetRoot   PublicWidget
}

// NewEventDispatcher creates a new event dispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		focusManager: NewFocusManager(),
	}
}

// SetWidgetRoot sets the root widget for hit-testing.
func (d *EventDispatcher) SetWidgetRoot(root PublicWidget) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.widgetRoot = root
	// Auto-build focus chain from widget tree
	d.focusManager.SetChainFromRoot(root)
}

// SetFocus sets focus to a specific widget.
func (d *EventDispatcher) SetFocus(w PublicWidget) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.focusManager.Focus(w)
}

// FocusedWidget returns the currently focused widget.
func (d *EventDispatcher) FocusedWidget() PublicWidget {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.focusManager.Focused()
}

// OnPointer registers a pointer event handler.
func (d *EventDispatcher) OnPointer(handler func(*PointerEvent)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pointerHandlers = append(d.pointerHandlers, handler)
}

// OnKey registers a keyboard event handler.
func (d *EventDispatcher) OnKey(handler func(*KeyEvent)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.keyHandlers = append(d.keyHandlers, handler)
}

// OnTouch registers a touch event handler.
func (d *EventDispatcher) OnTouch(handler func(*TouchEvent)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.touchHandlers = append(d.touchHandlers, handler)
}

// OnWindow registers a window event handler.
func (d *EventDispatcher) OnWindow(handler func(*WindowEvent)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.windowHandlers = append(d.windowHandlers, handler)
}

// OnCustom registers a custom event handler.
func (d *EventDispatcher) OnCustom(handler func(*CustomEvent)) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.customHandlers = append(d.customHandlers, handler)
}

// Dispatch routes an event to appropriate handlers.
func (d *EventDispatcher) Dispatch(evt Event) {
	if evt == nil {
		return
	}
	switch e := evt.(type) {
	case *PointerEvent:
		d.dispatchPointer(e)
	case *KeyEvent:
		d.dispatchKey(e)
	case *TouchEvent:
		d.dispatchTouch(e)
	case *WindowEvent:
		d.dispatchWindow(e)
	case *CustomEvent:
		d.dispatchCustom(e)
	}
}

func (d *EventDispatcher) dispatchPointer(evt *PointerEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.widgetRoot != nil {
		if target := d.hitTest(d.widgetRoot, int(evt.x), int(evt.y)); target != nil {
			target.HandleEvent(evt)
			if evt.Consumed() {
				return
			}
		}
	}

	for _, h := range d.pointerHandlers {
		if evt.Consumed() {
			return
		}
		h(evt)
	}
}

func (d *EventDispatcher) dispatchKey(evt *KeyEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.handleTabNavigation(evt) {
		return
	}
	if d.dispatchToFocusedWidget(evt) {
		return
	}
	d.dispatchToKeyHandlers(evt)
}

// handleTabNavigation processes Tab key for focus management.
func (d *EventDispatcher) handleTabNavigation(evt *KeyEvent) bool {
	if evt.EventType() != KeyPress || evt.Key() != KeyTab {
		return false
	}
	if evt.Modifiers()&ModShift != 0 {
		d.focusManager.FocusPrev()
	} else {
		d.focusManager.FocusNext()
	}
	evt.Consume()
	return true
}

// dispatchToFocusedWidget sends the key event to the currently focused widget.
func (d *EventDispatcher) dispatchToFocusedWidget(evt *KeyEvent) bool {
	focused := d.focusManager.Focused()
	if focused == nil {
		return false
	}
	if pw, ok := focused.(PublicWidget); ok {
		return pw.HandleEvent(evt)
	}
	return false
}

// dispatchToKeyHandlers sends the event to registered key handlers.
func (d *EventDispatcher) dispatchToKeyHandlers(evt *KeyEvent) {
	for _, h := range d.keyHandlers {
		if evt.Consumed() {
			return
		}
		h(evt)
	}
}

func (d *EventDispatcher) dispatchTouch(evt *TouchEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.widgetRoot != nil {
		if target := d.hitTest(d.widgetRoot, int(evt.x), int(evt.y)); target != nil {
			target.HandleEvent(evt)
			if evt.Consumed() {
				return
			}
		}
	}

	for _, h := range d.touchHandlers {
		if evt.Consumed() {
			return
		}
		h(evt)
	}
}

func (d *EventDispatcher) dispatchWindow(evt *WindowEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.windowHandlers {
		if evt.Consumed() {
			return
		}
		h(evt)
	}
}

func (d *EventDispatcher) dispatchCustom(evt *CustomEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.customHandlers {
		if evt.Consumed() {
			return
		}
		h(evt)
	}
}

func (d *EventDispatcher) hitTest(w PublicWidget, x, y int) PublicWidget {
	// Get widget position and bounds
	var px, py, width, height int
	if positioner, ok := w.(interface{ Position() (int, int) }); ok {
		px, py = positioner.Position()
	}
	width, height = w.Bounds()

	// Check if point is within widget bounds
	if x < px || x >= px+width || y < py || y >= py+height {
		return nil
	}

	// Check container children (bottom-up, so last child is checked first)
	if container, ok := w.(Container); ok {
		children := container.Children()
		for i := len(children) - 1; i >= 0; i-- {
			if child := d.hitTest(children[i], x, y); child != nil {
				return child
			}
		}
	}

	return w
}

// FocusManager manages keyboard focus and tab order.
type FocusManager struct {
	mu         sync.RWMutex
	chain      []PublicWidget
	focusedIdx int
}

// NewFocusManager creates a new focus manager.
func NewFocusManager() *FocusManager {
	return &FocusManager{focusedIdx: -1}
}

// SetChain sets the focus chain (tab order).
func (fm *FocusManager) SetChain(widgets []PublicWidget) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.chain = widgets
	fm.focusedIdx = -1
}

// SetChainFromRoot builds a focus chain from a widget tree.
// It collects all Focusable widgets in depth-first order.
func (fm *FocusManager) SetChainFromRoot(root PublicWidget) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	fm.chain = nil
	fm.focusedIdx = -1
	fm.collectFocusable(root)
}

// collectFocusable recursively collects focusable widgets.
func (fm *FocusManager) collectFocusable(w PublicWidget) {
	if focusable, ok := w.(Focusable); ok && focusable.CanTakeFocus() {
		fm.chain = append(fm.chain, w)
	}
	if container, ok := w.(Container); ok {
		for _, child := range container.Children() {
			fm.collectFocusable(child)
		}
	}
}

// Focus sets focus to a specific widget.
func (fm *FocusManager) Focus(w PublicWidget) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	// Clear previous focus
	if fm.focusedIdx >= 0 && fm.focusedIdx < len(fm.chain) {
		fm.chain[fm.focusedIdx].SetFocused(false)
	}
	// Find and focus the target widget
	for i, widget := range fm.chain {
		if widget == w {
			fm.focusedIdx = i
			widget.SetFocused(true)
			return
		}
	}
}

// FocusNext moves focus to the next widget in the chain.
func (fm *FocusManager) FocusNext() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	if len(fm.chain) == 0 {
		return
	}
	if fm.focusedIdx >= 0 && fm.focusedIdx < len(fm.chain) {
		fm.chain[fm.focusedIdx].SetFocused(false)
	}
	fm.focusedIdx = (fm.focusedIdx + 1) % len(fm.chain)
	fm.chain[fm.focusedIdx].SetFocused(true)
}

// FocusPrev moves focus to the previous widget in the chain.
func (fm *FocusManager) FocusPrev() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	if len(fm.chain) == 0 {
		return
	}
	if fm.focusedIdx >= 0 && fm.focusedIdx < len(fm.chain) {
		fm.chain[fm.focusedIdx].SetFocused(false)
	}
	fm.focusedIdx--
	if fm.focusedIdx < 0 {
		fm.focusedIdx = len(fm.chain) - 1
	}
	fm.chain[fm.focusedIdx].SetFocused(true)
}

// Focused returns the currently focused widget.
func (fm *FocusManager) Focused() PublicWidget {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	if fm.focusedIdx >= 0 && fm.focusedIdx < len(fm.chain) {
		return fm.chain[fm.focusedIdx]
	}
	return nil
}

// ClearFocus removes focus from all widgets.
func (fm *FocusManager) ClearFocus() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	if fm.focusedIdx >= 0 && fm.focusedIdx < len(fm.chain) {
		fm.chain[fm.focusedIdx].SetFocused(false)
	}
	fm.focusedIdx = -1
}
