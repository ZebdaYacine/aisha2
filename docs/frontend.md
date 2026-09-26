# AISHA MVP Frontend Development Guide

## 1. Product Goal

Build a premium, visually distinctive, secure, multilingual e-commerce storefront for Algerian artisanal products.

The platform must allow customers to discover products, understand their cultural origin, learn about the artisan, and complete a purchase through a small number of clear steps.

AISHA should feel like a modern international fashion and lifestyle store while preserving an authentic Algerian identity.

The visual experience should combine:

* The editorial minimalism of Zara.
* The bold product presentation of Nike.
* Algerian craftsmanship, patterns, materials, landscapes, and cultural storytelling.
* A calm and premium shopping experience.
* Strong product photography.
* Clear typography.
* Generous whitespace.
* Smooth but restrained interactions.

The interface must not directly copy Zara, Nike, or another brand. It should use similar design principles while maintaining an original AISHA identity.

---

## 2. Frontend Principles

* Next.js App Router.
* Server Components by default.
* Client Components only when interactivity requires them.
* Strict TypeScript.
* shadcn/ui components.
* Tailwind CSS.
* React Hook Form and Zod for forms.
* Translation keys for all visible text.
* Full Arabic RTL support.
* Responsive and accessible design.
* Mobile-first implementation.
* No business-critical rules only in the frontend.
* Avoid unnecessary visual clutter.
* Prefer reusable design-system components over page-specific styling.
* Maintain consistent spacing, typography, borders, icons, and interaction states.

---

## 3. Visual Direction

### 3.1 Brand Personality

AISHA should feel:

* Premium.
* Authentic.
* Contemporary.
* Cultural.
* Editorial.
* Warm.
* Trustworthy.
* Handmade but professionally presented.

Avoid:

* Generic marketplace layouts.
* Excessive gradients.
* Too many rounded cards.
* Large amounts of unrelated colour.
* Heavy shadows.
* Crowded product grids.
* Decorative animations that slow down shopping.
* Dashboard styling on customer-facing pages.
* Using every available shadcn/ui component without a clear purpose.

### 3.2 Design Inspiration

Use Zara-inspired principles for:

* Minimal navigation.
* Large editorial imagery.
* Neutral backgrounds.
* Strong typography.
* Full-width product sections.
* Spacious layouts.
* Simple product cards.
* Discreet interface controls.

Use Nike-inspired principles for:

* Strong campaign sections.
* Bold headlines.
* High-quality product galleries.
* Clear calls to action.
* Product storytelling.
* Responsive mobile interactions.
* Visually engaging category sections.
* Sticky purchase controls on product pages.

AISHA must remain visually original by incorporating:

* Algerian visual references.
* Local materials and textures.
* Artisan portraits.
* Regional origin labels.
* Cultural stories.
* Subtle geometric patterns inspired by Algerian crafts.
* Earthy and natural colours.

---

## 4. Design System

Create a central design system before building complete pages.

Recommended structure:

```text
components/
├── ui/
├── layout/
├── commerce/
├── editorial/
├── feedback/
└── forms/

styles/
├── globals.css
├── tokens.css
└── utilities.css
```

### 4.1 Colour Palette

Use a restrained palette.

Suggested base colours:

```text
Background:        warm off-white
Surface:           white
Primary text:      near black
Secondary text:    warm grey
Muted surface:     light stone
Border:            soft neutral grey
Primary accent:    deep terracotta
Secondary accent:  olive green
Premium accent:    muted gold
Error:             deep red
Success:           forest green
```

Example CSS variables:

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

Do not apply accent colours to every element. Most of the storefront should remain neutral so that product imagery remains dominant.

### 4.2 Typography

Use an editorial typography hierarchy.

Recommended approach:

* A clean sans-serif font for navigation, forms, prices, and interface controls.
* A refined display font for large campaign titles and cultural storytelling.
* An Arabic font with excellent readability and balanced weight.

Possible font pairing:

```text
Latin UI:      Inter or Geist
Latin Display: Cormorant Garamond or Playfair Display
Arabic UI:     IBM Plex Sans Arabic or Noto Sans Arabic
Arabic Display: Noto Kufi Arabic or Readex Pro
```

Typography rules:

* Use large headings with compact line height.
* Avoid excessive font-weight variation.
* Product names should remain readable and understated.
* Prices must be immediately visible.
* Use uppercase labels sparingly.
* Arabic headings must preserve readable spacing.
* Avoid justified text.
* Limit long paragraphs to readable line widths.

Suggested scale:

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

### 4.3 Spacing

Use generous spacing similar to premium fashion stores.

Recommended spacing principles:

