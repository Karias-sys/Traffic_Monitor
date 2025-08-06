# Epic 3: Real-time Dashboard & Visualization

**Epic Goal:** Deliver the primary user interface with real-time network visibility through live dashboard, charts, and WebSocket-powered updates. This epic provides network administrators with immediate visual insight into bandwidth utilization, protocol distribution, and top network talkers, fulfilling the core value proposition of instant network intelligence.

## Story 3.1: Static Web Server & Asset Management
As a **network administrator**,  
I want **a web interface served from the Netwatch binary**,  
so that **I can access the monitoring dashboard through any browser without additional infrastructure**.

**Acceptance Criteria:**
1. HTTP server serves static web assets (HTML, CSS, JavaScript) embedded in the Go binary
2. Web assets embedded using Go's embed functionality for single-binary deployment
3. Matrix-themed CSS framework with dark theme and cybersecurity color palette (green, amber, red, cyan)
4. Responsive design supporting desktop and tablet viewports
5. Basic HTML structure with placeholder areas for dashboard components
6. Static asset caching headers for optimal browser performance
7. Development mode supports live asset reloading for rapid iteration

## Story 3.2: WebSocket Real-time Communication
As a **system**,  
I want **WebSocket connections for real-time data streaming between backend and frontend**,  
so that **dashboard updates reflect network changes within sub-second latency**.

**Acceptance Criteria:**
1. WebSocket endpoint (/ws) with proper HTTP upgrade handling
2. JSON-based message protocol for sending flow metrics and statistics
3. Client-side WebSocket connection with automatic reconnection logic
4. Connection management handles multiple concurrent dashboard users
5. Message queuing and throttling to prevent WebSocket flooding under high traffic
6. Heartbeat/ping mechanism to detect and handle connection failures
7. Graceful degradation when WebSocket connections are unavailable

## Story 3.3: Live Bandwidth Charts & Metrics Display
As a **network administrator**,  
I want **real-time bandwidth visualization with live updating charts**,  
so that **I can immediately see current network utilization trends and spikes**.

**Acceptance Criteria:**
1. Real-time line chart displaying bandwidth over time (bps and pps)
2. Chart.js integration with streaming data updates via WebSocket
3. Time window controls (1min, 5min, 15min, 60min views) with smooth transitions
4. Y-axis auto-scaling based on current traffic levels with manual override option
5. Visual indicators for data freshness and WebSocket connection status
6. Chart performance optimization for smooth rendering during high update rates
7. Keyboard shortcuts for chart navigation and time window switching

## Story 3.4: Protocol Breakdown & Top Talkers Visualization
As a **network administrator**,  
I want **visual breakdown of protocol distribution and top bandwidth consumers**,  
so that **I can quickly identify what types of traffic dominate my network and which hosts are heaviest users**.

**Acceptance Criteria:**
1. Protocol pie chart or bar chart showing traffic distribution (TCP, UDP, ICMP, others)
2. Top talkers table with auto-refresh showing source/destination IPs and byte counts
3. Sortable columns in top talkers (by bytes, packets, duration) with visual sorting indicators
4. Color-coded protocol categories using matrix theme colors for immediate recognition
5. Percentage calculations for protocol distribution relative to total traffic
6. Configurable refresh intervals for top talkers list with manual refresh option
7. Click-through navigation from top talkers to detailed flow information (preparation for Epic 4)
