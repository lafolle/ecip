package ecip

/*
`ecip` is a Coredns plugin that emits,  number of DNS queries per client address, to prometheus.
*/

import (
	"context"
	"fmt"
	"net"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var cipqc = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: plugin.Namespace,
	Name:      "ecip_count",
	Help:      "Counter for queries emitted per client addr",
}, []string{"client_addr"})

type Ecip struct {
	Next plugin.Handler
}

func (ecip Ecip) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	addr := ecip.getClientIP(w)

	names, err := net.LookupAddr(addr)
	if err != nil {
		fmt.Println("ecip: failed to lookup for:", addr, "err:", err)
		return 0, err
	}
	if len(names) == 0 {
		fmt.Println(" ecip: names empty:", addr, "err:", err)
		return 0, err
	}

	chosen_name := names[0]

	cipqc.WithLabelValues(chosen_name)

	return plugin.NextOrFailure(ecip.Name(), ecip.Next, ctx, w, r)
}

func (ecip Ecip) Name() string { return "ecip" }

func (ecip Ecip) getClientIP(w dns.ResponseWriter) string {

	if addr, ok := w.RemoteAddr().(*net.UDPAddr); ok {
		return addr.IP.String()
	}

	addr, _ := w.RemoteAddr().(*net.TCPAddr)
	return addr.IP.String()
}
