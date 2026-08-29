package gameloop

import (
	"context"
	"errors"
	"testing"
	"time"

	"game/internal/gameplay/world"
)

type testAction int

const (
	testActionA testAction = iota
	testActionB
)

var errStageSentinel = errors.New("stage boom")

func TestBeforeAndAfterRunOncePerFrame(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](100 * time.Millisecond)

	var beforeCnt, tickCnt, afterCnt int

	l.Before(func(_ context.Context, _ *Frame[int, testAction]) error {
		beforeCnt++
		return nil
	})
	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		tickCnt++
		return nil
	})
	l.After(func(_ context.Context, _ *Frame[int, testAction]) error {
		afterCnt++
		return nil
	})

	w := &world.World{}

	// First frame: 60ms accumulated, no tick yet.
	_, err := l.Run(context.Background(), w, 60*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 1: %v", err)
	}
	// Second frame: accumulator reaches 120ms, one tick fires.
	_, err = l.Run(context.Background(), w, 60*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 2: %v", err)
	}

	if beforeCnt != 2 {
		t.Errorf("before stages should run once per frame, got %d", beforeCnt)
	}
	if afterCnt != 2 {
		t.Errorf("after stages should run once per frame, got %d", afterCnt)
	}
	if tickCnt != 1 {
		t.Errorf("tick stages should run once in total, got %d", tickCnt)
	}
}

func TestEachTickFixedStepAccumulator(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](100 * time.Millisecond)

	var tickCnt int
	var lastAlpha float64

	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		tickCnt++
		return nil
	})
	l.After(func(_ context.Context, f *Frame[int, testAction]) error {
		lastAlpha = f.Alpha
		return nil
	})

	w := &world.World{}

	// Frame 1: 250ms -> 2 ticks, 50ms left.
	_, err := l.Run(context.Background(), w, 250*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 1: %v", err)
	}
	if tickCnt != 2 {
		t.Errorf("frame 1: expected 2 ticks, got %d", tickCnt)
	}
	if lastAlpha != 0.5 {
		t.Errorf("frame 1: expected alpha 0.5, got %v", lastAlpha)
	}

	// Frame 2: +60ms -> 110ms accumulated -> 1 tick, 10ms left.
	_, err = l.Run(context.Background(), w, 60*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 2: %v", err)
	}
	if tickCnt != 3 {
		t.Errorf("frame 2: expected 3 ticks in total, got %d", tickCnt)
	}
	if lastAlpha != 0.1 {
		t.Errorf("frame 2: expected alpha 0.1, got %v", lastAlpha)
	}
}

func TestTickStagesSeeTickDtOthersSeeFrameDt(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](50 * time.Millisecond)

	var beforeDt, tickDt, afterDt time.Duration

	l.Before(func(_ context.Context, f *Frame[int, testAction]) error {
		beforeDt = f.Dt
		return nil
	})
	l.EachTick(func(_ context.Context, f *Frame[int, testAction]) error {
		tickDt = f.Dt
		return nil
	})
	l.After(func(_ context.Context, f *Frame[int, testAction]) error {
		afterDt = f.Dt
		return nil
	})

	_, err := l.Run(context.Background(), &world.World{}, 120*time.Millisecond)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if beforeDt != 120*time.Millisecond {
		t.Errorf("before stage should see real dt, got %v", beforeDt)
	}
	if tickDt != 50*time.Millisecond {
		t.Errorf("tick stage should see fixed tick, got %v", tickDt)
	}
	if afterDt != 120*time.Millisecond {
		t.Errorf("after stage should see real dt, got %v", afterDt)
	}
}

func TestFrameDtClamped(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](50 * time.Millisecond)

	var tickCnt int
	var frameDt time.Duration

	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		tickCnt++
		return nil
	})
	l.Before(func(_ context.Context, f *Frame[int, testAction]) error {
		frameDt = f.Dt
		return nil
	})

	// 10 seconds stall must be clamped to 250ms -> at most 5 ticks.
	_, err := l.Run(context.Background(), &world.World{}, 10*time.Second)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if frameDt != maxFrameDt {
		t.Errorf("dt should be clamped to %v, got %v", maxFrameDt, frameDt)
	}
	if tickCnt != maxTicksPerFrame {
		t.Errorf("expected %d ticks after clamp, got %d", maxTicksPerFrame, tickCnt)
	}
}

