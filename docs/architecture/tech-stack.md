# Tech Stack

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| Backend Language | Go | 1.21+ | High-performance packet processing and web server | Superior concurrency, memory safety, single-binary compilation, and excellent networking libraries |
| Backend Framework | Go Standard Library + Gorilla WebSocket | stdlib + v1.5.0 | HTTP server and WebSocket handling | Minimal dependencies, maximum performance, proven stability for network applications |
| Frontend Language | JavaScript (Vanilla) | ES2020+ | Browser-based dashboard interface | Eliminates framework overhead, reduces bundle size, optimal for specialized real-time interface |
| Frontend Framework | None (Vanilla Components) | N/A | Component organization without framework | Maximum performance for real-time updates, minimal complexity for focused use case |
| UI Component Library | Custom Components + Chart.js | Chart.js v4.4.0 | Real-time data visualization | Lightweight charting for bandwidth/protocol visualization, custom matrix-themed components |
| State Management | Native JavaScript Classes | ES6+ | Application state coordination | Simple state management sufficient for dashboard app, avoids framework complexity |
| API Style | REST + WebSocket | HTTP/1.1, WebSocket | Data access and real-time streaming | RESTful APIs for programmatic access, WebSocket for sub-second dashboard updates |
| Database | In-Memory (Go Maps + Sync) | Native Go | Flow data storage and indexing | 60-minute retention requirement, maximum query speed, eliminates database overhead |
| Cache | In-Memory LRU | Custom implementation | Flow aggregation optimization | Built-in caching for flow lookup performance, no external cache needed |
| File Storage | Embedded Assets (Go embed) | Go 1.16+ embed | Static web asset delivery | Single-binary deployment, eliminates external file dependencies |
| Authentication | Optional Token-based | Custom JWT-like | API access control | Optional security layer, localhost-first design with configurable auth |
| Frontend Testing | Jest + jsdom | Jest v29.0+ | Unit testing for JS components | Lightweight testing without browser overhead, sufficient for component logic |
| Backend Testing | Go testing + testify | Go stdlib + v1.8.0 | Unit and integration testing | Native Go testing with assertion library, mock packet generation |
| E2E Testing | Playwright | v1.40+ | Full application workflow testing | Comprehensive testing of WebSocket updates and dashboard interactions |
| Build Tool | Make + Go build | Native tooling | Build automation and cross-compilation | Simple, reliable build process without additional dependencies |
| Bundler | None (Direct serving) | N/A | Asset delivery optimization | Embedded assets served directly, no bundling complexity needed |
| IaC Tool | None (Single Binary) | N/A | Infrastructure management | Deployment is binary placement, no infrastructure to manage |
| CI/CD | GitHub Actions | Latest | Automated testing and releases | Free for open source, excellent Go support, cross-platform builds |
| Monitoring | Structured Logging + Metrics Endpoint | Go stdlib + slog | Application monitoring | Built-in metrics exposure, structured logging for operational insight |
| Logging | slog (structured logging) | Go 1.21+ | Application logging and debugging | High-performance structured logging, configurable levels |
| CSS Framework | Custom Matrix Theme | CSS3 | Cybersecurity-themed styling | Matrix aesthetic requirements, specialized for SOCC environments |
