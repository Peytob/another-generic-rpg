package client

import (
	"context"
	"fmt"
	"game/internal/engine/fsm"
	"game/internal/engine/fsm/state"
	"game/pkg/utils/logger"
)

type Client struct {
	fsm fsm.Machine
}

func NewClient() *Client {
	client := &Client{}

	client.fsm = initializeClientMachine()

	return client
}

func (c *Client) Run(ctx context.Context) error {
	for {
		var err error

		select {
		case <-ctx.Done():
			// Mark window as should close
			logger.FromCtx(ctx).Info("root context done, closing window")
			err = c.fsm.Event(fsm.StoppedEvent)
			if err != nil {
				return fmt.Errorf("failed to stop machine on context close: %w", err)
			}
		default:
			// nothing, keep running
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
	return nil
}

func initializeClientMachine() fsm.Machine {
	return fsm.NewBuilder().
		RegisterState(state.NewInitialState(), fsm.NoTransitions()).
		InitialState(state.InitialStateIdentifier).
		RegisterState(state.NewStoppedState(), fsm.NoTransitions()).
		GlobalTransitions(fsm.Transitions{
			fsm.StoppedEvent: state.StoppedStateIdentifier,
		}).
		FinalStates([]fsm.StateIdentifier{
			state.StoppedStateIdentifier,
		}).
		ShouldBuild()
}
