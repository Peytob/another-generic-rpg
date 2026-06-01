package window

import (
	"fmt"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Opts struct {
	Width  int
	Height int
	Title  string
}

type Window struct {
	window *glfw.Window
}

func Init(opts Opts) (*Window, error) {
	var err error
	w := &Window{}

	err = glfw.Init()
	if err != nil {
		return w, fmt.Errorf("failed to initialize GLFW: %w", err)
	}

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	if w.window, err = glfw.CreateWindow(opts.Width, opts.Height, opts.Title, nil, nil); err != nil {
		return w, fmt.Errorf("failed to create GLFW w: %w", err)
	}
	w.window.MakeContextCurrent()
	glfw.SwapInterval(int(time.Second.Milliseconds() / 60))

	return w, nil
}

func (w Window) Terminate() {
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
