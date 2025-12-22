package api

import (
	"context"
)

type Emitter interface {
	Emit(ctx context.Context, id string, message *Message) error
}

// EmitterFunc allows a plain function to be used as an Emitter
type EmitterFunc func(ctx context.Context, id string, message *Message) error

func (e EmitterFunc) Emit(ctx context.Context, id string, message *Message) error {
	return e(ctx, id, message)
}
