package domaintechnik

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/schuellerf/libdns-domaintechnik"
)

// Provider wraps the libdns Domaintechnik provider as a Caddy module.
type Provider struct{ *domaintechnik.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.domaintechnik",
		New: func() caddy.Module { return &Provider{new(domaintechnik.Provider)} },
	}
}

// Provision expands placeholders in the API token and optional endpoint.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.APIToken = repl.ReplaceAll(p.Provider.APIToken, "")
	p.Provider.APIEndpoint = repl.ReplaceAll(p.Provider.APIEndpoint, "")
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens. Syntax:
//
//	domaintechnik [<api_token>] {
//	    api_token <api_token>
//	    api_endpoint <url>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			p.Provider.APIToken = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "api_token":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.APIToken != "" {
					return d.Err("API token already set")
				}
				p.Provider.APIToken = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "api_endpoint":
				if !d.NextArg() {
					return d.ArgErr()
				}
				if p.Provider.APIEndpoint != "" {
					return d.Err("API endpoint already set")
				}
				p.Provider.APIEndpoint = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	if p.Provider.APIToken == "" {
		return d.Err("missing API token")
	}
	return nil
}

// Interface guards
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
