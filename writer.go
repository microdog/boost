package boost

import (
	"github.com/coredns/coredns/plugin/pkg/response"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"strings"
	"time"
)

type ResponseWriter struct {
	dns.ResponseWriter
	state request.Request
	*Boost
	server string
}

func skip(m *dns.Msg) bool {
	if m.Truncated || len(m.Question) > 1 {
		return true
	}

	if mt, _ := response.Typify(m, time.Now().UTC()); mt != response.NoError {
		return true
	}

	return false
}

func (w *ResponseWriter) Write(buf []byte) (int, error) {
	log.Warning("boost called with Write(): not handling reply")
	n, err := w.ResponseWriter.Write(buf)
	return n, err
}

func (w *ResponseWriter) WriteMsg(res *dns.Msg) error {
	skip := skip(res)
	if skip {
		return w.ResponseWriter.WriteMsg(res)
	}

	answersMap := make(map[string]Ans)
	for _, rr := range res.Answer {
		switch rr.Header().Rrtype {
		case dns.TypeA:
			rr := rr.(*dns.A)
			key := rr.A.String()
			if ans, ok := answersMap[key]; !ok || ans.Header().Ttl < rr.Hdr.Ttl {
				answersMap[key] = Ans{
					IP: rr.A,
					RR: rr,
				}
			}
		case dns.TypeAAAA:
			rr := rr.(*dns.AAAA)
			key := rr.AAAA.String()
			if ans, ok := answersMap[key]; !ok || ans.Header().Ttl < rr.Hdr.Ttl {
				answersMap[key] = Ans{
					IP: rr.AAAA,
					RR: rr,
				}
			}
		}
	}

	answers := make([]Ans, 0, len(answersMap))
	for _, ans := range answersMap {
		answers = append(answers, ans)
	}

	if w.Debug {
		lines := make([]string, 0, len(answers))
		for _, ans := range answers {
			lines = append(lines, ans.RR.String())
		}
		log.Debug("collected answers:\n\t" + strings.Join(lines, "\n\t"))
	}

	if len(answers) > 0 {
		var best Ans
		ok := methods[w.Boost.Method](w.Boost, answers, &best)
		if ok {
			res.Answer = []dns.RR{best.RR}
		}
	}

	if w.Debug {
		log.Debug("final response:\n\t" + strings.ReplaceAll(strings.TrimSpace(res.String()), "\n", "\n\t"))
	}

	return w.ResponseWriter.WriteMsg(res)
}
