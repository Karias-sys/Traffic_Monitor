# Unified Project Structure

```plaintext
netwatch/
├── .github/                    # CI/CD workflows
│   └── workflows/
│       ├── ci.yaml
│       ├── release.yaml
│       └── security.yaml
├── cmd/                        # Application entry points
│   └── netwatch/
│       └── main.go            # Main application entry
├── internal/                   # Private application code
│   ├── capture/               # Packet capture engine
│   │   ├── engine.go
│   │   ├── ring_buffer.go
│   │   ├── packet_parser.go
│   │   └── interface_manager.go
│   ├── flow/                  # Flow processing and storage
│   │   ├── aggregator.go
│   │   ├── storage.go
│   │   ├── indexer.go
│   │   ├── cleaner.go
│   │   └── repository.go
│   ├── api/                   # HTTP API layer
│   │   ├── router.go
│   │   ├── handlers/
│   │   │   ├── flows.go
│   │   │   ├── metrics.go
│   │   │   ├── health.go
│   │   │   └── websocket.go
│   │   └── middleware/
│   │       ├── cors.go
│   │       ├── auth.go
│   │       ├── rate_limit.go
│   │       └── logging.go
│   ├── websocket/             # Real-time streaming
│   │   ├── manager.go
│   │   ├── broadcaster.go
│   │   └── client.go
│   ├── metrics/               # System monitoring
│   │   ├── collector.go
│   │   ├── storage.go
│   │   └── calculator.go
│   └── config/                # Configuration management
│       ├── config.go
│       ├── validation.go
│       └── defaults.go
├── pkg/                       # Public library code
│   ├── types/                 # Shared data types
│   │   ├── flow.go
│   │   ├── metrics.go
│   │   └── api.go
│   ├── logger/                # Logging utilities
│   │   ├── logger.go
│   │   └── structured.go
│   └── utils/                 # Common utilities
│       ├── ip.go              # IP address utilities
│       ├── time.go            # Time formatting
│       └── validation.go      # Input validation
├── web/                       # Frontend assets (embedded)
│   ├── assets/
│   │   ├── js/
│   │   │   ├── components/    # UI components
│   │   │   │   ├── Dashboard/
│   │   │   │   │   ├── BandwidthChart.js
│   │   │   │   │   ├── ProtocolBreakdown.js
│   │   │   │   │   ├── TopTalkers.js
│   │   │   │   │   └── SystemStatus.js
│   │   │   │   ├── FlowTable/
│   │   │   │   │   ├── FlowTable.js
│   │   │   │   │   ├── FlowFilters.js
│   │   │   │   │   ├── FlowRow.js
│   │   │   │   │   └── FlowExport.js
│   │   │   │   ├── Common/
│   │   │   │   │   ├── TabNavigation.js
│   │   │   │   │   ├── WebSocketClient.js
│   │   │   │   │   ├── LoadingSpinner.js
│   │   │   │   │   └── StatusIndicator.js
│   │   │   │   └── Health/
│   │   │   │       ├── MetricsGrid.js
│   │   │   │       ├── PerformanceChart.js
│   │   │   │       └── AlertPanel.js
│   │   │   ├── services/      # API services
│   │   │   │   ├── apiClient.js
│   │   │   │   ├── websocketService.js
│   │   │   │   └── dataFormatter.js
│   │   │   ├── utils/         # Frontend utilities
│   │   │   │   ├── timeFormatter.js
│   │   │   │   ├── ipValidator.js
│   │   │   │   └── keyboardHandler.js
│   │   │   └── app.js         # Main application
│   │   ├── css/
│   │   │   ├── matrix-theme.css
│   │   │   ├── components.css
│   │   │   └── dashboard.css
│   │   └── vendor/            # Third-party libraries
│   │       ├── chart.js       # Chart.js for visualizations
│   │       └── normalize.css  # CSS normalization
│   └── index.html             # SPA shell
├── tests/                     # Test files
│   ├── integration/           # Integration tests
│   │   ├── api_test.go
│   │   ├── capture_test.go
│   │   └── websocket_test.go
│   ├── unit/                  # Unit tests
│   │   ├── flow/
│   │   │   ├── aggregator_test.go
│   │   │   ├── storage_test.go
│   │   │   └── repository_test.go
│   │   ├── capture/
│   │   │   ├── engine_test.go
│   │   │   └── parser_test.go
│   │   └── api/
│   │       ├── handlers_test.go
│   │       └── middleware_test.go
│   ├── mocks/                 # Mock implementations
│   │   ├── packet_generator.go
│   │   ├── mock_storage.go
│   │   └── mock_interfaces.go
│   └── e2e/                   # End-to-end tests
│       ├── playwright/
│       │   ├── dashboard.spec.js
│       │   ├── flow-table.spec.js
│       │   └── websocket.spec.js
│       └── fixtures/
│           ├── test-packets.pcap
│           └── flow-data.json
├── scripts/                   # Build and utility scripts
│   ├── build.sh              # Cross-platform build script
│   ├── test.sh               # Test runner
│   ├── generate-mocks.sh     # Mock generation
│   └── release.sh            # Release packaging
├── configs/                   # Configuration files
│   ├── netwatch.yaml.example # Example configuration
│   ├── docker/               # Docker configurations
│   │   ├── Dockerfile
│   │   ├── Dockerfile.dev
│   │   └── docker-compose.yml
│   └── systemd/              # Systemd service files
│       └── netwatch.service
├── docs/                     # Documentation
│   ├── prd.md
│   ├── front-end-spec.md
│   ├── architecture.md       # This document
│   ├── api/                  # API documentation
│   │   ├── openapi.yaml
│   │   └── websocket.md
│   ├── deployment/           # Deployment guides
│   │   ├── linux.md
│   │   ├── docker.md
│   │   └── systemd.md
│   └── development/          # Development guides
│       ├── setup.md
│       ├── testing.md
│       └── contributing.md
├── .env.example              # Environment template
├── .gitignore
├── .dockerignore
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── Makefile                  # Build automation
└── README.md
```
