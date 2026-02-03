# System Architecture

This project implements a simplified Content Delivery Network (CDN) consisting of three primary components:

- **Control Plane** – Coordination & intelligence
- **Edge Nodes** – Traffic serving & caching
- **Origin Server** – Source of truth for content

The architecture is modular and service‑oriented, allowing each component to scale independently.

---

## High‑Level Architecture

    +---------+
    | Client  |
    +----+----+
         |
         v
   +-----+------+
   |   Edge     |
   |  (Cache)   |
   +-----+------+
         |
Cache Miss v
+-----+------+
| Origin |
+------------+

Edge <---- Heartbeat / Register ----> Control Plane


---

## Components

### 1. Control Plane (The Brain)
Responsible for system coordination and decision making.

**Responsibilities**
- Edge registration
- Heartbeat monitoring
- Health status tracking
- Routing metadata (future)
- Central configuration source

**Key Modules**
- Registry
- Health Monitor
- Routing Logic
- API Layer

---

### 2. Edge Nodes (The Workers)
Distributed nodes that serve client traffic and cache content.

**Responsibilities**
- Register with Control Plane at startup
- Send periodic heartbeats
- Serve client requests
- Cache files locally
- Fetch from Origin on cache miss

**Edge Lifecycle**
1. Start service
2. Load configuration
3. Register with Control Plane
4. Begin heartbeat loop
5. Accept client traffic
6. Cache & serve files

---

### 3. Origin Server (The Source)
Central storage of original files.

**Responsibilities**
- Provide authoritative file copies
- Serve data to edges on cache miss
- No caching logic
- Simple and stable

---

## Request Flow

### Cache Hit
Client → Edge → Cached File → Client


### Cache Miss
Client → Edge → Origin → Edge Cache → Client


---

## Heartbeat Flow

Edge → Control Plane (/heartbeat)
Control Plane → Update Health Status


If heartbeats stop:
- Edge marked **unhealthy**
- Removed from routing pool (future)

---

## Failure Handling

| Scenario | Behavior |
|--------|---------|
| Edge Down | Marked unhealthy |
| Origin Down | Edge returns error |
| Control Plane Down | Edges continue serving cached data |

---

## Deployment Model

- Each component runs as an independent service
- Dockerized for isolation and portability
- Multiple edge instances supported
- Environment variables control runtime behavior

---

## Scalability Strategy

- Horizontal scaling of edges
- Stateless control plane APIs
- Disk‑based caching at edges
- Future: load‑aware routing & geo distribution

---

## Future Architectural Enhancements

- Intelligent Resolver Service
- Geo‑based routing
- Cache TTL & invalidation
- Metrics & observability
- Authentication & TLS
- Multi‑region deployment

---

## Summary

| Layer | Role | Scaling Pattern |
|------|------|----------------|
| Control Plane | Coordination | Horizontal |
| Edge Nodes | Traffic & Cache | Horizontal |
| Origin | Content Source | Vertical / Replicated |

The architecture evolves from **static local CDN → intelligent distributed platform** as routing, metrics, and automation are introduced.