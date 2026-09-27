# AISHA MVP Requirements

## 1. Product Vision

AISHA promotes and sells Algerian handmade and artistic products to local and international customers. The product must place the artisan and the provenance of each item at the centre of the shopping experience.

AISHA is not a basic listing website. The core business model requires:

1. Artisan or workshop identification.
2. Rich product information.
3. Physical reception of stock by AISHA.
4. Quality and conformity inspection.
5. Managed stock.
6. Packing and shipment from AISHA.
7. Optional made-to-order requests.

## 2. MVP Scope

### Included

- Public multilingual storefront.
- Customer registration, authentication, and profile.
- Administrator-managed roles and permissions.
- Artisan onboarding and profile management.
- Product drafts and multilingual content.
- Product moderation.
- Warehouse reception.
- Quality inspection.
- Stock ledger and reservations.
- Shopping cart.
- Checkout.
- Payment abstraction.
- Trusted payment confirmation.
- Order management.
- Shipment recording and tracking status.
- In-app notifications.
- Audit records.
- Docker Compose development and deployment baseline.
- Jenkins CI validation.

### Excluded unless explicitly requested

- Advanced promotion engine.
- Marketplace syndication to Amazon or other marketplaces.
- Multi-warehouse routing.
- Automated artisan payouts.
- Complex tax engine.
- Advanced recommendation engine.
- Real-time chat.
- Loyalty programme.
- Multi-vendor split payment.
- Full return logistics.
- Native mobile application.
- Detailed BI dashboards.

## 3. Roles

### Visitor

- Browse public catalogue.
- Search and filter.
- View product and artisan details.
- Change language and theme.
- Register and sign in.

### Customer

- Manage personal profile and addresses.
- Use cart and checkout.
- Pay and view payment status.
- View order history.
- Track shipments.
- Submit a made-to-order request if enabled.

### Artisan

For the MVP, artisans may be created by an administrator or may submit an application for approval.

- Manage artisan profile.
- Create and edit product drafts.
- Submit products for moderation.
- View moderation decisions.
- View stock status for owned products.
- View requests for made-to-order products.

### Moderator

- Review submitted products.
- Approve, reject, or request changes.
- Record mandatory reasons.

### Warehouse Agent

- Register incoming products.
- Perform inspection.
- Accept, reject, quarantine, or record damaged quantities.
- Prepare paid orders.
- Record shipment handover.
- Read user and artisan-application context, including submitted application
  files, without user-role, artisan-decision, or media-deletion permissions.

### Moderator

- Review, approve, request changes for, reject, suspend, activate, and archive
  product submissions according to the moderation state machine.
- Read user and product media for moderation context.
- Cannot manage user roles, artisan applications, warehouse stock, or media deletion.

### Administrator

- Manage users and roles.
- Review artisan applications.
- Manage all moderator, warehouse, artisan, media, and audit operations.
- Review audit events.
- Suspend users, artisans, and products.
- Configure categories and operational data.

## 4. Functional Requirements

## REQ-AUTH-001 - Registration

The system must allow a visitor to create a customer account using approved identity fields.

Acceptance criteria:

- Email is unique when email authentication is enabled.
- Password is hashed.
- Raw passwords are never stored or logged.
- Input is validated on the backend.
- Repeated registration requests are rate-limited.

## REQ-AUTH-002 - Authentication

The system must authenticate active users and create a secure session.

Acceptance criteria:

- Invalid credentials return a generic error.
- Suspended users cannot sign in.
- Access and refresh credentials are short-lived and rotated when token-based sessions are used.
- Logout revokes the current session.

## REQ-RBAC-001 - Authorisation

Every protected backend route must enforce Casbin policies.

Acceptance criteria:

- Unauthenticated requests return `401`.
- Authenticated but forbidden requests return `403`.
- Resource ownership is checked in addition to role.
- Frontend hiding is never treated as authorisation.

## REQ-ART-001 - Artisan Profile

The system must store:

- Public display name.
- Legal or internal name where required.
- Workshop name.
- Biography in supported languages.
- Wilaya or location.
- Craft categories.
- Contact data with public/private visibility.
- Profile image and optional presentation media.
- Approval and suspension status.

