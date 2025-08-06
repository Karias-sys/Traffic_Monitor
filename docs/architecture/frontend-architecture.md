# Frontend Architecture

## Component Architecture

### Component Organization

```text
web/assets/
├── js/
│   ├── components/           # Reusable UI components
│   │   ├── Dashboard/
│   │   │   ├── BandwidthChart.js      # Real-time line chart
│   │   │   ├── ProtocolBreakdown.js   # Protocol pie chart  
│   │   │   ├── TopTalkers.js          # Live updating table
│   │   │   └── SystemStatus.js        # Health indicators
│   │   ├── FlowTable/
│   │   │   ├── FlowTable.js           # Main data table
│   │   │   ├── FlowFilters.js         # Search/filter controls
│   │   │   ├── FlowRow.js             # Individual table row
│   │   │   └── FlowExport.js          # CSV export functionality
│   │   ├── Common/
│   │   │   ├── TabNavigation.js       # Tab switching component
│   │   │   ├── WebSocketClient.js     # WebSocket management
│   │   │   ├── LoadingSpinner.js      # Loading states
│   │   │   └── StatusIndicator.js     # Connection/health status
│   │   └── Health/
│   │       ├── MetricsGrid.js         # System metrics display
│   │       ├── PerformanceChart.js    # Historical performance
│   │       └── AlertPanel.js          # System alerts
│   ├── services/             # API and data services
│   │   ├── apiClient.js      # REST API wrapper
│   │   ├── websocketService.js # Real-time data streaming
│   │   └── dataFormatter.js  # Data transformation utilities
│   ├── utils/                # Utility functions
│   │   ├── timeFormatter.js  # Timestamp formatting
│   │   ├── ipValidator.js    # IP address validation
│   │   └── keyboardHandler.js # Keyboard shortcut management
│   └── app.js               # Main application entry point
├── css/
│   ├── matrix-theme.css     # Matrix cybersecurity styling
│   ├── components.css       # Component-specific styles
│   └── dashboard.css        # Layout and grid styles
└── index.html              # Single page application shell
```

### Component Template

```javascript
// Example: BandwidthChart.js - Real-time chart component
class BandwidthChart {
    constructor(containerId, options = {}) {
        this.container = document.getElementById(containerId);
        this.chart = null;
        this.data = [];
        this.options = {
            updateInterval: 1000,
            maxDataPoints: 300, // 5 minutes at 1s intervals
            ...options
        };
        
        this.initChart();
        this.bindEvents();
    }
    
    initChart() {
        const ctx = this.container.querySelector('canvas').getContext('2d');
        this.chart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'Bandwidth (bps)',
                    data: [],
                    borderColor: '#00FF41',
                    backgroundColor: 'rgba(0, 255, 65, 0.1)',
                    tension: 0.4
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: { duration: 0 }, // Disable for real-time
                scales: {
                    x: { type: 'time' },
                    y: { beginAtZero: true }
                }
            }
        });
    }
    
    updateData(newData) {
        const { timestamp, bytesPerSecond } = newData;
        
        // Add new data point
        this.data.push({ x: timestamp, y: bytesPerSecond });
        
        // Maintain sliding window
        if (this.data.length > this.options.maxDataPoints) {
            this.data.shift();
        }
        
        // Update chart
        this.chart.data.datasets[0].data = this.data;
        this.chart.update('none'); // Skip animation for performance
    }
    
    bindEvents() {
        // Keyboard shortcuts
        document.addEventListener('keydown', (e) => {
            if (e.ctrlKey && e.key === 'r') {
                e.preventDefault();
                this.refreshChart();
            }
        });
    }
    
    destroy() {
        if (this.chart) {
            this.chart.destroy();
        }
    }
}
```

## State Management Architecture

### State Structure

```javascript
// AppState.js - Centralized application state
class AppState {
    constructor() {
        this.state = {
            // Connection state
            websocket: {
                connected: false,
                reconnecting: false,
                lastUpdate: null
            },
            
            // Dashboard data
            dashboard: {
                metrics: null,
                flows: [],
                topTalkers: [],
                protocolBreakdown: {}
            },
            
            // Flow table state  
            flowTable: {
                flows: [],
                filters: {
                    searchTerm: '',
                    protocol: 'all',
                    timeRange: '5m',
                    sortBy: 'lastSeen',
                    sortOrder: 'desc'
                },
                pagination: {
                    offset: 0,
                    limit: 100,
                    total: 0
                }
            },
            
            // UI state
            ui: {
                activeTab: 'dashboard',
                loading: false,
                error: null,
                alerts: []
            }
        };
        
        this.subscribers = new Map();
        this.history = []; // For debugging
    }
    
    subscribe(path, callback) {
        if (!this.subscribers.has(path)) {
            this.subscribers.set(path, new Set());
        }
        this.subscribers.get(path).add(callback);
        
        // Return unsubscribe function
        return () => this.subscribers.get(path).delete(callback);
    }
    
    setState(path, value) {
        const oldValue = this.getState(path);
        this.setNestedProperty(this.state, path, value);
        
        // Store history for debugging
        this.history.push({
            timestamp: Date.now(),
            path,
            oldValue,
            newValue: value
        });
        
        // Notify subscribers
        this.notifySubscribers(path, value, oldValue);
    }
    
    getState(path) {
        return this.getNestedProperty(this.state, path);
    }
    
    // Helper methods for nested property access
    setNestedProperty(obj, path, value) {
        const keys = path.split('.');
        const lastKey = keys.pop();
        const target = keys.reduce((o, k) => o[k], obj);
        target[lastKey] = value;
    }
    
    getNestedProperty(obj, path) {
        return path.split('.').reduce((o, k) => o && o[k], obj);
    }
    
    notifySubscribers(path, newValue, oldValue) {
        const subscribers = this.subscribers.get(path);
        if (subscribers) {
            subscribers.forEach(callback => callback(newValue, oldValue));
        }
    }
}

// Global state instance
const appState = new AppState();
```

### State Management Patterns

- **Immutable Updates:** State changes create new objects rather than mutating existing ones for predictable updates
- **Path-based Subscriptions:** Components subscribe to specific state paths to minimize unnecessary re-renders
- **Action-based Updates:** State modifications through explicit actions for debugging and consistency
- **Local Component State:** Transient UI state (hover, focus) managed locally to reduce global state complexity

## Routing Architecture

### Route Organization

```text
Application Routes (Hash-based for SPA):
├── #/dashboard           # Default route - real-time dashboard
├── #/flows              # Flow analysis table
├── #/flows/:flowId      # Individual flow details (future)
├── #/health             # System health monitoring
└── #/settings           # Configuration (future)

Tab Navigation:
- Ctrl+1 → Dashboard
- Ctrl+2 → Flows  
- Ctrl+3 → Health
```