func TestLagDroppedAfterTickCap(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](10 * time.Millisecond)

	var tickCnt int

	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		tickCnt++
		return nil
	})

	// 250ms with a 10ms tick wants 25 ticks; the cap allows 5 and the
	// remaining lag must be dropped, so the next small frame produces
	// exactly one tick instead of a catch-up burst.
	_, err := l.Run(context.Background(), &world.World{}, 250*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 1: %v", err)
	}
	if tickCnt != maxTicksPerFrame {
		t.Errorf("frame 1: expected %d ticks, got %d", maxTicksPerFrame, tickCnt)
	}

	_, err = l.Run(context.Background(), &world.World{}, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("Run 2: %v", err)
	}
	if tickCnt != maxTicksPerFrame+1 {
		t.Errorf("frame 2: expected %d ticks in total, got %d", maxTicksPerFrame+1, tickCnt)
	}
}

func TestZeroTickRunsEachTickOncePerFrame(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](0)

	var tickCnt int
	var tickDt time.Duration
	var alpha float64

	l.EachTick(func(_ context.Context, f *Frame[int, testAction]) error {
		tickCnt++
		tickDt = f.Dt
		return nil
	})
	l.After(func(_ context.Context, f *Frame[int, testAction]) error {
		alpha = f.Alpha
		return nil
	})

	for i := 0; i < 3; i++ {
		_, err := l.Run(context.Background(), &world.World{}, 33*time.Millisecond)
		if err != nil {
			t.Fatalf("Run %d: %v", i, err)
		}
	}

	if tickCnt != 3 {
		t.Errorf("tick stages should run once per frame, got %d", tickCnt)
	}
	if tickDt != 33*time.Millisecond {
		t.Errorf("tick stage should see real dt, got %v", tickDt)
	}
	if alpha != 0 {
		t.Errorf("alpha should stay zero in variable-step mode, got %v", alpha)
	}
}

func TestTransitionRequestedBeforeFrameCompletes(t *testing.T) {
	t.Parallel()

	const stopEvent = 7

	l := NewLoop[int, testAction](0)

	var afterCnt int

	l.Before(func(_ context.Context, f *Frame[int, testAction]) error {
		f.RequestTransition(stopEvent)
		return nil
	})
	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		return nil
	})
	l.After(func(_ context.Context, _ *Frame[int, testAction]) error {
		afterCnt++
		return nil
	})

	event, err := l.Run(context.Background(), &world.World{}, time.Millisecond)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	if afterCnt != 1 {
		t.Errorf("frame should complete after transition request, after stages ran %d", afterCnt)
	}
	if event != stopEvent {
		t.Errorf("expected transition %d, got %d", stopEvent, event)
	}
}

func TestNoTransitionReturnsZeroValue(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](0)
	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		return nil
	})

	event, err := l.Run(context.Background(), &world.World{}, time.Millisecond)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if event != 0 {
		t.Errorf("expected zero event, got %d", event)
	}
}

func TestBeforeStageErrorPropagates(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](0)

	var tickCnt, afterCnt int

	l.Before(func(_ context.Context, _ *Frame[int, testAction]) error {
		return errStageSentinel
	})
	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		tickCnt++
		return nil
	})
	l.After(func(_ context.Context, _ *Frame[int, testAction]) error {
		afterCnt++
		return nil
	})

	_, err := l.Run(context.Background(), &world.World{}, time.Millisecond)
	if !errors.Is(err, errStageSentinel) {
		t.Fatalf("expected errStageSentinel, got %v", err)
	}
	if tickCnt != 0 || afterCnt != 0 {
		t.Errorf("no further stages should run on error, got ticks=%d after=%d", tickCnt, afterCnt)
	}
}

func TestTickStageErrorPropagates(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](100 * time.Millisecond)

	l.EachTick(func(_ context.Context, _ *Frame[int, testAction]) error {
		return errStageSentinel
	})

	// 2 seconds get clamped to 250ms, which is still enough for ticks
	// to fire and fail.
	_, err := l.Run(context.Background(), &world.World{}, 2*time.Second)
	if !errors.Is(err, errStageSentinel) {
		t.Fatalf("expected errStageSentinel, got %v", err)
	}
}

func TestWorldPassedToStages(t *testing.T) {
	t.Parallel()

	l := NewLoop[int, testAction](0)

	var seen *world.World

	l.EachTick(func(_ context.Context, f *Frame[int, testAction]) error {
		seen = f.World
		return nil
	})

	w := &world.World{}
	_, err := l.Run(context.Background(), w, time.Millisecond)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if seen != w {
		t.Error("stages should receive the world passed to Run")
	}
}
