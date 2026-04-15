# Fiber Apotek analysis

## Executive summary
- Existing app is a monolithic Fiber + GORM REST API.
- The API surface is broad, around 159 controller-backed routes, plus export PDF/Excel routes and a few inline utility routes.
- Business rules are embedded mostly inside controllers.
- Infra concerns, database access, token parsing, report syncing, and stock mutations are mixed into handler logic.
- This is exactly the kind of codebase that benefits from a Clean Architecture rewrite.

## Existing technical shape
- Entry point: `main.go`
- HTTP framework: Fiber v2
- Persistence: PostgreSQL via GORM
- Cache/session-ish infra: Redis, now mostly for auth token blacklist
- Supporting concerns: schedulers, exports, report syncing, env writer, backup/restore helpers

## High level domains

### System / auth
- login, logout, profile, menus
- branch selection and user-branch authority
- user, branch, user-branch management

### Master data
- products
- product categories
- units
- unit conversions
- suppliers
- supplier categories
- members
- member categories

### Transactions
- purchases + purchase items
- sales + sale items
- buy returns
- sale returns
- duplicate receipts + items + details
- expenses
- another incomes

### Audits / stock control
- first stocks + items
- opnames + items
- mobile opname views
- defecta + defecta items

### Reporting / dashboard
- transaction reports
- daily profit reports
- daily asset
- dashboard reports
- neraca saldo / profit graph

### Exports
- excel exports for most resources
- pdf exports for most resources and labels

## Existing route inventory

### Authentication and system basics
- POST `/api/login`
- POST `/api/logout`
- GET `/api/profile`
- GET `/api/menus`
- GET `/api/list_branches`
- POST `/api/set_branch`
- GET `/api/branches`
- GET `/api/branches/:id`
- POST `/api/branches`
- PUT `/api/branches/:id`
- DELETE `/api/branches/:id`
- GET `/api/users`
- GET `/api/users/:user_id`
- POST `/api/users`
- PUT `/api/users/:user_id`
- DELETE `/api/users/:user_id`
- GET `/api/user-branches`
- GET `/api/user-branches/:user_id/:branch_id`
- POST `/api/user-branches`
- PUT `/api/user-branches/:user_id/:branch_id`
- DELETE `/api/user-branches/:user_id/:branch_id`
- GET `/api/detail-users/:id`

### Master data
- Product categories
  - GET `/api/product-categories`
  - POST `/api/product-categories`
  - GET `/api/product-categories/:id`
  - PUT `/api/product-categories/:id`
  - DELETE `/api/product-categories/:id`
  - GET `/api/product-categories-combo`
- Products
  - GET `/api/products`
  - POST `/api/products`
  - GET `/api/products/:id`
  - PUT `/api/products/:id`
  - DELETE `/api/products/:id`
  - GET `/api/sales-products-combo`
  - GET `/api/purchase-products-combo`
- Units
  - GET `/api/units`
  - GET `/api/units/:id`
  - POST `/api/units`
  - PUT `/api/units/:id`
  - DELETE `/api/units/:id`
  - GET `/api/cmb-units`
  - GET `/api/units-combo`
- Unit conversions
  - GET `/api/unit-conversions`
  - GET `/api/unit-conversions/:id`
  - POST `/api/unit-conversions`
  - PUT `/api/unit-conversions/:id`
  - DELETE `/api/unit-conversions/:id`
  - GET `/api/conversion-products-combo`
- Suppliers
  - GET `/api/suppliers`
  - GET `/api/suppliers/:id`
  - POST `/api/suppliers`
  - PUT `/api/suppliers/:id`
  - DELETE `/api/suppliers/:id`
  - GET `/api/suppliers-combo`
- Supplier categories
  - GET `/api/supplier-categories`
  - DELETE `/api/supplier-categories/:id`
  - GET `/api/supplier-categories-combo`
- Members
  - GET `/api/members`
  - GET `/api/members/:id`
  - POST `/api/members`
  - PUT `/api/members/:id`
  - DELETE `/api/members/:id`
  - GET `/api/members-combo`
- Member categories
  - GET `/api/member-categories`
  - GET `/api/member-categories/:id`
  - POST `/api/member-categories`
  - PUT `/api/member-categories/:id`
  - DELETE `/api/member-categories/:id`
  - GET `/api/member-categories-combo`

