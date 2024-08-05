package core

import (
	"context"
)

type TracerCallbacks struct {
	OnStartSection func(name string, metadata map[string]string)
	OnEndSection   func()
}
type TEigenKey string

const EIGEN_KEY TEigenKey = "com.eigen.tracer"

func ContextWithTracing(ctx context.Context, callbacks *TracerCallbacks) context.Context {
	return context.WithValue(ctx, EIGEN_KEY, callbacks)
}

func GetContextTracingCallbacks(ctx context.Context) *TracerCallbacks {
	tracing, ok := ctx.Value(EIGEN_KEY).(*TracerCallbacks)
	if !ok || tracing == nil {
		return &TracerCallbacks{
			OnStartSection: func(name string, meta map[string]string) {},
			OnEndSection:   func() {},
		}
	}

	return tracing
}