* Large separation between editorial sections.
* Smaller and consistent spacing inside forms and product metadata.
* Avoid placing every section inside a bordered card.
* Use whitespace to define hierarchy.

Example section spacing:

```text
Mobile:  py-16
Tablet:  py-20
Desktop: py-28 or py-32
```

Recommended page container:

```tsx
<div className="mx-auto w-full max-w-[1600px] px-4 sm:px-6 lg:px-10 xl:px-14">
  {children}
</div>
```

### 4.4 Corners, Borders, and Shadows

Use:

* Small border radii.
* Thin neutral borders.
* Very subtle shadows.
* Square or nearly square product imagery.
* Full-width image sections without unnecessary containers.

Avoid:

* Large rounded cards everywhere.
* Strong drop shadows.
* Floating glassmorphism panels.
* Excessively pill-shaped buttons.

Suggested defaults:

```text
Cards: rounded-sm or rounded-none
Inputs: rounded-sm
Buttons: rounded-sm
Badges: rounded-full only when semantically appropriate
```

---

## 5. Interaction and Motion

Motion should communicate quality, not decoration.

Use:

* Subtle image zoom on hover.
* Fade and translate transitions for menus and dialogs.
* Smooth cart drawer transitions.
* Underline animations on navigation links.
* Crossfade between product gallery images.
* Skeleton loaders that match final content dimensions.
* Sticky purchase panels where appropriate.

Avoid:

* Large page entrance animations.
* Continuous motion.
* Parallax on every section.
* Animations that delay interaction.
* Bouncing buttons.
* Excessive use of Framer Motion.

Recommended timings:

```text
Fast interaction: 150ms
Standard transition: 250ms
Drawer or modal: 300ms
Editorial reveal: 400ms maximum
```

Respect:

```css
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    scroll-behavior: auto !important;
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## 6. Routes

Suggested MVP routes:

```text
/[locale]
/[locale]/products
/[locale]/products/[slug]
/[locale]/artisans
/[locale]/artisans/[slug]
/[locale]/categories/[slug]
/[locale]/login
/[locale]/register
/[locale]/cart
/[locale]/checkout
/[locale]/payment/[paymentId]
/[locale]/account
/[locale]/account/orders
/[locale]/account/orders/[id]
/[locale]/custom-orders
/[locale]/artisan
/[locale]/artisan/profile
/[locale]/artisan/products
/[locale]/artisan/products/new
/[locale]/artisan/products/[id]
/[locale]/admin
/[locale]/admin/artisans
/[locale]/admin/moderation
/[locale]/admin/warehouse
/[locale]/admin/inventory
/[locale]/admin/media
/[locale]/admin/artisan-applications (application Details modal includes authorized document/profile-media previews)
/[locale]/admin/orders
/[locale]/admin/audit
```

---

## 7. Global Storefront Layout

### 7.1 Announcement Bar

Optional, compact announcement bar for:

* Delivery information.
* New collections.
* Artisan campaigns.
* Important service updates.

Requirements:

* Maximum one short message.
* Dismissible only when necessary.
* Fully translated.
* Must not consume excessive vertical space.

### 7.2 Header

The header should be minimal and editorial.

Desktop layout:

```text
Left:    navigation
Centre:  AISHA logo
Right:   search, locale, account, cart
```

RTL layout must mirror naturally.

Recommended navigation:

```text
New
Products
Categories
Artisans
Our Story
```

Behaviour:

* Transparent over the home hero when contrast is sufficient.
* Solid background after scrolling.
* Sticky header.
* Compact height.
* Search may open as a full-width overlay.
* Cart should open as a drawer.
* Mobile navigation should use a full-height sheet or menu.
* Show cart item count without visual noise.

### 7.3 Footer

Include:

* AISHA story.
* Shop links.
* Artisan information.
* Customer support.
* Terms and privacy.
* Social links.
* Newsletter only when implemented.
* Language selector.
* Country or delivery information when required.

Use a spacious multi-column layout on desktop and accordions on mobile.

---

## 8. Public Pages

## 8.1 Home Page

The home page should feel like a premium campaign landing page rather than a dashboard.

Recommended section order:

```text
1. Full-width hero
2. Featured collection
3. Category editorial grid
4. Brand and cultural story
5. Featured artisans
6. Product carousel or grid
7. Values and provenance
8. Regional craftsmanship feature
9. Final discovery call to action
```

### Hero

Use:

* Full viewport or near-full viewport image.
* Strong headline.
* One short supporting sentence.
* One or two clear actions.
* High-quality Algerian craft imagery.
* Text placement that remains readable on mobile.

Example:

```text
Headline:
Crafted in Algeria.
Made to be remembered.

