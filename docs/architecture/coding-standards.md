# Coding Standards

## Critical Fullstack Rules

- **Performance First:** All code changes must maintain <5% CPU usage target - benchmark before merging
- **Memory Management:** Flow storage must never exceed 1GB - implement proper cleanup and monitoring
- **Error Handling:** All packet processing errors must be non-blocking - system continues operation
- **Real-time Priority:** WebSocket updates take precedence over HTTP API responses during high load
- **Security Validation:** All IP address and port inputs must be validated before processing
- **Logging Standards:** Use structured logging with consistent field names across frontend and backend

## Naming Conventions

| Element | Frontend | Backend | Example |
|---------|----------|---------|---------|
| Components | PascalCase | - | `BandwidthChart.js` |
| Functions | camelCase | camelCase | `updateFlowData()`, `ProcessPacket()` |
| API Endpoints | - | kebab-case | `/api/v1/top-talkers` |
| Go Packages | - | lowercase | `internal/capture` |
| Constants | UPPER_CASE | UPPER_CASE | `MAX_FLOWS`, `DEFAULT_PORT` |
