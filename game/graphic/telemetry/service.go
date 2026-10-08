package telemetry

import (
	"sync"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Telemetry interface {
	About() About
}

type telemetry struct {
	about     About
	aboutOnce sync.Once
}

func NewTelemetry() Telemetry {
	return &telemetry{}
}

func (t *telemetry) About() About {
	t.aboutOnce.Do(func() {
		t.about = About{
			Name:     gl.GoStr(gl.GetString(gl.VENDOR)),
			Version:  gl.GoStr(gl.GetString(gl.VERSION)),
			Renderer: gl.GoStr(gl.GetString(gl.RENDERER)),
		}
	})

	return t.about
}
