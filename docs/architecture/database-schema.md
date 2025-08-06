# Database Schema

**In-Memory Data Structures (No Traditional Database)**

Netwatch uses specialized in-memory data structures rather than traditional database schemas to achieve the performance requirements. This section documents the internal storage organization and indexing strategies.

## Primary Flow Storage Structure

```go
// FlowTable - Primary storage container
type FlowTable struct {
    // Primary hash table for O(1) flow lookup
    flows       sync.Map                    // map[string]*NetworkFlow
    
    // Secondary indexes for query optimization  
    ipIndex     *IPRangeIndex              // Source/destination IP lookup
    portIndex   *PortRangeIndex            // Port-based filtering
    timeIndex   *TimeSeriesIndex           // Time-range queries
    protocolIdx map[uint8][]*NetworkFlow   // Protocol filtering
    
    // Memory management
    flowCount   int64                       // Current flow count
    memoryUsage int64                       // Estimated memory usage
    oldestFlow  time.Time                   // Aging boundary
    
    // Statistics
    totalFlows  uint64                      // Cumulative flow count
    lastCleanup time.Time                   // Last cleanup operation
}

// NetworkFlow storage structure optimized for memory efficiency
type NetworkFlow struct {
    // 5-tuple key (24 bytes total)
    SrcIP    [16]byte    // IPv6-compatible, 16 bytes
    DstIP    [16]byte    // IPv6-compatible, 16 bytes  
    SrcPort  uint16      // 2 bytes
    DstPort  uint16      // 2 bytes
    Protocol uint8       // 1 byte
    _        [7]byte     // Padding for alignment
    
    // Flow statistics (32 bytes)
    Bytes      uint64    // Total bytes
    Packets    uint64    // Total packets  
    FirstSeen  int64     // Unix timestamp nanoseconds
    LastSeen   int64     // Unix timestamp nanoseconds
    
    // Flow state (8 bytes)
    Status     uint8     // Active/Idle/Closed
    Direction  uint8     // Bidirectional tracking
    _          [6]byte   // Reserved for future use
    
    // Total: 64 bytes per flow for cache line efficiency
}
```

## IP Range Index Structure

```go
// IPRangeIndex - Optimized for CIDR queries and IP range lookups
type IPRangeIndex struct {
    // Separate trees for IPv4 and IPv6 for efficiency
    ipv4Tree *RadixTree[uint32, []*NetworkFlow]
    ipv6Tree *RadixTree[[16]byte, []*NetworkFlow]
    
    // Reverse lookup for flow removal
    flowToNodes map[string][]RadixNode
    
    mutex sync.RWMutex
}

// RadixTree implementation for efficient IP prefix matching
type RadixTree[K comparable, V any] struct {
    root *RadixNode[K, V]
    size int
}

type RadixNode[K comparable, V any] struct {
    key      K
    mask     int        // CIDR prefix length
    value    V          // Slice of flows matching this prefix
    children []*RadixNode[K, V]
}
```

## Time Series Index Structure  

```go
// TimeSeriesIndex - Optimized for time-range queries
type TimeSeriesIndex struct {
    // Time-based buckets for efficient range queries
    buckets map[int64]*TimeBucket  // Unix timestamp / bucket_size
    
    // Ring buffer for 60-minute retention
    bucketRing []int64             // Circular buffer of bucket timestamps
    ringIndex  int                 // Current position in ring
    
    bucketSize time.Duration       // Default: 1 minute buckets
    retention  time.Duration       // 60 minutes
    
    mutex sync.RWMutex
}

type TimeBucket struct {
    timestamp time.Time
    flows     []*NetworkFlow       // Flows active in this time bucket
    flowCount int
    totalBytes uint64
    totalPackets uint64
}
```

## System Metrics Storage

```go
// MetricsStorage - Rolling window for historical metrics
type MetricsStorage struct {
    // Circular buffer for time-series metrics
    metrics    []SystemMetrics     // Fixed size ring buffer
    capacity   int                 // Buffer size (e.g., 3600 for 1-hour at 1s intervals)
    head       int                 // Current write position
    size       int                 // Current number of metrics stored
    
    // Aggregated statistics  
    current    SystemMetrics       // Most recent metrics
    
    mutex      sync.RWMutex
}

// SystemMetrics optimized for minimal memory usage
type SystemMetrics struct {
    Timestamp        int64    // Unix timestamp nanoseconds (8 bytes)
    CPUUsage         uint16   // CPU % * 100 for precision (2 bytes)
    MemoryUsage      uint32   // Memory in KB (4 bytes)  
    ActiveFlows      uint32   // Flow count (4 bytes)
    PacketsPerSecond uint32   // Packet rate (4 bytes)
    BytesPerSecond   uint64   // Bandwidth (8 bytes)
    DroppedPackets   uint32   // Cumulative drops (4 bytes)
    CaptureErrors    uint32   // Cumulative errors (4 bytes)
    // Total: 40 bytes per metric point
}
```
