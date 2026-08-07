# AISHA MVP User Stories and Workflows

## 1. Purpose and Business Rules

AISHA is a managed marketplace for Algerian handmade products.

The platform uses one user account model:

- Every registered user has customer capabilities.
- A registered user can activate artisan membership using the same account.
- Artisan membership adds seller capabilities; it does not replace customer capabilities.
- An artisan can browse, order, pay, cancel eligible orders, and track delivery like any other customer.
- An artisan must own at least one workshop.
- Creating the first workshop is mandatory during artisan onboarding.
- Supporting documents are optional during onboarding.
- One artisan can own multiple workshops.
- Every product must belong to exactly one workshop.
- One workshop can contain multiple products.
- A product is not sellable merely because it was created or approved.
- Received stock must pass warehouse inspection before it becomes available for sale.
- All sensitive operations must be authorised by the backend and recorded when auditing is required.

## 2. Core Business Workflows

### 2.1 Customer purchase workflow

```text
Visitor browses catalogue
-> visitor registers or signs in
-> user acts as customer
-> customer adds product to cart
-> customer checks out
-> backend validates product, price, delivery, and stock
-> stock is reserved
-> payment attempt is created
-> payment is confirmed by trusted provider event
-> order becomes paid
-> warehouse picks and packs
-> shipment is created
-> customer tracks delivery
-> order is delivered
```

### 2.2 Become an artisan workflow

```text
Registered user acts as customer
-> user clicks "Become an Artisan"
-> artisan onboarding page opens
-> user completes artisan profile information
-> user may upload supporting documents
-> user creates the first mandatory workshop
-> backend validates the artisan profile and workshop
-> artisan membership is activated on the existing account
-> user is redirected to artisan dashboard
-> user retains all customer capabilities
```

### 2.3 Workshop and product workflow

```text
Artisan opens artisan dashboard
-> artisan creates or selects an owned workshop
-> artisan creates product draft for that workshop
-> artisan submits product for moderation
-> moderator approves, rejects, or requests changes
-> approved product waits for accepted warehouse stock
-> warehouse receives product items
-> warehouse inspects received items
-> accepted quantity becomes available
-> product becomes sellable when all activation rules are satisfied
```

### 2.4 Managed-stock rule

AISHA uses a managed-stock reseller model. The system must never make received stock sellable before warehouse quality inspection. Product approval and inventory acceptance are separate controls.

## 3. User, Membership, and Ownership Model

### 3.1 User capabilities

Every authenticated user can:

- Manage their account.
- Manage delivery addresses.
- Browse active products and public artisan profiles.
- Add products to the cart.
- Place orders.
- Pay for orders.
- View their order history.
- Track shipments.
- Cancel eligible orders.
- Submit made-to-order requests.
- Apply for or activate artisan membership.

### 3.2 Artisan capabilities

A user with active artisan membership can additionally:

- Manage their artisan profile.
- Create and manage workshops.
- Create products for owned workshops.
- Submit products for moderation.
- View seller-side order items related to their products.
- Respond to made-to-order requests.
- View permitted workshop, product, sales, and fulfilment information.

### 3.3 Domain relationships

```text
User
├── has customer capabilities
├── has zero or one ArtisanMembership
├── has many CustomerOrders
├── has many Addresses
├── has one Cart
└── may own many Workshops through ArtisanMembership

ArtisanMembership
├── belongs to one User
├── has one ArtisanProfile
├── has many Workshops
└── has an optional VerificationStatus

Workshop
├── belongs to one ArtisanMembership
├── has one owner User through ArtisanMembership
└── has many Products

Product
├── belongs to one Workshop
├── belongs indirectly to one ArtisanMembership
└── may have many inventory batches and order items
```

### 3.4 Membership states

```text
ACTIVE
SUSPENDED
CLOSED
```

Artisan membership becomes `ACTIVE` after successful onboarding and creation of the first valid workshop.

### 3.5 Optional verification states

Document verification is separate from artisan membership activation.

```text
NOT_SUBMITTED
PENDING
VERIFIED
CHANGES_REQUESTED
REJECTED
```

Optional documents may be submitted during onboarding or later. A verification state must not silently remove the user's customer capabilities.

## 4. Visitor Stories

### US-VIS-001 - Browse catalogue

**As a visitor**, I want to browse active products so that I can discover Algerian handmade goods.

#### Main flow

1. Visitor opens the home page.
2. Visitor selects a supported language.
3. Visitor browses featured categories, products, workshops, and artisans.
4. Visitor applies filters or search criteria.
5. Visitor opens product details.

#### Acceptance criteria

- Only public and active products appear.
- Product image, name, workshop, artisan, price, currency, category, and availability are shown.
- Empty, loading, and error states exist.
- Arabic pages use RTL.
- Private artisan, workshop, legal, and contact information is not exposed.

### US-VIS-002 - View artisan profile

**As a visitor**, I want to view a public artisan profile so that I can learn about the artisan and their work.

#### The visitor can see

- Artisan public name.
- Biography.
- Public profile image.
- Craft categories.
- Public location information.
- Public workshops.
- Published products.
- Product-making stories when provided.
- Verification badge when the artisan is verified.

#### Acceptance criteria

