package ecip

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("ecip", setup) }

func setup(c *caddy.Controller) error {
	c.Next()
	if c.NextArg() {
		return plugin.Error("ecip", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin((func(next plugin.Handler) plugin.Handler {
		return Ecip{Next: next}
	}))

	return nil
}
