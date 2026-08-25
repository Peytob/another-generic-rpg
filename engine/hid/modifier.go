package hid

type Modifier int32

const (
	ModNone    Modifier = 0
	ModShift   Modifier = 1 << 0
	ModControl Modifier = 1 << 1
	ModAlt     Modifier = 1 << 2
	ModSuper   Modifier = 1 << 3
)
