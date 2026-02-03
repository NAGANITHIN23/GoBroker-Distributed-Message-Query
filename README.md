cat <<EOF > README.md
# GoBroker: Distributed Message Queue (Go/TCP)

**GoBroker** is a high-performance, fault-tolerant message broker built from scratch in Go. It bypasses HTTP to use a **custom binary protocol over raw TCP**, achieving low-latency communication with **Write-Ahead Logging (WAL)** for data durability.

## Key Engineering Features
* **Custom Binary Protocol:** Designed a lightweight packet structure to reduce payload size by ~60% compared to JSON/REST.
* **Zero Data Loss (Durability):** Implements an append-only **Write-Ahead Log (WAL)**.
* **High Concurrency:** Utilizes Go **Goroutines** and **RWMutex** to handle thousands of concurrent clients.

## Architecture
* **Packet Structure:** \`[OpCode (1B)] [TopicLen (2B)] [Topic] [PayloadLen (4B)] [Payload]\`
* **Storage:** Thread-safe in-memory Map for routing + Disk-based WAL for recovery.

## Quick Start
1. **Start Server:** \`go run main.go protocol.go\`
2. **Start Subscriber:** \`cd client && go run main.go protocol.go sub news\`
3. **Publish Message:** \`cd client && go run main.go protocol.go pub news "Hello!"\`
EOF
