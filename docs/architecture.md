# AISHA Clean Feature-Based Architecture

## 1. Purpose

This document defines the implementation architecture of the AISHA monorepo.

AISHA uses:

```text
Monorepo
+ Feature-Based Architecture
+ Clean Architecture
+ Modular Monolith Backend
+ Next.js App Router Frontend
+ Ports and Adapters
+ Dependency Injection
+ Optional MVVM-style ViewModel layer on the frontend
```

Business logic must not be placed in Next.js route files, React presentational components, frontend API clients, Fiber handlers, route registration files, SQLBoiler generated models, migrations, or generic utility folders.

This is the structural source of truth for the monorepo. `docs/frontend.md` and `docs/backend.md` carry the full prose rationale (visual direction, UX behaviour, domain reasoning); every concrete requirement from both has been bound into a folder, a file, or a rule below, so this document alone is buildable.

- **Part I — Frontend** (§4–§31): Next.js App Router, feature-based, ViewModel-optional, fully bound to concrete design tokens, routing, component inventory, forms, and CI commands — this is the MVP-ready version of the frontend, not just the conceptual layering.
- **Part II — Backend** (§32–§60): Go Fiber v3, feature-based Clean Architecture, fully bound to concrete packages, code shapes, domain rules, and CI commands — the MVP-ready version of the backend.

## 2. Dependency Rules

```text
Outer layers may depend on inner layers.
Inner layers must never depend on outer layers.
```

Backend:

```text
Server / Data / Infrastructure
             ↓
        Application
             ↓
           Domain
```

Frontend:

```text
App Routes
    ↓
Views and ViewModels
    ↓
Application / Domain
    ↓
Repository Contracts
    ↑
Data Repository Implementations
    ↓
Shared API Infrastructure
```

The frontend communicates with the backend only through stable API contracts.

## 3. Monorepo Tree

```text
aisha/
├── apps/
│   ├── web/
│   └── api/
├── packages/
│   ├── contracts/
│   ├── eslint-config/
│   ├── typescript-config/
│   └── testing/
├── infrastructure/
│   ├── docker/
│   ├── compose/
│   ├── scripts/
│   └── monitoring/
├── docs/
│   ├── architecture.md
│   ├── backend.md
│   ├── frontend.md
│   ├── workflows.md
│   ├── requirements.md
│   ├── jira_project.md
│   ├── deployment.md
│   ├── security.md
│   ├── ui_design.md
│   ├── userstory.md
├── .github/workflows/
├── .env.example
├── Makefile
├── README.md
```

# Part I — Frontend

## 4. Frontend Goal and Principles

Build a premium, visually distinctive, secure, multilingual e-commerce storefront for Algerian artisanal products. Customers discover products, understand their cultural origin, learn about the artisan, and complete a purchase in a small number of clear steps. AISHA should feel like a modern international fashion and lifestyle store while preserving an authentic Algerian identity — the editorial minimalism of Zara, the bold product presentation of Nike, and Algerian craftsmanship, patterns, and cultural storytelling, without directly copying any of them (§6).

Principles:

```text
Next.js App Router
Server Components by default; Client Components only when interactivity requires them
Strict TypeScript
shadcn/ui components
Tailwind CSS
React Hook Form + Zod for forms
Translation keys for all visible text
Full Arabic RTL support
Responsive and accessible design
Mobile-first implementation
No business-critical rules only in the frontend
Reusable design-system components over page-specific styling
```

Required frontend packages:

```text
next                        → App Router, Server/Client Components
react / react-dom
typescript                  → strict mode
tailwindcss                 → utility classes, consumes tokens from styles/tokens.css (§6)
shadcn/ui (+ Radix primitives) → core/components/ui (§10)
react-hook-form + zod       → features/<feature>/schemas, components/forms (§17)
next-intl (recommended)     → satisfies the /[locale]/ routing + messages/*.json convention in §8, §26; substitute an equivalent i18n router only if the routing and message-file shape below are preserved
lucide-react                → icon set paired with shadcn/ui
```

## 5. Frontend Base Structure

```text
apps/web/
├── app/
├── core/
├── features/
├── messages/
├── public/
├── tests/
├── next.config.ts
├── middleware.ts
├── components.json
├── package.json
├── tsconfig.json
└── .env.example
```

### Responsibilities

```text
app/       Routes, layouts, metadata, loading, errors, page composition
core/      Shared technical infrastructure and reusable UI
features/  Business feature modules
messages/  Translation files
public/    Static assets and PWA files
tests/     Unit, integration, E2E, and visual-regression tests (§27)
```

`components.json` is the shadcn/ui config (component output paths, alias mapping) — required because §10, §25 assume shadcn primitives live at `core/components/ui`. `middleware.ts` owns locale detection and redirect to `/[locale]/...` (§8); it must not implement authentication or authorization decisions — those stay server-side and authoritative (§19, backend §49).

## 6. Design System and Visual Tokens

Create the design system before building complete pages. AISHA should feel premium, authentic, contemporary, cultural, editorial, warm, trustworthy, and handmade-but-professionally-presented. Avoid generic marketplace layouts, excessive gradients, heavy shadows, crowded grids, decorative animations that slow down shopping, dashboard styling on customer-facing pages, and using every shadcn/ui component without a clear purpose.

```text
core/components/
├── ui/          shadcn/ui primitives (§25)
├── layout/      header, footer, nav (§9)
├── commerce/    Money, AvailabilityLabel, VerificationBadge (§10)
├── editorial/   reusable storytelling primitives (§10)
├── feedback/    EmptyState, ErrorState, LoadingSkeleton (§24)
└── forms/       shared form primitives (§17)

apps/web/app/
├── globals.css      imports tokens.css + utilities.css, Tailwind base
└── styles/
    ├── tokens.css   CSS custom properties below
    └── utilities.css
```

Global CSS stays under `app/` (not `core/`) because Next.js requires the root layout to import it directly — `app/[locale]/layout.tsx` is that root layout (§8).

### 6.1 Colour Palette

Restrained palette — most of the storefront stays neutral so product imagery remains dominant. Do not apply accent colours to every element.

```text
Background:        warm off-white
Surface:            white
Primary text:       near black
Secondary text:     warm grey
Muted surface:      light stone
Border:             soft neutral grey
Primary accent:     deep terracotta
Secondary accent:   olive green
Premium accent:     muted gold
Error:               deep red
Success:             forest green
```

```css
:root {
  --background: 40 33% 98%;
  --foreground: 20 10% 10%;

  --card: 0 0% 100%;
  --card-foreground: 20 10% 10%;

  --muted: 35 18% 93%;
  --muted-foreground: 25 8% 42%;

  --border: 30 12% 86%;
  --input: 30 12% 86%;

  --primary: 14 48% 38%;
  --primary-foreground: 40 33% 98%;

  --secondary: 75 18% 30%;
  --secondary-foreground: 40 33% 98%;

  --accent: 40 40% 56%;
  --accent-foreground: 20 10% 10%;

  --destructive: 0 62% 42%;
  --destructive-foreground: 0 0% 100%;

  --radius: 0.25rem;
}
```

### 6.2 Typography

```text
Latin UI:       Inter or Geist
Latin Display:  Cormorant Garamond or Playfair Display
Arabic UI:      IBM Plex Sans Arabic or Noto Sans Arabic
Arabic Display: Noto Kufi Arabic or Readex Pro
```

```text
Display XL:  clamp(3rem, 8vw, 7rem)
Display LG:  clamp(2.5rem, 6vw, 5rem)
Heading 1:   clamp(2rem, 4vw, 3.5rem)
Heading 2:   clamp(1.75rem, 3vw, 2.75rem)
Heading 3:   1.5rem
Body large:  1.125rem
Body:        1rem
Small:       0.875rem
Label:       0.75rem
```

Large headings use compact line height; avoid excessive font-weight variation; prices must be immediately visible; uppercase labels sparingly (never for Arabic unless specifically appropriate, §25); avoid justified text; limit paragraphs to readable line widths.

### 6.3 Spacing

```text
Mobile:  py-16
Tablet:  py-20
Desktop: py-28 or py-32
```

```tsx
<div className="mx-auto w-full max-w-[1600px] px-4 sm:px-6 lg:px-10 xl:px-14">
  {children}
</div>
```

Use whitespace to define hierarchy rather than wrapping every section in a bordered card.

### 6.4 Corners, Borders, Shadows

```text
Cards:   rounded-sm or rounded-none
Inputs:  rounded-sm
Buttons: rounded-sm
Badges:  rounded-full only when semantically appropriate
```

Thin neutral borders, very subtle shadows, square or nearly square product imagery, full-width image sections without unnecessary containers. Avoid large rounded cards everywhere, strong drop shadows, glassmorphism, and pill-shaped buttons.

## 7. Interaction and Motion

Motion communicates quality, not decoration.

```text
Fast interaction:    150ms
Standard transition: 250ms
Drawer or modal:     300ms
Editorial reveal:    400ms maximum
```

