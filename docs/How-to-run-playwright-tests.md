# How to Run Playwright Tests

Use these steps to execute the browser end-to-end test suite in `playwright_e2e_test/`.

## Prerequisites
- App is running and reachable (default: [http://localhost:8080](http://localhost:8080)); e.g., `make dev` in the project root.
- Node.js and npm installed.

## Install test dependencies (one time or after updates)
```bash
cd playwright_e2e_test
npm install
npx playwright install
```

## Run the suite
- Headless (CI-friendly): `npm test`
- Headed (visible browser): `npm run test:headed`

## Pointing at a different host/port
Set `BASE_URL` before running:
```bash
BASE_URL=http://localhost:3000 npm test
```

## Artifacts and traces
- On failure, Playwright keeps a trace (`trace.zip`) per test run (see `playwright-report/` by default).
- Use `npx playwright show-trace <path-to-trace.zip>` to inspect failures.
