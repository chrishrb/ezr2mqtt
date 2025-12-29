package api

import (
	"context"
)

type Emitter interface {
	Emit(ctx context.Context, name string, message *Message) error
	EmitHADiscovery(ctx context.Context, component HAComponent, message *HADiscovery) error
}