- Private identity information is not displayed.
- Private documents are never displayed.
- Private addresses and private contact details are not displayed.
- Suspended or closed artisan profiles do not allow new selling activity.

### US-VIS-003 - View workshop profile

**As a visitor**, I want to view a workshop so that I can understand where products are made or supplied.

#### The visitor can see

- Workshop name.
- Public description.
- Public location.
- Craft categories.
- Public media.
- Active products belonging to the workshop.
- Artisan public identity.

#### Acceptance criteria

- Only public workshops are visible.
- Only active products appear.
- Private workshop contact and legal information is hidden.

## 5. Authentication and Account Stories

### US-AUTH-001 - Register

**As a visitor**, I want to register so that I can purchase products and later become an artisan using the same account.

#### Main flow

1. Visitor opens registration.
2. Visitor enters required fields.
3. Visitor accepts required terms.
4. Backend validates the request.
5. Password is hashed using the configured secure password algorithm.
6. User account is created.
7. Customer capabilities are available by default.
8. User signs in or verifies the account depending on configuration.

#### Failures

- Duplicate email or other unique identifier.
- Weak password.
- Invalid form.
- Terms not accepted.
- Rate limit exceeded.
- Verification delivery failure.

#### Acceptance criteria

- Registration creates one user account.
- Registration does not create a separate customer record that prevents later artisan activation.
- A newly registered user can use customer features.

### US-AUTH-002 - Sign in

**As a user**, I want to sign in so that I can access my customer and artisan capabilities.

#### Main flow

1. User submits credentials.
2. Backend verifies credentials and account status.
3. Session is created.
4. Effective permissions are loaded.
5. User is redirected to the intended page.

#### Acceptance criteria

- An artisan signs in through the same authentication flow as any other user.
- The session exposes only backend-authorised capabilities.
- A suspended artisan can still use allowed customer features unless the entire user account is suspended.

#### Failures

- Invalid credentials.
- Suspended or disabled user account.
- Excessive attempts.
- Expired or invalid verification state when verification is mandatory for sign-in.

### US-AUTH-003 - Sign out and revoke session

**As a user**, I want to sign out so that the current session can no longer be used.

#### Acceptance criteria

- The active session is revoked.
- Repeating sign-out is safe.
- Revoking one session does not corrupt other valid sessions unless the user selects sign out from all devices.

### US-ACC-001 - View combined account dashboard

**As a user**, I want one account dashboard so that I can access customer features and artisan features without switching accounts.

#### Acceptance criteria

- Customer sections are available to all active users.
- Artisan sections appear only when artisan membership is active.
- The dashboard may provide a mode switch such as `Shopping` and `Artisan Dashboard`, but both modes use the same account and session.
- Switching interface mode does not change ownership or permissions.

## 6. Artisan Membership and Onboarding Stories

### US-ART-001 - Start artisan onboarding

**As a registered customer**, I want to click **Become an Artisan** so that I can add artisan capabilities to my account.

#### Preconditions

- User is authenticated.
- User account is active.
- User does not already have active artisan membership.

#### Main flow

1. User opens the account menu or customer dashboard.
2. User clicks **Become an Artisan**.
3. System opens the artisan onboarding page.
4. Existing account information is reused where appropriate.
5. User begins an onboarding draft.

#### Acceptance criteria

- The action does not create a second user account.
- Existing customer orders, cart, addresses, and profile remain unchanged.
- Reopening onboarding resumes the user's current draft when one exists.
- A user with active artisan membership is redirected to the artisan dashboard instead of creating a duplicate membership.

### US-ART-002 - Complete artisan profile

**As a customer becoming an artisan**, I want to provide artisan information so that the platform can create my artisan profile.

#### Required information

- Public artisan or business name.
- Artisan type: individual or organisation, when applicable.
- Biography or activity description.
- Main craft category or categories.
- Wilaya.
- Contact information required for platform operations.
- Consent and declarations.

#### Optional information

- Public profile image.
- Additional biography details.
- Social or public links.
- Supporting documents.
- Identity or organisation documents.
- Craft certificates.
- Other verification evidence.

#### Acceptance criteria

- Supporting documents are optional.
- Private data is separated from public profile data.
- The backend validates required fields.
- Uploaded files use allowed formats and size limits.
- Private uploads are stored securely and are not publicly accessible.

### US-ART-003 - Create first mandatory workshop

**As a customer becoming an artisan**, I want to create my first workshop so that my artisan membership has a valid production or business entity.

#### Required workshop information

- Workshop name.
- Workshop description.
- At least one craft category.
- Wilaya.
- Commune or city.
- Operational address.
- Contact information required for platform operations.

#### Optional workshop information

- Public workshop image or logo.
- Public gallery.
- Public location description.
- Opening information.
- Additional production details.

#### Main flow

1. User completes the artisan profile step.
2. User opens the first workshop step.
3. User enters required workshop information.
4. Backend validates the workshop.
5. Workshop is created and linked to the pending artisan membership.
6. Artisan membership is activated on the existing user account.
7. User is redirected to the artisan dashboard.

#### Acceptance criteria

- At least one workshop is mandatory.
- Artisan membership cannot become active without a valid first workshop.
- The first workshop belongs to the authenticated user through artisan membership.
- Workshop creation and membership activation occur atomically or are safely recoverable.
- Repeating the final request does not create duplicate memberships or duplicate workshops.
- The user retains all customer capabilities after activation.