## REQ-PROD-001 - Product Information

Each product must support:

- Artisan ownership.
- Product type: artisan-specific or standard/traditional.
- Multilingual name.
- Multilingual description.
- Product story and cultural context.
- Materials and composition.
- Production method.
- Intended use.
- Dimensions and weight.
- Category.
- Images and optional short video.
- Price in integer minor units.
- Planned workshop order quantity and total order price in integer minor units.
- Currency.
- Country and region of origin.
- Eco-friendly and fair-trade indicators when verified.
- Made-to-order eligibility.
- Moderation status.
- Publication status.
- Product submission supports one selected language (Arabic, English, French, or Spanish) and up to four presentation media files. Catalogue reads fall back to that selected translation when the visitor's language is unavailable.

## REQ-PROD-002 - Product Categories

The MVP must allow categories such as:

- Decoration and art.
- Domestic use.
- Jewellery and beauty accessories.
- Traditional copper products.
- Carpets and textiles.
- Glassware.
- Bamboo and halfa products.
- Traditional musical instruments.
- Souvenirs.
- Leather goods.
- Pottery and porcelain.
- Embroidery.

Categories must be data-driven, not hard-coded into the UI.

## REQ-MOD-001 - Product Moderation

A product cannot become public before approval.

Required states:

```text
DRAFT
PENDING_REVIEW
CHANGES_REQUESTED
APPROVED
ACTIVE
SUSPENDED
ARCHIVED
```

Rules:

- The artisan owns drafts.
- Submission freezes or versions the content under review.
- Rejection and change requests require a reason.
- Suspension prevents new purchases.
- Historical order snapshots remain unchanged.
- Approval generates a stable warehouse-facing product code; accepted warehouse stock may activate the approved product atomically when all publication gates pass.

## REQ-WH-001 - Reception

The warehouse must record incoming stock before inspection.

Reception must include:

- Product.
- Artisan or supplier.
- Received quantity.
- Reception reference.
- Date and receiving agent.
- Optional parcel or batch reference.
- Notes and evidence.

Received quantity is not sellable until accepted by inspection.

## REQ-WH-002 - Inspection

Inspection must classify quantities into:

- Accepted.
- Rejected.
- Quarantined.
- Damaged.

Invariant:

```text
accepted + rejected + quarantined + damaged = inspected quantity
```

Only accepted quantity can become available for sale.

For an approved product, a successful inspection with accepted quantity reevaluates the publication gates in the same transaction and publishes the product only when they pass.

## REQ-INV-001 - Stock Ledger

Inventory must use append-only movements and reservations rather than direct quantity editing.

Tracked values:

- On hand.
- Available.
- Reserved.
- Quarantined.
- Damaged.
- Rejected.
- Shipped.

Rules:

- Available stock cannot become negative.
- Every adjustment requires actor and reason.
- Multi-step stock changes use a database transaction.
- Redis is not the source of truth.

## REQ-CART-001 - Shopping Cart

Customers and visitors may add active products to a cart.

Rules:

- Adding to cart does not reserve stock.
- Cart price is informational until checkout.
- The backend recalculates all totals.
- Inactive products cannot remain purchasable.

## REQ-CHECKOUT-001 - Checkout

Checkout must:

1. Require authentication.
2. Reload products, prices, and availability.
3. Validate delivery address.
4. Calculate totals on the backend.
5. Create expiring inventory reservations.
6. Create a pending order.
7. Create a payment attempt.
8. Use an idempotency key.

The same request must not create duplicate orders or duplicate reservations.

## REQ-PAY-001 - Payment Adapter

Payment providers must be hidden behind an interface.

The MVP may start with:

- A sandbox provider.
- A manual development provider.
- One real provider only after credentials and exact business rules are confirmed.

Rules:

- The browser redirect is not trusted payment confirmation.
- Trusted webhook or server verification confirms payment.
- Amount and currency must match the order.
- Duplicate callbacks must not duplicate side effects.
- Raw payment-card data must never be stored.

The original business document mentions PayPal and Amazon Pay as possible tools, but provider selection is not considered confirmed until explicitly approved.