Actions:
Explore the collection
Meet the artisans
```

Hero requirements:

* Use art-directed responsive images.
* Support alternative mobile crops.
* Avoid carousels in the initial MVP.
* Keep visible copy concise.
* Do not cover important image content with text.
* Ensure sufficient contrast.

### Featured Collection

Use large editorial tiles rather than small cards.

Possible layout:

```text
Desktop:
- One large tile occupying two-thirds width.
- Two stacked smaller tiles occupying one-third width.

Mobile:
- Vertical full-width tiles.
```

Each tile should include:

* Collection name.
* Short story or label.
* Subtle call to action.
* Background image.
* Accessible text contrast.

### Categories

Use large image-based category blocks.

Examples:

* Jewellery.
* Textiles.
* Ceramics.
* Leather.
* Home décor.
* Traditional garments.

Do not present all categories as identical small icons.

### Featured Products

Use a clean grid with strong images.

Desktop:

```text
4 columns on large screens
3 columns on medium screens
2 columns on small tablets
2 or 1 columns on mobile depending on image width
```

Alternate selected product sections with larger editorial product blocks to avoid a generic marketplace appearance.

### Featured Artisans

Use:

* Artisan portrait.
* Name.
* Region.
* Primary craft.
* Short quotation or story.
* Link to profile.

Artisan cards should feel editorial and human, not administrative.

### Values Section

Present:

* Verified origin.
* Fair partnership.
* Quality inspection.
* Responsible materials.

Use concise copy and simple line icons.

Do not display verification badges unless the corresponding information is verified by the backend.

---

## 8.2 Product Listing

The product listing should prioritise imagery.

Page layout:

```text
Title and short introduction
Category navigation
Filter and sort toolbar
Product grid
Pagination or load-more control
```

Features:

* Search.
* Category filter.
* Artisan filter.
* Price filter if currencies are consistent.
* Availability.
* Pagination.
* Sort by newest, price, or relevance when supported.

### Filter Experience

Desktop:

* Compact horizontal toolbar or left filter panel.
* Sticky filter controls where useful.

Mobile:

* Filter button opens a bottom sheet or full-height sheet.
* Display active filter count.
* Provide clear apply and reset actions.

### Product Grid

Product cards should include:

* Primary product image.
* Optional alternate image on hover.
* Product name.
* Artisan name.
* Price and currency.
* Availability or made-to-order label.
* Verified badges only when applicable.
* Wishlist button only if the feature exists.

Product card styling:

* Minimal borders.
* No heavy card background.
* Image occupies most of the component.
* Use consistent image ratios.
* Product information appears beneath the image.
* Avoid large add-to-cart buttons inside every card.
* Entire product title and image area should link to the detail page.
* Do not show excessive metadata.

Hover behaviour:

* Slight image scale.
* Optional second image crossfade.
* Visible product action.
* Keyboard focus equivalent.

---

## 8.3 Product Detail

The product detail page is a primary conversion page.

Desktop layout:

```text
Left:  product gallery, approximately 60–65%
Right: sticky product information, approximately 35–40%
```

Mobile layout:

```text
Gallery
Product summary
Purchase controls
Story and details
Artisan information
Recommendations
```

Must show:

* Product gallery.
* Localised name and description.
* Artisan identity and profile link.
* Product history and cultural context.
* Materials.
* Manufacturing method.
* Usage.
* Dimensions and weight.
* Origin.
* Eco-friendly or fair-trade badges only when verified.
* Price and currency.
* Availability.
* Made-to-order option when enabled.
* Add-to-cart action.

### Gallery

Use:

* Large images.
* Swipe support on mobile.
* Image thumbnails on desktop where appropriate.
* Zoom or lightbox as an optional basic MVP feature.
* Loading placeholder.
* Proper aspect ratios.
* Accessible image labels.
* Video only when optimised and useful.

### Purchase Panel

Include:

* Product name.
* Artisan and region.
* Price.
* Tax or delivery note when applicable.
* Availability.
* Variant or quantity controls.
* Made-to-order information.
* Estimated production time when trusted.
* Add-to-cart button.
* Secondary custom-order action when supported.
* Delivery and return summary.
* Verification indicators.

Primary button should be visually prominent and full width.

On mobile, a sticky bottom purchase bar may show:

```text
Price
Add to cart
```

It must not hide important content or interfere with browser controls.

### Product Story

Use an editorial section with:

* Large image.
* Cultural context.
* Artisan quotation.
* Production process.
* Region of origin.
* Material details.

Avoid hiding all important information inside accordions. Accordions may be used for secondary technical information.

### Recommendations

Show:

* More from this artisan.
* Similar products.
* Related cultural collection.

Do not display recommendations when no valid data exists.

---

## 8.4 Artisan Listing

The artisan listing should be visually led by portraits and workshops.

Include:

* Search.
* Region filter.
* Craft filter.
* Featured artisan story.
* Artisan grid.

Each artisan preview should include:

* Portrait or workshop image.
* Name.
* Region.
* Craft.
* Short description.
* Number of active products when useful.

---

## 8.5 Artisan Detail

The artisan profile should feel like an editorial profile.

Recommended structure:

```text
Hero portrait or workshop image
Artisan name and region
Short biography
Craft history
Techniques and materials
Product collection
Gallery or workshop story
Verified information
```

Use long-form typography with readable line width.

The artisan’s story should not look like a seller dashboard.

---

## 9. Language and RTL

Supported locales:

```text
ar
fr
en
es
```

Recommended route strategy:

```text
/[locale]/...
```

Rules:

* Set `<html lang>` and `dir`.
* Arabic uses `dir="rtl"`.
* Icons that imply direction must mirror where appropriate.
* Form alignment and table layouts must work in RTL.
* Dates, numbers, and currencies must use locale-aware formatting.
* Do not concatenate translated fragments.
* Use logical CSS properties where possible.
* Avoid hardcoded `left` and `right` when `start` and `end` are appropriate.
* Test every major page in Arabic, not only the navigation.
* Product grids should maintain natural visual order in RTL.
* Carousels must support RTL navigation correctly.
* Breadcrumb arrows must mirror in Arabic.
* Drawer and sheet opening directions should feel natural for the active language.

Translation file layout:

```text
messages/
├── ar.json
├── fr.json
├── en.json
└── es.json
```

Example:

```tsx
const direction = locale === "ar" ? "rtl" : "ltr";

