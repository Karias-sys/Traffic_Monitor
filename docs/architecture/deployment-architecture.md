# Deployment Architecture

## Deployment Strategy

**Frontend Deployment:**
- **Platform:** Embedded in Go binary (no separate deployment)
- **Build Command:** `go build -ldflags="-X main.version=$(VERSION)" ./cmd/netwatch`
- **Output Directory:** Single binary executable
- **CDN/Edge:** Static assets served directly from embedded filesystem

**Backend Deployment:**
- **Platform:** Linux bare metal or VMs with CAP_NET_RAW capability
- **Build Command:** `make build-release` (cross-compilation for multiple architectures)
- **Deployment Method:** Direct binary placement with systemd service management

## CI/CD Pipeline

```yaml
name: Netwatch CI/CD
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Run tests
        run: |
          make test-unit
          make test-integration
      - name: Run linter
        run: make lint
  
  build:
    needs: test
    runs-on: ubuntu-latest
    strategy:
      matrix:
        goos: [linux]
        goarch: [amd64, arm64]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      - name: Build binary
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
        run: make build-release
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: netwatch-${{ matrix.goos }}-${{ matrix.goarch }}
          path: bin/netwatch*
  
  release:
    if: startsWith(github.ref, 'refs/tags/')
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: bin/*
```

## Environments

| Environment | Frontend URL | Backend URL | Purpose |
|-------------|--------------|-------------|---------|
| Development | http://localhost:8080 | http://localhost:8080/api/v1 | Local development |
| Staging | http://staging-netwatch.local | http://staging-netwatch.local/api/v1 | Pre-production testing |
| Production | http://netwatch.local | http://netwatch.local/api/v1 | Live environment |