#### Failures

- Missing required profile data.
- Missing workshop.
- Invalid workshop information.
- Unsupported location.
- Duplicate idempotency key.
- Invalid optional upload.
- User already has active artisan membership.

### US-ART-004 - Activate artisan membership

**As a customer who completed onboarding**, I want artisan membership activated so that I can manage workshops and products immediately.

#### Acceptance criteria

- The existing user receives artisan capabilities.
- Customer capabilities remain active.
- Artisan membership state becomes `ACTIVE`.
- First workshop ownership is recorded.
- An audit event records activation.
- The user can open the artisan dashboard without signing up again.
- Optional document verification may remain `NOT_SUBMITTED` or `PENDING` without creating another user account.

### US-ART-005 - Submit optional verification documents

**As an artisan**, I want to submit optional documents so that the platform can verify my identity, organisation, or craft claims.

#### Main flow

1. Artisan opens verification settings.
2. Artisan selects document type.
3. Artisan uploads one or more allowed files.
4. Backend validates and stores files privately.
5. Verification status becomes `PENDING`.
6. Authorised reviewer is notified.

#### Acceptance criteria

- Submission is optional unless a later business rule requires verification for a specific feature.
- Documents are private.
- Only the owner and authorised reviewers can access them.
- File access is audited where required.
- Uploading documents does not create a new artisan membership.

### US-ART-006 - Review optional verification

**As an authorised reviewer**, I want to review artisan documents so that I can verify supported claims.

#### Reviewer actions

- Mark verified.
- Request changes.
- Reject submitted evidence.
- Add an internal and user-visible reason where appropriate.

#### Acceptance criteria

- Every decision is audited.
- Rejection of verification does not delete historical data.
- Rejection of verification does not remove customer capabilities.
- Any restriction placed on selling must follow an explicit policy and audited state change.

### US-ART-007 - Edit artisan profile

**As an artisan**, I want to edit my artisan profile so that public and operational information remains accurate.

#### Acceptance criteria

- Only the owner or authorised administrator can edit the profile.
- Private and public fields are clearly separated.
- Sensitive changes may require re-verification.
- Changes are audited when required.

### US-ART-008 - Suspend artisan membership

**As an authorised administrator**, I want to suspend artisan membership so that new selling activity can be stopped without deleting the user or historical records.

#### Acceptance criteria

- New product submissions and new selling activity are blocked according to policy.
- Existing customer capabilities remain available unless the user account itself is suspended.
- Historical orders, products, workshops, payments, and audit data remain accessible to authorised actors.
- Suspension requires a reason.
- The action is audited.

## 7. Workshop Stories

### US-WS-001 - View my workshops

**As an artisan**, I want to view all my workshops so that I can manage each production or business location.

#### Acceptance criteria

- The artisan sees only workshops they own.
- Workshop status, product count, and required operational summary are shown.
- Empty, loading, and error states exist.

### US-WS-002 - Create additional workshop

**As an artisan**, I want to create additional workshops so that I can organise products under different locations or workshop identities.

#### Main flow

1. Artisan opens **My Workshops**.
2. Artisan clicks **Create Workshop**.
3. Artisan enters required information.
4. Backend verifies active artisan membership.
5. Backend validates the request.
6. Workshop is created and linked to the artisan.
7. Audit event is recorded.

#### Acceptance criteria

- One artisan can own many workshops.
- Every workshop belongs to exactly one artisan membership.
- A regular customer without active artisan membership cannot create a workshop.
- Repeating the request with the same idempotency key does not create duplicates.

### US-WS-003 - Update workshop

**As an artisan**, I want to update an owned workshop so that its information remains accurate.

#### Acceptance criteria

- The backend verifies ownership.
- Another artisan cannot update the workshop.
- Public and private fields are handled separately.
- Relevant changes are audited.

### US-WS-004 - Activate or deactivate workshop

**As an artisan**, I want to deactivate a workshop so that it stops accepting new product activity without deleting history.

#### Acceptance criteria

- Deactivation does not delete products or historical orders.
- New products cannot be created for a deactivated workshop.
- Existing products follow configured suspension or archival rules.
- Reactivation is allowed only when ownership and membership are valid.

### US-WS-005 - Delete empty workshop

**As an artisan**, I want to delete an unused workshop so that incorrect or abandoned drafts do not remain in my account.

#### Acceptance criteria

- Permanent deletion is allowed only when the workshop has no products, inventory, orders, financial records, or other protected history.
- A workshop with protected history must be deactivated or archived instead.
- The action is audited.

### US-WS-006 - Prevent cross-artisan access

**As the platform**, I want to prevent artisans from accessing other artisans' private workshops so that ownership boundaries are enforced.

#### Acceptance criteria

- Reading another artisan's private workshop returns `403` or a non-disclosing `404` according to policy.
- Updating or deleting another artisan's workshop returns `403`.
- Creating a product for another artisan's workshop returns `403`.
- Hiding frontend controls is not considered sufficient authorisation.

## 8. Product Stories

### US-PROD-001 - Create product draft for workshop

