package api

import (
	"context"
)

type Emitter interface {
	Emit(ctx context.Context, name string, message *Message) error
	EmitHADiscovery(ctx context.Context, component HAComponent, message HASensorDiscovery) error
}

// EmitterFunc allows a plain function to be used as an Emitter
type EmitterFunc func(ctx context.Context, name string, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, name string, message *Message) error {
	return e(ctx, name, message)
}
