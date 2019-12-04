package boost

import (
	"fmt"
	"github.com/microdog/go-ping"
	"sort"
	"strings"
	"sync"
	"time"
)

func makePinger(b *Boost, ans Ans, ch chan pingStats, wg *sync.WaitGroup) func() {
	finished := false

	pinger, err := ping.NewPinger(ans.IP.String())
	if err != nil {
		log.Errorf("cannot create pinger: %s", err)
		return nil
	}

	pinger.Count = b.PingCount
	if pinger.Count < 1 {
		pinger.Count = 1
	}
	pinger.Interval = b.PingInterval
	pinger.Timeout = b.PingTimeout
	pinger.OnFinish = func(stats *ping.Statistics) {
		if b.Debug {
			log.Debugf("ping stats: %+v", stats)
		}
		ch <- pingStats{
			Ans:      ans,
			Sent:     stats.PacketsSent,
			Received: stats.PacketsRecv,
			Rtts:     stats.Rtts,
			Best:     stats.MinRtt,
			Worst:    stats.MaxRtt,
			Avg:      stats.AvgRtt,
			StdDev:   stats.StdDevRtt,
			Loss:     stats.PacketLoss,
		}
		if !finished {
			wg.Done()
			finished = true
		}
	}
	pinger.SetPrivileged(true)

	wg.Add(1)
	return func() {
		pinger.Run()
		if pinger.PacketsSent == 0 {
			log.Warning("0 packets Sent, maybe you are running as an non-privileged user")
			if !finished {
				wg.Done()
				finished = true
			}
		}
	}
}

func measureByPing(b *Boost, in []Ans, out *Ans) bool {
	wg := &sync.WaitGroup{}
	resultsCh := make(chan pingStats, len(in))

	for _, ans := range in {
		pinger := makePinger(b, ans, resultsCh, wg)
		go pinger()
	}

	wg.Wait()
	close(resultsCh)

	results := make([]pingStats, 0, len(in))
	for result := range resultsCh {
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		ri, rj := results[i], results[j]
		if ri.Loss < rj.Loss {
			return true
		}
		if ri.Loss > rj.Loss {
			return false
		}
		if ri.Avg < rj.Avg {
			return true
		}
		if ri.Avg > rj.Avg {
			return false
		}
		return ri.StdDev < rj.StdDev
	})

	if b.Debug {
		lines := make([]string, 0, len(results))
		for _, result := range results {
			lines = append(lines, result.String())
		}
		log.Debug("sorted results\n\t" + strings.Join(lines, "\n\t"))
	}

	if len(results) > 0 {
		*out = results[0].Ans
		return true
	}
	return false
}

type pingStats struct {
	Ans
	Sent     int
	Received int
	Rtts     []time.Duration
	Best     time.Duration
	Worst    time.Duration
	Avg      time.Duration
	StdDev   time.Duration
	Loss     float64
}

func (r pingStats) String() string {
	return fmt.Sprintf("pingStats(IP=%s, Sent=%d, Loss=%f, Avg=%s, StdDev=%s)", r.IP.String(), r.Sent, r.Loss, r.Avg, r.StdDev)
}
