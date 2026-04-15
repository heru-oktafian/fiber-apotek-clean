# Clean Architecture notes for fiber-apotek-clean

## Uncle Bob essentials applied here
- Entities contain enterprise rules and core invariants.
- Use cases orchestrate application specific workflows.
- Interface adapters translate HTTP, DB, JWT, Redis, and file/export details.
- Frameworks and drivers stay outermost.
- Dependencies point inward only.

## Practical translation for this project

### Entities
Examples:
- Product
- Purchase
- PurchaseItem
- Sale
- SaleItem
- Opname
- OpnameItem
- User
- Branch

These should represent business concepts, not GORM model details.

### Use cases
Examples:
- Login user
- Select branch
- Create purchase transaction
- Create sale transaction
- Create opname and opname item
- List products combo for sale/purchase/opname

These own business flow and invariants.

### Interface adapters
Examples:
- Fiber handlers for request/response binding
- PostgreSQL repository implementations
- JWT token service implementation
- Redis token blacklist implementation

### Frameworks and drivers
Examples:
- Fiber
- GORM
- Redis client
- Postgres driver
- env loader

## Architecture rules for the new repo
- No Fiber import inside domain or usecase packages.
- No GORM model leakage into usecase signatures.
- Repository interfaces should be expressed in terms of domain entities or dedicated DTOs.
- Transaction boundaries should be owned by use cases via a unit-of-work or transaction manager port.
- Response formatting belongs in transport adapters, not in use cases.

## Preferred project skeleton
- `cmd/api/main.go`
- `internal/bootstrap/app.go`
- `internal/domain/...`
- `internal/usecase/...`
- `internal/ports/...`
- `internal/adapters/http/fiber/...`
- `internal/adapters/persistence/postgres/...`
- `internal/adapters/auth/jwt/...`
- `internal/adapters/cache/redis/...`
- `internal/shared/response`
- `internal/shared/pagination`
- `internal/shared/clock`
- `internal/shared/idgen`

## Contract strategy
Because the old app already has frontend expectations, keep:
- route names
- key request field names
- key response envelope shape where feasible

But improve internals:
- server-side totals
- explicit validation
- deterministic transactions
- better separation of concerns
