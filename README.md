# Domaintechnik module for Caddy

This package contains a DNS provider module for [Caddy](https://github.com/caddyserver/caddy). It solves the ACME **DNS-01** challenge using Domaintechnik’s `acme.present` / `acme.cleanup` API (wildcards, or when HTTP-01 is not possible). Caddy already supports HTTP-01 without this plugin.

This is not a general-purpose DNS manager. The HTTP client lives in [`github.com/schuellerf/libdns-domaintechnik`](https://github.com/schuellerf/libdns-domaintechnik).

## Caddy module name

```
dns.providers.domaintechnik
```

(The GitHub repo name is `caddy-dns-domaintechnik`; the Caddyfile still uses `domaintechnik`.)

## Building Caddy

You might need to install xcaddy first:
```bash
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
```

```bash
xcaddy build --with github.com/schuellerf/caddy-dns-domaintechnik
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

Global options (no hostname on this block; it applies to every site):

```
{
	acme_dns domaintechnik {env.DT_TOKEN}
}

example.com, *.example.com {
	reverse_proxy localhost:8080
}
```

Per site (hostname in front of the site block, not in front of `tls`):

```
example.com, *.example.com {
	tls {
		dns domaintechnik {env.DT_TOKEN}
	}
	reverse_proxy localhost:8080
}
```

You could set the `DT_TOKEN` environment variable in the systemd service file or
just replace `{env.DT_TOKEN}` with the actual token.