**As an artisan**, I want to create a product for one of my workshops so that I can prepare it for moderation and sale.

#### Preconditions

- User has active artisan membership.
- Selected workshop is active.
- Selected workshop belongs to the authenticated artisan.

#### Artisan enters

- Workshop.
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

#### Acceptance criteria

- Every product belongs to exactly one workshop.
- The product owner is derived from workshop ownership.
- Product is created as `DRAFT`.
- Product remains invisible while `DRAFT`.
- Another artisan cannot create a product for the selected workshop.

### US-PROD-002 - Edit product draft

**As an artisan**, I want to edit my draft product so that I can complete or correct it before submission.

#### Acceptance criteria

- Only the owner can edit the draft.
- Submitted immutable versions are not overwritten.
- Changes are validated.
- A product cannot be moved to a workshop owned by another artisan.

### US-PROD-003 - Move product between owned workshops

**As an artisan**, I want to move an eligible product between my workshops so that product organisation remains accurate.

#### Acceptance criteria

- Source and destination workshops belong to the same artisan.
- Destination workshop is active.
- Product with protected inventory or order history follows explicit migration rules.
- Historical order items preserve the workshop identity captured at purchase time.
- The action is audited.

### US-PROD-004 - Submit product for moderation

**As an artisan**, I want to submit a complete product so that a moderator can review it.

#### Main flow

1. Artisan opens a complete draft.
2. Backend verifies product and workshop ownership.
3. Backend validates mandatory fields.
4. Product changes to `PENDING_REVIEW`.
5. Moderator is notified.
6. Submitted version is preserved.

#### Failures

- Missing mandatory field.
- Invalid price or currency.
- Missing required media.
- Inactive workshop.
- Suspended artisan membership.
- Ownership failure.

### US-PROD-005 - Moderate product

**As a moderator**, I want to review submitted products so that only acceptable listings can proceed toward sale.

#### Moderator actions

- Approve.
- Reject.
- Request changes.
- Suspend an active product.

#### Acceptance criteria

- Reasons are mandatory for rejection, requested changes, and suspension.
- Decision is audited.
- Approval does not create stock.
- Approval alone does not make the product sellable.

### US-PROD-006 - Respond to requested changes

**As an artisan**, I want to update a product after requested changes so that I can resubmit it.

#### Acceptance criteria

- Artisan sees the reviewer reason.
- Artisan edits a new draft based on the reviewed version.
- Previous submitted versions remain preserved.
- Resubmission returns the product to `PENDING_REVIEW`.

### US-PROD-007 - Activate sellable product

**As the platform**, I want to activate only products that satisfy all selling rules so that unavailable or unapproved products cannot be purchased.

A product becomes sellable only when:

- Product is approved.
- Product is not suspended or archived.
- Owning workshop is active.
- Artisan membership is active.
- At least one accepted stock unit is available, unless an explicit made-to-order path applies.
- Price and currency are valid.
- Required public media exists.

### US-PROD-008 - Archive product

**As an artisan or authorised administrator**, I want to archive an eligible product so that it is no longer offered while history remains preserved.

#### Acceptance criteria

- Archived products are not sellable.
- Historical order items remain readable.
- Inventory is handled according to warehouse policy.
- Permanent deletion is prevented when protected history exists.

## 9. Warehouse and Inventory Stories

### US-WH-001 - Receive stock

**As a warehouse agent**, I want to receive product stock so that delivered items can enter inspection.

#### Main flow

1. Agent identifies product, workshop, and supplier or artisan.
2. Agent records batch or parcel.
3. Agent records received quantity.
4. System creates reception record.
5. Quantity becomes `RECEIVED_PENDING_INSPECTION`.

#### Acceptance criteria

- Received quantity is not immediately sellable.
- Product and workshop ownership context is preserved.
- Duplicate reception requests are safely handled.

### US-WH-002 - Inspect batch

**As a warehouse agent**, I want to inspect a received batch so that acceptable items become available and unacceptable items are separated.

#### Main flow

1. Agent opens pending batch.
2. Agent enters accepted, rejected, quarantined, and damaged quantities.
3. Agent adds required reasons and evidence.
4. Backend validates the quantity equation.
5. Transaction creates inventory movements.
6. Accepted stock becomes available.

#### Acceptance criteria

```text
received quantity
= accepted
+ rejected
+ quarantined
+ damaged
```

- Quantities cannot be negative.
- Accepted stock becomes available only after the transaction succeeds.
- Product activation rules are reevaluated after accepted stock changes.

### US-WH-003 - Correct inventory

**As an authorised inventory manager**, I want to correct inventory using explicit movements so that stock remains traceable.

#### Required data

- Product or batch.
- Adjustment quantity.
- Reason.
- Actor.
- Timestamp.
- Correlation or reference information when applicable.

#### Acceptance criteria

- Stock totals cannot be edited silently.
- Every correction creates an inventory movement and audit event.
- Available stock cannot become negative.
- Permissions are enforced by the backend.

### US-WH-004 - View inventory by workshop

**As an authorised artisan or warehouse user**, I want to view permitted inventory grouped by workshop so that stock can be understood in its ownership context.

#### Acceptance criteria

- Artisan sees only inventory for products in owned workshops.
- Warehouse roles see inventory according to assigned permissions.
- Private warehouse details are not exposed to customers.

