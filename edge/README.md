# CDN Edge Node

This directory contains the **Edge Node** implementation for our custom CDN system.

An edge node represents a geographically distributed server responsible for:
- Registering itself with the control plane
- Maintaining liveness via heartbeats
- Serving requests close to users
- (Later) caching and fetching content from origin

---

## 📌 Role of the Edge in the CDN

The CDN is split into two major parts:

- **Control Plane** → decides *where* traffic should go
- **Edge Nodes** → actually *serve* the traffic

This repository implements the **Edge Node**, which acts as a lightweight, autonomous service that continuously reports its health to the control plane and serves user requests.

---

## 🧠 High‑Level Responsibilities

An edge node performs the following operations:

1. **Startup Initialization**
   - Loads configuration (edge ID, region, ports, control‑plane URL)
   - Identifies itself uniquely in the CDN

2. **Edge Registration**
   - Registers with the control plane via `/register`
   - Shares metadata like:
     - Edge ID
     - Region
     - Listening port

3. **Heartbeat Reporting**
   - Sends periodic heartbeats to `/heartbeat`
   - Allows the control plane to:
     - Detect failures
     - Mark edges healthy/unhealthy

4. **HTTP Server**
   - Runs a local HTTP server
   - Exposes `/health` endpoint
   - (Later) serves cached or origin content

5. **Autonomous Operation**
   - Continues running independently
   - Requires no constant control‑plane interaction beyond heartbeats

---

## 🗂 Folder Structure (MVP)

edge/
├── main.go # Orchestrates edge startup
├── config.go # Loads edge configuration
├── controlplane.go # HTTP client for control-plane communication
├── heartbeat.go # Periodic heartbeat logic
├── server.go # Edge HTTP server
├── config.json # Edge runtime configuration
├── go.mod
└── README.md

---

## 🔁 Runtime Flow

Edge starts
↓
Load config.json
↓
POST /register → Control Plane
↓
Start heartbeat loop
↓
Start HTTP server
↓
Edge marked ALIVE


If heartbeats stop:
- Control plane marks edge **UNHEALTHY**
- Edge is removed from routing decisions

---

## 🔌 Control Plane Interaction

| Endpoint       | Purpose                         |
|---------------|----------------------------------|
| `/register`   | Initial edge registration        |
| `/heartbeat`  | Liveness signal                  |
| `/resolve`   | (Indirect) client routing target |

---

## ⚙ Configuration (`config.json`)

Each edge is configured independently.

Example:

```json
{
  "edge_id": "edge-mumbai",
  "region": "india-west",
  "port": 9001,
  "control_plane_url": "http://localhost:8080",
  "heartbeat_interval": 5
}
This allows multiple edges to run simultaneously with different identities.

✅ Current Capabilities
✔ Self‑registration
✔ Heartbeat‑based health tracking
✔ Independent HTTP server
✔ Multi‑edge simulation (local or Docker)

🚀 Planned Enhancements
Future features will be added incrementally:

Caching
In‑memory cache

TTL handling

Disk cache

Content Serving
Fetch from origin on cache miss

Serve static and dynamic content

Metrics
Request counts

Latency reporting

Load metrics sent via heartbeat

Security
Authenticated control‑plane communication

TLS support

🧩 How This Fits the Big Picture
This edge node mirrors how real CDN providers (Cloudflare, Fastly, Akamai) design their edge infrastructure:

Small

Fast

Disposable

Region‑aware

Control‑plane driven

🛠 Usage (Local)
go run .
Ensure:

Control plane is running

config.json is present

📄 License
Internal learning project. Not production‑ready.