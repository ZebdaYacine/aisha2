# AISHA MVP - Development Documentation

AISHA is an e-commerce platform dedicated to Algerian handmade and artistic products. Its main differentiators are:

- Every product is linked to its artisan, artist, workshop, origin, materials, story, and production method.
- AISHA operates primarily as a reseller with managed stock, not as a pure dropshipping marketplace.
- Products are received by AISHA, inspected, accepted into stock, packed, and shipped to customers.
- Customers may request made-to-order products from selected artisans.
- The public experience must support Arabic, French, English, and Spanish.
- The initial commercial target is Europe, followed by wider international markets.

## MVP Objective

Deliver a secure, multilingual platform that proves the complete business loop:

```text
Artisan identified
-> artisan profile created
-> product submitted
-> product moderated
-> physical stock received
-> quality inspection completed
-> stock becomes sellable
-> customer discovers product
-> customer checks out
-> payment confirmed
-> warehouse prepares parcel
-> shipment is recorded
-> order is delivered
```

## Approved Technical Stack

### Frontend

- Next.js 16 App Router
- React 19
- TypeScript
- Tailwind CSS 4
- shadcn/ui
- React Hook Form
- Zod
- next-themes
- Jest and React Testing Library
- Playwright

### Backend

- Go
- Fiber
- Casbin
- SQLBoiler
- Google Wire
- PostgreSQL
- Redis
- MinIO or S3-compatible storage
- OpenAPI 3
- golang-migrate

### Infrastructure

- Docker
- Docker Compose
- Jenkins
- Nginx or Traefik
- Persistent PostgreSQL and MinIO volumes

## Documentation Reading Order

1. `requirements.md`
2. `userstory.md`
3. `architecture.md`
4. `backend.md`
5. `frontend.md`
6. `security.md`
7. `deployment.md`

## Source-of-Truth Priority

1. Explicit instruction for the current task.
2. `requirements.md`.
3. Approved architecture decisions.
4. Existing tested behaviour.
5. Other documentation.

When documents conflict, Codex must stop and report the conflict. It must not silently invent business rules.

## Repository Target Layout

```text
/
├── AGENTS.md
├── README.md
├── Makefile
├── docker-compose.yml
├── Jenkinsfile
├── .env.example
├── docs/
├── backend/
│   ├── cmd/
│   ├── features/
│   ├── server/
│   ├── core/
│   ├── db/
│   ├── openapi/
│   ├── test/
│   └── tools/
├── frontend/
│   ├── app/
│   ├── components/
│   ├── features/
│   ├── lib/
│   ├── messages/
│   └── tests/
└── deploy/
    ├── nginx/
    ├── scripts/
    └── monitoring/
```

## MVP Definition of Done

The MVP is complete only when:

- Artisan and product data can be created and reviewed.
- Stock cannot become available before warehouse inspection.
- Customers can browse, add to cart, check out, and place an order.
- Payment confirmation is handled through a trusted callback or a clearly marked manual adapter for development.
- Inventory reservation prevents overselling.
- Order and shipment states are tracked.
- Backend authorisation is enforced.
- Arabic RTL and all four planned languages are structurally supported.
- PostgreSQL, Redis, and MinIO run through Docker Compose.
- Migrations, tests, builds, health checks, and Jenkins validation succeed.
- No secrets are committed.
