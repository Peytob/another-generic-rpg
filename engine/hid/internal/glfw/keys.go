package glfw

import (
	"engine/hid"

	"github.com/go-gl/glfw/v3.3/glfw"
)

// hidToGlfwKey maps engine key codes to GLFW ones, element index = hid.Key.
var hidToGlfwKey = [...]glfw.Key{
	hid.KeyUnknown: glfw.KeyUnknown,

	hid.KeyA: glfw.KeyA,
	hid.KeyB: glfw.KeyB,
	hid.KeyC: glfw.KeyC,
	hid.KeyD: glfw.KeyD,
	hid.KeyE: glfw.KeyE,
	hid.KeyF: glfw.KeyF,
	hid.KeyG: glfw.KeyG,
	hid.KeyH: glfw.KeyH,
	hid.KeyI: glfw.KeyI,
	hid.KeyJ: glfw.KeyJ,
	hid.KeyK: glfw.KeyK,
	hid.KeyL: glfw.KeyL,
	hid.KeyM: glfw.KeyM,
	hid.KeyN: glfw.KeyN,
	hid.KeyO: glfw.KeyO,
	hid.KeyP: glfw.KeyP,
	hid.KeyQ: glfw.KeyQ,
	hid.KeyR: glfw.KeyR,
	hid.KeyS: glfw.KeyS,
	hid.KeyT: glfw.KeyT,
	hid.KeyU: glfw.KeyU,
	hid.KeyV: glfw.KeyV,
	hid.KeyW: glfw.KeyW,
	hid.KeyX: glfw.KeyX,
	hid.KeyY: glfw.KeyY,
	hid.KeyZ: glfw.KeyZ,

	hid.Key0: glfw.Key0,
	hid.Key1: glfw.Key1,
	hid.Key2: glfw.Key2,
	hid.Key3: glfw.Key3,
	hid.Key4: glfw.Key4,
	hid.Key5: glfw.Key5,
	hid.Key6: glfw.Key6,
	hid.Key7: glfw.Key7,
	hid.Key8: glfw.Key8,
	hid.Key9: glfw.Key9,

	hid.KeyF1:  glfw.KeyF1,
	hid.KeyF2:  glfw.KeyF2,
	hid.KeyF3:  glfw.KeyF3,
	hid.KeyF4:  glfw.KeyF4,
	hid.KeyF5:  glfw.KeyF5,
	hid.KeyF6:  glfw.KeyF6,
	hid.KeyF7:  glfw.KeyF7,
	hid.KeyF8:  glfw.KeyF8,
	hid.KeyF9:  glfw.KeyF9,
	hid.KeyF10: glfw.KeyF10,
	hid.KeyF11: glfw.KeyF11,
	hid.KeyF12: glfw.KeyF12,

	hid.KeyUp:    glfw.KeyUp,
	hid.KeyDown:  glfw.KeyDown,
	hid.KeyLeft:  glfw.KeyLeft,
	hid.KeyRight: glfw.KeyRight,

	hid.KeyInsert:   glfw.KeyInsert,
	hid.KeyDelete:   glfw.KeyDelete,
	hid.KeyHome:     glfw.KeyHome,
	hid.KeyEnd:      glfw.KeyEnd,
	hid.KeyPageUp:   glfw.KeyPageUp,
	hid.KeyPageDown: glfw.KeyPageDown,

	hid.KeyEnter:     glfw.KeyEnter,
	hid.KeySpace:     glfw.KeySpace,
	hid.KeyTab:       glfw.KeyTab,
	hid.KeyBackspace: glfw.KeyBackspace,
	hid.KeyEscape:    glfw.KeyEscape,

	hid.KeyLeftShift:    glfw.KeyLeftShift,
	hid.KeyLeftControl:  glfw.KeyLeftControl,
	hid.KeyLeftAlt:      glfw.KeyLeftAlt,
	hid.KeyLeftSuper:    glfw.KeyLeftSuper,
	hid.KeyRightShift:   glfw.KeyRightShift,
	hid.KeyRightControl: glfw.KeyRightControl,
	hid.KeyRightAlt:     glfw.KeyRightAlt,
	hid.KeyRightSuper:   glfw.KeyRightSuper,
}