## 10. Cart and Checkout Stories

### US-CART-001 - Add item to cart

**As a user**, I want to add a product to my cart so that I can purchase it later.

Adding to cart does not reserve stock.

#### Cart displays

- Product.
- Workshop.
- Artisan.
- Quantity.
- Current display price.
- Currency.
- Subtotal.
- Availability warning.

#### Acceptance criteria

- An artisan can add products to the cart like any other user.
- A user may add their own product to the cart only if business policy permits; the backend must not assume artisan status blocks customer actions.
- Current product state is revalidated at checkout.

### US-CART-002 - Update or remove cart item

**As a user**, I want to change quantities or remove items so that my cart matches my intended purchase.

#### Acceptance criteria

- Quantity is positive and within configured limits.
- Removing a missing item is idempotent.
- Cart updates do not reserve stock.

### US-CHECK-001 - Checkout

**As a customer**, I want to check out so that I can create and pay for an order.

#### Main flow

1. User signs in.
2. Backend reloads current products, workshop state, artisan state, prices, and availability.
3. User chooses delivery address.
4. Backend validates delivery.
5. Backend calculates trusted totals.
6. Backend creates stock reservations.
7. Backend creates pending order.
8. Backend creates payment attempt.
9. User continues to payment.

#### Failures

- Product inactive.
- Workshop inactive.
- Artisan membership suspended.
- Price changed.
- Insufficient stock.
- Invalid address.
- Unsupported destination.
- Duplicate request.

### US-CHECK-002 - Concurrent purchase of final item

**As the platform**, I want concurrent checkout protected so that the same final stock unit cannot be sold twice.

Given one available item:

1. Customer A and Customer B submit checkout concurrently.
2. Database transaction or lock protects stock.
3. Only one reservation succeeds.
4. The other request receives `OUT_OF_STOCK`.
5. Available stock never becomes negative.

### US-CHECK-003 - Artisan purchases as customer

**As an artisan**, I want to purchase products and track my order using the same account so that becoming an artisan does not remove customer functionality.

#### Acceptance criteria

- Artisan can use cart and checkout.
- Artisan can use saved delivery addresses.
- Artisan can pay and receive customer notifications.
- Artisan can view the order in the customer order history.
- Seller-side visibility is separate from customer-side ownership.

## 11. Payment Stories

### US-PAY-001 - Start payment

**As a customer**, I want to start payment using the trusted order total so that the payment amount cannot be manipulated by the frontend.

#### Acceptance criteria

- Backend sends trusted amount and currency to the payment provider adapter.
- Frontend never supplies a trusted final total.
- Payment attempt is linked to one order.
- Repeated requests are idempotent where required.

### US-PAY-002 - Return from provider

**As a customer**, I want to return to the platform after payment so that I can see the current payment state.

#### Acceptance criteria

- Browser return page may display `PENDING`.
- A success query parameter does not confirm payment.
- Final confirmation depends on trusted provider verification or webhook processing.

### US-PAY-003 - Process provider webhook

**As the platform**, I want to process trusted payment events so that orders change state safely.

#### Main flow

1. Verify signature and timestamp.
2. Reject unknown or invalid event.
3. Detect duplicate event.
4. Match order, amount, and currency.
5. Validate state transition.
6. Mark payment captured.
7. Mark order paid.
8. Convert reservation into committed allocation.
9. Queue notification through outbox.

#### Acceptance criteria

- Duplicate webhook does not duplicate payment, allocation, or notification.
- Amount and currency mismatch prevents capture processing.
- State transition is transactional.

### US-PAY-004 - Payment failure or expiry

**As the platform**, I want failed or expired payments handled so that reserved stock is released.

#### Main flow

1. Payment becomes failed or expired.
2. Order becomes `PAYMENT_FAILED` or remains eligible for retry according to policy.
3. Stock reservation is released.
4. Customer is notified.

### US-PAY-005 - Refund payment

**As an authorised actor**, I want to refund eligible payments so that cancellations and returns can be settled correctly.

#### Acceptance criteria

- Refund amount is validated against captured and previously refunded amounts.
- Payment becomes `PARTIALLY_REFUNDED` or `REFUNDED`.
- Order state is updated according to policy.
- The action is audited.

## 12. Order, Seller Visibility, and Fulfilment Stories

### US-ORD-001 - View customer order

**As a customer**, I want to view my order so that I can understand its payment, fulfilment, and shipment status.

#### Acceptance criteria

- Customer sees only orders owned by their user account.
- An artisan sees their own purchases in the same customer order history.
- Customer view includes order items, totals, payment status, delivery address summary, and tracking information as allowed.
- Customer cannot read another user's order.

### US-ORD-002 - View artisan order items

**As an artisan**, I want to view order items related to products from my workshops so that I can understand sales and permitted fulfilment information.

#### Acceptance criteria

- Artisan sees only order items whose products belong to their workshops.
- Artisan does not automatically see unrelated items in the same customer order.
- Artisan sees only minimum customer information required by the business workflow.
- Sensitive payment data is not exposed.
- Seller-side access does not make the artisan the owner of the customer order.

### US-ORD-003 - Separate buyer and seller contexts

