package boost

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

func (b *Boost) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	zone := plugin.Zones(b.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(b.Name(), b.Next, ctx, w, r)
	}

	rw := &ResponseWriter{ResponseWriter: w, Boost: b, state: state, server: metrics.WithServer(ctx)}
	return plugin.NextOrFailure(b.Name(), b.Next, ctx, rw, r)
}

func (b *Boost) Name() string {
	return "boost"
}