return (
  <html lang={locale} dir={direction}>
    <body>{children}</body>
  </html>
);
```

---

## 10. Feature Structure

```text
features/
├── catalogue/
│   ├── api/
│   ├── components/
│   ├── schemas/
│   ├── types/
│   └── utils/
├── cart/
├── checkout/
├── artisans/
├── authentication/
├── account/
├── moderation/
├── warehouse/
└── custom-orders/
```

Keep API calls in a defined API layer.

Presentational components must not call random endpoints directly.

Recommended component separation:

```text
ProductCard
ProductGrid
ProductGallery
ProductPurchasePanel
ProductStory
ProductSpecifications
ArtisanCard
ArtisanProfileHero
CategoryEditorialCard
Money
AvailabilityLabel
VerificationBadge
EmptyState
ErrorState
LoadingSkeleton
```

---

## 11. Data Fetching

Use Server Components for:

* Catalogue pages.
* Product detail.
* Artisan public pages.
* Category pages.
* Initial account pages when authentication design supports server access.
* SEO metadata.
* Initial navigation data.

Use Client Components for:

* Cart interactions.
* Filters with client state.
* Forms.
* Upload progress.
* Product gallery interaction.
* Search overlays.
* Cart drawer.
* Rich interactive moderation or warehouse controls.

Every fetch must handle:

* Loading.
* Empty result.
* Recoverable error.
* Authentication required.
* Forbidden.
* Not found.

Use page-specific skeletons rather than a generic spinner.

Examples:

* Product grid skeleton.
* Product gallery skeleton.
* Artisan profile skeleton.
* Order detail skeleton.

Avoid full-page loading spinners for catalogue navigation.

---

## 12. Forms

Use React Hook Form and Zod.

Required forms:

* Registration.
* Login.
* Address.
* Artisan application.
* Artisan profile.
* Product creation and edit.
* Product submission.
* Moderation decision.
* Warehouse reception.
* Inspection.
* Inventory adjustment.
* Checkout.
* Custom-order request.
* Quote response.

Rules:

* Reuse Zod schemas for client feedback where practical.
* Backend remains authoritative.
* Display stable API error messages mapped to translated user text.
* Move focus to error summary after failed submit.
* Preserve user input after recoverable errors.
* Disable duplicate submissions while pending.
* Display field-level and form-level errors.
* Do not rely on placeholder text as a label.
* Use clear required and optional indicators.
* Use inline help only when useful.
* Group related fields into visual sections.
* Use a sticky action footer for long admin or artisan forms when beneficial.

Form styling should remain clean:

* Labels above fields.
* Comfortable vertical spacing.
* Thin borders.
* Strong focus states.
* Minimal decorative containers.

---

## 13. Cart

Cart behaviour:

* May be local for anonymous visitors and server-backed after login.
* Anonymous items are merged into the authenticated owner cart on login; matching quantities are summed and capped at 100.
* Display current price as provisional.
* Revalidate when opening the cart and checkout.
* Clearly mark unavailable items.
* Do not claim that cart quantity is reserved.
* Wishlist saves/removes use the authenticated owner API; unauthenticated visitors keep browsing without a fake saved state.

### Cart Drawer

The cart icon should open a side drawer for quick review.

Include:

* Product image.
* Product name.
* Artisan.
* Quantity.
* Price.
* Remove action.
* Provisional subtotal.
* View cart.
* Checkout.

The drawer must:

* Trap focus.
* Be keyboard accessible.
* Restore focus when closed.
* Work correctly in RTL.
* Provide empty-cart content.
* Show loading and revalidation states.

### Full Cart Page

Use a clean two-column desktop layout:

```text
Left:  cart items
Right: order summary
```

On mobile:

```text
Cart items
Order summary
Checkout action
```

The order summary may be sticky on desktop.

---

## 14. Checkout UX

Recommended steps:

```text
1. Review items
2. Delivery address
3. Shipping method
4. Payment
5. Confirmation
```

Use a clear progress indicator.

The checkout should be visually quieter than the main storefront.

Avoid:

* Promotional distractions.
* Large unrelated navigation menus.
* Product recommendations during payment.
* Hidden fees.
* Automatically selected optional services.

At confirmation:

* Display trusted server totals.
* Display currency.
* Display delivery address.
* Display shipping method.
* Require explicit final confirmation.
* Send an idempotency key.
* Handle price and stock changes gracefully.

Use a two-column layout on desktop:

```text
Left:  current checkout step
Right: persistent order summary
```

On mobile, show a collapsible order summary near the top.

---

## 15. Payment Return Page

States:

* `pending`: provider redirect occurred but trusted confirmation is incomplete.
* `paid`: backend confirms payment.
* `failed`: backend confirms failure.
* `expired`: payment session expired.

The page may poll the backend for a limited period.

It must not use a query string such as `success=true` as proof of payment.

Visual states:

### Pending

* Calm loading indicator.
* Clear explanation.
* Do not encourage duplicate payment.
* Provide safe navigation.

### Paid

* Strong confirmation heading.
* Order reference.
* Summary.
* Link to order details.
* Continue shopping action.

### Failed

* Explain the failure safely.
* Allow retry only through a backend-authorised flow.
* Provide support guidance.

### Expired

* Explain that the session expired.
* Offer a safe return to the cart or checkout.

---

## 16. Authentication UI

Include:

* Login.
* Registration.
* Protected route layout.
* Forbidden page.
* Session-expired handling.
* Logout.
* Role-aware navigation for convenience only.

Backend remains authoritative.

Authentication page style:

* Split-screen editorial image on desktop.
* Focused form panel.
* Single-column layout on mobile.
* Minimal distractions.
* Clear password requirements.
* Visible error feedback.
* Links between login and registration.

Do not present authentication pages as generic admin screens.

---

## 17. Account Area

The account area may be more functional than the public storefront, but it should remain visually consistent.

Recommended navigation:

```text
Overview
Orders
Addresses
Profile
Security
Custom orders
Logout
```

Desktop:

* Side navigation.
* Main content area.

Mobile:

* Compact tab selector or account menu.

Order cards should show:

* Order reference.
* Date.
* Status.
* Total.
* Product preview.
* Detail action.

---

## 18. Artisan Dashboard

MVP dashboard:

* Profile completion.
* Application status.

Artisan onboarding is a staged customer flow. The form first saves the
workshop/profile details as a draft, then uploads at least one supporting
document and one profile media item, and finally submits the application. The
submit action is unavailable until those files are selected or already exist;
the API repeats the same requirement so a client cannot bypass it.
* Product counts by status.
* Create product.
* Edit draft.
* Submit for review.
* View moderation reason.
* View accepted, reserved, and shipped stock for owned products.
* View custom-order requests when included.

Workshop and product management use responsive data tables on desktop and a
horizontal scroll container on narrow screens. Each table puts row actions in
the first column, provides one primary Add/New action above the table, and
uses read-only Details dialogs plus Update dialogs for editing. Workshop
deletion is confirmation-protected and remains subject to ownership/history
rules; products use the history-preserving Archive action.

The approved artisan workspace groups these areas into an accessible tab panel:
`Workshops`, `Private files`, and `Product authoring`. The Workshops tab is
the default; private files and product authoring remain hidden when seller
capabilities are unavailable or membership is suspended.
Approved artisans do not see the onboarding/application form in this area;
that form remains available only while an application is being started or
completed.

Private files use horizontally scrollable tables. Each profile-media row shows
the media name, media type, locale-formatted upload date, and View, Update, and
Delete actions. Add and Update open modal forms with a media-kind combobox and
local file picker; Delete is confirmation-protected. Application documents use
the same row treatment and an Add document modal. Product authoring keeps its
table as the index and opens New draft and Update inside a scrollable modal.

The header and footer language controls replace only the locale segment of the
current pathname and preserve query strings and hashes. The header displays a
flag plus the two-letter locale abbreviation; the option list retains the full
language names. Theme switching uses Light, Dark, and System modes, closes on
Escape or outside click, and shares the `aisha-theme` storage key with the
pre-hydration theme bootstrap.

The first visit in a browser starts with the AISHA brand splash: a short
animated story sequence, progress bar, and reduced-motion-safe styling. It is
stored in local storage after completion so route changes do not replay it.

Use a distinct application shell separate from the public storefront.

Recommended structure:

```text
Desktop sidebar
Top header
Page title and primary action
Status summary
Main workspace
```

Dashboard design rules:

* Do not imitate the editorial home page.
* Prioritise clarity and task completion.
* Use cards only for meaningful grouped information.
* Use data tables on desktop.
* Use stacked rows or cards on mobile.
* Keep the AISHA design tokens and typography.
* Provide clear status labels.
* Status must not be communicated by colour alone.

---

## 19. Admin and Warehouse UI

Admin and warehouse pages should prioritise efficiency.

The shared admin presentation uses a consistent control language across the
dashboard, user management, moderation, warehouse, artisan review, audit,
artisan workspace, and account surfaces:

* Searchable comboboxes are used for supported language, role, country, wilaya,
  category, workshop, product type, and workflow status selections.
* Multi-value role assignments render as removable chips and submit the exact
  selected role list to the API.
* Workflow states render as semantic badges with text and colour so the state
  is never communicated by colour alone.
* Administrative and customer-facing dates use locale-aware long date formats;
  audit timestamps include the local time.
* Native locale links remain available as a crawlable and accessible fallback
  below the interactive language combobox.

### Moderation Queue

Display:

* Product.
* Artisan.
* Submission date.
* Translation completeness.
* Media.
* Current status.
* Actions.

Actions require confirmation and reason.

Recommended desktop layout:

* Filter toolbar.
* Dense but readable table.
* Product preview drawer.
* Moderation action panel.

On mobile:

* Stacked moderation cards.
* Full-screen review sheet.

### Reception

Form fields:

* Product.
* Artisan.
* Quantity.
* Batch or parcel reference.
* Notes.
* Date.

Support:

* Product search.
* Barcode or reference entry when available.
* Confirmation summary.
* Clear success state.
* Prevention of duplicate submission.

### Inspection

Display received quantity and require the sum of outcomes to match the inspected quantity.

Show the calculation visibly:

```text
Accepted + Rejected + Damaged + Pending = Received quantity
```

Prevent submission when totals do not match.

### Order Fulfilment

Display only paid and eligible orders.

Include:

* Order reference.
* Payment status.
* Fulfilment status.
* Delivery method.
* Required items.
* Picking state.
* Packing state.
* Shipment reference.

---

## 20. API Client

Create one typed client with:

* Base URL.
* Cookie or token strategy.
* Correlation ID handling.
* JSON parsing.
* Standard error parsing.
* Refresh behaviour where applicable.
* Abort signal.
* Retry only for safe idempotent reads unless specifically designed.

Do not silently retry checkout or payment creation without the same idempotency key.

Suggested structure:

```text
lib/api/
├── client.ts
├── errors.ts
├── request.ts
├── response.ts
├── correlation.ts
└── endpoints/
```

Example error model:

```ts
export interface ApiError {
  code: string;
  message: string;
  correlationId?: string;
  fieldErrors?: Record<string, string[]>;
}
```

UI components should map stable backend error codes to translation keys.

Do not display raw backend stack traces or internal error messages.

---

## 21. Accessibility

* Use native semantic controls.
* Every input has a label.
* Error text is associated with its input.
* Keyboard navigation works.
* Focus is visible.
* Dialogs trap and restore focus.
* Images have useful alt text.
* Video has captions or a transcript when practical.
* Status is not communicated by colour alone.
* Touch targets are sufficiently large.
* Skip links are available.
* Heading hierarchy is valid.
* Product galleries are keyboard accessible.
* Carousels do not auto-advance unexpectedly.
* Drawers and overlays announce their purpose.
* Loading and success states use appropriate live regions.
* Text contrast meets WCAG requirements.
* Focus indicators must remain visible on image backgrounds.
* RTL must not break keyboard or reading order.

---

## 22. Media

Use optimised image components.

Product media requirements:

* Correct aspect ratios.
* Placeholder while loading.
* Basic zoom or gallery for the MVP.
* Do not expose private draft media with permanent public URLs.
* Handle signed URLs with expiry.
* Provide different responsive image sizes.
* Avoid loading full-resolution images in product grids.
* Prevent layout shifts.
* Support portrait, square, and landscape editorial media.
* Use consistent crop rules by component.

Recommended ratios:

```text
Product card:       4:5
Product gallery:    4:5 or original constrained ratio
Artisan portrait:   3:4
Editorial banner:   16:9 or 3:2
Category tile:      4:5 or 1:1
Mobile hero:        4:5
Desktop hero:       16:9
```

---

## 23. SEO and Discovery

MVP SEO:

* Localised title and description.
* Product and artisan metadata.
* Canonical URLs.
* Open Graph images.
* Sitemap for active public content.
* Structured product data only when accurate.
* Localised alternate URLs.
* Meaningful product slugs.
* No indexing of draft, moderation, cart, checkout, or account pages.
* Server-rendered critical product information.
* Share images that match the AISHA visual identity.

---

## 24. Responsive Behaviour

Define behaviour explicitly for:

```text
Mobile: 320px and above
Small tablet: 640px and above
Tablet: 768px and above
Desktop: 1024px and above
Large desktop: 1280px and above
Editorial wide: 1536px and above
```

Key requirements:

* Header becomes a compact mobile navigation.
* Filters move into a sheet.
* Product detail columns stack.
* Sticky purchase actions remain usable.
* Tables transform into mobile-friendly rows.
* Large display typography scales through `clamp`.
* Editorial images use mobile-specific crops.
* Checkout summary becomes collapsible.
* Product card text does not overflow.
* Arabic content is tested at every breakpoint.

---

## 25. Empty, Loading, and Error States

Every important page must have designed states.

### Empty States

Examples:

* Empty cart.
* No search results.
* No products in a category.
* No artisan products.
* No customer orders.
* No moderation submissions.
* No warehouse receptions.

Each empty state should include:

* Clear heading.
* Short explanation.
* One relevant action.
* Optional restrained illustration.

### Loading States

Use layout-matched skeletons.

Do not use a large centred spinner for all pages.

### Error States

Include:

* Human-readable message.
* Retry action when safe.
* Navigation alternative.
* Correlation ID for support when available.
* Distinct handling for forbidden and not found.

---

## 26. Component Styling Standards

### Buttons

Variants:

```text
Primary
Secondary
Outline
Ghost
Destructive
Text link
```

Primary storefront buttons should generally use:

* Strong contrast.
* Medium or large height.
* Minimal radius.
* Clear hover, focus, active, disabled, and pending states.

Example:

```tsx
<Button
  size="lg"
  className="h-12 w-full rounded-sm px-8 text-sm font-medium uppercase tracking-wide"
