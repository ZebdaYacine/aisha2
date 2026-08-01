# AISHA MVP User Stories and Workflows

## 1. Core Business Workflow

```text
Admin identifies or approves artisan
-> artisan profile is published
-> artisan creates product draft
-> moderator reviews product
-> approved product waits for warehouse stock
-> warehouse receives items
-> warehouse inspects items
-> accepted quantity becomes available
-> customer browses product
-> customer adds item to cart
-> customer checks out
-> stock is reserved
-> payment is confirmed
-> order becomes paid
-> warehouse picks and packs
-> shipment is created
-> customer tracks delivery
-> order is delivered
```

This workflow reflects AISHA's managed-stock reseller model. The system must never make received stock sellable before quality inspection.

## 2. Visitor Stories

### US-VIS-001 - Browse catalogue

**As a visitor**, I want to browse active products so that I can discover Algerian handmade goods.

Main flow:

1. Open home page.
2. Select language.
3. Browse featured categories, products, and artisans.
4. Apply filters.
5. Open product details.

Acceptance:

- Only active products appear.
- Product image, name, artisan, price, currency, category, and availability are shown.
- Empty, loading, and error states exist.
- Arabic uses RTL.

### US-VIS-002 - View artisan profile

The visitor can see:

- Artisan name or workshop.
- Craft category.
- Location.
- Biography.
- Public media.
- Published products.
- Product-making story when provided.

Private legal and contact information is not exposed.

## 3. Authentication Stories

### US-AUTH-001 - Register

Main flow:

1. Visitor opens registration.
2. Enters required fields.
3. Accepts terms.
4. Backend validates.
5. Password is hashed.
6. Account is created.
7. User signs in or verifies account depending on configuration.

Failures:

- Duplicate email.
- Weak password.
- Invalid form.
- Rate limit.
- Verification delivery failure.

### US-AUTH-002 - Sign in

Main flow:

1. User submits credentials.
2. Backend verifies password and account status.
3. Session is created.
4. User is redirected to intended page.

Failures:

- Invalid credentials.
- Suspended user.
- Disabled account.
- Excessive attempts.

### US-AUTH-003 - Sign out and revoke session

The system revokes the active session. Repeating sign-out is safe.

## 4. Artisan Onboarding

### US-ART-001 - Create artisan application

The applicant submits:

- Identity or organisation information.
- Workshop name.
- Wilaya and address.
- Craft categories.
- Biography.
- Contact data.
- Required documents.
- Profile image.
- Consent and declarations.

States:

```text
DRAFT
SUBMITTED
UNDER_REVIEW
CHANGES_REQUESTED
APPROVED
REJECTED
SUSPENDED
```

### US-ART-002 - Review application

Admin flow:

1. Open submitted applications.
2. Review information and private documents.
3. Approve, reject, or request changes.
4. Add reason.
5. Notify applicant.
6. On approval, grant artisan role.

Security:

- Documents are private.
- Only authorised reviewers may read them.
- Every decision is audited.

## 5. Product Stories

### US-PROD-001 - Create draft

Artisan enters:

- Names in supported languages.
- Descriptions.
- Product story.
- Materials.
- Production method.
- Intended use.
- Cultural origin.
- Category.
- Dimensions and weight.
- Price and currency.
- Images.
- Optional video.
- Made-to-order eligibility.
- Eco-friendly or fair-trade claims for later verification.

The product remains invisible while `DRAFT`.

### US-PROD-002 - Submit for moderation

1. Artisan opens complete draft.
2. Backend validates mandatory fields.
3. Product changes to `PENDING_REVIEW`.
4. Moderator is notified.
5. Submitted version is preserved.

### US-PROD-003 - Moderate product

Moderator actions:

- Approve.
- Reject.
- Request changes.
- Suspend an active product.

Reasons are mandatory for rejection, requested changes, and suspension.

### US-PROD-004 - Activate sellable product

Approval alone is insufficient.

A product becomes sellable only when:

- Product is approved.
- Product is not suspended.
- At least one accepted stock unit is available.
- Price and currency are valid.
- Required public media exists.

## 6. Warehouse Stories

### US-WH-001 - Receive stock

1. Agent identifies product and supplier.
2. Records batch or parcel.
3. Records received quantity.
4. System creates reception record.
5. Quantity becomes `RECEIVED_PENDING_INSPECTION`.

### US-WH-002 - Inspect batch

1. Agent opens pending batch.
2. Enters accepted, rejected, quarantined, and damaged quantities.
3. Adds reasons and evidence.
4. Backend validates the quantity equation.
5. Transaction creates inventory movements.
6. Accepted stock becomes available.

### US-WH-003 - Correct inventory

Only an authorised manager can adjust inventory.

Required:

- Product or batch.
- Adjustment quantity.
- Reason.
- Actor.
- Timestamp.
- Audit event.

Stock totals cannot be edited silently.

## 7. Cart and Checkout

### US-CART-001 - Add item to cart

Adding to cart does not reserve stock.

The cart displays:

- Product.
- Artisan.
- Quantity.
- Current display price.
- Currency.
- Subtotal.
- Availability warning.

### US-CHECK-001 - Checkout

Main flow:

1. Customer signs in.
2. Backend reloads current products and prices.
3. Customer chooses address.
4. Backend validates delivery.
5. Backend calculates totals.
6. Backend creates stock reservations.
7. Backend creates pending order.
8. Backend creates payment attempt.
9. Customer continues to payment.

