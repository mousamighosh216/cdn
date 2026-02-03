# Routing & Traffic Steering

Routing is the process of deciding **which Edge Node should serve a client request**.

In a CDN, this is critical for:
- Low latency
- Load balancing
- High availability
- Fault tolerance

---

## Current Implementation (Basic)

At the current stage, routing is **manual / static**:

- Client directly calls a known Edge URL
- Example: `http://localhost:8001/file.txt`
- No automatic decision logic exists yet

This is suitable for:
- Local development
- Testing cache behavior
- Verifying edge registration

---

## Planned Routing Strategies

### 1. Least Latency Routing
Route client to the closest edge based on:
- Geographic region
- Ping / response time
- IP location lookup

**Goal:** Faster response times for end users.

---

### 2. Round‑Robin Load Balancing
Requests are distributed evenly across edges.

Example:
Request 1 → Edge A
Request 2 → Edge B
Request 3 → Edge C

**Goal:** Prevent overload on a single edge.

---

### 3. Health‑Aware Routing
Only route traffic to edges marked **healthy** by the Control Plane.

If an edge stops sending heartbeats:
- It is marked *dead*
- It is removed from routing pool

**Goal:** High availability and fault tolerance.

---

### 4. Weighted Routing (Future)
Edges can have weights based on:
- CPU capacity
- Region priority
- Network bandwidth

Example:
- Edge A weight = 2
- Edge B weight = 1  
Edge A receives twice the traffic.

---

## Control Plane Role

The Control Plane acts as the **decision engine** by:

- Maintaining a registry of active edges
- Tracking heartbeat status
- Exposing routing APIs
- Providing edge metadata (region, load, health)

---

## Future Enhancements

- Geo‑DNS routing
- Anycast IP simulation
- Adaptive load balancing
- Real‑time metrics based routing
- Failover clusters

---

## Summary

| Stage | Routing Type | Purpose |
|------|-------------|--------|
| Current | Static | Development & testing |
| Next | Round‑Robin | Load distribution |
| Future | Latency + Health | Production‑grade performance |

Routing evolves from **manual → automated → intelligent** as the platform matures.