**As the platform**, I want buyer and seller contexts separated so that an artisan purchasing products does not receive unauthorised seller information.

#### Acceptance criteria

- Customer ownership is based on `order.user_id` or equivalent buyer reference.
- Artisan visibility is based on workshop and product ownership of individual order items.
- An artisan buying their own or another artisan's product sees the purchase in customer context.
- Seller context never grants access to unrelated buyer orders.

### US-FUL-001 - Pick and pack

**As a warehouse agent**, I want to pick and pack paid allocations so that orders can be shipped.

#### Main flow

1. Paid order enters warehouse queue.
2. Agent picks reserved or allocated stock.
3. Agent records packed parcel.
4. Order becomes `READY_TO_SHIP` when fulfilment rules are satisfied.

#### Acceptance criteria

- Picked quantity cannot exceed allocated quantity.
- Partial fulfilment follows explicit policy.
- Workshop and product ownership data remains traceable.

### US-SHIP-001 - Create shipment

**As a fulfilment agent**, I want to create a shipment so that the customer can track delivery.

#### Main flow

1. Agent chooses carrier or manual shipment.
2. Tracking reference is stored.
3. Order becomes `SHIPPED` when applicable.
4. Customer receives notification.

### US-SHIP-002 - Record delivery event

**As the platform**, I want to record shipment events so that customers receive accurate tracking information.

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

#### Acceptance criteria

- Duplicate events do not duplicate notifications.
- Invalid backwards transitions are rejected or explicitly handled.
- Customer can view permitted tracking events.
- Artisan can track purchases made through their customer account.

## 13. Cancellation and Return Stories

### US-CAN-001 - Cancel unpaid order

**As a customer**, I want to cancel an eligible unpaid order so that I am not charged and reserved stock is released.

#### Acceptance criteria

- Customer may cancel only their own eligible order.
- Reservation is released.
- Payment attempt is cancelled or expired where supported.
- Repeating cancellation is idempotent.

### US-CAN-002 - Cancel paid order before shipment

**As a customer or authorised support actor**, I want to cancel an eligible paid order before shipment so that a refund can begin.

#### Acceptance criteria

- Eligibility is checked.
- Refund starts when required.
- Stock is returned through inventory movement when applicable.
- Action is audited.
- Shipped orders cannot be silently cancelled.

### US-RET-001 - Record return

**As an authorised fulfilment actor**, I want to record a returned shipment so that refund and inventory decisions can be processed.

#### Acceptance criteria

- Return event is linked to shipment and order.
- Returned items do not automatically become sellable.
- Returned stock requires inspection before becoming available again.
- Refund follows explicit policy.

## 14. Made-to-Order Stories

### US-CUSTOM-001 - Submit request

**As a customer**, I want to submit a made-to-order request so that an artisan can evaluate custom work.

#### Customer provides

- Artisan, workshop, or product.
- Description.
- Quantity.
- Desired date.
- Reference images or documents.
- Contact preference.

#### Acceptance criteria

- An artisan may submit a request as a customer using the same account.
- Request is visible only to authorised participants.
- Uploaded references follow file security rules.

### US-CUSTOM-002 - Artisan responds

**As an artisan**, I want to respond to a request associated with my workshop so that I can clarify, decline, or quote the work.

#### Artisan actions

- Request clarification.
- Decline.
- Send a quote with price, scope, expiry, estimated production time, and workshop.

#### Acceptance criteria

- Artisan can respond only when the request targets their profile, workshop, or eligible product.
- Quote versions are immutable.
- Private customer information is limited to what is necessary.

### US-CUSTOM-003 - Customer accepts quote

**As a customer**, I want to accept a valid quote so that trusted checkout and production can begin.

#### Acceptance criteria

- Expired quote cannot be accepted.
- Accepted quote becomes trusted checkout data.
- Payment and production workflow starts only after explicit acceptance.
- Quote cannot be changed after acceptance without a new version and renewed acceptance.

## 15. Notification Stories

### US-NOTIF-001 - Receive account and membership notifications

**As a user**, I want notifications about important account and artisan membership events so that I know when action is required.

Events may include:

- Account verification.
- Artisan membership activation.
- Verification changes requested.
- Verification decision.
- Artisan suspension.
- Workshop status change.

### US-NOTIF-002 - Receive commerce notifications

**As a customer**, I want notifications about orders and shipments so that I can track my purchase.

Events may include:

- Payment pending.
- Payment confirmed.
- Payment failed.
- Order cancelled.
- Order packed.
- Shipment created.
- Delivery event.
- Refund completed.

### US-NOTIF-003 - Receive artisan notifications

**As an artisan**, I want notifications about my workshops and products so that I can respond to platform actions.

Events may include:

- Product submitted.
- Changes requested.
- Product approved or rejected.
- Product suspended.
- New made-to-order request.
- Quote accepted.
- Seller-side order item status where applicable.

#### Acceptance criteria

- Notifications are generated through an outbox or equivalent reliable mechanism where required.
- Duplicate events do not create duplicate notifications.
- Notifications do not expose private data to unauthorised recipients.

## 16. Administration and Moderation Stories

### US-ADMIN-001 - Manage permissions

**As an authorised administrator**, I want to manage roles and permissions so that administrative access is controlled independently from customer and artisan membership.