### Transactions
- Purchases
  - GET `/api/purchases`
  - GET `/api/purchases/:id`
  - POST `/api/purchases`
  - PUT `/api/purchases/:id`
  - DELETE `/api/purchases/:id`
  - GET `/api/purchase-items/all/:id`
  - POST `/api/purchase-items`
  - PUT `/api/purchase-items/:id`
  - DELETE `/api/purchase-items/:id`
- Sales
  - GET `/api/sales`
  - GET `/api/sales/:id`
  - POST `/api/sales`
  - PUT `/api/sales/:id`
  - DELETE `/api/sales/:id`
  - GET `/api/sale-items/all/:id`
  - POST `/api/sale-items`
  - PUT `/api/sale-items/:id`
  - DELETE `/api/sale-items/:id`
  - GET `/api/sales-details`
- Buy returns
  - GET `/api/buy-returns`
  - GET `/api/buy-returns/:id`
  - POST `/api/buy-returns`
  - GET `/api/cmb-prod-buy-returns`
  - GET `/api/cmb-purchases`
- Sale returns
  - GET `/api/sale-returns`
  - GET `/api/sale-returns/:id`
  - POST `/api/sale-returns`
  - GET `/api/cmb-prod-sale-returns`
  - GET `/api/cmb-sales`
- Duplicate receipts
  - GET `/api/duplicate-receipts`
  - GET `/api/duplicate-receipts/:id`
  - POST `/api/duplicate-receipts`
  - PUT `/api/duplicate-receipts/:id`
  - DELETE `/api/duplicate-receipts/:id`
  - GET `/api/duplicate-receipts-items/all/:id`
  - POST `/api/duplicate-receipts-items`
  - PUT `/api/duplicate-receipts-items/:id`
  - DELETE `/api/duplicate-receipts-items/:id`
  - GET `/api/duplicate-receipts-details`
- Other financial transactions
  - GET `/api/expenses`
  - POST `/api/expenses`
  - PUT `/api/expenses/:id`
  - DELETE `/api/expenses/:id`
  - GET `/api/another-incomes`
  - POST `/api/another-incomes/`
  - PUT `/api/another-incomes/:id`
  - DELETE `/api/another-incomes/:id`

### Audits and stock control
- First stocks
  - GET `/api/first-stocks`
  - POST `/api/first-stocks`
  - PUT `/api/first-stocks/:id`
  - DELETE `/api/first-stocks/:id`
  - GET `/api/first-stock-with-items/:id`
  - GET `/api/first-stock-items/:id`
  - POST `/api/first-stock-items`
  - PUT `/api/first-stock-items/:id`
  - DELETE `/api/first-stock-items/:id`
- Opnames
  - GET `/api/opnames`
  - POST `/api/opnames`
  - GET `/api/opnames/:id`
  - PUT `/api/opnames/:id`
  - DELETE `/api/opnames/:id`
  - GET `/api/opname-items`
  - POST `/api/opname-items-all`
  - POST `/api/opname-items`
  - PUT `/api/opname-items`
  - DELETE `/api/opname-items`
  - GET `/api/cmb-product-opname`
- Mobile opname support
  - GET `/api/mobile-opnames`
  - GET `/api/mobile-opnames-active`
  - GET `/api/mobile-opnames-item-details`
  - GET `/api/mobile-opnames-items-glimpse`
- Defecta
  - GET `/api/sys-defectas`
  - GET `/api/sys-defectas/:id`
  - POST `/api/sys-defectas`
  - PUT `/api/sys-defectas/:id`
  - DELETE `/api/sys-defectas/:id`
  - GET `/api/sys-defecta-items`
  - POST `/api/sys-defecta-items`
  - PUT `/api/sys-defecta-items/:id`
  - DELETE `/api/sys-defecta-items/:id`

### Dashboard and reports
- GET `/api/daily_asset`
- GET `/api/dashboard/monthly-profit-report`
- GET `/api/dashboard/daily-profit-report`
- GET `/api/dashboard/weekly-profit-report`
- GET `/api/dashboard/profit-today-by-user`
- GET `/api/dashboard/top-selling-report`
- GET `/api/dashboard/least-selling-report`
- GET `/api/dashboard/neared-report`
- GET `/api/report/neraca-saldo`
- GET `/api/report/profit-by-month`

