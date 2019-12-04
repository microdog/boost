package boost

import (
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"net"
	"time"
)

var log = clog.NewWithPlugin("boost")

type Ans struct {
	net.IP
	dns.RR
}

type MeasureMethod func(b *Boost, in []Ans, out *Ans) (ok bool)
type Method int

const (
	MethodPing Method = iota
)

var methods = map[Method]MeasureMethod{
	MethodPing: measureByPing,
}

type RankOptions struct {
	Method

	// Ping related options
	PingCount    int
	PingInterval time.Duration
	PingTimeout  time.Duration
}

type Boost struct {
	Debug bool
	Next  plugin.Handler
	Zones []string

	RankOptions
}

const (
	defaultMethod       = MethodPing
	defaultPingCount    = 3
	defaultPingInterval = 25 * time.Millisecond
	defaultPingTimeout  = 500 * time.Millisecond
)

func New() *Boost {
	return &Boost{
		Debug: false,
		Next:  nil,
		Zones: []string{"."},
		RankOptions: RankOptions{
			Method:       defaultMethod,
			PingCount:    defaultPingCount,
			PingInterval: defaultPingInterval,
			PingTimeout:  defaultPingTimeout,
		},
	}
}
