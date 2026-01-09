package polling

import (
	"context"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type Poller struct {
	name     string
	client   transport.Client
	emitter  api.Emitter
	runEvery time.Duration
	store    store.Store
}

func NewPoller(
	name string,
	client transport.Client,
	emitter api.Emitter,
	runEvery time.Duration,
	store store.Store,
) *Poller {
	return &Poller{
		name:     name,
		client:   client,
		emitter:  emitter,
		runEvery: runEvery,
		store:    store,
	}
}

func (r *Poller) Run(ctx context.Context) {
	go r.pollOnce(ctx)
	go r.pollPeriodic(ctx)
}
