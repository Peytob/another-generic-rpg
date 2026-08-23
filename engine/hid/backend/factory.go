package backend

import (
	"engine/hid"
	"engine/hid/internal/glfw"
)

func NewGlfwHid() (hid.Hid, error) {
	return glfw.NewHid()
}
