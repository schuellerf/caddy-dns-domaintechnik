package domaintechnik

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/schuellerf/libdns-domaintechnik"
)

func TestUnmarshalCaddyfileTokenArg(t *testing.T) {
	d := caddyfile.NewTestDispenser(`domaintechnik {env.DT_TOKEN}`)
	p := &Provider{Provider: new(domaintechnik.Provider)}
	if err := p.UnmarshalCaddyfile(d); err != nil {
		t.Fatal(err)
	}
	if p.APIToken != "{env.DT_TOKEN}" {
		t.Errorf("token = %q", p.APIToken)
	}
}

func TestUnmarshalCaddyfileBlock(t *testing.T) {
	d := caddyfile.NewTestDispenser(`domaintechnik {
		api_token dt1_literal
		api_endpoint https://example.test/dns-api/v1
	}`)
	p := &Provider{Provider: new(domaintechnik.Provider)}
	if err := p.UnmarshalCaddyfile(d); err != nil {
		t.Fatal(err)
	}
	if p.APIToken != "dt1_literal" {
		t.Errorf("token = %q", p.APIToken)
	}
	if p.APIEndpoint != "https://example.test/dns-api/v1" {
		t.Errorf("endpoint = %q", p.APIEndpoint)
	}
}

func TestUnmarshalCaddyfileMissingToken(t *testing.T) {
	d := caddyfile.NewTestDispenser(`domaintechnik`)
	p := &Provider{Provider: new(domaintechnik.Provider)}
	if err := p.UnmarshalCaddyfile(d); err == nil {
		t.Fatal("expected error")
	}
}

func TestUnmarshalCaddyfileUnknownDirective(t *testing.T) {
	d := caddyfile.NewTestDispenser(`domaintechnik {
		nope x
	}`)
	p := &Provider{Provider: new(domaintechnik.Provider)}
	if err := p.UnmarshalCaddyfile(d); err == nil {
		t.Fatal("expected error")
	}
}
