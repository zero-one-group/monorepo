# Go Modular

Template for a modular monolith backend: one deployable composed of feature Modules with shared infrastructure.

## Language

**Module**:
A self-contained feature area of the application (for example User or Auth) that owns its use cases, persistence, and HTTP surface.
_Avoid_: package, plugin, microservice

**Modular Monolith**:
A single deployable application composed of Modules that share one process and one primary datastore boundary.
_Avoid_: distributed system, service mesh
