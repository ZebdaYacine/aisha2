# PRD-06 · Public Catalogue & Discovery

| | |
| --- | --- |
| **Epic** | 6 — Public Catalogue & Discovery |
| **Stories** | US-VIS-001 … US-VIS-003 |
| **Priority** | High |
| **Status** | ✅ Implemented |
| **Surfaces** | `apps/web/app/[locale]/products`, `categories`, `artisans`, `search` |
| **API** | `GET /api/v1/categories`, `GET /api/v1/products`, `GET /api/v1/products/:id`, `GET /api/v1/artisans`, `GET /api/v1/artisans/:id`, `GET /api/v1/workshops`, `GET /api/v1/workshops/:id` |

## 1. Summary

Help visitors discover trustworthy Algerian products, categories, artisans, and workshops through localized, fast, editorial public pages.

## 2. Problem

Customers cannot evaluate an artisan product from a flat listing alone. Discovery must combine product information, provenance, artisan context, availability, and safe publication rules.

## 3. Goals

- Show only public, approved, active content.
- Make product, artisan, workshop, origin, and availability information discoverable.
- Support localized search, filtering, pagination, and responsive public pages.

### Non-goals

- Personalized recommendations or advanced BI dashboards in the MVP.

## 4. Users

- Visitor browsing without an account.
- Customer comparing products before purchase.
- Artisan whose approved public profile and products are discoverable.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | List only public and active products. | US-VIS-001 |
| R2 | Show product image, name, workshop, artisan, price, currency, category, and availability. | US-VIS-001 |
| R3 | Provide localized search, category filtering, pagination, and empty/error states. | US-VIS-001 |
| R4 | Provide public artisan profiles with biography, craft, public location, workshops, products, and verified information only when backed by verification. | US-VIS-002 |
| R5 | Provide public workshop profiles and authoritative stock availability when those capabilities are available. | US-VIS-003 |
| R6 | Never expose private documents, addresses, contact data, or unpublished entities. | US-VIS-002..003 |

## 6. Flow

```text
Landing page → category/search → product or artisan profile → product detail → cart
```

## 7. Technical notes

- Public APIs apply publication filtering server-side and support locale, pagination, category, workshop, and product search parameters.
- Public workshop read models are backed by approved artisan profiles and expose only active, public workshops.
- Product availability is derived from the append-only inventory movement ledger; the catalogue exposes stock state without allowing public mutation.
- Server Components handle initial catalogue content and metadata; interactive filters and pagination preserve URL state.
- Public artisan craft categories come from approved profile-category relations; verification is not shown unless a trusted verification field exists.

## 8. Success metrics

- Catalogue-to-product-detail engagement.
- Search/filter completion and zero-result rate.
- Public pages with no privacy or unpublished-content leaks.

## 9. Risks & open questions

- Full workshop ownership, lifecycle, and membership activation remain in PRD-07; this PRD consumes the public read model only.
- Inventory reception, inspection, reservation, and adjustment commands remain in PRD-10/PRD-11; catalogue availability is read-only.
- Collection and region discovery require confirmed data models before implementation.

## Source traceability

`REQ-PROD-001`, `REQ-PROD-002`, `REQ-I18N-001`, `US-VIS-001..003`, `docs/frontend.md`.
