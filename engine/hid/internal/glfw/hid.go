package glfw

import (
	"engine/hid"
)

type Hid struct {
}

func NewHid() (Hid, error) {
	return Hid{}, nil
}

func (h Hid) Keyboard() hid.Keyboard {
	//TODO implement me
	panic("implement me")
}

func (h Hid) Mouse() hid.Mouse {
	//TODO implement me
	panic("implement me")
}