Use: subtle image zoom on hover, fade/translate transitions for menus and dialogs, smooth cart-drawer transitions, underline animations on nav links, crossfade between gallery images, skeleton loaders matching final content dimensions, sticky purchase panels where appropriate.

Avoid: large page-entrance animations, continuous motion, parallax on every section, animations that delay interaction, bouncing buttons, excessive Framer Motion.

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    scroll-behavior: auto !important;
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

Shared motion utilities live in `core/lib/motion.ts` (timing constants, shared variants) and `core/hooks/use-prefers-reduced-motion.ts` — features consume these rather than hardcoding durations.

## 8. Frontend Routes

Route strategy is locale-prefixed (§26): `/[locale]/...`. `app/[locale]/layout.tsx` is the root layout — Next.js allows the dynamic segment to own the root `<html>`/`<body>` when it is the only top-level route, so no separate non-dynamic `app/layout.tsx` is needed; `[locale]/layout.tsx` imports `app/globals.css` and sets `lang`/`dir` (§26).

Route files stay thin: they may read params, define metadata, and compose feature views. They must not implement repositories, perform scattered raw fetch calls, duplicate backend authorization, or contain trusted calculations.

```text
apps/web/app/
├── globals.css
├── styles/
│   ├── tokens.css
│   └── utilities.css
└── [locale]/
    ├── layout.tsx
    ├── page.tsx                        # home — composes features/home (§15)
    ├── login/
    ├── register/
    ├── init-email/
    ├── confirmOTP/
    ├── init-password/
    ├── products/
    │   ├── page.tsx
    │   └── [slug]/
    ├── categories/[slug]/
    ├── collections/[slug]/
    ├── regions/[slug]/
    ├── artisans/
    │   ├── page.tsx
    │   └── [slug]/
    ├── search/
    ├── wishlist/
    ├── cart/
    ├── checkout/
    │   ├── page.tsx
    │   └── success/
    ├── payment/[paymentId]/
    ├── account/
    │   ├── layout.tsx                  # side nav shell, §19
    │   ├── page.tsx
    │   ├── profile/
    │   ├── addresses/
    │   ├── orders/
    │   │   ├── page.tsx
    │   │   └── [id]/
    │   ├── wishlist/
    │   └── reviews/
    ├── custom-orders/
    │   ├── page.tsx
    │   ├── new/
    │   └── [id]/
    ├── artisan/
    │   ├── layout.tsx                  # dashboard shell, §20
    │   ├── page.tsx
    │   ├── profile/
    │   ├── products/
    │   │   ├── page.tsx
    │   │   ├── new/
    │   │   └── [id]/
    │   ├── orders/
    │   └── custom-orders/
    └── admin/
        ├── layout.tsx                  # dashboard shell, §20
        ├── page.tsx
        ├── users/
        ├── artisans/
        ├── catalogue/
        ├── moderation/
        ├── warehouse/
        ├── inventory/
        ├── orders/
        ├── payments/
        ├── shipments/
        ├── custom-orders/
        └── audit/
```

MVP route list, for reference against the tree above (all under `/[locale]`): `/`, `/products`, `/products/[slug]`, `/artisans`, `/artisans/[slug]`, `/categories/[slug]`, `/login`, `/register`, `/cart`, `/checkout`, `/payment/[paymentId]`, `/account`, `/account/orders`, `/account/orders/[id]`, `/custom-orders`, `/artisan`, `/artisan/profile`, `/artisan/products`, `/artisan/products/new`, `/artisan/products/[id]`, `/admin`, `/admin/artisans`, `/admin/moderation`, `/admin/warehouse`, `/admin/orders`, `/admin/audit`.

## 9. Global Storefront Layout

```text
core/components/layout/
├── AnnouncementBar.tsx
├── StorefrontHeader.tsx
├── StorefrontFooter.tsx
├── MobileNav.tsx
├── SearchOverlay.tsx
├── LocaleSwitcher.tsx
└── Breadcrumbs.tsx
```

**Announcement bar** — optional, one short message maximum, dismissible only when necessary, fully translated, must not consume excessive vertical space.

**Header** — minimal and editorial, sticky, compact height; transparent over the home hero when contrast is sufficient, solid after scrolling.

```text
Desktop layout:
Left:    navigation (New, Products, Categories, Artisans, Our Story)
Centre:  AISHA logo
Right:   search, locale, account, cart
```

RTL mirrors this naturally (§26). Search opens as a full-width overlay (`SearchOverlay.tsx`); cart opens as a drawer (`features/cart/components/containers/CartDrawer.tsx`, §18); mobile nav uses a full-height sheet (`MobileNav.tsx`); cart item count shows without visual noise.

**Footer** — AISHA story, shop links, artisan information, customer support, terms/privacy, social links, newsletter (only when implemented), language selector, delivery information when required. Spacious multi-column layout on desktop, accordions on mobile.

## 10. Frontend Core

```text
apps/web/core/
├── api/
│   ├── client.ts
│   ├── server-client.ts
│   ├── request.ts
│   ├── response.ts
│   ├── errors.ts            # ApiError model, §21
│   ├── correlation.ts
│   ├── token-manager.ts
│   └── endpoints/
├── context/
│   ├── AuthContext.tsx
│   ├── CartContext.tsx
│   └── AppProviders.tsx
├── components/
│   ├── ui/                  # shadcn primitives, §25
│   ├── layout/               # §9
│   ├── commerce/
│   │   ├── Money.tsx
│   │   ├── AvailabilityLabel.tsx
│   │   └── VerificationBadge.tsx
│   ├── editorial/
│   │   ├── EditorialSection.tsx
│   │   └── EditorialTile.tsx
│   ├── feedback/
│   │   ├── EmptyState.tsx
│   │   ├── ErrorState.tsx
│   │   └── LoadingSkeleton.tsx
│   ├── guards/
│   │   ├── ProtectedRoute.tsx
│   │   └── RoleGuard.tsx
│   └── forms/
│       ├── FormField.tsx
│       └── FormErrorSummary.tsx
├── hooks/
│   ├── use-media-query.ts
│   └── use-prefers-reduced-motion.ts
├── lib/
│   ├── motion.ts
│   └── utils.ts
├── providers/
├── zod/
├── translation.ts
└── types.ts
```

