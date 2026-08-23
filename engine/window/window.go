package window

import (
	"context"
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Opts struct {
	Width   int
	Height  int
	Title   string
	GLMajor int
	GLMinor int
	Visible bool
}

type Window struct {
	window *glfw.Window
}

func Init(opts Opts) (*Window, error) {
	var err error
	w := &Window{}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.ContextVersionMajor, opts.GLMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, opts.GLMinor)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, toGlfwBool(opts.Visible))
	if w.window, err = glfw.CreateWindow(opts.Width, opts.Height, opts.Title, nil, nil); err != nil {
		return w, fmt.Errorf("failed to create GLFW w: %w", err)
	}
	w.window.MakeContextCurrent()
	glfw.SwapInterval(1)

	return w, nil
}

func (w Window) Terminate(_ context.Context) {
	w.window.Destroy()
}

func (w Window) Close() {
	w.window.SetShouldClose(true)
}

func (w Window) ShouldClose() bool {
	return w.window.ShouldClose()
}

func (w Window) PoolEvents() {
	glfw.PollEvents()
}

func (w Window) OnClose(callback func(Window)) {
	w.window.SetCloseCallback(func(_ *glfw.Window) {
		callback(w)
	})
}

func (w Window) OnSizeChanged(callback func(width int, height int)) {
	w.window.SetSizeCallback(func(_ *glfw.Window, width int, height int) {
		callback(width, height)
	})
}

func (w Window) Clear() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func (w Window) Show() {
	w.window.SwapBuffers()
}

func (w Window) Size() (int, int) {
	return w.window.GetSize()
}

// Raw returns underlying GLFW window, intended for GLFW-specific
// engine modules only (e.g. hid backend).
func (w Window) Raw() *glfw.Window {
	return w.window
}

func toGlfwBool(b bool) int {
	if b {
		return glfw.True
	}
	return glfw.False
}
