# Data Models

## NetworkFlow

**Purpose:** Core entity representing aggregated network traffic between two endpoints, tracking connection metadata and traffic statistics for real-time analysis and historical querying.

**Key Attributes:**
- flowKey: string - Unique identifier based on 5-tuple hash for efficient lookup
- srcIP: net.IP - Source IP address in binary format for memory efficiency  
- dstIP: net.IP - Destination IP address in binary format for memory efficiency
- srcPort: uint16 - Source port number for transport layer identification
- dstPort: uint16 - Destination port number for transport layer identification
- protocol: uint8 - IP protocol number (TCP=6, UDP=17, ICMP=1, etc.)
- bytes: uint64 - Total bytes transferred in both directions
- packets: uint64 - Total packet count in both directions
- firstSeen: time.Time - Timestamp of first packet in flow for aging calculations
- lastSeen: time.Time - Timestamp of most recent packet for idle detection
- status: FlowStatus - Current flow state (active, idle, closed)

### TypeScript Interface
```typescript
interface NetworkFlow {
  flowKey: string;
  srcIP: string;
  dstIP: string;
  srcPort: number;
  dstPort: number;
  protocol: number;
  protocolName: string; // Derived field for UI display
  bytes: number;
  packets: number;
  firstSeen: string; // ISO 8601 timestamp
  lastSeen: string;  // ISO 8601 timestamp
  duration: number;  // Calculated field in milliseconds
  status: 'active' | 'idle' | 'closed';
}
```

### Relationships
- One-to-many with PacketSample (for detailed analysis)
- Aggregated into TopTalker summaries
- Referenced by SystemMetrics for bandwidth calculations

## SystemMetrics

**Purpose:** Real-time system performance and network statistics for monitoring system health and capture effectiveness.

**Key Attributes:**
- timestamp: time.Time - Metric collection timestamp for time series data
- cpuUsage: float64 - Current CPU utilization percentage
- memoryUsage: uint64 - Current memory consumption in bytes
- activeFlows: uint32 - Number of currently tracked flows
- packetsPerSecond: uint64 - Current packet capture rate
- bytesPerSecond: uint64 - Current bandwidth utilization
- droppedPackets: uint64 - Cumulative dropped packet count
- captureErrors: uint64 - Cumulative capture error count

### TypeScript Interface
```typescript
interface SystemMetrics {
  timestamp: string;
  cpuUsage: number;
  memoryUsage: number;
  activeFlows: number;
  packetsPerSecond: number;
  bytesPerSecond: number;
  droppedPackets: number;
  captureErrors: number;
  interfaceStatus: 'up' | 'down' | 'error';
}
```

### Relationships
- Aggregates data from NetworkFlow collection
- Used by WebSocket streaming for dashboard updates
- Historical data maintained in rolling window

## TopTalker

**Purpose:** Aggregated view of highest bandwidth consumers for dashboard visualization and rapid network analysis.

**Key Attributes:**
- identifier: string - IP address or flow identifier for ranking
- displayName: string - Human-readable identifier for UI
- totalBytes: uint64 - Aggregated byte count across all flows
- totalPackets: uint64 - Aggregated packet count across all flows
- flowCount: uint32 - Number of active flows for this talker
- rank: uint8 - Current ranking position (1-N)
- percentageOfTotal: float64 - Percentage of total network traffic

### TypeScript Interface
```typescript
interface TopTalker {
  identifier: string;
  displayName: string;
  totalBytes: number;
  totalPackets: number;
  flowCount: number;
  rank: number;
  percentageOfTotal: number;
  flows: NetworkFlow[]; // Associated flows for drill-down
}
```

### Relationships
- Aggregates multiple NetworkFlow entities
- Updated in real-time via WebSocket streams
- Used for click-through navigation to detailed flow analysis