`core` contains generic reusable infrastructure only. Feature-specific rules stay in their feature. `core/api/endpoints/` holds one file per backend feature group (mirroring backend §48's routing table) rather than a single growing `endpoints.ts`.

## 11. Frontend Feature Structure

```text
apps/web/features/<feature>/
├── index.ts
├── domain/
│   ├── entities/
│   ├── value-objects/
│   ├── rules/
│   ├── errors/
│   ├── repositories/
│   └── services/
├── application/
│   ├── commands/
│   ├── queries/
│   ├── dto/
│   ├── mappers/
│   └── use-cases/
├── data/
│   ├── api/
│   ├── repositories/
│   ├── mappers/
│   ├── cache/
│   └── storage/
├── viewmodel/
│   ├── hooks/
│   ├── state/
│   ├── actions/
│   └── types/
├── view/
├── components/
│   ├── presentational/
│   ├── containers/
│   ├── forms/
│   └── feedback/
├── schemas/
├── hooks/
├── types/
├── utils/
├── constants/
└── tests/
    ├── domain/
    ├── application/
    ├── data/
    ├── viewmodel/
    └── view/
```

Not every simple feature needs every folder. In particular: `data/api` + `components/` + `schemas/` + `types/` + `utils/` is the MVP-minimum subset — a feature can ship with just those (presentation plus a typed API call) and grow into `domain/`, `application/`, and `viewmodel/` only once it accumulates real business rules or multi-source orchestration. `features/home` (§15) is the clearest example: view and components only, no domain/application/data, because it composes other features' already-fetched data rather than owning any of its own.

## 12. Frontend Layer Responsibilities

### Domain

Contains entities, value objects, framework-independent rules, repository contracts, and domain errors.

The frontend domain must not import React, Next.js, TanStack Query, Axios, `fetch`, browser APIs, or UI components.

### Application

Optional layer for commands, queries, use cases, DTOs, and reusable orchestration shared across several ViewModels.

### Data

Contains API calls, request/response types, repository implementations, API-to-domain mappers, cache adapters, and storage adapters.

### ViewModel

Coordinates loading, empty, error, filter, pagination, form, selection, and action state (§16 defines exactly which states every fetch must handle).

Allowed forms:

```text
Custom React hook
Plain TypeScript class
Reducer with actions
TanStack Query wrapper
State-machine adapter
```

The ViewModel must not render JSX, use database concepts, call random endpoints when a repository exists, or duplicate trusted backend rules (backend §50 is authoritative; this layer is UX only).

### View

Consumes ViewModels and composes the feature screen. Views do not implement repositories, call raw endpoints, or calculate trusted business values.

### Components

```text
components/
├── presentational/
├── containers/
├── forms/
└── feedback/
```

Presentational components only receive props and emit callbacks; they must not call APIs directly (§20 of the checklist).

## 13. Frontend Features

```text
apps/web/features/
├── home/          # NEW — landing-page composition only, see §15
├── auth/
├── account/
├── address/
├── catalogue/
├── category/
├── collection/
├── region/
├── product/
├── artisan/
├── search/
├── cart/
├── checkout/
├── order/
├── payment/
├── shipment/
├── wishlist/
├── review/
├── custom-order/
├── moderation/
├── warehouse/
├── inventory/
└── admin/
```

This is the canonical feature list, matching backend §46's Feature Mapping table one-to-one. `docs/frontend.md` names a shorter, looser-cased subset for its MVP examples (`authentication`, `artisans`, `custom-orders`); those map directly to `auth`, `artisan`, `custom-order` above — use the folder names in this list, not the guide's prose casing.

## 14. Example Frontend Product Feature

```text
features/product/
├── index.ts
├── domain/
│   ├── entities/product.ts
│   ├── value-objects/product-price.ts
│   ├── rules/product-rules.ts
│   ├── errors/product-errors.ts
│   └── repositories/product-repository.ts
├── application/
│   ├── queries/get-products.ts
│   ├── queries/get-product-by-slug.ts
│   ├── commands/
│   ├── dto/
│   └── mappers/
├── data/
│   ├── api/product-api.ts
│   ├── api/product-api.types.ts
│   ├── api/endpoints.ts
│   ├── repositories/api-product-repository.ts
│   └── mappers/product-mapper.ts
├── viewmodel/
│   ├── product-list-viewmodel.ts
│   ├── product-details-viewmodel.ts
│   ├── product-state.ts
│   ├── product-actions.ts
│   └── hooks/
│       ├── use-product-list-viewmodel.ts
│       └── use-product-details-viewmodel.ts
├── view/
│   ├── ProductListingView.tsx
│   ├── ProductDetailsView.tsx
│   └── index.ts
├── components/
│   ├── presentational/
│   │   ├── ProductCard.tsx
│   │   ├── ProductGrid.tsx
│   │   ├── ProductGallery.tsx
│   │   ├── ProductPrice.tsx
│   │   ├── ProductHeading.tsx
│   │   ├── ProductStory.tsx
│   │   ├── ProductSpecifications.tsx
│   │   └── RelatedProducts.tsx
│   ├── containers/
│   │   ├── ProductPurchasePanel.tsx    # price, availability, quantity, add-to-cart, sticky on mobile
│   │   ├── ProductOptions.tsx           # variant selection state
│   │   └── AddToCartButton.tsx
│   ├── forms/
│   └── feedback/
├── schemas/
├── hooks/
├── types/
├── utils/
└── tests/
```

## 15. Home Feature and Page Component Trees

`features/home/` composes the landing page from other features' data and is view/components-only (§11):

```text
features/home/
├── index.ts
├── view/HomeView.tsx
├── components/
│   ├── HomeHero.tsx
│   ├── FeaturedCollection.tsx
│   ├── EditorialCategoryGrid.tsx
│   ├── BrandStory.tsx
│   ├── FeaturedArtisans.tsx
│   ├── FeaturedProductGrid.tsx
│   ├── ProvenanceValues.tsx
│   ├── RegionalCraftFeature.tsx
│   └── FinalDiscoveryBanner.tsx
└── tests/view/
```

### Home page composition

```tsx
<StorefrontLayout>
  <AnnouncementBar />
  <StorefrontHeader />

  <main>
    <HomeHero />
    <FeaturedCollection />
    <EditorialCategoryGrid />
    <BrandStory />
    <FeaturedArtisans />
    <FeaturedProductGrid />
    <ProvenanceValues />
    <RegionalCraftFeature />
    <FinalDiscoveryBanner />
  </main>

  <StorefrontFooter />
  <CartDrawer />
  <SearchOverlay />
</StorefrontLayout>
```

Hero: full/near-full viewport image, strong headline, one supporting sentence, one or two actions, art-directed responsive images with mobile crops, no carousel in MVP. Featured Collection: large editorial tiles (desktop: one two-thirds tile + two stacked one-third tiles; mobile: vertical full-width). Categories: large image-based blocks (jewellery, textiles, ceramics, leather, home décor, traditional garments) — never identical small icons. Featured Products: 4 columns large / 3 medium / 2 small tablet / 1–2 mobile, alternated with larger editorial blocks so the page doesn't read as a generic marketplace. Values section: verified origin, fair partnership, quality inspection, responsible materials — concise copy, simple line icons; never show a verification badge unless the backend has actually verified it.

### Product page composition

```tsx
<ProductPage>
  <Breadcrumbs />

  <section className="product-layout">
    <ProductGallery />

    <ProductPurchasePanel>
      <ProductHeading />
      <ArtisanLink />
      <Money />
      <Availability />
      <ProductOptions />
      <AddToCartButton />
      <DeliverySummary />
    </ProductPurchasePanel>
  </section>

  <ProductStory />
  <ProductSpecifications />
  <ArtisanFeature />
  <RelatedProducts />
</ProductPage>
```

Desktop: gallery ~60–65% width, sticky purchase panel ~35–40%. Mobile order: gallery → summary → purchase controls → story/details → artisan info → recommendations, with an optional sticky bottom bar showing price + add-to-cart that never hides content or interferes with browser controls. `ArtisanLink`, `Availability`, and `DeliverySummary` live in `features/product/components/presentational/`; `ArtisanFeature` lives in `features/artisan/components/presentational/ArtisanFeature.tsx`. Recommendations (more from this artisan, similar products, related collection) render nothing when no valid data exists — never an empty section.

## 16. Data Fetching Rules

**Server Components** for: catalogue pages, product detail, artisan public pages, category pages, initial account pages where the auth design supports server access, SEO metadata, initial navigation data.

**Client Components** for: cart interactions, filters with client state, forms, upload progress, product gallery interaction, search overlays, cart drawer, rich interactive moderation/warehouse controls.

Every fetch must handle: loading, empty result, recoverable error, authentication required, forbidden, not found (§24 defines the shared components for these). Use page-specific skeletons (product grid skeleton, product gallery skeleton, artisan profile skeleton, order detail skeleton) rather than a generic spinner — never a full-page spinner for catalogue navigation.

## 17. Forms

React Hook Form + Zod throughout. Zod schemas live in each feature's `schemas/`; reuse them for client-side feedback, but the backend remains authoritative (§50) — client validation is UX only.

| Form | Feature | File |
|---|---|---|
| Registration | `auth` | `features/auth/components/forms/RegisterForm.tsx` |
| Login | `auth` | `features/auth/components/forms/LoginForm.tsx` |
| Address | `address` | `features/address/components/forms/AddressForm.tsx` |
| Artisan application | `artisan` | `features/artisan/components/forms/ArtisanApplicationForm.tsx` |
| Artisan profile | `artisan` | `features/artisan/components/forms/ArtisanProfileForm.tsx` |
| Product creation / edit | `product` | `features/product/components/forms/ProductForm.tsx` |
| Product submission | `product` | `features/product/components/forms/ProductSubmitForm.tsx` |
| Moderation decision | `moderation` | `features/moderation/components/forms/ModerationDecisionForm.tsx` |
| Warehouse reception | `warehouse` | `features/warehouse/components/forms/ReceptionForm.tsx` |
| Inspection | `warehouse` | `features/warehouse/components/forms/InspectionForm.tsx` |
| Inventory adjustment | `inventory` | `features/inventory/components/forms/InventoryAdjustmentForm.tsx` |
| Checkout (per step) | `checkout` | `features/checkout/components/forms/` |
| Custom-order request | `custom-order` | `features/custom-order/components/forms/CustomOrderRequestForm.tsx` |
| Quote response | `custom-order` | `features/custom-order/components/forms/QuoteResponseForm.tsx` |

Rules: display stable API error messages mapped to translated user text (§21); move focus to the error summary after a failed submit (`core/components/forms/FormErrorSummary.tsx`); preserve user input after recoverable errors; disable duplicate submissions while pending; show field-level and form-level errors; never use placeholder text as a label; label fields above the input with comfortable vertical spacing, thin borders, and strong focus states.

## 18. Cart, Checkout, and Payment Return

**Cart** — local for anonymous visitors, server-backed after login, with an explicit merge policy. Current price shown as provisional; revalidate on opening the cart and on checkout; clearly mark unavailable items; never claim cart quantity is reserved.

```text
features/cart/
├── components/containers/CartDrawer.tsx   # focus trap, keyboard accessible, restores focus on close, RTL-correct
└── ...(standard feature structure, §11)
```

Drawer contents: product image, name, artisan, quantity, price, remove action, provisional subtotal, view-cart, checkout. Full cart page: desktop two-column (items / sticky order summary), mobile stacked (items → summary → checkout action).

**Checkout** — `features/checkout/view/CheckoutView.tsx` drives five steps with a progress indicator, visually quieter than the main storefront:

```text
features/checkout/components/containers/
├── ReviewItemsStep.tsx
├── DeliveryAddressStep.tsx
├── ShippingMethodStep.tsx
├── PaymentStep.tsx
└── ConfirmationStep.tsx
```

Avoid promotional distractions, large nav menus, product recommendations during payment, hidden fees, and pre-selected optional services. At confirmation: display trusted server totals, currency, delivery address, shipping method; require explicit final confirmation; send an `Idempotency-Key` (backend §52); handle price/stock changes gracefully. Desktop: two-column (current step / persistent summary). Mobile: collapsible summary near the top.

**Payment return** — `features/payment/view/PaymentReturnView.tsx`, one component per state:

```text
features/payment/components/presentational/
├── PaymentPending.tsx    # calm loading, no duplicate-payment encouragement, safe navigation
├── PaymentPaid.tsx       # confirmation heading, order reference, summary, link to order, continue shopping
├── PaymentFailed.tsx     # safe explanation, retry only via a backend-authorised flow, support guidance
└── PaymentExpired.tsx    # explains session expiry, safe return to cart/checkout
```

The page may poll the backend for a limited period; it must never treat a query string like `success=true` as proof of payment — only a backend-confirmed status counts.

## 19. Authentication and Account UI

**Auth** — login, registration, protected route layout (`core/components/guards/ProtectedRoute.tsx`), forbidden page, session-expired handling, logout, role-aware nav (convenience only — backend §49 is authoritative). Style: split-screen editorial image on desktop, focused single-column form on mobile, minimal distractions, clear password requirements, visible error feedback. Never present auth pages as generic admin screens.

**Account area** (`app/[locale]/account/layout.tsx`) — side navigation on desktop, compact tab selector on mobile:

```text
Overview · Orders · Addresses · Profile · Security · Custom orders · Logout
```

Order cards show: reference, date, status, total, product preview, detail action.

## 20. Artisan Workspace and Admin/Warehouse UI

Both use a distinct application shell, separate from the public storefront — do not imitate the editorial home page here; prioritise clarity and task completion.

```text
app/[locale]/artisan/layout.tsx   # sidebar + top header + page title/primary action + status summary + workspace
app/[locale]/admin/layout.tsx     # same shell pattern
```

**Artisan dashboard** (MVP): profile completion, application status, product counts by status, create/edit-draft/submit-for-review, moderation reason display, accepted/reserved/shipped stock for owned products, custom-order requests.

**Moderation queue** (`features/moderation/components/containers/`): `ModerationQueueTable.tsx` (product, artisan, submission date, translation completeness, media, status, actions — desktop: filter toolbar + dense table + preview drawer + action panel; mobile: stacked cards + full-screen review sheet) and `ModerationActionPanel.tsx` (actions require confirmation and reason).

**Reception** (`features/warehouse/components/forms/ReceptionForm.tsx`): product, artisan, quantity, batch/parcel reference, notes, date; product search, barcode/reference entry where available, confirmation summary, duplicate-submission prevention.

**Inspection** (`features/warehouse/components/forms/InspectionForm.tsx` + `InspectionTotalsCheck.tsx`): received quantity must equal the sum of outcomes, shown visibly —

```text
Accepted + Rejected + Damaged + Pending = Received quantity
```

— submission is blocked when totals don't match.

**Order fulfilment** (`features/order/components/containers/` under `app/[locale]/admin/orders/`): only paid, eligible orders; reference, payment status, fulfilment status, delivery method, required items, picking state, packing state, shipment reference.

Dashboard rules throughout: cards only for meaningful grouped information, data tables on desktop / stacked rows or cards on mobile, keep AISHA's tokens and typography (§6), status never communicated by colour alone.

## 21. API Client

```text
core/api/
├── client.ts
├── server-client.ts
├── request.ts
├── response.ts
├── errors.ts
├── correlation.ts
├── token-manager.ts
└── endpoints/
```

```ts
export interface ApiError {
  code: string;
  message: string;
  correlationId?: string;
  fieldErrors?: Record<string, string[]>;
}
```

UI components map stable backend error codes (backend §51) to translation keys — never display raw backend stack traces or internal error messages. The client never silently retries checkout or payment creation without reusing the same idempotency key (backend §52); retries are safe only for idempotent reads unless a flow is specifically designed otherwise.

## 22. Accessibility, Media, and SEO

**Accessibility** — native semantic controls; every input labelled and every error associated with its input; full keyboard navigation with visible focus (including over image backgrounds); dialogs trap and restore focus; meaningful alt text; status never colour-only; sufficient touch targets; skip links; valid heading hierarchy; keyboard-accessible galleries; carousels that don't auto-advance unexpectedly; drawers/overlays announce their purpose; RTL never breaks keyboard or reading order.

**Media** — optimised image components, correct aspect ratios, loading placeholders, basic zoom/gallery for MVP, signed URLs with expiry (never permanent public URLs for private draft media), responsive sizes, no full-resolution images in grids, no layout shift.

```text
Product card:      4:5
Product gallery:    4:5 or original constrained ratio
Artisan portrait:   3:4
Editorial banner:   16:9 or 3:2
Category tile:      4:5 or 1:1
Mobile hero:        4:5
Desktop hero:        16:9
```

**SEO** — localised title/description, product and artisan metadata, canonical URLs, Open Graph images matching AISHA's visual identity, sitemap for active public content only, accurate structured product data, localised alternate URLs, meaningful slugs, no indexing of draft/moderation/cart/checkout/account pages, server-rendered critical product info.

## 23. Responsive Behaviour

```text
Mobile:          320px and above
Small tablet:    640px and above
Tablet:          768px and above
Desktop:         1024px and above
Large desktop:   1280px and above
Editorial wide:  1536px and above
```

Header collapses to compact mobile nav; filters move into a sheet; product-detail columns stack; sticky purchase actions remain usable; tables become mobile-friendly rows; display typography scales via `clamp` (§6.2); checkout summary becomes collapsible; product-card text never overflows; Arabic content is tested at every breakpoint (§26).

## 24. Empty, Loading, and Error States

```text
core/components/feedback/
├── EmptyState.tsx     # heading, short explanation, one action, optional restrained illustration
├── ErrorState.tsx     # human-readable message, retry when safe, nav alternative, correlation ID, distinct forbidden/not-found handling
└── LoadingSkeleton.tsx # layout-matched — never a large centred spinner for a whole page
```

Examples requiring an empty state: empty cart, no search results, no products in a category, no artisan products, no customer orders, no moderation submissions, no warehouse receptions. Page-specific skeleton variants (product grid, gallery, artisan profile, order detail) live alongside each feature's own `components/feedback/`.

## 25. Component Styling Standards

```text
core/components/ui/   (shadcn primitives)
├── button.tsx    # Primary, Secondary, Outline, Ghost, Destructive, Text link
├── card.tsx
├── badge.tsx
├── dialog.tsx    # focused confirmation, irreversible actions (alert dialog variant)
├── sheet.tsx     # filters, mobile menus, cart, detail panels
└── skeleton.tsx
```

```tsx
<Button
  size="lg"
  className="h-12 w-full rounded-sm px-8 text-sm font-medium uppercase tracking-wide"
>
  {t("product.addToCart")}
</Button>
```

Primary storefront buttons: strong contrast, medium/large height, minimal radius, clear hover/focus/active/disabled/pending states; never uppercase for Arabic text unless specifically appropriate. Cards only when content needs visual grouping — product cards avoid visible backgrounds, heavy borders, large shadows, excessive padding. Badges only for verified provenance, fair-trade status, eco-friendly certification, new products, made-to-order, and stock state — never for ordinary descriptive information.

## 26. Internationalization and RTL

Supported locales:

```text
ar
fr
en
es
```

```text
apps/web/messages/
├── ar.json
├── fr.json
├── en.json
└── es.json
```

```tsx
const direction = locale === "ar" ? "rtl" : "ltr";

return (
  <html lang={locale} dir={direction}>
    <body>{children}</body>
  </html>
);
```

All visible text uses translation keys. Arabic uses `dir="rtl"`; icons implying direction mirror where appropriate; form alignment and table layouts work in RTL; dates, numbers, and currencies use locale-aware formatting; never concatenate translated fragments; use logical CSS properties (`start`/`end`, not hardcoded `left`/`right`); test every major page in Arabic, not only navigation. Product grids keep natural visual order in RTL; carousels support RTL navigation; breadcrumb arrows mirror; drawer/sheet opening direction feels natural for the active language.

## 27. Frontend Testing

```text
apps/web/tests/
├── unit/          # ProductCard, ProductGallery, Money, LanguageSwitch, RTL layout, forms/validation,
│                  # OrderStatus, AvailabilityLabel, CartItem, VerificationBadge, filter controls,
│                  # mobile nav, PaymentStatus — mirrors each feature's own tests/{domain,application,viewmodel}
├── integration/   # API error mapping, protected layouts, cart state, anonymous-to-authenticated cart merge,
│                  # checkout steps, locale routing, server-rendered product data, signed media URL handling,
│                  # role-aware navigation, moderation action validation, warehouse quantity validation
├── e2e/           # Playwright: browse products, change locale, verify Arabic RTL, search/filter, open a
│                  # product, add to cart, register/login, complete checkout with the dev payment adapter,
│                  # handle payment pending→paid, view an order, artisan creates/submits a product, moderator
│                  # approves a product, warehouse receives/inspects stock, mobile nav, keyboard nav
└── visual/        # Home hero, product listing, product detail, cart drawer, checkout, Arabic RTL layouts,
                   # mobile navigation, artisan profile, moderation queue
```

Each feature's own `tests/{domain,application,data,viewmodel,view}/` (§11) covers that feature in isolation; `apps/web/tests/` above is the cross-cutting suite referenced by backend §56.

## 28. Performance Requirements

```text
LCP: below 2.5 seconds
CLS: below 0.1
INP: below 200 milliseconds
```

Optimise hero images; preload only essential fonts/media; avoid unnecessary client-side JS; lazy-load below-the-fold content; Server Components for product content; avoid heavy animation libraries unless justified; limit third-party scripts; dynamic imports for rich admin controls; prevent cumulative layout shift; appropriate cache/revalidation strategy; lightweight product-card components; avoid loading full product objects when summary data is sufficient.

## 29. Frontend Implementation Phases

```text
1. Design Foundation    — colours, typography, spacing scale, CSS variables, button variants,
                           form controls, header, footer, container, responsive rules, RTL foundation
2. Public Storefront     — home, product listing, product detail, category page, artisan listing,
                           artisan profile, search, cart drawer
3. Authentication & Checkout — login, registration, cart page, address, checkout steps, payment
                           return, confirmation
4. Customer Account       — overview, orders, order detail, addresses, custom orders
5. Artisan Workspace       — dashboard, profile, product management, media upload, submission/
                           moderation feedback, stock overview
6. Admin and Warehouse     — moderation queue, reception, inspection, inventory, paid-order
                           fulfilment, audit views
7. Quality                — accessibility audit, RTL audit, performance audit, visual regression,
                           Playwright coverage, cross-browser testing, responsive review
```

## 30. Per-Page Definition of Done

A page is not complete until it has: desktop, tablet, mobile, and Arabic RTL layouts; loading, empty (where applicable), and error states; keyboard accessibility with visible focus; translated text with correct locale formatting; optimised images; metadata where public; component tests where appropriate; no raw backend error messages; no broken layout under long translated content; and no business-critical rule implemented only in the frontend.

## 31. Required Frontend Commands

```bash
npm ci
npm run lint
npm run type-check
npm test
npm run build
npm run test:e2e
```

Recommended additional commands:

```bash
npm run format:check
npm run test:visual
npm run test:a11y
npm run analyse
```

These map to `apps/web/package.json` scripts and are the frontend half of the CI command set — see backend §58 for the Go equivalent.

# Part II — Backend

Part II is written to be directly buildable: every layer below names the exact
package it is implemented with, the exact code shape handlers/use
cases/repositories take, and the exact commands CI runs. It supersedes any
more abstract "framework-independent" phrasing from earlier drafts wherever
the two disagree — concrete wins.

## 32. Backend Goal and Required Packages

Build a Go Fiber API that implements the business workflow without placing
business logic in handlers or SQLBoiler models.

```text
github.com/gofiber/fiber/v3            → server/, feature server/ handlers
github.com/casbin/casbin/v2            → internal/pkg/authorization
github.com/volatiletech/sqlboiler/v4   → features/<feature>/data/models, data/repositories
github.com/google/wire                 → internal/container (compile-time DI)
github.com/jackc/pgx/v5                → internal/pkg/database (Postgres driver under SQLBoiler)
github.com/redis/go-redis/v9           → rate limiting, OTP state, short-lived cache
github.com/minio/minio-go/v7           → internal/pkg/storage
github.com/golang-migrate/migrate/v4   → db/migrations tooling (Makefile target, never in-process)
github.com/google/uuid                 → domain IDs
github.com/go-playground/validator/v10 → server/request struct tags
```

Choose exact versions compatible with the repository and current Go version.
Do not mix multiple HTTP routers or ORMs, and do not substitute an equivalent
(no Gin, no GORM, no manual JWT parsing) — these choices are fixed.

Fiber v3 uses value-receiver contexts (`fiber.Ctx`, not `*fiber.Ctx`) —
handler signatures across every feature follow this consistently (§40).

### Coding rules

- Run `gofmt`.
- Handle every error; wrap errors with context.
- Use `context.Context` in application and repository methods.
- Keep interfaces small and consumer-owned.
- Use constructor injection; avoid global mutable state.
- Use UTC timestamps.
- Use structured logs with correlation IDs.
- Use integer minor units for money.
- Never expose SQLBoiler models as API responses.

## 33. Backend Base Structure

```text
apps/api/
├── cmd/
│   └── main.go
├── internal/
│   ├── bootstrap/
│   ├── config/
│   ├── container/
│   ├── server/
│   ├── pkg/
│   └── features/
├── db/
│   ├── schema/
│   ├── migrations/
│   ├── seeds/
│   ├── queries/
│   └── sqlboiler.toml
├── tests/
│   ├── unit/
│   ├── integration/
│   ├── concurrency/
│   └── e2e/
├── docs/
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
└── .env.example
```

`cmd/main.go` contains startup only — a single call into `bootstrap`:

```go
func main() {
    app := bootstrap.NewApp()
    defer app.Close()
    app.Server.Run()
}
```

### `internal/bootstrap/` (owns startup sequencing)

```text
apps/api/internal/bootstrap/
├── app.go        # Assembles config + container + server, returns a runnable App
├── database.go   # Opens the Postgres pool, runs health check, wires it into container
└── env.go        # Loads .env / process env into internal/config structs, fails fast on missing keys
```

Distinct from `internal/container/` (DI wiring of feature modules) and
`internal/config/` (typed config values). Rule: `bootstrap` may import
`config`, `container`, and `server`. Nothing in `container`, `server`, `pkg`,
or `features` may import `bootstrap` — the dependency arrow points one way.
`bootstrap` is also the only package allowed to call `os.Exit` or panic on
startup misconfiguration.

## 34. Backend Shared Structure

```text
apps/api/internal/
├── config/
│   ├── config.go
│   ├── database.go
│   ├── redis.go
│   ├── minio.go
│   ├── jwt.go
│   ├── mail.go
│   ├── payment.go
│   ├── shipping.go
│   └── validation.go
├── container/
│   ├── wire.go            # //go:build wireinject — provider sets, e.g. ProductProviderSet
│   ├── wire_gen.go         # generated by `wire`, committed to VCS, never hand-edited
│   ├── infrastructure.go
│   ├── features.go
│   └── lifecycle.go
├── server/
│   ├── server.go
│   ├── router.go
│   ├── middleware/
│   ├── routes/
│   └── error_handler/
├── pkg/
│   ├── apperror/
│   ├── auth/
│   ├── authorization/
│   ├── clock/
│   ├── database/
│   ├── email/
│   ├── events/
│   ├── idempotency/
│   ├── money/
│   ├── pagination/
│   ├── response/
│   ├── storage/
│   ├── transaction/
│   ├── validator/
│   └── utils/
└── features/
```

## 35. Backend Feature Structure

```text
apps/api/internal/features/<feature>/
├── module.go
├── domain/
│   ├── entity.go
│   ├── aggregate.go
│   ├── value_objects.go
│   ├── repository.go       # interfaces declared here — see §36, §40
│   ├── services.go
│   ├── policy.go
│   ├── rules.go
│   ├── events.go
│   └── errors.go
├── application/
│   ├── commands/
│   │   ├── command.go
│   │   └── handler.go
│   ├── queries/
│   │   ├── query.go
│   │   └── handler.go
│   ├── usecases/
│   ├── ports/                # external-provider interfaces, e.g. PaymentGateway
│   ├── dto.go
│   ├── mapper.go
│   └── service.go
├── data/
│   ├── repositories/
│   │   └── sql_repository.go
│   ├── mappers/
│   │   └── persistence_mapper.go
│   ├── models/
│   ├── queries/
│   ├── cache/
│   └── adapters/
├── server/
│   ├── router.go
│   ├── routes.go
│   ├── handler.go
│   ├── middleware/
│   ├── request/
│   ├── response/
│   └── presenter/
└── mocks/                    # generated — see §57, never hand-edited
    ├── mock_<feature>_repository.go
    └── mock_<port>.go
```

## 36. Backend Layer Responsibilities

### Domain

Owns entities, aggregates, value objects, invariants, state transitions,
repository contracts, policies, services, events, and errors.

Every repository interface a feature's data layer implements is declared
directly in that feature's `domain/repository.go` — not inferred from naming
convention:

```go
// features/product/domain/repository.go
package domain

type ProductRepository interface {
    GetByID(ctx context.Context, id string) (*Product, error)
    GetBySlug(ctx context.Context, slug string) (*Product, error)
    Create(ctx context.Context, product *Product) error
    Update(ctx context.Context, product *Product) error
    List(ctx context.Context, f ProductFilter) ([]*Product, int, error)
}
```

Method naming is fixed: `GetByID` / `Create` / `Update` / `List`, consistently
across every feature — not `FindByID`/`Save` in one feature and `GetByID`/
`Create` in another.

Domain must not import Fiber, SQLBoiler, PostgreSQL drivers, Redis, MinIO,
Casbin implementation, provider SDKs, or HTTP structs.

### Application

Owns commands, queries, handlers, use cases, DTOs, transaction boundaries,
authorization orchestration, ownership checks, idempotency, audit, and
outbox coordination. Execution order inside a use case: authorise → validate
command and domain invariants → load required data → open transaction where
needed → call repositories → write audit/outbox event → commit → return DTO.

### Data

Owns SQLBoiler repositories, persistence mappers, optimized queries, Redis
adapters, MinIO adapters, cache, and external provider adapters. SQLBoiler
models remain inside data and are never exposed as API responses.

```go
// features/product/data/repositories/sql_product_repository.go
package repositories

type sqlProductRepository struct{ db *sql.DB }

func NewSQLProductRepository(db *sql.DB) domain.ProductRepository {
    return &sqlProductRepository{db: db}
}
```

`container/features.go` (generated via `wire`, §34) wires the concrete type
behind the interface at startup — nothing above `data/` ever imports
`sqlProductRepository` directly, only `domain.ProductRepository`. Transaction-
aware ports may receive a `Tx` abstraction, or a repository constructed from
a transaction scope — `*sql.Tx`/pgx equivalents must never leak above `data/`.

### Server

Owns router, route registration, handlers, feature middleware, request
structs, response structs, validation, mapping, HTTP status handling, and
presenters. See §40 for the concrete handler shape.

Handlers do not: open SQL queries, decide product state transitions,
calculate trusted prices, implement Casbin policy manually, or call MinIO
directly.

## 37. Global Server Structure

```text
apps/api/internal/server/
├── server.go
├── router.go
├── middleware/
│   ├── request_id.go
│   ├── recovery.go
│   ├── logging.go
│   ├── cors.go
│   ├── authentication.go
│   ├── rate_limit.go
│   └── security_headers.go
├── routes/
│   ├── health.go
│   ├── public.go
│   ├── authenticated.go
│   └── admin.go
└── error_handler/
    └── error_handler.go
```

Global middleware belongs here. JWT verification happens exactly once, in
`middleware/authentication.go` — feature middleware only adds
authorization/ownership checks on top of an already-verified principal.

## 38. Feature Server Structure

```text
features/<feature>/server/
├── router.go
├── routes.go
├── handler.go
├── middleware/
│   ├── authentication.go
│   ├── authorization.go
│   └── validation.go
├── request/
│   ├── create_request.go
│   ├── update_request.go
│   ├── filter_request.go
│   └── mapper.go
├── response/
│   ├── response.go
│   ├── list_response.go
│   ├── detail_response.go
│   └── mapper.go
└── presenter/
    └── presenter.go
```

Feature middleware contains only feature-specific concerns. Do not duplicate
global authentication middleware.

## 39. Router, Request, and Response Rules

### Router

Groups feature routes and delegates to handlers.

### Request

Owns JSON names, validation tags, path/query/header parsing, and mapping to
application commands or queries. It contains no business logic.

### Response

Defines the public JSON contract and maps from application DTOs. It never
exposes SQLBoiler models, secrets, or provider internals.

### Handler

Extracts context and principal, parses input, validates transport data, maps
to a command/query, calls the application layer, maps the result, and
returns the correct status.

## 40. Handler / Use Case / Repository Pattern (concrete shapes)

### Handler (`features/<feature>/server/handler.go`)

```go
func (h *ProductHandler) Create(c fiber.Ctx) error {
    principal, err := h.auth.RequirePrincipal(c)
    if err != nil {
        return h.errors.Write(c, err)
    }

    var req CreateProductRequest
    if err := c.Bind().Body(&req); err != nil {
        return h.errors.Write(c, ErrInvalidJSON)
    }
    if err := h.validator.Struct(req); err != nil {
        return h.errors.Write(c, NewValidationError(err))
    }

    result, err := h.createProduct.Execute(c.Context(), principal, req.ToCommand())
    if err != nil {
        return h.errors.Write(c, err)
    }
    return c.Status(fiber.StatusCreated).JSON(ProductResponseFrom(result))
}
```

### Use case (`features/<feature>/application/usecases/create_product.go`)

```go
type CreateProductCommand struct {
    ArtisanID    string
    CategoryID   string
    PriceMinor   int64
    Currency     string
    Translations []ProductTranslationInput
}

type CreateProduct interface {
    Execute(ctx context.Context, actor Principal, cmd CreateProductCommand) (ProductDTO, error)
}
```

### Repository

Method names and shape are fixed by §36 (`GetByID` / `Create` / `Update` /
`List`). See the `sqlProductRepository` example in §36.

### Request flow — public endpoint (no auth)

```text
Client
  → router.go (feature)                     [route registration only]
  → middleware/ (global: request_id, recovery, logging, cors, rate_limit)
  → handler.go                              [parse + validate request struct]
  → application/queries/handler.go          [orchestration, no I/O logic itself]
  → data/repositories/sql_repository.go     [SQLBoiler query, mapped via persistence_mapper]
  → domain/entity.go                        [returned entity, business invariants already enforced]
  ← application/dto.go + mapper.go          [entity → DTO]
  ← server/response/mapper.go               [DTO → public JSON contract]
  ← Client
```

### Request flow — authenticated endpoint (JWT access token)

```text
Client (Authorization: Bearer <accessToken>)
  → router.go
  → middleware/authentication.go (global)   [verify signature + expiry, extract principal]
  → middleware/authorization.go (feature)   [role/policy check via internal/pkg/authorization]
  → handler.go                              [principal now in context]
  → application/commands/handler.go         [ownership check, transaction boundary, idempotency]
  → domain/policy.go                        [authoritative business rule — frontend guard is UX only]
  → data/repositories/sql_repository.go
  ← ... same response path as above
```

Refresh flow: access tokens are short-lived, refresh tokens are stored as
session records (Postgres) and rotated on use; `internal/pkg/auth` owns
signing/verification, `features/auth` owns the issue/refresh/revoke use
cases. Redis is used only for OTP state and rate limiting, never as the
token source of truth.

## 41. Backend Features

```text
apps/api/internal/features/
├── auth/
├── user/
├── address/
├── artisan/
├── category/
├── collection/
├── region/
├── product/
├── product_media/
├── moderation/
├── warehouse/
├── inventory/
├── reservation/
├── cart/
├── checkout/
├── order/
├── payment/
├── shipment/
├── custom_order/
├── wishlist/
├── review/
├── notification/
├── audit/
└── admin/
```

## 42. Example Backend Product Feature

```text
features/product/
├── module.go
├── domain/
│   ├── product.go
│   ├── product_id.go
│   ├── product_status.go
│   ├── product_repository.go   # GetByID / GetBySlug / Create / Update / List
│   ├── product_policy.go
│   ├── product_rules.go        # submission / activation invariants, see §43
│   ├── product_events.go
│   └── product_errors.go
├── application/
│   ├── commands/
│   │   ├── create_product.go
│   │   ├── update_product.go
│   │   ├── submit_product.go
│   │   └── suspend_product.go
│   ├── queries/
│   │   ├── get_product.go
│   │   └── list_products.go
│   ├── ports/
│   ├── product_dto.go
│   └── product_mapper.go
├── data/
│   ├── repositories/sql_product_repository.go
│   ├── mappers/product_persistence_mapper.go
│   ├── queries/product_queries.go
│   └── models/
├── server/
│   ├── router.go
│   ├── handler.go
│   ├── middleware/product_access.go
│   ├── request/
│   │   ├── create_product_request.go
│   │   ├── update_product_request.go
│   │   ├── product_filter_request.go
│   │   └── request_mapper.go
│   ├── response/
│   │   ├── product_response.go
│   │   ├── product_list_response.go
│   │   └── response_mapper.go
│   └── presenter/product_presenter.go
└── mocks/
    ├── mock_product_repository.go
    └── mock_payment_gateway.go   # if product feature consumes a port
```

## 43. Domain Rules by Feature

Each feature's `domain/rules.go` and `domain/policy.go` carry these
invariants. This is the pattern every feature follows, illustrated on the
features currently in MVP scope.

**`features/product/domain/`** — submission is allowed only when: artisan is
approved and active, required translation fields exist, category exists,
price is positive, currency is supported by config, required media exists,
and current status permits submission. Activation is allowed only when:
moderation is approved, product is not suspended, and accepted inventory is
available.

**`features/inventory/domain/`** — append-only movement log, never a
mutable balance edited in place:

```text
RECEPTION · INSPECTION_ACCEPT · INSPECTION_REJECT · INSPECTION_QUARANTINE
INSPECTION_DAMAGE · RESERVATION · RESERVATION_RELEASE · ORDER_COMMIT
SHIPMENT · ADJUSTMENT_IN · ADJUSTMENT_OUT · RETURN
```

An `inventory_balances` table may exist in `data/` purely as a read
optimisation; every value in it must be reconstructable from the movement
log, and the movement log — not the balance table — is the source of truth.

**`features/reservation/domain/`** — required fields: ID, product/stock item
reference, order or checkout reference, quantity, status, `expires_at`,
`created_at`, `consumed_at`/`released_at`. A background job expires stale
active reservations; this expiry job lives in `internal/pkg` or a
feature-owned worker, not inside a request handler.

**`features/order/domain/`** — order lines are immutable snapshots taken at
checkout time. Never join current product price to render a historical
order total — that is a correctness bug, not a style preference.

**`features/payment/domain/`** — persisted fields: internal payment ID,
order ID, provider, provider payment reference, amount minor, currency,
status, idempotency key, created/updated timestamps. Provider webhook events
are stored separately from the payment record itself, keyed by a unique
external event ID, so a provider retry cannot double-apply.

## 44. Database Migrations

Migration rules:

- Every schema change has `up` and `down`.
- Prefer additive changes.
- Add explicit foreign keys.
- Add unique constraints.
- Add indexes for frequent lookups.
- Define deletion behaviour explicitly.
- Production migrations never run invisibly inside API startup — `bootstrap/`
  (§33) opens connections and health-checks, it does not migrate.
- Regenerate SQLBoiler models after schema changes — a stale generated model
  is treated as a build break, same policy as a stale mock (§57).

Migration sequence:

```text
000001_extensions
000002_users_sessions
000003_artisans
000004_catalogue
000005_moderation
000006_warehouse
000007_inventory
000008_cart_addresses
000009_orders
000010_payments
000011_shipments
000012_custom_orders
000013_notifications_outbox
000014_audit_idempotency
```

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

## 45. Cross-Feature Communication

A feature must never import another feature's concrete data implementation.

Allowed mechanisms:

```text
Application port
Exported application query
Domain service interface
Stable shared contract
Domain event
Integration event
Outbox event
```

## 46. Feature Mapping

| Capability | Frontend | Backend |
|---|---|---|
| Authentication | `auth` | `auth` |
| Profile | `account` | `user` |
| Addresses | `address` | `address` |
| Artisan | `artisan` | `artisan` |
| Catalogue | `catalogue` | `category`, `collection`, `region` |
| Products | `product` | `product`, `product_media` |
| Search | `search` | catalogue queries |
| Moderation | `moderation` | `moderation` |
| Warehouse | `warehouse` | `warehouse` |
| Inventory | `inventory` | `inventory`, `reservation` |
| Cart | `cart` | `cart` |
| Checkout | `checkout` | `checkout`, `reservation` |
| Orders | `order` | `order` |
| Payments | `payment` | `payment` |
| Shipments | `shipment` | `shipment` |
| Wishlist | `wishlist` | `wishlist` |
| Reviews | `review` | `review` |
| Custom orders | `custom-order` | `custom_order` |
| Administration | `admin` | `admin`, `audit` |
| Landing composition | `home` | — (composes existing public endpoints only) |

## 47. API Contract

Every endpoint defines method, path, authentication, permission, request
schema, response schema, stable error codes, pagination, idempotency, and
examples.

```json
{
  "data": {},
  "meta": {},
  "requestId": "uuid"
}
```

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "The request is invalid.",
    "fields": {
      "email": "INVALID_EMAIL"
    }
  },
  "requestId": "uuid"
}
```

## 48. API Endpoints and Routing

Every endpoint group below registers in its feature's
`features/<feature>/server/routes.go`, then is mounted from the global
`internal/server/routes/{public,authenticated,admin}.go` (§37).

### Authentication

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/logout
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/me
```