>
  {t("product.addToCart")}
</Button>
```

Do not use uppercase for Arabic text unless specifically appropriate.

### Cards

Use cards only when content needs visual grouping.

Product cards should generally not use:

* Visible card backgrounds.
* Heavy borders.
* Large shadows.
* Excessive internal padding.

### Badges

Use badges for:

* Verified provenance.
* Fair-trade status.
* Eco-friendly certification.
* New products.
* Made-to-order.
* Stock state.

Do not create badges for ordinary descriptive information.

### Dialogs and Sheets

Use:

* Dialogs for focused confirmation.
* Sheets for filters, mobile menus, cart, and detail panels.
* Alert dialogs for destructive or irreversible actions.

---

## 27. Testing

### Component Tests

* Product card.
* Product gallery.
* Money display.
* Language switch.
* RTL layout.
* Forms and validation.
* Order status.
* Availability label.
* Cart item.
* Verification badge.
* Filter controls.
* Mobile navigation.
* Payment status.

### Integration Tests

* API error mapping.
* Protected layouts.
* Cart state.
* Anonymous-to-authenticated cart merge.
* Checkout steps.
* Locale routing.
* Server-rendered product data.
* Signed media URL handling.
* Role-aware navigation.
* Moderation action validation.
* Warehouse quantity validation.

### Playwright Tests

* Browse products.
* Change locale.
* Verify Arabic RTL.
* Search and filter products.
* Open a product.
* Add an item to the cart.
* Register and log in.
* Complete checkout with the development payment adapter.
* Handle payment pending and then paid.
* View an order.
* Artisan creates a product.
* Artisan submits a product for review.
* Moderator approves a product.
* Warehouse receives and inspects stock.
* Verify mobile storefront navigation.
* Verify keyboard navigation for critical flows.

### Visual Regression

Add visual tests for:

* Home hero.
* Product listing.
* Product detail.
* Cart drawer.
* Checkout.
* Arabic RTL layouts.
* Mobile navigation.
* Artisan profile.
* Moderation queue.

---

## 28. Performance Requirements

Target strong Core Web Vitals.

Requirements:

* Optimise hero images.
* Preload only essential fonts and media.
* Avoid unnecessary client-side JavaScript.
* Lazy-load below-the-fold content.
* Use Server Components for product content.
* Avoid heavy animation libraries unless justified.
* Limit third-party scripts.
* Use dynamic imports for rich admin controls.
* Prevent cumulative layout shift.
* Use appropriate cache and revalidation strategies.
* Keep product card components lightweight.
* Avoid loading complete product objects when summary data is sufficient.

Suggested performance targets:

```text
LCP: below 2.5 seconds
CLS: below 0.1
INP: below 200 milliseconds
```

---

## 29. Suggested Home Page Component Tree

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

---

## 30. Suggested Product Page Component Tree

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

---

## 31. Tailwind Styling Guidance

Prefer reusable primitives and readable class composition.

Example editorial section:

```tsx
<section className="py-16 sm:py-20 lg:py-28">
  <div className="mx-auto max-w-[1600px] px-4 sm:px-6 lg:px-10 xl:px-14">
    <div className="mb-10 flex items-end justify-between gap-6 lg:mb-14">
      <div className="max-w-2xl">
        <p className="mb-3 text-xs font-medium uppercase tracking-[0.2em] text-muted-foreground">
          {t("home.featured.label")}
        </p>

        <h2 className="text-3xl font-medium tracking-tight sm:text-4xl lg:text-5xl">
          {t("home.featured.title")}
        </h2>
      </div>

      <Link
        href={`/${locale}/products`}
        className="hidden border-b border-foreground pb-1 text-sm font-medium md:inline-flex"
      >
        {t("common.viewAll")}
      </Link>
    </div>

    <FeaturedProductGrid />
  </div>
