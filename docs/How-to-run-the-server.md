WhoKnows Variations

# How to Run the Server

---

## Development (with auto-reload via Docker + Air)

This setup provides hot reload (automatic rebuild on file changes) and uses Docker volumes for persistent database and build cache.

### 1. Prerequisites
- Install Docker and Docker Compose  
  [https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/)

### 2. Clone the repository
```bash
git clone https://github.com/5-uptime-gang/whoknows_variations.git
cd whoknows_variations
```

### 2.1 Create .air.toml file

- Create an .air.toml file in the root of the project:
```bash
# .air.toml — local dev config (add to .gitignore)
root = "."
tmp_dir = "/tmp"

[build]
cmd = "go build -buildvcs=false -o /tmp/main ./cmd"
bin = "/tmp/main"
include_ext = ["go", "tpl", "tmpl", "html", "js", "css"]
exclude_dir = ["vendor"]
delay = 1000

[log]
time = true
```

### 3. Run the development environment
```bash
docker compose -f docker-compose.dev.yml up --build
```

**What happens:**
- Builds the development image using `Dockerfile.dev`.
- Runs the app with Air for live reload.
- Creates persistent Docker volumes for:
  - Database → `whoknows_dev_data`
  - Go build cache → `whoknows_dev_cache`
- Runs as a non-root user (`appuser`).

Access the app at:  
[http://localhost:8080](http://localhost:8080)

### 4. Edit and develop
- Any change to Go, HTML, JS, or CSS files triggers an automatic rebuild inside the container.
- Logs are streamed live in the terminal.

### 5. Stop the environment
```bash
docker compose -f docker-compose.dev.yml down
```

The database and cache volumes are preserved unless you use the `-v` flag.

---

## Production (optimized build)

This setup builds a minimal, non-root production image and runs it persistently.

### 1. Build and run the container
From the project root:
```bash
sudo docker-compose up -d --build
```

**What happens:**
- Builds the application using the `Dockerfile`.
- Runs it as a lightweight container (`whoknows_variations`).
- Mounts a persistent database volume (`appdata` → `/usr/src/app/data`).

Access the app at:  
[http://localhost:8080](http://localhost:8080)

### 2. View logs
```bash
docker-compose logs -f
```

### 3. Stop or restart
```bash
sudo docker-compose down        # stop the container
sudo docker-compose restart     # restart quickly
```

---

## Cleanup

If you want to reset everything, including database volumes:

**Development:**
```bash
docker compose -f docker-compose.dev.yml down -v
```

**Production:**
```bash
sudo docker-compose down -v
```

---

## Summary

| Environment   | Command                                           | Description                              |
|----------------|---------------------------------------------------|------------------------------------------|
| Development    | `docker compose -f docker-compose.dev.yml up --build` | Hot reload with persistent volumes       |
| Production     | `sudo docker-compose up -d --build`                   | Optimized static build, persistent DB    |
| Stop dev       | `docker compose -f docker-compose.dev.yml down`  | Stops container, keeps DB                |
| Reset dev      | `docker compose -f docker-compose.dev.yml down -v` | Stops and deletes DB volume              |



How to run the server: 
// Deprecated 
1. Install Go (https://go.dev/doc/install, version 1.25.0.)
2. Clone the repository: git clone https://github.com/5-uptime-gang/whoknows_variations.git 
3. Navigate to the backend folder: cd /go_rewrite/src/backend 
4. Run the server: go run . (or go run main.go to specify the file)
4.1 Alternatively you can install air for development auto reloading.
    - go install github.com/air-verse/air@latest
    - go init
4.2 run "Air" to run the project with auto reload


Documentation: This project follows the Go REST API with Gin tutorial (https://go.dev/doc/tutorial/web-service-gin).
