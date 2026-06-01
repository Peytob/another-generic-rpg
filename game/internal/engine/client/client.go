package client

import (
	"context"
	"fmt"
	"game/internal/engine/fsm"
	"game/pkg/engine"
	"game/pkg/utils/logger"
	"game/pkg/window"
)

type Client struct {
	fsm    engine.Machine
	window *window.Window
}

func (c *Client) Run(ctx context.Context) error {
	for {
		var err error

		select {
		case <-ctx.Done():
			// Mark window as should close
			logger.FromCtx(ctx).Info("root context done, closing window")
			c.window.Close()
		default:
			// nothing, keep running
		}

		c.window.PoolEvents()

		if c.window.ShouldClose() {
			logger.FromCtx(ctx).Debug("stopping game machine")
			err = c.fsm.Event(fsm.StoppedEvent)
			if err != nil {
				return fmt.Errorf("failed to stop machine on window close: %w", err)
			}
			break
		}

		event, err := c.fsm.State().Update(ctx)
		if err != nil {
			return fmt.Errorf("error while executing FSM update: %w", err)
		}

		if event != fsm.NoEvent {
			err = c.fsm.Event(event)
			if err != nil {
				return fmt.Errorf("error while changing FSM state: %w", err)
			}
		}

		if !c.fsm.IsRunning() {
			logger.FromCtx(ctx).Info("game FSM is completed, closing window")
			break
		}
	}

	return nil
}

func (c *Client) Shutdown(ctx context.Context) error {
	logger.FromCtx(ctx).Info("shutting down client")
	c.window.Terminate()
	return nil
}
