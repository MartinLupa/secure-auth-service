

## Justify selection of the following
- REST vs gRPC
- Service discovery: etcd, Consul, etc.
- API Gateway: Kong, Envoy, etc.
- Security: circuit breakers (using Hystrix-like patterns in Go)
- Social (OAuth via Goth), MFA (TOTP)


## Design
- Error handling: fmt.Errorf --> log.Fatal --> panic

## Architecture
The pattern I applied is **Clean Architecture** (also called Layered/Hexagonal Architecture), which organizes code into concentric layers where dependencies flow inward. Each layer has a specific responsibility and only knows about layers below it, never above.

**The layers work like this:**

- **Models** (`internal/models/`) - Pure data structures with no dependencies
- **Repository** (`internal/repository/`) - Data access layer that talks to the database, depends only on models
- **Service** (`internal/service/`) - Business logic layer that orchestrates operations, depends on repository interfaces
- **Handlers** (`internal/handlers/`) - HTTP layer that handles requests/responses, depends on service interfaces