package polling

import (
	"context"
	"time"

	"github.com/chrishrb/ezr2mqtt/api"
	"github.com/chrishrb/ezr2mqtt/store"
	"github.com/chrishrb/ezr2mqtt/transport"
)

type Poller struct {
	name              string
	client            transport.Client
	emitter           api.Emitter
	runEvery          time.Duration
	discoveryRunEvery time.Duration
	store             store.Store
}

type PollerOption func(*Poller)

func WithDiscoveryInterval(interval time.Duration) PollerOption {
	return func(p *Poller) {
		p.discoveryRunEvery = interval
	}
}

func NewPoller(
	name string,
	client transport.Client,
	emitter api.Emitter,
	runEvery time.Duration,
	store store.Store,
	opts ...PollerOption,
) *Poller {
	p := &Poller{
		name:              name,
		client:            client,
		emitter:           emitter,
		runEvery:          runEvery,
		discoveryRunEvery: 10 * time.Minute,
		store:             store,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (r *Poller) Run(ctx context.Context) {
	go r.runDiscovery(ctx)
	go r.pollPeriodic(ctx)
}