</section>
```

Example product image:

```tsx
<div className="group relative aspect-[4/5] overflow-hidden bg-muted">
  <Image
    src={product.image.url}
    alt={product.image.alt}
    fill
    sizes="(max-width: 640px) 50vw, (max-width: 1024px) 33vw, 25vw"
    className="object-cover transition-transform duration-500 ease-out group-hover:scale-[1.025]"
  />
</div>
```

---

## 32. Implementation Phases

### Phase 1: Design Foundation

* Brand colours.
* Typography.
* Spacing scale.
* CSS variables.
* Button variants.
* Form controls.
* Header.
* Footer.
* Container.
* Responsive rules.
* RTL foundation.

### Phase 2: Public Storefront

* Home.
* Product listing.
* Product detail.
* Category page.
* Artisan listing.
* Artisan profile.
* Search.
* Cart drawer.

### Phase 3: Authentication and Checkout

* Login.
* Registration.
* Cart page.
* Address.
* Checkout steps.
* Payment return.
* Confirmation.

### Phase 4: Customer Account

* Account overview.
* Orders.
* Order detail.
* Addresses.
* Custom orders.

### Phase 5: Artisan Workspace

* Artisan dashboard.
* Profile.
* Product management.
* Media upload.
* Submission and moderation feedback.
* Stock overview.

### Phase 6: Admin and Warehouse

* Moderation queue.
* Reception.
* Inspection.
* Inventory.
* Paid-order fulfilment.
* Audit views.

### Phase 7: Quality

* Accessibility audit.
* RTL audit.
* Performance audit.
* Visual regression.
* Playwright coverage.
* Cross-browser testing.
* Responsive review.

---

## 33. Definition of Done for Every Page

A page is not complete until it has:

* Desktop layout.
* Tablet layout.
* Mobile layout.
* Arabic RTL layout.
* Loading state.
* Empty state where applicable.
* Error state.
* Keyboard accessibility.
* Visible focus states.
* Translated text.
* Correct locale formatting.
* Optimised images.
* Metadata where public.
* Component tests where appropriate.
* No raw backend error messages.
* No broken layout with long translated content.
* No business-critical rules implemented only in the frontend.

---

## 34. Required Commands

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

---

## 35. Final UI Quality Rules

The storefront must:

* Feel premium before adding complex functionality.
* Use product photography as the main visual element.
* Maintain generous whitespace.
* Use consistent image ratios.
* Avoid generic marketplace styling.
* Avoid excessive cards, badges, shadows, and rounded corners.
* Keep calls to action clear.
* Make prices and availability easy to find.
* Tell the artisan and product story without interrupting the purchase flow.
* Work equally well in Arabic, French, English, and Spanish.
* Look carefully designed on mobile, not merely compressed from desktop.
* Remain fast, accessible, and secure.
* Preserve an original AISHA visual identity rather than copying an existing store.
