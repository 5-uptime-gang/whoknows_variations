# Security Breach Notes

## Potential breach points
- Plain HTTP backend plus non-secure auth cookie (`util/cookie.go`) makes session hijack easy over sniffed traffic.
- Login/register lack rate limiting or password policy, inviting brute-force and credential-stuffing.
- `/metrics` is public and leaks user/browser counts and paths; good recon fodder.
- DB access depends on `DATABASE_URL` (often `sslmode=disable`) and env secrets; exposed logs or sniffed links could yield credentials.
- Swagger UI pulls JS/CSS from a CDN without pinning, creating a small supply-chain/XSS risk.

## Quick reflections on mitigation
- Force HTTPS end-to-end (HSTS) and set cookies as `Secure`, `HttpOnly`, `SameSite=Strict`, with short lifetimes and rotation on login.
- Add rate limiting + basic lockout/backoff on auth endpoints; enforce stronger passwords; alert on abnormal failures.
- Protect `/metrics` with auth or IP allowlists; avoid PII in labels; sample if needed.
- Require TLS to Postgres, keep secrets in a vault/runner env vars, rotate DB creds, and enable frequent backups.
- Bundle Swagger assets locally or pin versions with integrity/CSP headers.
