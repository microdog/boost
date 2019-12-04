package boost

import (
	"fmt"
	"github.com/caddyserver/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"strconv"
	"time"
)

func init() {
	plugin.Register("boost", setup)
}

func setup(c *caddy.Controller) error {
	boost, err := parse(c)
	if err != nil {
		return plugin.Error("boost", err)
	}

	config := dnsserver.GetConfig(c)
	config.AddPlugin(func(next plugin.Handler) plugin.Handler {
		boost.Debug = config.Debug
		boost.Next = next
		return boost
	})

	return nil
}

func parse(c *caddy.Controller) (*Boost, error) {
	boost := New()

	// cache [zones..]
	c.Next()
	origins := make([]string, len(c.ServerBlockKeys))
	copy(origins, c.ServerBlockKeys)
	args := c.RemainingArgs()

	if len(args) > 0 {
		copy(origins, args)
	}

	for i := range origins {
		origins[i] = plugin.Host(origins[i]).Normalize()
	}
	boost.Zones = origins

	for c.NextBlock() {
		switch c.Val() {
		case "method":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			method, ok := argsMethods[args[0]]
			if !ok {
				return nil, c.ArgErr()
			}
			boost.Method = method
		case "ping_count":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			pingCount, err := strconv.Atoi(args[0])
			if err != nil {
				return nil, err
			}
			if pingCount < 1 {
				return nil, fmt.Errorf("ping_count can not be less than 1")
			}
			boost.PingCount = pingCount
		case "ping_interval":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			pingInterval, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, err
			}
			if pingInterval <= 0.0 {
				return nil, fmt.Errorf("ping_interval must be greater than 0")
			}
			boost.PingInterval = time.Duration(pingInterval * float64(time.Second))
		case "ping_timeout":
			args := c.RemainingArgs()
			if len(args) != 1 {
				return nil, c.ArgErr()
			}
			pingTimeout, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return nil, err
			}
			if pingTimeout <= 0.0 {
				return nil, fmt.Errorf("ping_timeout must be greater than 0")
			}
			boost.PingTimeout = time.Duration(pingTimeout * float64(time.Second))
		}
	}

	return boost, nil
}

var (
	argsMethods = map[string]Method{
		"ping": MethodPing,
	}
)
