# Domaintechnik module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It solves the ACME **DNS-01** challenge using Domaintechnik’s `acme.present` / `acme.cleanup` API (wildcards, or when HTTP-01 is not possible). Caddy already supports HTTP-01 without this plugin.

This is not a general-purpose DNS manager. The HTTP client lives in [`github.com/schuellerf/libdns-domaintechnik`](https://github.com/schuellerf/libdns-domaintechnik).

## Caddy module name

```
dns.providers.domaintechnik
```

(The GitHub repo name is `caddy-dns-domaintechnik`; the Caddyfile still uses `domaintechnik`.)

## Building Caddy

```bash
xcaddy build --with github.com/schuellerf/caddy-dns-domaintechnik
```

If the libdns module is not tagged yet, also pass it explicitly:

```bash
xcaddy build \
  --with github.com/schuellerf/caddy-dns-domaintechnik \
  --with github.com/schuellerf/libdns-domaintechnik
```

From local checkouts:

```bash
xcaddy build \
  --with github.com/schuellerf/caddy-dns-domaintechnik=/path/to/caddy-dns \
  --with github.com/schuellerf/libdns-domaintechnik=/path/to/libdns
```

## Testing without creating a DNS record

None of these talk to Domaintechnik’s `acme.present` / `acme.cleanup` endpoints. They do not issue a certificate.

1. **Unit tests** (httptest mock API):

   ```bash
   go test ./...                    # in libdns-domaintechnik
   go test ./...                    # in caddy-dns-domaintechnik
   ```

2. **Compile Caddy with the plugin** (xcaddy “dry-run”):

   ```bash
   xcaddy build --with github.com/schuellerf/caddy-dns-domaintechnik
   ./caddy list-modules | grep domaintechnik
   # expect: dns.providers.domaintechnik
   ```

3. **Validate a Caddyfile** (loads the module, no ACME client):

   ```bash
   cat > /tmp/Caddyfile <<'EOF'
   {
   	acme_dns domaintechnik dummy-token
   }
   :80 {
   	respond "ok"
   }
   EOF
   ./caddy validate --config /tmp/Caddyfile
   ```

   `caddy run` is what starts certificate issuance and would create `_acme-challenge` records. Skip that until you want a real cert.

Optional live check that still **does not** change DNS: Domaintechnik `ping` (curl or `Provider.Ping`). That only verifies the bearer token.

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