## REQ-ORDER-001 - Order Snapshot

An order must preserve:

- Product name.
- Product description summary.
- Variant or dimensions if applicable.
- Unit price.
- Currency.
- Quantity.
- Artisan identity.
- Delivery address.
- Shipping fee.
- Applied discount if any.
- Tax if any.
- Customer and order timestamps.

Product edits after purchase must not alter order history.

## REQ-ORDER-002 - Order States

Minimum states:

```text
PENDING_PAYMENT
PAID
PREPARING
READY_TO_SHIP
SHIPPED
DELIVERED
CANCELLED
PAYMENT_FAILED
PARTIALLY_REFUNDED
REFUNDED
```

State transitions must be validated in the application layer.

## REQ-SHIP-001 - Shipping

The MVP must permit manual or adapter-based shipment creation.

Required data:

- Carrier name.
- Tracking reference.
- Shipment status.
- Handover date.
- Delivery estimate when available.
- Shipment events.

Possible carriers referenced by the business document include international delivery companies, but no carrier is hard-coded as the exclusive provider.

## REQ-CUSTOM-001 - Made-to-Order Request

The MVP may include a basic custom request workflow:

- Customer selects an artisan or eligible product.
- Customer describes requirements.
- Customer attaches safe reference files.
- Artisan requests clarification, declines, or sends a quote.
- Customer accepts or rejects the quote.

Recommended states:

```text
SUBMITTED
CLARIFICATION_REQUESTED
QUOTED
QUOTE_ACCEPTED
QUOTE_REJECTED
CANCELLED
```

Payment and production automation for custom orders may remain post-MVP if time is limited.

## REQ-I18N-001 - Languages

The frontend architecture must support:

- Arabic.
- French.
- English.
- Spanish.

Arabic must use RTL layout.

Visible text must use translation keys.

## REQ-MEDIA-001 - Product Media

Media must be stored in MinIO or compatible object storage.

Rules:

- Files are private by default during draft and review.
- File type, signature, and size are validated.
- Storage names are generated by the server.
- Product media is public only after publication.
- Private artisan documents use authorised or signed access.

## REQ-AUDIT-001 - Audit

Append-only audit records are required for:

- Role changes.
- Artisan approval.
- Product moderation.
- Stock reception and inspection.
- Inventory adjustments.
- Payment status changes.
- Refunds.
- Order cancellation.
- Shipment status changes.
- Security-sensitive administration.

## 5. Data and Financial Rules

- Money uses integer minor units and currency code.
- PostgreSQL timestamps use UTC.
- Unique constraints prevent duplicate external references.
- Foreign keys have explicit deletion behaviour.
- Important status columns use constrained values or validated enums.
- Permanent data is stored in PostgreSQL.
- Redis is used for caching, rate limiting, locks where appropriate, and short-lived data.

## 6. Non-Functional Requirements

### Security

- Backend authorisation.
- Secure password hashing.
- Rate limits on sensitive endpoints.
- Secret management outside Git.
- Secure headers.
- Safe file upload.
- Redacted structured logs.
- No stack traces or SQL errors in public responses.

### Reliability

- Idempotency for checkout, payment, refund, and shipment creation.
- Database transactions for critical state changes.
- Outbox pattern for notifications and external side effects.
- Persistent volumes are preserved during deployment.

### Accessibility

- Semantic HTML.
- Keyboard access.
- Visible focus.
- Accessible labels and errors.
- Proper RTL behaviour.

### Performance

- Pagination.
- Indexed frequent queries.
- Optimised images.
- Cache non-authoritative public data.
- Avoid N+1 SQL queries.

## 7. Open Questions Codex Must Not Invent

- Exact payment provider.
- Exact shipping providers and tariffs.
- Supported countries at launch.
- Currencies accepted at launch.
- Tax rules.
- Commission rates.
- Artisan payout mechanism.
- Refund window.
- Whether cash on delivery is supported.
- Whether one order may include multiple artisans.
- Whether orders may be split into multiple shipments.
- Maximum file sizes.
- Required legal documents for artisan approval.
- Exact certification rules for fair-trade, origin, and eco-friendly labels.