### Exports
- Broad excel and pdf export surface exists for nearly every resource.
- This should be postponed until the core transactional rewrite is stable.

## Behavior and business rule notes

### Auth flow
1. Login validates active user and returns short lived token with `sub`.
2. Set branch consumes that token and returns a richer token with:
   - `sub`
   - `name`
   - `branch_id`
   - `user_role`
   - `default_member`
   - `quota`
   - `subscription_type`
   - `real_asset`
3. Role middleware authorizes by `user_role` claim.
4. Logout blacklists JWT in Redis until token expiry.

### Product and combo behavior
- Product is branch scoped.
- Product create defaults stock to 0.
- Empty SKU is replaced with product ID.
- Sale combo reads sale price from product.
- Purchase combo reads purchase price from product.
- Opname combo reads purchase price plus live stock.

### Purchase transaction behavior
- Header and items are created inside a DB transaction.
- Purchase total is server calculated.
- Unit conversion may alter actual stock addition and normalized price.
- Product stock is increased.
- Product expired date may move earlier if purchased batch expires earlier.
- A transaction report entry is also created.

### Sale transaction behavior
- Header and items are created inside a DB transaction.
- Product stock is reduced.
- Daily profit report is updated.
- Branch quota may be reduced for quota subscriptions.
- Member points may be updated for non default members.
- Known legacy bug: item price and subtotal are still trusted from client payload.

### Opname behavior
- Opname header is created first.
- Adding opname item mutates product stock directly to the counted quantity.
- Existing stock and old purchase price are captured into `qty_exist` and `sub_total_exist`.
- Total opname is recalculated from item deltas.
- Report sync is triggered after header or total recalculation.
- Known legacy inconsistency: detail response and item-list response do not expose `qty_exist` / `sub_total_exist` consistently.

## Structural pain points in the old code
- Controllers own parsing, validation, transaction orchestration, stock math, report sync, and persistence.
- Global mutable infra via `configs.DB` and `configs.RDB`.
- Domain rules are not isolated into use cases.
- Framework helpers are tightly coupled with core behavior.
- Mixed synchronous and asynchronous stock/report updates make reasoning harder.
- Exports live directly beside core API and increase startup coupling.

## Clean Architecture rewrite recommendation

### Core dependency rule
- Domain and use case layers must not import Fiber, GORM, Redis, or env/config packages.
- HTTP, DB, Redis, scheduler, and export concerns belong to outer adapters.

### Suggested target structure
- `cmd/api`
- `internal/bootstrap`
- `internal/domain`
- `internal/usecase`
- `internal/ports`
- `internal/adapters/http`
- `internal/adapters/persistence/postgres`
- `internal/adapters/cache/redis`
- `internal/adapters/auth/jwt`
- `internal/adapters/reporting`
- `internal/shared`

### Proposed domain grouping
- `auth`
- `branch`
- `user`
- `product`
- `unit`
- `supplier`
- `member`
- `purchase`
- `sale`
- `buyreturn`
- `salereturn`
- `opname`
- `firststock`
- `defecta`
- `report`
- `dashboard`
- `export`

### Migration strategy
1. Freeze and document old API contracts first.
2. Build clean architecture foundation.
3. Reimplement auth and master data first.
4. Reimplement stock moving transactions next.
5. Reimplement reports and exports last.
6. Keep contract tests against old endpoints wherever feasible.

## Rewrite sequencing recommendation
Phase 1:
- bootstrap app
- config loader
- postgres and redis adapters
- JWT auth
- health endpoint
- response envelope
- auth endpoints
- branches, users, products, suppliers, units, members basic CRUD/combo

Phase 2:
- purchases
- sales
- buy returns
- sale returns
- opnames
- first stocks
- stock and report consistency tests

Phase 3:
- dashboard and report endpoints
- defecta
- exports
- scheduler jobs

## Important correctness notes for new implementation
- Preserve API path and payload compatibility where intentionally required.
- Do not preserve legacy bugs unless explicitly asked.
- For financial transactions, always calculate totals and subtotals on server side.
- Keep auth blacklist in Redis if logout invalidation remains required.
- Make stock mutation use cases explicit and testable.
- Avoid hidden async writes for correctness critical paths unless there is a strong reason.