Failures:

- Product inactive.
- Price changed.
- Insufficient stock.
- Invalid address.
- Unsupported destination.
- Duplicate request.

### US-CHECK-002 - Concurrent last item

Given one available item:

1. Customer A and B submit checkout concurrently.
2. Database transaction or lock protects the stock.
3. Only one reservation succeeds.
4. Other request receives `OUT_OF_STOCK`.
5. Available stock never becomes negative.

## 8. Payment Stories

### US-PAY-001 - Start payment

The backend sends the trusted amount and currency to the provider adapter. The frontend never supplies a trusted total.

### US-PAY-002 - Return from provider

The browser return page may display `PENDING`. A success query parameter does not confirm payment.

### US-PAY-003 - Process provider webhook

1. Verify signature and timestamp.
2. Reject unknown or invalid event.
3. Detect duplicate event.
4. Match order, amount, and currency.
5. Validate state transition.
6. Mark payment captured.
7. Mark order paid.
8. Convert reservation into committed allocation.
9. Queue notification through outbox.

### US-PAY-004 - Payment failure or expiry

1. Payment becomes failed or expired.
2. Order becomes `PAYMENT_FAILED` or remains eligible for retry.
3. Stock reservation is released.
4. Customer is notified.

## 9. Order and Fulfilment

### US-ORD-001 - View order

Customer sees only their own order.

Artisan sees only their own order items and minimum required customer information.

Warehouse sees fulfilment information.

### US-FUL-001 - Pick and pack

1. Paid order enters warehouse queue.
2. Agent picks reserved stock.
3. Agent records packed parcel.
4. Order becomes `READY_TO_SHIP`.

Picked quantity cannot exceed allocated quantity.

### US-SHIP-001 - Create shipment

1. Agent chooses carrier or manual shipment.
2. Tracking reference is stored.
3. Order becomes `SHIPPED`.
4. Customer receives notification.

### US-SHIP-002 - Record delivery event

Shipment events may include:

```text
PENDING
SHIPPED
IN_TRANSIT
OUT_FOR_DELIVERY
DELIVERED
DELIVERY_FAILED
RETURNED
```

Duplicate events do not duplicate notifications.

## 10. Cancellation

### US-ORD-002 - Cancel unpaid order

- Customer may cancel eligible unpaid order.
- Reservation is released.
- Payment attempt is cancelled or expired.

### US-ORD-003 - Cancel paid order before shipment

- Eligibility is checked.
- Refund starts if required.
- Stock is returned through a movement.
- Action is audited.

Shipped orders cannot be silently cancelled.

## 11. Made-to-Order

### US-CUSTOM-001 - Submit request

Customer provides:

- Artisan or product.
- Description.
- Quantity.
- Desired date.
- Reference images or documents.
- Contact preference.

### US-CUSTOM-002 - Artisan responds

Artisan can:

- Request clarification.
- Decline.
- Send a quote with price, scope, expiry, and estimated production time.

Quote versions are immutable.

### US-CUSTOM-003 - Customer accepts quote

- Expired quote cannot be accepted.
- Accepted quote becomes trusted checkout data.
- Payment and production workflow is started only after explicit acceptance.

## 12. Administration

### US-ADMIN-001 - Manage roles

Role change requires backend permission and creates an audit event.

### US-ADMIN-002 - Suspend artisan or product

Suspension prevents new selling activity but does not delete historical orders.

### US-ADMIN-003 - Review audits

Authorised users can search audit events by:

- Actor.
- Target.
- Event type.
- Date.
- Correlation ID.

## 13. Critical State Machines

### Product

```text
DRAFT
-> PENDING_REVIEW
-> CHANGES_REQUESTED
-> PENDING_REVIEW
-> APPROVED
-> ACTIVE
-> SUSPENDED
-> ARCHIVED
```

### Reception

```text
RECEIVED_PENDING_INSPECTION
-> INSPECTING
-> ACCEPTED
-> PARTIALLY_ACCEPTED
-> REJECTED
-> QUARANTINED
-> CLOSED
```

### Reservation

```text
ACTIVE
-> CONSUMED
-> RELEASED
-> EXPIRED
```

### Payment

```text
CREATED
-> PENDING
-> CAPTURED
-> FAILED
-> CANCELLED
-> PARTIALLY_REFUNDED
-> REFUNDED
```

### Order

```text
PENDING_PAYMENT
-> PAID
-> PREPARING
-> READY_TO_SHIP
-> SHIPPED
-> DELIVERED

PENDING_PAYMENT -> PAYMENT_FAILED
PENDING_PAYMENT -> CANCELLED
PAID -> CANCELLED
PAID -> PARTIALLY_REFUNDED
PAID -> REFUNDED
```

## 14. Required End-to-End Tests

1. Artisan onboarding approval.
2. Product moderation with requested changes.
3. Warehouse partial acceptance.
4. Product not sellable before inspection.
5. Standard purchase.
6. Concurrent purchase of final item.
7. Payment webhook delayed.
8. Duplicate payment webhook.
9. Payment failure releases stock.
10. Paid order fulfilment and shipment.
11. Customer cannot read another customer's order.
12. Artisan cannot edit another artisan's product.
13. Direct API authorisation bypass attempt returns `403`.
14. Arabic RTL page smoke test.