### Public Catalogue

```text
GET    /api/v1/categories
GET    /api/v1/products
GET    /api/v1/products/:id
GET    /api/v1/artisans
GET    /api/v1/artisans/:id
```

### Artisan

```text
POST   /api/v1/artisan-applications
POST   /api/v1/artisan-applications/draft
POST   /api/v1/artisan-applications/me/submit
GET    /api/v1/artisan-applications/me
PATCH  /api/v1/artisan/profile
POST   /api/v1/artisan/products
PATCH  /api/v1/artisan/products/:id
POST   /api/v1/artisan/products/:id/media
POST   /api/v1/artisan/products/:id/submit
GET    /api/v1/artisan/products
GET    /api/v1/artisan/inventory
```

### Moderation

```text
GET    /api/v1/admin/product-submissions
POST   /api/v1/admin/product-submissions/:id/approve
POST   /api/v1/admin/product-submissions/:id/request-changes
POST   /api/v1/admin/product-submissions/:id/reject
POST   /api/v1/admin/products/:id/suspend
```

### Warehouse and Inventory

```text
POST   /api/v1/warehouse/receptions
GET    /api/v1/warehouse/receptions
POST   /api/v1/warehouse/receptions/:id/inspect
GET    /api/v1/warehouse/inventory
POST   /api/v1/warehouse/inventory/:productId/adjust
GET    /api/v1/warehouse/orders
POST   /api/v1/warehouse/orders/:id/prepare
POST   /api/v1/warehouse/orders/:id/ship
```

