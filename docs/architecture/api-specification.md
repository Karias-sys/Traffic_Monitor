# API Specification

## REST API Specification

```yaml
openapi: 3.0.0
info:
  title: Netwatch Network Monitoring API
  version: 1.0.0
  description: Real-time network flow monitoring and analysis API for high-performance packet capture and traffic analysis
servers:
  - url: http://localhost:8080/api/v1
    description: Local development server
  - url: https://netwatch.local/api/v1
    description: Production deployment

paths:
  /health:
    get:
      summary: System health check
      description: Returns current system status and capture health
      responses:
        '200':
          description: System operational
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "healthy"
                  timestamp:
                    type: string
                    format: date-time
                  capture_active:
                    type: boolean
                  interface_status:
                    type: string
                    enum: [up, down, error]

  /flows:
    get:
      summary: Query network flows
      description: Retrieve filtered and paginated network flow data
      parameters:
        - name: src_ip
          in: query
          description: Filter by source IP (supports CIDR)
          schema:
            type: string
            example: "192.168.1.0/24"
        - name: dst_ip
          in: query
          description: Filter by destination IP (supports CIDR)
          schema:
            type: string
            example: "10.0.0.1"
        - name: protocol
          in: query
          description: Filter by protocol number or name
          schema:
            type: string
            example: "tcp"
        - name: port
          in: query
          description: Filter by port number (src or dst)
          schema:
            type: integer
            example: 443
        - name: time_start
          in: query
          description: Start time for flow query (ISO 8601)
          schema:
            type: string
            format: date-time
        - name: time_end
          in: query
          description: End time for flow query (ISO 8601)
          schema:
            type: string
            format: date-time
        - name: limit
          in: query
          description: Maximum number of flows to return
          schema:
            type: integer
            default: 100
            maximum: 1000
        - name: offset
          in: query
          description: Number of flows to skip
          schema:
            type: integer
            default: 0
        - name: sort_by
          in: query
          description: Sort field
          schema:
            type: string
            enum: [bytes, packets, duration, first_seen, last_seen]
            default: "last_seen"
        - name: sort_order
          in: query
          description: Sort direction
          schema:
            type: string
            enum: [asc, desc]
            default: "desc"
      responses:
        '200':
          description: Flow data retrieved successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  flows:
                    type: array
                    items:
                      $ref: '#/components/schemas/NetworkFlow'
                  pagination:
                    type: object
                    properties:
                      total_count:
                        type: integer
                      has_more:
                        type: boolean
                      limit:
                        type: integer
                      offset:
                        type: integer
        '400':
          description: Invalid query parameters
        '500':
          description: Internal server error

  /metrics:
    get:
      summary: Current system metrics
      description: Real-time system performance and network statistics
      responses:
        '200':
          description: Current metrics
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/SystemMetrics'

  /metrics/history:
    get:
      summary: Historical metrics
      description: Time-series metrics data within retention window
      parameters:
        - name: duration
          in: query
          description: Duration for historical data (5m, 15m, 1h)
          schema:
            type: string
            default: "15m"
        - name: interval
          in: query
          description: Data point interval (1s, 5s, 30s)
          schema:
            type: string
            default: "5s"
      responses:
        '200':
          description: Historical metrics data
          content:
            application/json:
              schema:
                type: object
                properties:
                  metrics:
                    type: array
                    items:
                      $ref: '#/components/schemas/SystemMetrics'
                  interval:
                    type: string
                  duration:
                    type: string

  /top-talkers:
    get:
      summary: Top bandwidth consumers
      description: Current top network talkers by traffic volume
      parameters:
        - name: limit
          in: query
          description: Number of top talkers to return
          schema:
            type: integer
            default: 10
            maximum: 50
        - name: duration
          in: query
          description: Time window for calculation (5m, 15m, 1h)
          schema:
            type: string
            default: "5m"
      responses:
        '200':
          description: Top talkers data
          content:
            application/json:
              schema:
                type: object
                properties:
                  top_talkers:
                    type: array
                    items:
                      $ref: '#/components/schemas/TopTalker'
                  time_window:
                    type: string
                  generated_at:
                    type: string
                    format: date-time

  /ws:
    get:
      summary: WebSocket endpoint for real-time updates
      description: Upgrade to WebSocket for streaming real-time flow and metrics data
      responses:
        '101':
          description: WebSocket connection established
        '400':
          description: WebSocket upgrade failed

components:
  schemas:
    NetworkFlow:
      type: object
      properties:
        flowKey:
          type: string
        srcIP:
          type: string
        dstIP:
          type: string
        srcPort:
          type: integer
        dstPort:
          type: integer
        protocol:
          type: integer
        protocolName:
          type: string
        bytes:
          type: integer
          format: int64
        packets:
          type: integer
          format: int64
        firstSeen:
          type: string
          format: date-time
        lastSeen:
          type: string
          format: date-time
        duration:
          type: integer
          description: Duration in milliseconds
        status:
          type: string
          enum: [active, idle, closed]

    SystemMetrics:
      type: object
      properties:
        timestamp:
          type: string
          format: date-time
        cpuUsage:
          type: number
          format: float
        memoryUsage:
          type: integer
          format: int64
        activeFlows:
          type: integer
        packetsPerSecond:
          type: integer
          format: int64
        bytesPerSecond:
          type: integer
          format: int64
        droppedPackets:
          type: integer
          format: int64
        captureErrors:
          type: integer
          format: int64
        interfaceStatus:
          type: string
          enum: [up, down, error]

    TopTalker:
      type: object
      properties:
        identifier:
          type: string
        displayName:
          type: string
        totalBytes:
          type: integer
          format: int64
        totalPackets:
          type: integer
          format: int64
        flowCount:
          type: integer
        rank:
          type: integer
        percentageOfTotal:
          type: number
          format: float
        flows:
          type: array
          items:
            $ref: '#/components/schemas/NetworkFlow'

  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      description: Optional token-based authentication when security mode enabled

security:
  - BearerAuth: []
  - {} # No authentication for localhost
```