// glfwToHidKey maps GLFW key codes to engine ones, element index = glfw.Key.
var glfwToHidKey [glfw.KeyLast + 1]hid.Key

// glfwMouseButtonInvalid mirrors GLFW_MOUSE_BUTTON_INVALID which is not
// exposed by go-gl bindings.
const glfwMouseButtonInvalid = glfw.MouseButton(-1)

// hidToGlfwMouseButton maps engine mouse buttons to GLFW ones, element index = hid.MouseButton.
var hidToGlfwMouseButton = [...]glfw.MouseButton{
	hid.MouseButtonUnknown: glfwMouseButtonInvalid,
	hid.MouseButtonLeft:    glfw.MouseButtonLeft,
	hid.MouseButtonRight:   glfw.MouseButtonRight,
	hid.MouseButtonMiddle:  glfw.MouseButtonMiddle,
	hid.MouseButton4:       glfw.MouseButton4,
	hid.MouseButton5:       glfw.MouseButton5,
	hid.MouseButton6:       glfw.MouseButton6,
	hid.MouseButton7:       glfw.MouseButton7,
	hid.MouseButton8:       glfw.MouseButton8,
}

// glfwToHidMouseButton maps GLFW mouse buttons to engine ones, element index = glfw.MouseButton.
var glfwToHidMouseButton [glfw.MouseButtonLast + 1]hid.MouseButton

func init() {
	for hidKey, glfwKey := range hidToGlfwKey {
		if glfwKey > glfw.KeyUnknown && glfwKey <= glfw.KeyLast {
			glfwToHidKey[glfwKey] = hid.Key(hidKey)
		}
	}

	for hidBtn, glfwBtn := range hidToGlfwMouseButton {
		if glfwBtn > glfwMouseButtonInvalid && glfwBtn <= glfw.MouseButtonLast {
			glfwToHidMouseButton[glfwBtn] = hid.MouseButton(hidBtn)
		}
	}
}

func toGlfwKey(key hid.Key) glfw.Key {
	if key < 0 || int(key) >= len(hidToGlfwKey) {
		return glfw.KeyUnknown
	}
	return hidToGlfwKey[key]
}

func fromGlfwKey(key glfw.Key) hid.Key {
	if key <= glfw.KeyUnknown || key > glfw.KeyLast {
		return hid.KeyUnknown
	}
	return glfwToHidKey[key]
}

func toGlfwMouseButton(button hid.MouseButton) glfw.MouseButton {
	if button < 0 || int(button) >= len(hidToGlfwMouseButton) {
		return glfwMouseButtonInvalid
	}
	return hidToGlfwMouseButton[button]
}

func fromGlfwMouseButton(button glfw.MouseButton) hid.MouseButton {
	if button <= glfwMouseButtonInvalid || button > glfw.MouseButtonLast {
		return hid.MouseButtonUnknown
	}
	return glfwToHidMouseButton[button]
}

// hid.Action ordering (Pressed=0, Released=1) does not match GLFW
// (Release=0, Press=1), so actions are converted explicitly instead of casting.
func toGlfwAction(action hid.Action) glfw.Action {
	switch action {
	case hid.Released:
		return glfw.Release
	case hid.Repeat:
		return glfw.Repeat
	case hid.Pressed:
		return glfw.Press
	default:
		return glfw.Release
	}
}

func fromGlfwAction(action glfw.Action) hid.Action {
	switch action {
	case glfw.Release:
		return hid.Released
	case glfw.Repeat:
		return hid.Repeat
	case glfw.Press:
		return hid.Pressed
	default:
		return hid.Released
	}
}

// hid.Modifier bitmask values match GLFW ModifierKey values, plain cast is enough.
func toGlfwMods(mods hid.Modifier) glfw.ModifierKey {
	return glfw.ModifierKey(mods)
}

func fromGlfwMods(mods glfw.ModifierKey) hid.Modifier {
	return hid.Modifier(mods)
}
