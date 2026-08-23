package glfw

import (
	"testing"

	"engine/hid"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func TestKeyRoundTrip(t *testing.T) {
	for i := range hidToGlfwKey {
		hidKey := hid.Key(i)
		restored := fromGlfwKey(toGlfwKey(hidKey))
		if restored != hidKey {
			t.Fatalf("round trip failed for hid.Key(%d): got %d", hidKey, restored)
		}
	}
}

func TestGlfwKeyRoundTrip(t *testing.T) {
	for key := glfw.KeySpace; key <= glfw.KeyLast; key++ {
		hidKey := fromGlfwKey(key)
		if hidKey == hid.KeyUnknown {
			continue
		}
		if toGlfwKey(hidKey) != key {
			t.Fatalf("round trip failed for glfw.Key(%d): got %d", key, toGlfwKey(hidKey))
		}
	}
}

func TestGlfwKeyMapping(t *testing.T) {
	cases := map[hid.Key]glfw.Key{
		hid.KeyA:      glfw.KeyA,
		hid.KeyZ:      glfw.KeyZ,
		hid.Key0:      glfw.Key0,
		hid.KeyF12:    glfw.KeyF12,
		hid.KeyEscape: glfw.KeyEscape,
		hid.KeyUp:     glfw.KeyUp,
	}
	for hidKey, glfwKey := range cases {
		if toGlfwKey(hidKey) != glfwKey {
			t.Fatalf("expected hid.Key(%d) -> glfw.Key(%d), got %d", hidKey, glfwKey, toGlfwKey(hidKey))
		}
		if fromGlfwKey(glfwKey) != hidKey {
			t.Fatalf("expected glfw.Key(%d) -> hid.Key(%d), got %d", glfwKey, hidKey, fromGlfwKey(glfwKey))
		}
	}
}

func TestKeyUnknownMapping(t *testing.T) {
	if toGlfwKey(hid.KeyUnknown) != glfw.KeyUnknown {
		t.Fatal("expected hid.KeyUnknown -> glfw.KeyUnknown")
	}
	if fromGlfwKey(glfw.KeyUnknown) != hid.KeyUnknown {
		t.Fatal("expected glfw.KeyUnknown -> hid.KeyUnknown")
	}
	if fromGlfwKey(glfw.Key(-100)) != hid.KeyUnknown {
		t.Fatal("expected invalid glfw key -> hid.KeyUnknown")
	}
	if toGlfwKey(hid.Key(-100)) != glfw.KeyUnknown {
		t.Fatal("expected invalid hid key -> glfw.KeyUnknown")
	}
}

func TestMouseButtonRoundTrip(t *testing.T) {
	for i := range hidToGlfwMouseButton {
		hidBtn := hid.MouseButton(i)
		restored := fromGlfwMouseButton(toGlfwMouseButton(hidBtn))
		if restored != hidBtn {
			t.Fatalf("round trip failed for hid.MouseButton(%d): got %d", hidBtn, restored)
		}
	}
	for btn := glfw.MouseButton1; btn <= glfw.MouseButtonLast; btn++ {
		hidBtn := fromGlfwMouseButton(btn)
		if toGlfwMouseButton(hidBtn) != btn {
			t.Fatalf("round trip failed for glfw.MouseButton(%d)", btn)
		}
	}
}

func TestActionConversion(t *testing.T) {
	cases := map[hid.Action]glfw.Action{
		hid.Pressed:  glfw.Press,
		hid.Released: glfw.Release,
		hid.Repeat:   glfw.Repeat,
	}
	for hidAction, glfwAction := range cases {
		if toGlfwAction(hidAction) != glfwAction {
			t.Fatalf("expected hid.Action(%d) -> glfw.Action(%d)", hidAction, glfwAction)
		}
		if fromGlfwAction(glfwAction) != hidAction {
			t.Fatalf("expected glfw.Action(%d) -> hid.Action(%d)", glfwAction, hidAction)
		}
	}
}

func TestModsConversion(t *testing.T) {
	mods := hid.ModShift | hid.ModControl | hid.ModAlt | hid.ModSuper
	glfwMods := toGlfwMods(mods)
	if glfwMods != glfw.ModShift|glfw.ModControl|glfw.ModAlt|glfw.ModSuper {
		t.Fatalf("unexpected glfw mods: %d", glfwMods)
	}
	if fromGlfwMods(glfwMods) != mods {
		t.Fatalf("unexpected hid mods: %d", fromGlfwMods(glfwMods))
	}
}
