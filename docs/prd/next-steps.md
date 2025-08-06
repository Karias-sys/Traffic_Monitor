# Next Steps

## UX Expert Prompt

"Please create the UX architecture for Netwatch based on the attached PRD. Focus on the matrix-themed SOCC interface with tabbed navigation, real-time dashboard components, and flow table design. The target users are network administrators and SOCC analysts who need rapid access to technical network data. Prioritize information density and keyboard navigation over visual polish."

## Architect Prompt

"Please create the technical architecture for Netwatch based on the attached PRD. This is a high-performance Go monolith targeting 1 Gbps packet capture with AF_PACKET/TPACKETv3, real-time WebSocket streaming, and single-binary deployment. Focus on the packet processing pipeline, flow aggregation engine, memory management for 60-minute retention, and concurrent data structures for high-throughput operation. The system must maintain <5% CPU utilization and <1GB memory usage while providing sub-second dashboard updates."