# Epic 2: Flow Processing & Data Management

**Epic Goal:** Transform raw packet capture into structured network flow data with efficient in-memory storage and management. This epic delivers the core intelligence layer that aggregates packets into flows, manages memory usage within the 60-minute constraint, and provides the data foundation for all visualization and analysis features.

## Story 2.1: Flow Data Structures & Key Generation
As a **system**,  
I want **efficient data structures for representing network flows with unique identification**,  
so that **I can aggregate packets into flows and track connection metadata accurately**.

**Acceptance Criteria:**
1. Flow struct defined with source/destination IPs, ports, protocol, byte/packet counters, and timestamps
2. Flow key generation implemented using 5-tuple (src_ip, dst_ip, src_port, dst_port, protocol)
3. Bidirectional flow support with consistent key generation regardless of packet direction
4. Efficient serialization/deserialization for flow data structures
5. Flow state tracking (active, idle, closed) with appropriate state transitions
6. Memory-efficient storage with appropriate data types for counters and addresses
7. Unit tests validating flow key generation and data structure operations

## Story 2.2: Packet-to-Flow Aggregation Engine
As a **system**,  
I want **real-time aggregation of captured packets into network flows**,  
so that **I can provide flow-based network analysis instead of individual packet data**.

**Acceptance Criteria:**
1. Packet processing pipeline extracts 5-tuple and payload information from captured packets
2. Flow lookup mechanism efficiently finds existing flows or creates new ones
3. Packet statistics aggregated into flow counters (bytes, packets, first/last seen timestamps)
4. Protocol-specific handling for TCP (connection state), UDP (stateless), and other protocols
5. Performance optimization for high-throughput packet processing (target: 1 Gbps)
6. Error handling for malformed packets without disrupting flow processing
7. Flow creation and update operations are thread-safe for concurrent access

## Story 2.3: In-Memory Flow Storage & Indexing
As a **system**,  
I want **efficient in-memory storage and indexing of active network flows**,  
so that **I can quickly retrieve and analyze flow data for real-time monitoring**.

**Acceptance Criteria:**
1. In-memory flow table implemented using concurrent-safe data structures
2. Primary index by flow key for O(1) flow lookup and updates
3. Secondary indexes for common query patterns (by IP, by port, by protocol, by bytes/packets)
4. Time-based indexing for efficient retrieval of flows within time ranges
5. Memory usage monitoring and reporting for flow storage
6. Lock-free or minimal-locking design for high-performance concurrent access
7. Comprehensive testing of concurrent flow operations and index consistency

## Story 2.4: Flow Lifecycle & Memory Management
As a **network administrator**,  
I want **automatic flow cleanup and memory management within the 60-minute constraint**,  
so that **the system maintains stable memory usage during extended monitoring periods**.

**Acceptance Criteria:**
1. Flow aging mechanism automatically removes flows older than 60 minutes
2. Idle flow detection with configurable timeout for inactive connections
3. Memory pressure handling with graceful flow eviction when approaching limits
4. Flow cleanup preserves recent high-traffic flows over older low-traffic flows
5. Memory usage statistics exposed through metrics endpoints
6. Background cleanup process with minimal impact on packet processing performance
7. Configuration options for memory limits and flow retention policies
