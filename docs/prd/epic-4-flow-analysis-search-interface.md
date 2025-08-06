# Epic 4: Flow Analysis & Search Interface

**Epic Goal:** Provide comprehensive flow-level analysis capabilities with advanced search, filtering, and drill-down functionality for detailed network investigation. This epic enables network administrators to move beyond high-level dashboard metrics to investigate specific network flows, troubleshoot issues, and perform detailed traffic analysis.

## Story 4.1: Tabbed Navigation & Flow Table Interface
As a **network administrator**,  
I want **tabbed navigation between Dashboard and Flow Table views**,  
so that **I can switch between high-level monitoring and detailed flow analysis without losing context**.

**Acceptance Criteria:**
1. Tab-based navigation implemented with Dashboard and Flow Table as primary tabs
2. Browser-standard tab behavior with keyboard navigation (Tab/Shift+Tab, Ctrl+1/Ctrl+2)
3. Tab state persistence across browser refreshes using localStorage
4. Loading states displayed when switching between tabs with different data requirements
5. Active tab styling using matrix theme with appropriate visual indicators
6. Smooth transition animations between tab content areas
7. URL routing supports deep linking to specific tabs (/dashboard, /flows)

## Story 4.2: Comprehensive Flow Data Table
As a **network administrator**,  
I want **a detailed table showing all network flows with essential metadata**,  
so that **I can examine individual connections and identify specific network patterns**.

**Acceptance Criteria:**
1. Flow table displays: Source IP, Destination IP, Source Port, Destination Port, Protocol, Bytes, Packets, Duration, First Seen, Last Seen
2. Responsive table design with horizontal scrolling on smaller screens
3. Real-time flow updates via WebSocket without disrupting user interaction (scroll position, selections)
4. Row highlighting on hover with matrix-themed styling (green/cyan accent colors)
5. Monospace font for IP addresses and port numbers for consistent alignment
6. Flow status indicators (active, idle, closed) with appropriate visual encoding
7. Infinite scrolling or pagination for handling large flow datasets efficiently

## Story 4.3: Advanced Filtering & Search Capabilities
As a **network administrator**,  
I want **powerful filtering and search functionality across all flow attributes**,  
so that **I can quickly isolate flows of interest during network investigation and troubleshooting**.

**Acceptance Criteria:**
1. Search bar with intelligent parsing for IP addresses, ports, and protocol names
2. Filter controls for common patterns: source/destination IP ranges, port ranges, protocol types
3. Time-based filtering with presets (last 5min, 15min, 1hr) and custom time range selection
4. Byte/packet threshold filtering to focus on high-volume or low-volume flows
5. Real-time search results with debounced input to prevent performance issues
6. Filter combination logic (AND/OR operations) with visual query builder
7. Saved filter presets for common investigation patterns
8. Clear all filters functionality with keyboard shortcut (Esc key)

## Story 4.4: Flow Sorting & Data Export
As a **network administrator**,  
I want **flexible sorting options and data export capabilities**,  
so that **I can organize flow data for analysis and share findings with team members**.

**Acceptance Criteria:**
1. Column-based sorting for all flow table columns (ascending/descending) with visual sort indicators
2. Multi-column sorting capability with priority indicators (1st sort, 2nd sort, etc.)
3. Default sorting by "Last Seen" timestamp for most recent activity first
4. Sort state persistence across page refreshes and tab switches
5. Export filtered flow data to CSV format with proper header row
6. Export functionality respects current filters and sort order
7. Keyboard shortcuts for common sorting operations (click + Shift for multi-sort)
8. Performance optimization for sorting large flow datasets without blocking UI
