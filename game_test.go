//go:build windows || darwin || android || ios || linux

package wayne

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestEbitenKeyToWayne(t *testing.T) {
	tests := []struct {
		name      string
		ebitenKey ebiten.Key
		wantKey   Key
	}{
		// Special keys
		{"Escape", ebiten.KeyEscape, KeyEscape},
		{"Enter", ebiten.KeyEnter, KeyReturn},
		{"Tab", ebiten.KeyTab, KeyTab},
		{"Backspace", ebiten.KeyBackspace, KeyBackspace},
		{"Delete", ebiten.KeyDelete, KeyDelete},
		{"Left Arrow", ebiten.KeyArrowLeft, KeyLeft},
		{"Up Arrow", ebiten.KeyArrowUp, KeyUp},
		{"Right Arrow", ebiten.KeyArrowRight, KeyRight},
		{"Down Arrow", ebiten.KeyArrowDown, KeyDown},
		{"Home", ebiten.KeyHome, KeyHome},
		{"End", ebiten.KeyEnd, KeyEnd},
		{"PageUp", ebiten.KeyPageUp, KeyPageUp},
		{"PageDown", ebiten.KeyPageDown, KeyPageDown},
		{"Space", ebiten.KeySpace, KeySpace},

		// Modifier keys
		{"ShiftLeft", ebiten.KeyShiftLeft, KeyShiftL},
		{"ShiftRight", ebiten.KeyShiftRight, KeyShiftR},
		{"ControlLeft", ebiten.KeyControlLeft, KeyControlL},
		{"ControlRight", ebiten.KeyControlRight, KeyControlR},
		{"Alt", ebiten.KeyAlt, KeyAltL},
		{"AltRight", ebiten.KeyAltRight, KeyAltR},
		{"MetaLeft", ebiten.KeyMetaLeft, KeySuperL},
		{"MetaRight", ebiten.KeyMetaRight, KeySuperR},

		// Letter keys
		{"A", ebiten.KeyA, Key('a')},
		{"B", ebiten.KeyB, Key('b')},
		{"M", ebiten.KeyM, Key('m')},
		{"Z", ebiten.KeyZ, Key('z')},

		// Digit keys
		{"Digit0", ebiten.KeyDigit0, Key('0')},
		{"Digit5", ebiten.KeyDigit5, Key('5')},
		{"Digit9", ebiten.KeyDigit9, Key('9')},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ebitenKeyToWayne(tt.ebitenKey)
			if got != tt.wantKey {
				t.Errorf("ebitenKeyToWayne(%v) = %v, want %v", tt.ebitenKey, got, tt.wantKey)
			}
		})
	}
}

func TestEbitenKeyToWayneUnmapped(t *testing.T) {
	// Test that unmapped keys fall through to the default cast behavior
	unmappedKey := ebiten.Key(9999)
	got := ebitenKeyToWayne(unmappedKey)
	want := Key(9999)
	if got != want {
		t.Errorf("ebitenKeyToWayne(unmapped) = %v, want %v", got, want)
	}
}

func TestEbitenToWayneKeyMapCompleteness(t *testing.T) {
	// Verify all expected keys are in the map
	expectedKeys := []ebiten.Key{
		ebiten.KeyEscape, ebiten.KeyEnter, ebiten.KeyTab, ebiten.KeyBackspace,
		ebiten.KeyDelete, ebiten.KeyArrowLeft, ebiten.KeyArrowUp, ebiten.KeyArrowRight,
		ebiten.KeyArrowDown, ebiten.KeyHome, ebiten.KeyEnd, ebiten.KeyPageUp,
		ebiten.KeyPageDown, ebiten.KeySpace, ebiten.KeyShiftLeft, ebiten.KeyShiftRight,
		ebiten.KeyControlLeft, ebiten.KeyControlRight, ebiten.KeyAlt, ebiten.KeyAltRight,
		ebiten.KeyMetaLeft, ebiten.KeyMetaRight,
		ebiten.KeyA, ebiten.KeyB, ebiten.KeyC, ebiten.KeyD, ebiten.KeyE, ebiten.KeyF,
		ebiten.KeyG, ebiten.KeyH, ebiten.KeyI, ebiten.KeyJ, ebiten.KeyK, ebiten.KeyL,
		ebiten.KeyM, ebiten.KeyN, ebiten.KeyO, ebiten.KeyP, ebiten.KeyQ, ebiten.KeyR,
		ebiten.KeyS, ebiten.KeyT, ebiten.KeyU, ebiten.KeyV, ebiten.KeyW, ebiten.KeyX,
		ebiten.KeyY, ebiten.KeyZ,
		ebiten.KeyDigit0, ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3,
		ebiten.KeyDigit4, ebiten.KeyDigit5, ebiten.KeyDigit6, ebiten.KeyDigit7,
		ebiten.KeyDigit8, ebiten.KeyDigit9,
	}

	for _, key := range expectedKeys {
		if _, ok := ebitenToWayneKeyMap[key]; !ok {
			t.Errorf("ebitenToWayneKeyMap is missing key: %v", key)
		}
	}

	// Verify the map has exactly the expected number of keys
	expectedCount := len(expectedKeys)
	actualCount := len(ebitenToWayneKeyMap)
	if actualCount != expectedCount {
		t.Errorf("ebitenToWayneKeyMap has %d keys, expected %d", actualCount, expectedCount)
	}
}