#### Acceptance criteria

- Customer capability is the default authenticated-user capability.
- Artisan membership is an additive business membership.
- Admin, moderator, warehouse, finance, and support permissions are separate privileged permissions.
- Permission changes require backend authorisation.
- Permission changes create audit events.

### US-ADMIN-002 - Suspend user account

**As an authorised administrator**, I want to suspend a user account so that all access can be blocked when necessary.

#### Acceptance criteria

- User account suspension is separate from artisan membership suspension.
- Suspending the user account blocks sign-in according to policy.
- Historical records remain preserved.
- Reason and audit event are required.

### US-ADMIN-003 - Suspend artisan, workshop, or product

**As an authorised administrator or moderator**, I want to suspend a selling entity so that new sales can be prevented without deleting history.

#### Acceptance criteria

- Artisan suspension blocks new selling activity across owned workshops according to policy.
- Workshop suspension blocks eligible selling activity for that workshop.
- Product suspension blocks that product.
- Historical orders are not deleted.
- Suspension requires a reason.
- Action is audited.

### US-ADMIN-004 - Review audits

**As an authorised auditor**, I want to search audit events so that important actions can be traced.

#### Search criteria

- Actor.
- Target.
- Event type.
- Date range.
- Correlation ID.
- User.
- Artisan membership.
- Workshop.
- Product.
- Order.

### US-ADMIN-005 - Review optional artisan verification

**As an authorised reviewer**, I want a review queue for submitted documents so that optional verification can be processed without changing the one-account model.

#### Acceptance criteria

- Reviewer can access only assigned or permitted submissions.
- Documents remain private.
- Decision, reason, actor, and timestamp are audited.
- Verification decisions do not create another artisan account.

## 17. Authorisation and Security Rules

### 17.1 General rules

- All permissions are enforced by the backend.
- Frontend visibility is not authorisation.
- Resource ownership is checked for every read and write operation.
- Sensitive operations use least privilege.
- Private files use protected storage and authorised access paths.
- Important state changes are audited.
- Idempotency is used for operations vulnerable to duplicate requests.

### 17.2 User and artisan rules

- Every active registered user can use customer features.
- Artisan membership adds capabilities to the same user.
- Becoming an artisan must not remove cart, order, payment, address, cancellation, or tracking access.
- Only the user can manage their own customer profile and addresses, except authorised administrative actions.
- Only the artisan owner or authorised administrator can manage the artisan profile.

### 17.3 Workshop rules

- Only an active artisan can create a workshop.
- Only the owner or authorised administrator can modify a workshop.
- An artisan cannot access another artisan's private workshop data.
- Workshop ownership cannot be trusted from frontend input alone.
- Backend derives and validates artisan ownership.

### 17.4 Product rules

- Every product belongs to one workshop.
- Product ownership is derived from workshop ownership.
- Artisan can manage only products in owned workshops.
- Moderator decisions require moderator permission.
- Product approval does not bypass inventory inspection.

### 17.5 Order rules

- Customer can read only orders they own.
- Artisan can read only permitted order-item information for products in owned workshops.
- Artisan seller access does not grant access to the entire customer order.
- Warehouse sees fulfilment information according to role.
- Finance sees payment information according to role.

### 17.6 File rules

- Optional artisan documents are private.
- Made-to-order references are private to authorised participants.
- Public product and workshop media are stored separately or exposed through explicitly public access.
- File type, size, and malware controls are applied according to platform policy.

## 18. Critical State Machines

### 18.1 Artisan membership

```text
ACTIVE
-> SUSPENDED
-> ACTIVE
-> CLOSED
```

Membership is created as `ACTIVE` only after the first valid workshop is created.

### 18.2 Artisan verification

```text
NOT_SUBMITTED
-> PENDING
-> VERIFIED

PENDING
-> CHANGES_REQUESTED
-> PENDING

PENDING
-> REJECTED
-> PENDING
```

### 18.3 Workshop

```text
ACTIVE
-> INACTIVE
-> ACTIVE
-> SUSPENDED
-> ARCHIVED
```

### 18.4 Product

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

`APPROVED -> ACTIVE` occurs only when all activation rules are satisfied.

### 18.5 Reception

```text
RECEIVED_PENDING_INSPECTION
-> INSPECTING
-> ACCEPTED
-> PARTIALLY_ACCEPTED
-> REJECTED
-> QUARANTINED
-> CLOSED
```

### 18.6 Reservation

```text
ACTIVE
-> CONSUMED
-> RELEASED
-> EXPIRED
```

### 18.7 Payment

```text
CREATED
-> PENDING
-> CAPTURED
-> FAILED
-> CANCELLED
-> PARTIALLY_REFUNDED
-> REFUNDED
```

### 18.8 Order

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
SHIPPED -> RETURNED
```

### 18.9 Shipment

```text
PENDING
-> SHIPPED
-> IN_TRANSIT
-> OUT_FOR_DELIVERY
-> DELIVERED