### Cart and Checkout

```text
GET    /api/v1/cart
POST   /api/v1/cart/items
PATCH  /api/v1/cart/items/:id
DELETE /api/v1/cart/items/:id
POST   /api/v1/checkout
GET    /api/v1/orders
GET    /api/v1/orders/:id
POST   /api/v1/orders/:id/cancel
```

### Payment and Shipping Callbacks

```text
POST   /api/v1/payments/:paymentId/retry
GET    /api/v1/payments/:paymentId
POST   /api/v1/webhooks/payments/:provider
POST   /api/v1/webhooks/shipments/:provider
```

### Custom Orders

```text
POST   /api/v1/custom-orders
GET    /api/v1/custom-orders
GET    /api/v1/custom-orders/:id
POST   /api/v1/custom-orders/:id/messages
POST   /api/v1/custom-orders/:id/quotes
POST   /api/v1/custom-orders/:id/quotes/:quoteId/accept
```

### Routing table

| Route group | Feature | Global route file |
|---|---|---|
| `/auth/*`, `/me` | `auth` | `public.go` (register/login/refresh) + `authenticated.go` (`/me`) |
| `/categories`, `/products*`, `/artisans*` (GET) | `catalogue`, `product`, `artisan` | `public.go` |
| `/artisan-applications*`, `/artisan/*` | `artisan` | `authenticated.go` |
| `/admin/product-submissions*`, `/admin/products/:id/suspend` | `moderation` | `admin.go` |
| `/warehouse/*` | `warehouse`, `inventory` | `admin.go` (warehouse_agent role) |
| `/cart*`, `/checkout`, `/orders*` | `cart`, `checkout`, `order` | `authenticated.go` |
| `/payments/*` (non-webhook) | `payment` | `authenticated.go` |
| `/webhooks/payments/:provider`, `/webhooks/shipments/:provider` | `payment`, `shipment` | `public.go`, with signature-verification middleware instead of JWT (§52) |
| `/custom-orders*` | `custom_order` | `authenticated.go` |

