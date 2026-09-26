# PRD-12 · Cart & Wishlist Persistence

| | |
| --- | --- |
| **Epic** | 12 — Cart & Wishlist Persistence |
| **Stories** | US-CART-001 … US-CART-002; wishlist capability |
| **Priority** | High |
| **Status** | ✅ Shipped |
| **Surfaces** | `apps/web/app/[locale]/cart`, cart drawer, wishlist |
| **API** | `GET /api/v1/cart`, `POST /api/v1/cart/items`, `POST /api/v1/cart/merge`, `PATCH /api/v1/cart/items/:productId`, `DELETE /api/v1/cart/items/:productId`, `GET /api/v1/wishlist`, `POST/DELETE /api/v1/wishlist/items/:productId` |

## 1. Summary

Let visitors and customers collect products, update quantities, remove items, and save products for later while keeping price and availability explicitly provisional until checkout.

## 2. Problem

A client-only cart can be lost at sign-in, show stale prices, imply stock is reserved, or expose another customer's cart when persistence is added without ownership checks.

## 3. Goals

- Support anonymous cart state and authenticated server persistence.
- Define safe merge-on-login behavior.
- Revalidate product state when the cart opens and before checkout.
- Provide wishlist persistence only for authenticated owners.

### Non-goals

- Reserving inventory when an item is added to the cart.
- Coupon, promotion, or recommendation engines before their business rules are confirmed.

## 4. Users

- Visitor building a temporary cart.
- Customer managing a persistent cart and wishlist.
- Artisan purchasing as a customer using the same account.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Add active products to a cart without reserving stock. | US-CART-001 |
| R2 | Show product, workshop, artisan, quantity, provisional price, currency, subtotal, and availability warning. | US-CART-001 |
| R3 | Update positive quantities and remove items idempotently. | US-CART-002 |
| R4 | Merge anonymous and authenticated carts according to one explicit policy. | Cart architecture rule |
| R5 | Revalidate current product state, price, and availability on cart open and checkout. | US-CHECK-001 |
| R6 | Allow an authenticated user to save and remove wishlist items with ownership checks. | Wishlist architecture capability |

## 6. Flow

```text
Browse → add to cart → local/server cart → sign in and merge → revalidate → checkout
Browse → save wishlist item → authenticated wishlist → remove or open product
```

## 7. Technical notes

- Anonymous cart state is held in local storage; authenticated sessions persist through owner-scoped APIs and merge local items on login.
- Client cart state is UX only; the cart API exposes current product price, currency, active state, and availability warnings, while checkout recalculates trusted totals.
- Cart and wishlist responses use DTOs and owner-scoped queries. Repeated wishlist saves and cart removals are idempotent.

## 8. Success metrics

- Cart persistence and merge success rate.
- Cart-to-checkout conversion after revalidation.
- Zero unauthorized cart or wishlist reads.

## 9. Risks & open questions

- Anonymous cart merge conflicts use the documented sum-and-cap-at-100 policy; changing that policy requires product approval.
- Wishlist inclusion in the MVP and back-in-stock notifications need confirmation.

## Source traceability

`REQ-CART-001`, `US-CART-001..002`, checkout rules in `docs/architecture.md`.
