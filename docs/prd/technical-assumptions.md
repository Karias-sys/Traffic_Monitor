# Technical Assumptions

## Repository Structure: Monorepo

Single repository containing all project components (backend Go application, web frontend, documentation, configuration) for simplified development, testing, and deployment of this focused single-binary solution.

## Service Architecture

**High-Performance Go Monolith** - Single Go binary containing integrated packet capture engine, flow aggregation service, WebSocket server, REST API, and static file serving. This monolithic approach aligns with the "zero-infrastructure deployment" goal and eliminates network latency between components that would be critical for real-time packet processing.

**CRITICAL RATIONALE**: Microservices would introduce unacceptable latency for real-time packet processing and contradict the single-binary deployment requirement. The performance demands (1 Gbps capture, sub-second updates) require tightly coupled components sharing memory efficiently.

## Testing Requirements

**Unit + Integration Testing Strategy** focusing on:
- Unit tests for core packet processing and flow aggregation logic
- Integration tests for WebSocket streaming and REST API endpoints  
- Performance benchmarking tests to validate 1 Gbps capture capability
- Mock packet injection for reliable testing without requiring actual network traffic
- Automated testing pipeline that can run without raw socket permissions

**CRITICAL RATIONALE**: Given the performance requirements and packet capture complexity, comprehensive testing is essential, but E2E testing would require complex network simulation environments that may not be practical for CI/CD.

## Additional Technical Assumptions and Requests

**Programming Language & Runtime:**
- **Go 1.21+** for high-performance packet processing, excellent concurrency primitives, and single-binary compilation
- **CGO disabled** where possible to ensure truly portable binaries

**Packet Capture Technology:**
- **AF_PACKET with TPACKETv3** on Linux for zero-copy packet capture
- **Fallback to raw sockets** for broader compatibility if needed
- **Ring buffer implementation** for efficient packet queue management

**Web Technology Stack:**
- **Go standard library HTTP server** with WebSocket upgrades for minimal dependencies
- **Vanilla JavaScript frontend** with Chart.js for visualizations to avoid framework overhead
- **Embedded static assets** using Go embed for single-binary deployment

**Performance & Memory Management:**
- **Lock-free data structures** where possible for packet processing hot paths
- **Memory pooling** for packet buffer management to reduce GC pressure
- **Structured logging** with configurable levels to minimize production overhead

**Security & Deployment:**
- **Capability-based permissions** (CAP_NET_RAW) instead of requiring root access
- **TLS support** for encrypted WebSocket/API access when needed
- **Configuration via environment variables** and command-line flags only

**Development & Build:**
- **Make-based build system** for simplicity and cross-platform compatibility
- **Docker containerization** for development environment consistency
- **GitHub Actions CI/CD** for automated testing and binary releases