Handlers for admin-only groups additionally pass through
`middleware/authorization.go` with the relevant Casbin role.

## 49. Authentication and Authorization

Authentication includes access tokens, refresh tokens, session records,
rotation, revocation, OTP, and rate limiting.

Recommended roles:

```text
customer
artisan
warehouse_agent
moderator
admin
```

Frontend guards improve UX. Backend application policies are authoritative.
See §40 for the concrete public/authenticated request-flow diagrams.

## 50. Validation

Validate at two levels — neither substitutes for the other:

### Request DTO (`server/request/*.go`, `validator/v10` struct tags)

- Required fields.
- String length.
- Numeric range.
- UUID format.
- Supported file metadata.
- ISO currency format.
- Pagination bounds.

This layer never sees a `Principal` or touches the database.

### Domain/Application (`domain/policy.go`, `application/commands/handler.go`)

- Actor permission.
- Resource ownership.
- Valid status transition.
- Stock availability.
- Product readiness.
- Amount consistency.
- Duplicate callback.
- Idempotency.

Do not rely on frontend validation (§12's ViewModel layer) — it is UX only
and is never trusted here.

## 51. Error Catalogue

```text
VALIDATION_ERROR · INVALID_CREDENTIALS · AUTHENTICATION_REQUIRED · FORBIDDEN
RESOURCE_NOT_FOUND · CONFLICT · ARTISAN_NOT_APPROVED · PRODUCT_NOT_EDITABLE
PRODUCT_NOT_SELLABLE · INVALID_STATE_TRANSITION · OUT_OF_STOCK
RESERVATION_EXPIRED · IDEMPOTENCY_CONFLICT · PAYMENT_AMOUNT_MISMATCH
INVALID_WEBHOOK_SIGNATURE · DUPLICATE_WEBHOOK · UNSUPPORTED_FILE_TYPE
FILE_TOO_LARGE · RATE_LIMITED · INTERNAL_ERROR
```

`internal/pkg/apperror/codes.go` defines these as typed constants;
`internal/server/error_handler/error_handler.go` is the single place that
maps a code to an HTTP status and the public error envelope from §47.
Feature code raises a code, never an HTTP status directly. The frontend API
client maps every one of these to a translation key (§21) — never a raw
message.

## 52. Idempotency

Checkout, payment creation, refund, shipment creation, and external webhook
processing all go through `internal/pkg/idempotency`:

1. Client sends an `Idempotency-Key` header (checkout) or the provider sends
   a unique external event ID (webhooks).
2. `pkg/idempotency` stores key, actor, route, request hash, status, and
   response — this table is part of the `000014_audit_idempotency` migration.
3. Same key + same request hash → return the saved response, no re-execution.
4. Same key + different request hash → `IDEMPOTENCY_CONFLICT`.

Webhook signature verification (`INVALID_WEBHOOK_SIGNATURE`) happens in
feature middleware *before* idempotency is even checked — an unverified
webhook is rejected outright, not deduplicated.

## 53. File Upload

Order of operations, enforced in that order:

1. Validate request size before full processing.
2. Validate MIME by content signature, not by trusting the extension.
3. Generate the object key server-side — the original filename is never
   used as or folded into the storage path.
4. Upload to a private MinIO bucket (§55: buckets private by default).
5. Store metadata in Postgres.
6. Scan if a scanner is configured.
7. Issue a signed URL only for an authorised read.
8. For product media specifically: publish to the public catalogue only
   after moderation approval, not at upload time.

Applies to `product_media`, artisan verification documents, custom-order
attachments, and moderation files alike. The frontend never renders these as
permanent public URLs — see §22's signed-URL-with-expiry requirement.

## 54. OpenAPI

Every route documents: authentication requirement, Casbin permission,
request schema, response schema, possible error codes (§51), pagination,
`Idempotency-Key` usage where applicable, and one example payload. The spec
is validated in CI — a route merged without a matching OpenAPI entry fails
the pipeline, not just a code review comment.

## 55. Data Responsibilities

PostgreSQL is the source of truth for users, sessions, artisans, catalogue,
products, moderation, inventory, reservations, carts, orders, payments,
shipments, custom orders, reviews, notifications, audit, idempotency, and
outbox.

Redis is used for rate limiting, short-lived OTP state, justified caching,
session assistance, and short-lived coordination.

MinIO stores product media, artisan documents, custom-order attachments, and
moderation files. Buckets are private by default.

## 56. Testing

```text
apps/api/tests/
├── unit/          # domain state transitions, money calc, inventory invariants,
│                  # order totals, permission policies, DTO mapping — no mocks needed
├── integration/   # real Postgres repositories, transaction rollback, Casbin
│                  # enforcement, MinIO storage, Redis rate limiting, migration up/down
├── concurrency/   # last-unit reservation race, duplicate checkout, duplicate
│                  # webhook, reservation-expiry-vs-payment-confirmation race
└── e2e/           # product → warehouse → sale; customer purchase; payment
                   # failure; full shipment flow
```

`unit/` mirrors `application/` tests using generated mocks (§57); `integration/`
and `concurrency/` require the real test-container stack and never use mocks.

Frontend testing — component, integration, Playwright E2E, and visual
regression — is defined in full in §27; it is the other half of the same CI
suite.

## 57. Mock Generation Workflow

Every interface in `domain/repository.go` and every port in
`application/ports/` must have a generated mock, so application and domain
tests never touch a real database, Redis instance, or provider SDK.

```bash
# Regenerate mocks for a single feature's domain + application ports
mockery --dir=internal/features/product/domain --output=internal/features/product/mocks --outpkg=mocks --all
mockery --dir=internal/features/product/application/ports --output=internal/features/product/mocks --outpkg=mocks --all
```

Add to `apps/api/Makefile`:

```makefile
mocks:
	@for f in $$(ls internal/features); do \
		mockery --dir=internal/features/$$f/domain --output=internal/features/$$f/mocks --outpkg=mocks --all 2>/dev/null; \
		mockery --dir=internal/features/$$f/application/ports --output=internal/features/$$f/mocks --outpkg=mocks --all 2>/dev/null; \
	done

wire:
	wire ./internal/container/...
```

Test layering rule:

- **Domain tests**: no mocks needed — pure functions/invariants, no I/O.
- **Application tests**: use generated mocks for repositories and ports;
  assert via `testify/assert` and `testify/require`; verify call expectations
  via `testify/mock`.
- **Data tests**: hit a real (test-container) Postgres/Redis instance —
  mocks are not used here.
- Regenerate mocks whenever a repository or port interface signature changes;
  a stale mock is a build break.

## 58. Required Commands

```bash
go mod download
gofmt -w .
go vet ./...
go test ./...
go build ./...
migrate -path db/migrations -database "$DATABASE_URL" up
```

Combined with the `mocks` and `wire` Makefile targets (§57), this is the full
local + CI command set for the backend — see §31 for the frontend equivalent.

## 59. Architecture Checklist

### Backend

- Domain imports no framework or infrastructure package.
- Application depends on interfaces.
- Server handlers contain no business rules.
- SQLBoiler models remain inside data.
- Transactions are controlled by application use cases.
- Request and response structs stay in server.
- Feature middleware is separate from global middleware.
- External callbacks are verified and idempotent.
- Money uses integer minor units.
- UTC is used internally.
- Every domain repository interface and application port has a
  corresponding entry in `mocks/`, regenerated on signature change.
- `bootstrap/` contains no business logic and is the only package allowed
  to call `os.Exit` or panic on startup misconfiguration.
- Application-layer tests use mocks; data-layer tests use real
  test-container instances — never the reverse.
- JWT verification happens once, in global `middleware/authentication.go`;
  feature middleware only adds authorization/ownership checks on top.
- Repository interfaces use `GetByID` / `Create` / `Update` naming
  consistently across every feature.
- Every inventory balance is explainable from the movement log, never
  edited independently of it.
- Every route has a matching OpenAPI entry validated in CI; every error
  path uses a code from §51, never a raw HTTP status string.
- Idempotency keys are checked before checkout/payment/refund/shipment
  logic runs, and webhook signatures are verified before idempotency is
  checked.
- File uploads never trust the client-provided filename as a storage path;
  product media is published only post-moderation.
- `go vet`, `gofmt -w .`, `go test ./...`, and `migrate ... up` all pass
  before merge.

### Frontend

- Route files remain thin; every route lives under `app/[locale]/`.
- Domain imports no React or HTTP library.
- Repository contracts live inward.
- Data implements repository contracts.
- Views consume ViewModels.
- Presentational components do not call APIs.
- ViewModels do not render JSX.
- Backend rules are not duplicated as trusted frontend rules.
- All visible text is translated; no concatenated translation fragments.
- Arabic RTL is tested on every major page, not only navigation.
- Loading, empty, error, forbidden, not-found, and offline states exist,
  using the shared `core/components/feedback/` components (§24).
- Design tokens from §6 are used rather than one-off colours or spacing.
- Every fetch handles the six states listed in §16.
- Checkout and payment creation reuse the same idempotency key on retry —
  never a silent retry with a new one (§21).
- Signed URLs, not permanent public URLs, are used for private media (§22).
- Every page satisfies the per-page Definition of Done in §30.

## 60. Definition of Done

The architecture is correctly applied when frontend features support domain,
repository contracts, data implementations, optional application use cases,
ViewModels, views, and components where needed, styled from the shared
design tokens in §6 and passing the per-page checklist in §30; backend
features support domain, application, data, server, and generated mocks; the
server includes router, middleware, request, response, handler, and
presenter responsibilities; DI wiring is generated via `wire` rather than
hand-written; framework dependencies stay outside domain; persistence models
do not leak; features communicate through stable public interfaces; every
route has a matching OpenAPI entry and stable error code that the frontend
maps to a translated message; idempotency and webhook signature verification
are in place for every external callback; Arabic RTL, accessibility, and the
Core Web Vitals targets in §28 are met on every public page; and the system
can evolve feature by feature, on both sides of the API contract, without
coupling the whole codebase.
