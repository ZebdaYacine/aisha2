# PRD-08 · Workshop-Linked Product Authoring & Media

| | |
| --- | --- |
| **Epic** | 8 — Workshop-Linked Product Authoring & Media |
| **Stories** | US-PROD-001 … US-PROD-004; US-PROD-006; US-PROD-008 |
| **Priority** | High |
| **Status** | ✅ Shipped |
| **Surfaces** | `apps/web/app/[locale]/artisan/products`, product editor and media controls |
| **API** | `GET /api/v1/artisan/workshops`, `POST/PATCH/DELETE /api/v1/artisan/workshops`, `POST /api/v1/artisan/workshops/:id/status`, `POST /api/v1/artisan/products`, `PATCH /api/v1/artisan/products/:id`, `POST /api/v1/artisan/products/:id/media`, `POST /api/v1/artisan/products/:id/submit`, `POST /api/v1/artisan/products/:id/archive` |

## 1. Summary

Give artisans a safe product-draft workspace tied to an owned workshop, with localized content and private media before moderation.

## 2. Problem

Product content can become untrustworthy when ownership is inferred from frontend input, drafts expose private media, or submitted versions are overwritten.

## 3. Goals

- Tie every product to exactly one owned workshop.
- Support localized, culturally rich product content and safe media upload.
- Preserve immutable submitted versions and clear draft/review states.

### Non-goals

- Direct public publication by an artisan without moderation.

## 4. Users

- Artisan creating or editing a product draft.
- Moderator receiving a submitted immutable version.
- Administrator auditing ownership and media actions.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Create products only for an active workshop owned by the artisan; every product has exactly one workshop. | US-PROD-001 |
| R2 | Capture localized names, descriptions, story, materials, method, use, dimensions, category, price, currency, origin, and eligibility flags. | US-PROD-001 |
| R3 | Allow owner-only draft editing and eligible moves between owned workshops. | US-PROD-002..003 |
| R4 | Validate required fields and media before submission. | US-PROD-004 |
| R5 | Store immutable submission snapshots and preserve prior versions after requested changes. | US-PROD-004, US-PROD-006 |
| R6 | Keep draft and review media private with generated storage keys. | REQ-MEDIA-001 |
| R7 | Preserve history when a product is archived. | US-PROD-008 |
| R8 | Present owned products in a responsive management table with first-column Details/Update/Archive actions and one New draft action. | US-PROD-001..008 |

## 6. Flow

```text
GET owned workshops → select active workshop → draft → localized content/media
       → edit or move between owned active workshops → validate → PENDING_REVIEW → moderation
                                                               ↘ CHANGES_REQUESTED → edit → resubmit
       → owner archive → ARCHIVED (history retained; no public publication)
```

## 7. Technical notes

- Product creation and edits re-check ACTIVE artisan membership, workshop ownership, and active workshop status in the repository transaction. Workshop lifecycle operations are now implemented by PRD-07 and are shared by the product editor.
- Product moves are allowed only between the same artisan’s active workshops and are blocked once inventory movement history exists.
- Price uses integer minor units and currency code.
- Submission snapshots include the workshop identity; later edits create a new snapshot version after requested changes.
- Product media is sent as multipart form data to `/artisan/products/:id/media`; the API validates signatures and size, generates a `products/{productId}/{uuid}` key, writes to the configured private MinIO bucket (`MINIO_PRIVATE_BUCKET`, default `product-private`), persists metadata only after the object write succeeds, removes the object on metadata failure, and returns a 15-minute signed URL. SQLBoiler models do not leave repositories.
- Artisan archive is owner-scoped, rejects products still in `PENDING_REVIEW`, clears `published_at`, and records an immutable audit/outbox event. The artisan workspace exposes that history-preserving action as Archive rather than destructive Delete. Admin activation and moderation decisions remain in PRD-09.
- The artisan workspace renders an owned-product table with Details, Update, and Archive actions in the first column. Details is read-only and includes localized product information and private media previews; New draft and Update open the authoring form in a scrollable modal while the table remains the index.

## 8. Success metrics

- Draft-to-submission completion rate.
- Submission rejection rate caused by missing fields, media, inactive workshops, or ownership violations.
- Zero cross-artisan product access or workshop assignment incidents.
- Share of drafts created and moved through an owned active workshop.

## 9. Risks & open questions

- Maximum upload sizes and certification rules remain configuration decisions under PRD-20.
- Admin activation and notification delivery remain PRD-09 concerns.
- Artisan onboarding provisions a default workshop before product creation is enabled; suspended or closed memberships cannot create, move, edit, submit, or publish products.
- Archived products retain their product, submission, inventory, and order-line history; workshop deletion remains blocked when protected history exists.

## Source traceability

`REQ-PROD-001`, `REQ-MEDIA-001`, `US-PROD-001..004`, `US-PROD-006`, `US-PROD-008`.
