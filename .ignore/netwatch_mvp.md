# Netwatch (MVP) — Live Traffic Web UI in Go

A single-host, ~1 Gbps, flow-aware traffic monitor. Captures from a live interface, aggregates flows and counters every second, and streams updates to a web UI. Implemented entirely in **Go**.

---

## Goal

- Show **bps/pps**, **protocol breakdown**, **top talkers**, and a **flow table** in near real time.
- Keep **~60 min** of in-memory history (flows + counters).
- Single binary written in **Go**.

---

## Stack (Go)

- **Capture/Decode:** Go + `gopacket/afpacket` (TPACKETv3), BPF filter, snaplen 128–160.
- **Aggregation:** Per-CPU sharded maps; 1 s rollups; TTL eviction.
- **Transport:** WebSocket for live updates; REST for paging/search.
- **UI:** Minimal HTML + JS (HTMX/Alpine) + Chart.js or ECharts.
- **(Optional later):** SQLite/Prometheus/ClickHouse, GeoIP.

---

## Architecture (minimal)

```
[ NIC ] -> [ AF_PACKET ring (Go) ] -> [ flow sharded table + 1s rollup ]
                                      |             |
                                      v             v
                               [ REST /api/* ]   [ WS /ws ]
                                      \             /
                                       \           /
                                        [ Static HTML/JS UI ]
```

---

## Data Model (Go types)

```go
// Flow key is bidirectional (canonicalized A<=B)
type FlowKey struct {
  AIP, BIP [16]byte
  APort, BPort uint16
  Proto uint8  // 6 TCP, 17 UDP, etc.
  VLAN  uint16 // 0 if none
}

type Flow struct {
  Key                 FlowKey
  FirstSeen, LastSeen int64 // unix ms
  Packets, Bytes      uint64
  AtoBBytes, BtoABytes uint64
  TCPFlags            uint16
  // Optional: SNI, DNSName
}

type Counters struct {
  TS int64
  PPS, BPS uint64
  ByProto map[string]uint64 // tcp/udp/icmp/other
  Drops uint64              // kernel->user drops (interval)
}
```

---

## API (minimal)

- **WebSocket** `GET /ws`
  - Client → `{"subscribe":["rates","top_talkers"]}`
  - Server (1/s):
    - `{"type":"rates","ts":..., "pps":..., "bps":..., "by_proto":{"tcp":...}}`
    - `{"type":"top_talkers","rows":[{"a":"10.0.0.1:1234","b":"8.8.8.8:53","bytes":12345}]}`

- **REST**
  - `GET /api/flows?limit=200&sort=bytes&dir=desc&proto=tcp&ip=10.0.0.0/24`
  - `GET /api/counters?from=...&to=...&step=1s`
  - `GET /api/stats` (iface, snaplen, ring/drops)

---

## Repo Layout

```
netwatch/
  cmd/netwatchd/main.go
  internal/capture/afpacket.go  // open ring, BPF, batch poll
  internal/flow/key.go          // canonical key
  internal/flow/table.go        // sharded map, TTL eviction
  internal/flow/rollup.go       // 1s counters, top-K
  internal/api/rest.go          // /api/*
  internal/api/ws.go            // /ws
  internal/web/server.go        // static, templates
  pkg/types/{flow.go,messages.go}
  build/systemd.service
```

---

## Implementation Steps

1. **Capture**
   - Open AF_PACKET (TPACKETv3) with ring buffers.
   - Apply BPF (default: `ip or ip6 and not port 22`).
   - Snaplen 160; parse L2→L4; extract 5-tuple + VLAN.

2. **Flow Table**
   - Canonicalize key (A<=B).
   - Shard by hash (numCPU shards), per-shard mutex.
   - Update packets/bytes and direction bytes; set `LastSeen`.
   - TTL eviction loop (TCP 120s, UDP 30s, ICMP 10s).

3. **Rollups**
   - Every 1 s: compute PPS/BPS/ByProto + Top-K (bounded heap).
   - Append to in-memory rings (≤60 min).

4. **Server**
   - `/ws`: broadcast 1 s frames; drop if client lags.
   - `/api/flows`: paginate from current table snapshot.
   - `/api/counters`: return ring slice.

5. **UI (one page)**
   - WebSocket connect on load; update two charts (bps/pps + proto bar).
   - Top-talkers list (auto-updating).
   - Flow table via `/api/flows` (HTMX swap + filter inputs).

---

## Build & Run

```bash
go build -o netwatchd ./cmd/netwatchd

# grant capture without root
sudo setcap 'cap_net_raw,cap_net_admin+eip' ./netwatchd

# run
./netwatchd   --iface=eth0   --bpf='ip or ip6 and not port 22'   --snaplen=160   --retention=60m
```

---

## Config (flags/env)

- `--iface` (string): e.g., `eth0`
- `--bpf` (string)
- `--snaplen` (int): default 160
- `--retention` (duration): default 60m
- `--bind` (addr): default `127.0.0.1:8080`
- `--auth-token` (string, optional)

---

## Performance Checklist (1 Gbps target)

- Enable NIC **RSS**; verify multiple RX queues.
- Run pollers = RX queues; increase ring size if drops.
- Keep snaplen small; pre-filter with BPF.
- Track and display kernel drop counters.

---

## Security (minimal)

- Bind UI to localhost by default.
- Optional bearer token for REST/WS.
- systemd unit with `AmbientCapabilities=cap_net_raw,cap_net_admin`.

---

## Next (post-MVP, optional)

- Persist counters to SQLite/Prometheus.
- ClickHouse for flow search >60 min.
- GeoIP/ASN, DNS/TLS SNI toggles.
- PCAP export for selected flows.
