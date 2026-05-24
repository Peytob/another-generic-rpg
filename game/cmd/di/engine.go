package di

import (
	"game/internal/engine/client"
)

func Client(ctx Context) *client.Client {
	return client.NewClient()
}
