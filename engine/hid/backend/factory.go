package backend

import (
	"engine/hid"
	hidglfw "engine/hid/internal/glfw"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func NewGlfwHid(win *glfw.Window) (hid.Hid, error) {
	return hidglfw.NewHid(win)
}
