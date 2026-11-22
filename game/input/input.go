package input

import (
	"maps"
)

type Action int

const (
	MoveUp Action = iota
	MoveDown
	MoveLeft
	MoveRight
	Exit
)

type Frame struct {
	pressed  map[Action]struct{}
	released map[Action]struct{}
	held     map[Action]struct{}
}

func (f Frame) Pressed(a Action) bool {
	_, ok := f.pressed[a]
	return ok
}

func (f Frame) Released(a Action) bool {
	_, ok := f.released[a]
	return ok
}

func (f Frame) Held(a Action) bool {
	_, ok := f.held[a]
	return ok
}

func (f Frame) Empty() bool {
	return len(f.pressed) == 0 && len(f.released) == 0 && len(f.held) == 0
}

type Collector struct {
	pressed  []Action
	released []Action
	held     map[Action]struct{}
}

func NewCollector() *Collector {
	return &Collector{held: make(map[Action]struct{})}
}

func (c *Collector) Press(a Action) {
	c.pressed = append(c.pressed, a)
	c.held[a] = struct{}{}
}

func (c *Collector) Release(a Action) {
	c.released = append(c.released, a)
	delete(c.held, a)
}

func (c *Collector) Drain() Frame {
	pressed := make(map[Action]struct{}, len(c.pressed))
	for _, a := range c.pressed {
		pressed[a] = struct{}{}
	}

	released := make(map[Action]struct{}, len(c.released))
	for _, a := range c.released {
		released[a] = struct{}{}
	}

	held := make(map[Action]struct{}, len(c.held))
	maps.Copy(held, c.held)

	frame := Frame{pressed: pressed, released: released, held: held}

	c.pressed = nil
	c.released = nil

	return frame
}
