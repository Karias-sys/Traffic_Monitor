# Testing Strategy

## Testing Pyramid

```text
                E2E Tests (10)
               /              \
          Integration Tests (50)
             /                    \
        Frontend Unit (100)  Backend Unit (200)
```

## Test Organization

### Frontend Tests

```text
tests/e2e/playwright/
├── dashboard.spec.js          # Real-time dashboard functionality
├── flow-table.spec.js         # Flow filtering and sorting
├── websocket.spec.js          # Real-time updates
└── keyboard-navigation.spec.js # SOCC analyst workflows
```

### Backend Tests

```text
tests/unit/
├── capture/
│   ├── engine_test.go         # Packet capture functionality
│   ├── parser_test.go         # Protocol parsing
│   └── ring_buffer_test.go    # Memory management
├── flow/
│   ├── aggregator_test.go     # Flow aggregation logic
│   ├── storage_test.go        # In-memory storage
│   └── indexer_test.go        # Query performance
└── api/
    ├── handlers_test.go       # HTTP endpoints
    └── websocket_test.go      # Real-time streaming
```

### E2E Tests

```text
tests/e2e/
├── full-workflow.spec.js      # End-to-end monitoring workflow
├── performance.spec.js        # Load testing with mock packets
└── failure-scenarios.spec.js  # Network interface failures
```
