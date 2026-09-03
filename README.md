# Domaintechnik module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It solves the ACME **DNS-01** challenge using Domaintechnik’s `acme.present` / `acme.cleanup` API (wildcards, or when HTTP-01 is not possible). Caddy already supports HTTP-01 without this plugin.

This is not a general-purpose DNS manager. The HTTP client lives in [`github.com/schuellerf/libdns-domaintechnik`](https://github.com/schuellerf/libdns-domaintechnik).

## Caddy module name

```
dns.providers.domaintechnik
```

(The GitHub repo name is `caddy-dns-domaintechnik`; the Caddyfile still uses `domaintechnik`.)

## Building Caddy

Once both repos are on GitHub, build with xcaddy from the module path:

```bash
xcaddy build --with github.com/schuellerf/caddy-dns-domaintechnik
```

Go will pull [`github.com/schuellerf/libdns-domaintechnik`](https://github.com/schuellerf/libdns-domaintechnik) as a dependency. Publish the libdns repo first (or in the same step), or xcaddy cannot resolve it.

From local checkouts before push:

```bash
xcaddy build \
  --with github.com/schuellerf/caddy-dns-domaintechnik=/path/to/domaintechnik/caddy-dns \
  --with github.com/schuellerf/libdns-domaintechnik=/path/to/domaintechnik/libdns
```

## Token

Caddy only reads a string from config. Prefer an environment variable on the **Caddy process**:

- systemd: `Environment=DT_TOKEN=dt1_…` or `EnvironmentFile=/etc/caddy/env` (mode `600`)
- shell: `export DT_TOKEN=dt1_…` then `caddy run`
- Docker: `environment:` / `--env-file`

`{env.DT_TOKEN}` is expanded at runtime. `{$DT_TOKEN}` is substituted when the Caddyfile is adapted. A literal `dt1_…` in the Caddyfile also works if that file is not in git and is readable only by the `caddy` user.

## Caddyfile

Global:

```
{
	acme_dns domaintechnik {env.DT_TOKEN}
}
```

Per site:

```
tls {
	dns domaintechnik {env.DT_TOKEN}
}
```

Block form (optional custom API base URL):

```
tls {
	dns domaintechnik {
		api_token {env.DT_TOKEN}
		api_endpoint https://manage.domaintechnik.at/dns-api/v1
	}
}
```

## JSON

```json
{
	"module": "acme",
	"challenges": {
		"dns": {
			"provider": {
				"name": "domaintechnik",
				"api_token": "{env.DT_TOKEN}"
			}
		}
	}
}
```
