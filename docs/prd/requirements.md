# Requirements

## Functional Requirements

**FR1**: The system shall capture network packets in real-time from a specified network interface using AF_PACKET with TPACKETv3 for efficient Linux packet processing.

**FR2**: The system shall aggregate captured packets into network flows, tracking source/destination IPs, ports, protocols, byte counts, and packet counts.

**FR3**: The system shall provide a real-time traffic dashboard displaying live bandwidth charts (bps/pps) with sub-second updates.

**FR4**: The system shall display protocol breakdown visualization showing traffic distribution across different network protocols.

**FR5**: The system shall maintain a top talkers list with automatic refresh showing highest bandwidth consumers.

**FR6**: The system shall provide a searchable and filterable flow table interface with pagination support for large datasets.

**FR7**: The system shall support sorting flows by bytes, packets, and duration in both ascending and descending order.

**FR8**: The system shall implement WebSocket-based live updates for sub-second metric streaming to connected clients.

**FR9**: The system shall provide automatic WebSocket reconnection handling for reliable real-time updates.

**FR10**: The system shall expose REST API endpoints for flow queries with filtering capabilities (by IP, port, protocol, time range).

**FR11**: The system shall provide REST API access to historical counter data within the 60-minute memory buffer.

**FR12**: The system shall expose system statistics and health status through REST API endpoints.

## Non-Functional Requirements

**NFR1**: The system shall handle up to 1 Gbps sustained network traffic capture with minimal packet loss.

**NFR2**: The system shall maintain less than 5% CPU utilization on the monitoring host under typical load conditions.

**NFR3**: The system shall use less than 1GB of memory for 60 minutes of flow history storage.

**NFR4**: The system shall respond to API queries within 100ms under normal operating conditions.

**NFR5**: The system shall maintain 99.9% uptime with continuous packet capture capability.

**NFR6**: The system shall bind to localhost by default for security, with optional token authentication.

**NFR7**: The system shall deploy as a single binary with no external dependencies or complex setup requirements.

**NFR8**: The system shall be compatible with x86_64 and ARM64 architectures on modern Linux distributions.

**NFR9**: The system shall maintain flow data only in memory with automatic eviction after 60 minutes (no persistent storage in MVP).

**NFR10**: The system shall support graceful shutdown with proper cleanup of network capture resources.
