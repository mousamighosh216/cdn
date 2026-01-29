# CDN Control Plane

This directory contains the **Control Plane** implementation for our custom CDN system.

The control plane is the **brain of the CDN**.  
It does **not** serve user traffic. Instead, it coordinates edge nodes, tracks their health, and makes routing decisions.

---

## 📌 Role of the Control Plane

In this CDN architecture:

- **Edge Nodes** serve content
- **Control Plane** decides *which* edge should serve a request

The control plane maintains **global state** and ensures traffic is routed only to **healthy, available edges**.

---

## 🧠 Core Responsibilities

The control plane performs the following high‑level functions:

### 1️⃣ Edge Lifecycle Management
- Accepts edge registrations
- Maintains a registry of all edges
- Stores metadata:
  - Edge ID
  - Region
  - Port
  - Last heartbeat time
  - Health status

### 2️⃣ Health Monitoring
- Receives periodic heartbeats from edges
- Detects failed or unreachable edges
- Marks edges healthy/unhealthy automatically

### 3️⃣ Routing Decisions
- Selects an appropriate edge for a client request
- Ensures only healthy edges are returned
- (Later) applies region, load, or latency logic

### 4️⃣ Configuration Authority
- Acts as the single source of truth
- Owns routing rules and system behavior
- (Later) distributes configs to edges

---

## 🗂 Folder Structure (MVP)

control-plane/
├── api/
│ ├── main.go # Service entry point
│ ├── models.go # Edge data structures & registry
│ ├── health.go # Health evaluation logic
│ └── config.go # Config loader
│
├── config/
│ └── config.json # Runtime configuration
│
├── registry/ # (future) persistent edge storage
├── routing/ # (future) routing algorithms
├── logs/ # (optional) service logs
└── README.md

> All executable Go code currently lives inside `api/`.

---

## 🔌 Public API Endpoints

### `POST /register`
Registers a new edge node.

Used when:
- An edge starts
- An edge restarts

---

### `POST /heartbeat`
Receives periodic liveness signals from edges.

Used to:
- Detect failures
- Update health status

---

### `GET /resolve`
Returns a healthy edge for request routing.

Used by:
- DNS layer (future)
- Clients / load balancers (simulation)

---

## 🔁 Control Plane Runtime Flow

Edge boots
↓
POST /register
↓
Edge added to registry
↓
Heartbeats received
↓
Health status updated
↓
GET /resolve returns healthy edge

If heartbeats stop:
- Edge marked UNHEALTHY
- Removed from routing decisions

---

## 🧠 Architecture Overview (ASCII)
    ┌────────────┐
            │   Client   │
            └─────┬──────┘
                  │
            (resolve edge)
                  │
          ┌───────▼────────┐
          │  Control Plane  │
          │-----------------│
          │ Edge Registry   │
          │ Health Monitor  │
          │ Routing Logic   │
          └───────┬────────┘
                  │
      ┌───────────┴───────────┐
      │                       │
┌─────▼─────┐           ┌─────▼─────┐
│  Edge A   │           │  Edge B   │
│ (Healthy) │           │ (Healthy) │
└───────────┘           └───────────┘

---

## 🧩 Architecture Diagram (Mermaid)

```mermaid
flowchart TD
    Client -->|Resolve| ControlPlane
    ControlPlane -->|Select Edge| Edge1
    ControlPlane -->|Select Edge| Edge2

    Edge1 -->|Register| ControlPlane
    Edge2 -->|Register| ControlPlane

    Edge1 -->|Heartbeat| ControlPlane
    Edge2 -->|Heartbeat| ControlPlane
🧠 Separation of Concerns (Very Important)
Responsibility	Control Plane	Edge
Edge discovery	✅	❌
Health evaluation	✅	❌
Routing decisions	✅	❌
Content serving	❌	✅
Caching	❌	✅
Origin fetch	❌	✅
✅ Current Capabilities
✔ Edge registration
✔ Heartbeat‑based health tracking
✔ In‑memory edge registry
✔ Simple routing (first healthy edge)

🚀 Planned Enhancements
Routing
Region‑aware routing

Latency‑based routing

Load‑based routing

Registry
Persistent storage (Redis / DB)

Edge metadata enrichment

Security
Authenticated edge registration

Mutual TLS

Rate limiting

Observability
Metrics ingestion

Edge health dashboards

Alerting

🧠 Design Philosophy
The control plane is designed to be:

Stateless where possible

Deterministic

Horizontally scalable

Independent from request traffic

This mirrors how real CDNs (Fastly, Cloudflare, Akamai) design their control infrastructure.

🛠 Usage (Local)
cd control-plane/api
go run .
Ensure:

Port configured in config.json

Edge nodes are running

📄 License
Internal learning project. Not production‑ready.