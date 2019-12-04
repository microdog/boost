package boost

import (
	"github.com/caddyserver/caddy"
	"strings"
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input        string
		shouldErr    bool
		err          string
		method       Method
		pingCount    int
		pingInterval time.Duration
		pingTimeout  time.Duration
	}{
		{"boost", false, "", defaultMethod, defaultPingCount, defaultPingInterval, defaultPingTimeout},
		{"boost {\nmethod ping\n}", false, "", defaultMethod, defaultPingCount, defaultPingInterval, defaultPingTimeout},
		{"boost {\nping_count 5\n}", false, "", defaultMethod, 5, defaultPingInterval, defaultPingTimeout},
		{"boost {\nping_interval 0.5\n}", false, "", defaultMethod, defaultPingCount, 500 * time.Millisecond, defaultPingTimeout},
		{"boost {\nping_timeout 1.0\n}", false, "", defaultMethod, defaultPingCount, defaultPingInterval, time.Second},
		{"boost {\nping_count 0\n}", true, "can not be less than 1", 0, 0, 0, 0},
		{"boost {\nping_interval 0\n}", true, "must be greater than 0", 0, 0, 0, 0},
		{"boost {\nping_timeout 0\n}", true, "must be greater than 0", 0, 0, 0, 0},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		b, err := parse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: expected no error but found one for input %s, got: %v", i, test.input, err)
			}

			if !strings.Contains(err.Error(), test.err) {
				t.Errorf("Test %d: expected error to contain: %v, found error: %v, input: %s", i, test.err, err, test.input)
			}
		}

		if !test.shouldErr {
			if test.method != b.Method {
				t.Errorf("Test %d: expected %d, got: %d", i, test.method, b.Method)
			}
			if test.pingCount != b.PingCount {
				t.Errorf("Test %d: expected %d, got: %d", i, test.pingCount, b.PingCount)
			}
			if test.pingInterval != b.PingInterval {
				t.Errorf("Test %d: expected %s, got: %s", i, test.pingInterval, b.PingInterval)
			}
			if test.pingTimeout != b.PingTimeout {
				t.Errorf("Test %d: expected %s, got: %s", i, test.pingTimeout, b.PingTimeout)
			}
		}
	}
}
