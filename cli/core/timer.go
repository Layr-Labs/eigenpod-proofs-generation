package core

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Event struct {
	desc  string
	start time.Time
}

type Timer struct {
	events []Event
	id     string
}

type TPerfTimerKey string

const PERF_TIMER_KEY TPerfTimerKey = "com.eigen.timer"

func NewTimer() *Timer {
	trace := rand.Int()
	return &Timer{
		id: fmt.Sprintf("request-%d", trace),
	}
}

func TimerFromContext(ctx context.Context) (*Timer, context.Context) {
	timer := ctx.Value(PERF_TIMER_KEY)
	if timer != nil {
		// existing ctx already has a timer.
		return timer.(*Timer), ctx
	}

	newTimer := NewTimer()
	return newTimer, context.WithValue(ctx, PERF_TIMER_KEY, newTimer)
}

func (t *Timer) Start(section string) {
	var e Event
	e.desc = section
	e.start = time.Now()

	t.events = append(t.events, e)
}

func (t *Timer) End() {
	e := t.events[len(t.events)-1]
	t.events = t.events[:len(t.events)-1]

	dur := time.Now().UnixMilli() - e.start.UnixMilli()
	fmt.Printf("[%s] %s (%dms)\n", t.id, e.desc, dur)
}