IN_TRANSIT -> DELIVERY_FAILED
OUT_FOR_DELIVERY -> DELIVERY_FAILED
DELIVERY_FAILED -> IN_TRANSIT
DELIVERY_FAILED -> RETURNED
DELIVERED -> RETURNED
```

## 19. Required End-to-End Tests

### 19.1 Authentication and customer tests

1. Visitor registers and receives customer capabilities.
2. User signs in and views the combined account dashboard.
3. User signs out and the active session is revoked.
4. Customer cannot read another customer's order.
5. Arabic RTL page smoke test.

### 19.2 Artisan membership tests

6. Registered customer clicks **Become an Artisan**.
7. Existing user information is reused without creating a second account.
8. Artisan onboarding fails when no workshop is created.
9. Artisan onboarding succeeds without supporting documents.
10. Artisan onboarding succeeds with optional valid documents.
11. Artisan membership activates after creation of the first valid workshop.
12. Customer capabilities remain available after artisan activation.
13. Duplicate onboarding completion does not create duplicate membership or workshop.
14. Existing customer orders remain visible after becoming an artisan.
15. Artisan can place and track a new order using the same account.
16. Suspended artisan membership does not remove allowed customer capabilities.

### 19.3 Workshop tests

17. Artisan creates an additional workshop.
18. Artisan views all owned workshops.
19. Artisan updates an owned workshop.
20. Artisan cannot read another artisan's private workshop.
21. Artisan cannot update another artisan's workshop.
22. Regular customer without artisan membership cannot create a workshop.
23. Deactivated workshop preserves products and historical orders.
24. Workshop with protected history cannot be permanently deleted.
25. Empty unused workshop can be deleted according to policy.

### 19.4 Product tests

26. Artisan creates a product for an owned workshop.
27. Artisan cannot create a product for another artisan's workshop.
28. Every product is linked to exactly one workshop.
29. Product remains invisible while `DRAFT`.
30. Product moderation requests changes and preserves submitted version.
31. Artisan resubmits after requested changes.
32. Product approval does not make the product sellable without accepted stock.
33. Product cannot become active when workshop is inactive.
34. Product cannot become active when artisan membership is suspended.
35. Product can move only between workshops owned by the same artisan when eligible.
36. Historical order items preserve original workshop information.

### 19.5 Warehouse and inventory tests

37. Warehouse receives stock as `RECEIVED_PENDING_INSPECTION`.
38. Product is not sellable before inspection.
39. Warehouse partially accepts a batch.
40. Accepted stock becomes available transactionally.
41. Invalid inspection quantity equation is rejected.
42. Inventory correction creates movement and audit event.
43. Available stock never becomes negative.

### 19.6 Cart, checkout, and order tests

44. Standard customer purchase succeeds.
45. Artisan purchases as a customer using the same account.
46. Cart does not reserve stock.
47. Checkout reloads trusted product and price data.
48. Concurrent purchase of final item allows only one reservation.
49. Seller-side artisan view exposes only owned order items.
50. Artisan seller view does not expose unrelated items from the same order.
51. Direct API authorisation bypass attempt against another order returns `403`.
52. Direct API authorisation bypass attempt against another workshop returns `403`.
53. Direct API authorisation bypass attempt against another artisan's product returns `403`.

### 19.7 Payment and fulfilment tests

54. Payment provider return does not confirm payment by query parameter alone.
55. Delayed payment webhook captures payment correctly.
56. Duplicate payment webhook is idempotent.
57. Payment amount or currency mismatch is rejected.
58. Payment failure releases stock reservation.
59. Paid order enters fulfilment.
60. Picked quantity cannot exceed allocation.
61. Shipment is created and tracking is visible to the customer.
62. Artisan can track purchases made with their account.
63. Duplicate shipment event does not duplicate notification.
64. Paid order fulfilment and delivery completes successfully.

### 19.8 Cancellation, return, and made-to-order tests

65. Customer cancels eligible unpaid order and reservation is released.
66. Eligible paid order cancellation starts refund.
67. Shipped order cannot be silently cancelled.
68. Returned stock is not sellable before reinspection.
69. Customer submits made-to-order request.
70. Artisan responds only to requests targeting owned artisan resources.
71. Expired quote cannot be accepted.
72. Accepted quote becomes trusted checkout data.

### 19.9 Verification and administration tests

73. Optional verification document is stored privately.
74. Unauthorised user cannot access verification document.
75. Reviewer verifies an artisan submission and action is audited.
76. Verification rejection does not remove customer capabilities.
77. Admin suspends artisan membership with a reason.
78. Admin suspends workshop or product without deleting history.
79. Role or permission change creates audit event.
80. Audit search returns events filtered by actor, target, event type, date, and correlation ID.

## 20. MVP Completion Criteria

The MVP user-story scope is complete when:

- One account supports both customer and artisan activity.
- A customer can become an artisan through **Become an Artisan**.
- The first workshop is mandatory during onboarding.
- Supporting documents are optional.
- Artisan membership activates on the existing user account.
- Artisan retains all customer capabilities.
- An artisan can own multiple workshops.
- Every product belongs to one workshop.
- Product ownership follows workshop ownership.
- Product moderation and warehouse inspection remain separate controls.
- Customer, artisan, moderator, warehouse, finance, support, and administrator permissions are enforced by the backend.
- Buyer order ownership and artisan seller visibility are separated.
- Payment, reservation, inventory, shipment, cancellation, and audit workflows are idempotent and traceable where required.
- Required end-to-end tests pass.
