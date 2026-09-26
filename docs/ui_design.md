# AISHA E-commerce UI Design Specification

## 1. Purpose

This document defines the visual language, page structure, reusable components, responsive behavior, and user experience of the AISHA e-commerce application.

AISHA is a multilingual marketplace dedicated to Algerian handmade and artistic products. The interface must combine:

- the clarity and conversion-focused structure of a modern e-commerce store;
- the warmth, authenticity, and storytelling of Algerian craftsmanship;
- a premium editorial appearance;
- simple purchasing flows;
- strong visibility for artisans, product origin, materials, and cultural history;
- responsive design for desktop, tablet, and mobile;
- full support for Arabic RTL, French, English, and Spanish.

The application must be built with **Next.js App Router**, **TypeScript**, **Tailwind CSS**, and **shadcn/ui**.

---

## 2. Product Identity

### Brand name

**AISHA — North Africa Treasures**

### Brand positioning

AISHA is not a generic marketplace. It is a curated destination for authentic Algerian handmade products.

The user interface must continuously communicate:

- authenticity;
- craftsmanship;
- trust;
- Algerian origin;
- fair trade;
- sustainability;
- premium quality;
- connection between the customer and the artisan.

### Main visual direction

The design should feel like a mix of:

- a premium lifestyle store;
- an editorial cultural magazine;
- an artisan marketplace;
- a modern international e-commerce application.

Avoid an overly decorative or folkloric interface. Algerian identity should appear through photography, textures, patterns, product stories, typography choices, and restrained visual accents.

---

## 3. Design Principles

### 3.1 Product-first

Product photography must remain the main visual focus. Cards, buttons, labels, and text should support the product instead of competing with it.

### 3.2 Artisan-first storytelling

Every major product surface should provide direct access to:

- artisan name;
- workshop name;
- region or wilaya;
- production technique;
- materials;
- product story;
- handmade or certified origin status.

### 3.3 Simple shopping flow

A customer should be able to discover, select, and buy a product with a minimum number of steps.

Primary actions must always be visible and understandable:

- Add to cart;
- Buy now;
- Select options;
- Continue to checkout;
- Confirm order.

### 3.4 Trust by design

The interface must make quality control, secure payment, stock availability, return conditions, artisan verification, and delivery information visible at the correct moment.

### 3.5 Calm visual hierarchy

Use generous whitespace, large imagery, limited colors, strong headings, readable body text, and clear spacing.

### 3.6 Multilingual by default

All visible content must use translation keys. Components must work in LTR and RTL layouts without custom page duplication.

---

## 4. Visual Style

## 4.1 Color palette

Use CSS variables compatible with shadcn/ui themes.

### Light theme

```css
:root {
  --background: 42 33% 97%;
  --foreground: 24 24% 14%;

  --card: 0 0% 100%;
  --card-foreground: 24 24% 14%;

  --popover: 0 0% 100%;
  --popover-foreground: 24 24% 14%;

  --primary: 24 43% 30%;
  --primary-foreground: 42 40% 98%;

  --secondary: 38 34% 90%;
  --secondary-foreground: 24 33% 22%;

  --muted: 38 24% 93%;
  --muted-foreground: 24 10% 42%;

  --accent: 148 20% 34%;
  --accent-foreground: 42 40% 98%;

  --destructive: 0 72% 51%;
  --destructive-foreground: 0 0% 98%;

  --border: 33 22% 84%;
  --input: 33 22% 84%;
  --ring: 24 43% 30%;

  --radius: 0.75rem;
}
```

### Dark theme

```css
.dark {
  --background: 24 22% 9%;
  --foreground: 42 30% 94%;

  --card: 24 18% 12%;
  --card-foreground: 42 30% 94%;

  --popover: 24 18% 12%;
  --popover-foreground: 42 30% 94%;

  --primary: 33 46% 62%;
  --primary-foreground: 24 24% 12%;

  --secondary: 24 14% 18%;
  --secondary-foreground: 42 25% 90%;

  --muted: 24 14% 16%;
  --muted-foreground: 33 12% 65%;

  --accent: 148 18% 46%;
  --accent-foreground: 24 24% 12%;

  --destructive: 0 62% 52%;
  --destructive-foreground: 0 0% 98%;

  --border: 24 14% 20%;
  --input: 24 14% 20%;
  --ring: 33 46% 62%;

  --radius: 0.75rem;
}
```

Dark mode is not an inverted color scheme; it keeps the same warm, artisanal identity. Deep brown becomes the background rather than white, brass/clay tones lighten slightly so they stay legible at low luminance, and olive is desaturated a touch to avoid vibrating against the dark surface. Product photography stays untouched — never dim, tint, or overlay images to "match" dark mode; only surfaces, text, and chrome change.

### Theme behavior

- Respect the operating system preference by default (`prefers-color-scheme`), applied via the `.dark` class on `<html>` (works alongside the existing `dir` attribute).
- Provide an explicit **`ThemeToggle`** control (sun/moon icon `Button`, `variant="ghost"`, `size="icon"`) in the header's right zone (desktop) and inside the mobile `Sheet`, with three states: `Light`, `Dark`, `System`.
- Persist the user's explicit choice (not just the system default) across sessions; `System` re-syncs with OS changes live.
- Apply the theme before first paint (no flash of the wrong theme) via an inline script or Next.js theme script pattern.
- Every custom composition (not just shadcn/ui primitives) must be built from the CSS variable tokens above, never hardcoded hex values, so both themes stay correct automatically.
- Badges, trust-signal colors (Section 26), and order-status colors (Section 17) keep their semantic hue but shift lightness/saturation for dark-mode contrast; status meaning never changes between themes.
- Maintain WCAG 2.2 AA contrast in both themes for text, icons, borders, and focus indicators.

### Recommended semantic colors

| Role | Suggested value | Usage |
|---|---:|---|
| Sand | `#F5F0E7` | Backgrounds and neutral sections |
| Clay | `#8A5A3B` | Primary actions and brand accents |
| Deep brown | `#2E231D` | Main text and premium headings |
| Olive | `#4F6958` | Sustainability and verified origin |
| Terracotta | `#C66B4E` | Promotional labels and editorial accents |
| Brass | `#B68A4A` | Premium badges and decorative details |
| White | `#FFFFFF` | Cards and content surfaces |

### Usage rule

Do not use more than one strong accent color in the same component. Most screens should remain neutral and image-led.

---

## 4.2 Typography

Use font families that support Latin and Arabic well.

### Suggested fonts

- Latin headings: `Cormorant Garamond`, `Playfair Display`, or `DM Serif Display`.
- Latin body: `Inter` or `Manrope`.
- Arabic headings: `Noto Kufi Arabic` or `IBM Plex Sans Arabic`.
- Arabic body: `Noto Sans Arabic` or `IBM Plex Sans Arabic`.

### Type scale

| Token | Desktop | Mobile | Usage |
|---|---:|---:|---|
| Display XL | 64px | 40px | Main hero title |
| Display L | 48px | 34px | Editorial section title |
| H1 | 40px | 30px | Page title |
| H2 | 32px | 26px | Section title |
| H3 | 24px | 21px | Card groups and content sections |
| Body L | 18px | 17px | Introductions and descriptions |
| Body | 16px | 16px | Default content |
| Small | 14px | 14px | Metadata and helper text |
| Caption | 12px | 12px | Labels, badges, technical metadata |

### Typography rules

- Use serif fonts only for editorial headings and selected product titles.
- Use sans-serif fonts for navigation, forms, buttons, prices, and operational content.
- Prices must use tabular numerals.
- **Dates must always render in full, localized, unabbreviated form** (e.g. `September 26, 2026`, `26 septembre 2026`, `٢٦ سبتمبر ٢٠٢٦`) anywhere a date appears next to a status — order status, review date, stock restock date, promotion end date. Never use relative time ("2 days ago") or numeric-only dates (`09/26/26`) as the primary display; a relative or short form may appear only as a secondary `Small`/`Caption` hint next to the full date.
- Keep body line length between 55 and 75 characters where possible.

---

## 4.3 Spacing and sizing

Use the Tailwind spacing scale consistently.

Recommended layout spacing:

- Page horizontal padding: `px-4 sm:px-6 lg:px-8`.
- Maximum content width: `max-w-7xl`.
- Editorial content width: `max-w-4xl`.
- Section vertical space: `py-12 md:py-16 lg:py-24`.
- Card gap: `gap-4 md:gap-6`.
- Form spacing: `space-y-5`.
- Primary button height: `h-11` or `h-12`.
- Icon button size: `size-10`.

---

## 4.4 Borders, radius, and shadows

- Product cards: `rounded-xl`.
- Hero media: `rounded-2xl` where the layout is contained.
- Buttons: `rounded-md` or `rounded-lg`.
- Inputs: `rounded-md`.
- Drawers and dialogs: `rounded-t-2xl` on mobile and `rounded-xl` on desktop.
- Use subtle shadows only for floating surfaces.
- Prefer borders and background contrast over heavy shadows.

Suggested shadows:

```css
--shadow-soft: 0 10px 30px rgb(46 35 29 / 0.08);
--shadow-floating: 0 18px 50px rgb(46 35 29 / 0.14);
```

---

## 4.5 Photography direction

Product images should be:

- high resolution;
- naturally lit;
- neutral or artisan-workshop backgrounds;
- consistent by category;
- available in portrait and square ratios;
- free from excessive text overlays.

Required image ratios:

- Product grid: `4:5`.
- Product gallery main image: `4:5` or `1:1`.
- Editorial banner: `16:9` or `3:2`.
- Artisan portrait: `4:5`.
- Category tile: `3:4`.

Use video or short looping visual content sparingly on artisan and story pages.

In dark mode, product and editorial images are shown exactly as authored — no dimming, tinting, or overlay filters. Only the surrounding card/page surface, borders, and text adapt to the dark palette; keep a `1px` `border` token around white-background product shots so they don't bleed into the darker page background.

---

## 5. Global Layout

## 5.1 Top announcement bar

A thin, optional bar above the header.

Content examples:

- Free delivery above a threshold;
- Handmade in Algeria;
- International shipping available;
- Secure payment.

### shadcn/ui usage

Use a simple semantic `<div>` rather than a Card.

### Behavior

- Can rotate between two or three messages.
- Hidden when empty.
- Text centered.
- Maximum height: `36px`.

---

## 5.2 Main header

### Desktop layout

Three-zone header:

1. Left: categories/menu trigger and optional location selector.
2. Center: AISHA logo.
3. Right: search, language, account, wishlist, cart.

### Mobile layout

- Menu button;
- centered logo;
- search or cart icon;
- remaining actions inside a Sheet.

### Components

Use:

- `Button` with `variant="ghost"`;
- `NavigationMenu`;
- `Sheet`;
- `DropdownMenu` for short, non-searchable menus (e.g. account quick menu);
- `Combobox` (Command + Popover) for **language selection**, since the list carries language name, native name and flag, and benefits from type-ahead;
- `Command` for advanced search;
- `ThemeToggle` (custom, built on `Button` + `DropdownMenu`) for switching between Light, Dark, and System — see 4.1 "Theme behavior";
- `Badge` for cart count and any status indicator in the header (e.g. "New" ribbon on a promoted nav item);
- `Separator`.

### Language switcher

The `LanguageSwitcher` component uses a `Combobox`:

- trigger shows current language flag/code;
- popover lists Arabic, French, English, Spanish with native name + English name, filterable by typing;
- selecting an option updates the `[locale]` route segment and persists the choice;
- fully operable via keyboard and screen reader (`aria-expanded`, `aria-controls`, `role="listbox"`).

### Sticky behavior

The header should become sticky after scroll:

```text
position: sticky
inset-block-start: 0
z-index: 50
backdrop blur
subtle bottom border
```

---

## 5.3 Mega menu

Desktop category navigation should open a large, editorial mega menu.

Suggested columns:

- Categories;
- Collections;
- Regions;
- Featured artisans;
- Promotional image or seasonal campaign.

Example categories:

- Pottery;
- Rugs and textiles;
- Traditional jewelry;
- Copper and metalwork;
- Leather goods;
- Halfa and basketry;
- Decorative art;
- Musical instruments;
- Souvenirs;
- Eco-friendly products.

Use `NavigationMenu`, but style its viewport as a wide custom surface.

---

## 5.4 Footer

The footer should include:

- brand statement;
- newsletter form;
- shopping links;
- artisan links;
- customer service;
- company information;
- language and currency selectors (`Combobox` for language, matching the header switcher; `Select` for currency, since the list is short);
- payment method icons;
- social links;
- legal links.

Use a dark brown background with warm neutral text.

---

## 6. Core Routes

```text
/[locale]
/[locale]/products
/[locale]/products/[slug]
/[locale]/categories/[slug]
/[locale]/collections/[slug]
/[locale]/artisans
/[locale]/artisans/[slug]
/[locale]/regions/[slug]
/[locale]/search
/[locale]/wishlist
/[locale]/cart
/[locale]/checkout
/[locale]/checkout/success
/[locale]/login
/[locale]/register
/[locale]/forgot-password
/[locale]/account
/[locale]/account/profile
/[locale]/account/addresses
/[locale]/account/orders
/[locale]/account/orders/[id]
/[locale]/account/wishlist
/[locale]/about
/[locale]/our-mission
/[locale]/fair-trade
/[locale]/sustainability
/[locale]/contact
```

---

## 7. Homepage

The homepage must combine commerce and storytelling.

## 7.1 Hero section

### Desktop

Use a split layout:

- left: headline, short brand statement, two CTAs;
- right: large lifestyle image or product composition.

### Mobile

Stack image first or text first depending on campaign content.

### Content structure

- Eyebrow: `Authentic Algerian craftsmanship`;
- H1: emotional and memorable;
- short explanation;
- Primary CTA: `Shop the collection`;
- Secondary CTA: `Meet the artisans`.

### Components

- `Button`;
- custom responsive image component;
- optional `Badge` (e.g. "Limited edition" — always paired with the campaign's full end date when time-bound).

Avoid using a carousel as the default hero. A static campaign is clearer and faster.

---

## 7.2 Featured categories

Display six to eight categories as image tiles.

### Desktop

- Four-column grid;
- first tile may span two columns and two rows.

### Mobile

- horizontally scrollable two-column-feel list or two-column grid.

Each tile contains:

- image;
- category name;
- optional short descriptor;
- arrow icon.

Use custom image cards with minimal overlays.

---

## 7.3 Featured products

Use a standard e-commerce product grid.

Recommended title:

`Selected treasures`

Include:

- section title;
- short introduction;
- View all link;
- four products desktop;
- two products tablet;
- two or one product mobile depending on screen width.

---

## 7.4 Artisan story banner

Large editorial split section:

- artisan portrait or workshop image;
- artisan name;
- region;
- short story;
- CTA to artisan profile.

This section differentiates AISHA from ordinary marketplaces.

---

## 7.5 Values section

Show four trust pillars:

- Verified Algerian origin;
- Quality checked;
- Fair artisan partnership;
- Secure international delivery.

Use icons with short descriptions. Avoid oversized cards.

---

## 7.6 Regional discovery

A visual section allowing users to explore products by region or cultural area.

Examples:

- Sahara;
- Kabylia;
- Aurès;
- Tlemcen;
- Ghardaïa;
- Algiers;
- Hoggar.

Use image-led cards rather than a complex map in the MVP.

---

## 7.7 Editorial collection

Create curated campaign blocks such as:

- Gifts from Algeria;
- Sustainable home;
- Handmade tableware;
- Traditional jewelry;
- Ramadan collection;
- Wedding gifts.

The layout may alternate between image-left and image-right.

---

## 7.8 Testimonials and press

Use:

- `Carousel` only when there are at least three testimonials;
- `Avatar`;
- `Card` with restrained borders;
- logos in grayscale.

---

## 7.9 Newsletter

A simple full-width section with:

- title;
- one-sentence benefit;
- email input;
- submit button;
- privacy note.

Use `Input`, `Button`, and `Form`.

---

## 8. Product Listing Page

## 8.1 Page header

Include:

- breadcrumb;
- page title;
- category description;
- product count;
- optional campaign image.

Use `Breadcrumb` and semantic heading elements.

---

## 8.2 Filter system

### Desktop

Use a left sidebar with collapsible filter groups.

### Mobile

Use a full-height `Sheet` opened by a `Filters` button.

### Filters

- Category;
- Product type;
- Artisan;
- Region;
- Price;
- Material;
- Color;
- Availability;
- Ready to ship;
- Made to order;
- Eco-friendly;
- Fair trade;
- Rating.

### Components

- `Accordion`;
- `Checkbox`;
- `Slider`;
- `RadioGroup`;
- `Combobox` for the **Artisan** and **Region** filters specifically, since both lists can grow large and benefit from type-ahead search; keep `Accordion` + `Checkbox` for short, fixed-length filters (Category, Availability, Material, Color);
- `ScrollArea`;
- `Sheet`;
- `Badge` for active filters — every active filter chip is a `Badge` with a remove (×) affordance;
- `Button` for clear all.

Filter state should appear in the URL.

---

## 8.3 Sort control

Options:

- Featured;
- Newest;
- Price: low to high;
- Price: high to low;
- Best selling;
- Highest rated.

Use `Select` (short, fixed list — a `Combobox` is not needed here).

---

## 8.4 Product grid

### Grid columns

```text
Mobile: 2
Tablet: 3
Desktop: 3 or 4 depending on sidebar
Wide desktop: 4
```

Use CSS Grid and stable card heights.

---

## 8.5 Empty state

Use a centered empty state with:

- simple illustration or icon;
- clear explanation;
- Clear filters button;
- Continue shopping link.

Use `Button` and optional `Card`.

---

## 9. Product Card

The Product Card is the most repeated component and must remain simple.

## 9.1 Structure

1. Product image;
2. badges;
3. wishlist button;
4. category or region;
5. product name;
6. artisan name;
7. price and previous price;
8. optional rating;
9. optional color swatches.

## 9.2 Required states

Every state below is communicated with a `Badge` (never color alone), and any state carrying a date (sale end, restock, back-order) shows the **full date**, not a relative label:

- Default;
- Hover;
- Out of stock — `Badge variant="secondary"`, optional full restock date;
- Sale — `Badge` showing discount, plus full sale-end date when time-limited;
- New — `Badge`, optionally with the full date the product was listed;
- Made to order — `Badge`, with full estimated-ready date;
- Low stock — `Badge variant="destructive"` or warning tone;
- Wishlist active;
- Loading skeleton.

## 9.3 Interaction

On desktop hover:

- subtle image zoom;
- optionally swap to second image;
- reveal quick add action;
- maintain keyboard accessibility.

On mobile:

- no hover dependency;
- wishlist remains visible;
- tap card to open product.

## 9.4 Suggested implementation

Use a custom component composed from:

- `Card` or semantic `<article>`;
- `Badge`;
- `Button`;
- `Tooltip`;
- `Skeleton`.

Do not use a heavy bordered Card. The product image should visually define the card.

---

## 10. Product Detail Page

## 10.1 Desktop layout

Two-column layout:

- left: product gallery occupying approximately 60%;
- right: sticky purchase panel occupying approximately 40%.

## 10.2 Mobile layout

- image carousel;
- product title and price;
- variants;
- purchase actions;
- details accordions;
- sticky bottom add-to-cart bar.

---

## 10.3 Product gallery

Use:

- thumbnails on desktop;
- touch carousel on mobile;
- image zoom dialog;
- video thumbnail where available.

Components:

- `Carousel`;
- `Dialog`;
- `Button`;
- custom image zoom.

---

## 10.4 Product information panel

Required content order:

1. Breadcrumb;
2. region or collection label;
3. product title;
4. rating and review count;
5. price;
6. installment or tax information if applicable;
7. short product summary;
8. variant selectors;
9. quantity selector;
10. Add to cart;
11. Buy now;
12. delivery estimate (shown as a full date, e.g. "Arrives by October 4, 2026");
13. stock state — as a `Badge` (In stock / Low stock / Out of stock / Made to order);
14. trust reassurance;
15. wishlist and share actions.

---

## 10.5 Variant selectors

Use:

- `RadioGroup` for size and finish;
- custom swatches for color;
- `Select` when a short list of options exists;
- `Combobox` instead of `Select` when a variant list is long (e.g. a rug size list spanning many dimensions, or a broad finish/material catalog);
- inline validation when a required option is missing.

Selected options must be visually strong and accessible.

---

## 10.6 Purchase buttons

Primary action:

`Add to cart`

Secondary action:

`Buy now`

Rules:

- full width;
- minimum height `48px`;
- visible loading state;
- disabled state with explanation;
- success feedback through `Sonner` toast or cart drawer update.

---

## 10.7 Product details

Use **`Tabs`** on desktop (where horizontal space allows a persistent panel switch) with sections:

- Story;
- Materials;
- Dimensions;
- Care instructions;
- Production method;
- Origin and certification;
- Shipping and returns.

On mobile, render the same sections as `Accordion` instead, since a tab strip with seven labels does not fit comfortably on small screens. The first tab/section is open or selected by default.

---

## 10.8 Artisan profile preview

A dedicated card below the main product information:

- artisan image;
- artisan or workshop name;
- location;
- short bio;
- verified badge (`Badge`, e.g. "Verified artisan", with the full verification date on hover/tooltip);
- View profile CTA.

Use `Avatar`, `Badge`, `Button`, and a custom bordered section.

---

## 10.9 Made-to-order products

For custom products, replace the normal immediate purchase flow with:

- customization options;
- production time (stated as a full estimated completion date once options are chosen);
- personalization text;
- file upload where required;
- artisan message field;
- price estimate;
- request confirmation.

Use:

- `Textarea`;
- `Input`;
- `Form`;
- `Calendar` if date selection is relevant;
- `Dialog` for guidance;
- `Alert` for production constraints.

---

## 10.10 Recommendations

Show:

- More from this artisan;
- Similar products;
- Complete the collection;
- Recently viewed.

Use separate horizontal product sections instead of a large mixed carousel.

---

## 11. Artisan Directory

## 11.1 Directory page

The page should allow discovery by:

- craft;
- region — use a `Combobox` for region selection when the region list spans all Algerian wilayas;
- name;
- verified status (`Badge` filter chip).

Each artisan card contains:

- portrait or workshop image;
- artisan name;
- workshop name;
- location;
- primary craft;
- number of products;
- short statement.

Use a three-column desktop grid and one-column mobile layout.

---

## 11.2 Artisan profile page

Sections, presented as **`Tabs`** on desktop (Biography / Craft & Techniques / Products / Gallery / Certifications / Reviews) and as a stacked scroll with anchor links on mobile:

1. Cover image;
2. portrait and identity;
3. biography;
4. craft and techniques;
5. region;
6. short video;
7. products by the artisan;
8. workshop gallery;
9. certifications or awards (each with a `Badge` and full award date);
10. customer reviews (each review timestamped with a full date).

The design should feel editorial and human, not like a seller dashboard.

---

## 12. Search Experience

## 12.1 Search trigger

Desktop:

- search icon opens a wide Command dialog.

Mobile:

- search opens a dedicated full-screen Sheet.

## 12.2 Search content

Before typing:

- recent searches;
- popular searches;
- featured categories.

After typing:

- product matches;
- artisan matches;
- category matches;
- region matches;
- View all results action.

Use:

- `Command`;
- `Dialog`;
- `Sheet`;
- `ScrollArea`;
- `Avatar`;
- `Badge`.

---

## 13. Cart

## 13.1 Cart drawer

Adding a product should open or update a right-side `Sheet` on desktop and a bottom/full-height Sheet on mobile.

Content:

- cart title and count;
- progress to free shipping;
- line items;
- quantity controls;
- remove action;
- subtotal;
- checkout button;
- View cart link;
- trust note.

Use `Sheet`, `Separator`, `Button`, and `Progress`.

---

## 13.2 Full cart page

Desktop:

- left: cart items;
- right: sticky summary.

Mobile:

- stacked list;
- summary below;
- sticky checkout button.

Summary includes:

- subtotal;
- discounts;
- estimated shipping (shown as a full delivery-window date);
- taxes note;
- total;
- coupon input.

Use `Input`, `Button`, `Separator`, and `Alert`.

---

## 14. Checkout

Checkout must feel minimal, secure, and distraction-free.

## 14.1 Checkout header

Use a simplified header with:

- logo;
- secure checkout label;
- return to cart link.

Do not show the full navigation menu.

## 14.2 Desktop layout

- left: checkout steps;
- right: sticky order summary.

## 14.3 Mobile layout

- collapsible order summary at top;
- checkout form below;
- sticky final action when appropriate.

## 14.4 Steps

1. Contact information;
2. Delivery address;
3. Shipping method;
4. Payment method;
5. Review and confirmation.

Use a single-page progressive form for the MVP unless the payment provider requires separate screens.

## 14.5 Address fields — country and province

The **Country** field and the **State / Province / Wilaya** field both use a `Combobox`:

- `Country`: full, searchable list of all countries, showing flag + localized country name; typing filters instantly.
- `State / Province / Wilaya`: a dependent `Combobox` scoped to the selected country — for Algeria it lists the 58 wilayas; for other countries it lists that country's states/provinces/regions. The field is disabled until a country is chosen, and repopulates whenever the country changes.
- Both comboboxes remain fully keyboard-operable and screen-reader labeled, and mirror correctly in RTL.

## 14.6 Components

- `Form`;
- `Input`;
- `Combobox` for country and state/province/wilaya fields (see 14.5);
- `Select` for short, fixed lists only (e.g. shipping method, payment method);
- `RadioGroup`;
- `Checkbox`;
- `Alert`;
- `Separator`;
- `Button`;
- `Skeleton`;
- `Sonner`.

## 14.7 Validation

Use React Hook Form and Zod.

Validation messages must:

- appear below the field;
- explain how to fix the problem;
- be translated;
- not rely only on color.

---

## 15. Authentication

## 15.1 Login and registration layout

Use a split screen on desktop:

- form side;
- lifestyle or artisan image side.

On mobile, show only the form with a small brand image or logo.

## 15.2 Form content

Login:

- email;
- password;
- forgot password;
- submit;
- social login if supported;
- register link.

Registration:

- name;
- email;
- password;
- confirmation;
- terms agreement;
- submit;
- optional country field, using the same `Combobox` pattern as checkout (14.5) if collected at signup.

Use `Card`, `Form`, `Input`, `Checkbox`, `Button`, and `Separator`.

---

## 16. Customer Account

## 16.1 Desktop layout

Use a left navigation sidebar and right content panel. Where a section itself has sub-views (e.g. Orders: Active / Completed / Cancelled), render those sub-views as **`Tabs`** within the content panel.

## 16.2 Mobile layout

Use a top Select or list-based navigation; sub-views within a section still use `Tabs` when the content panel is wide enough, otherwise a `Combobox`/`Select` switch above the list.

## 16.3 Account sections

- Overview;
- Profile;
- Addresses (country/province fields use the `Combobox` pattern from 14.5);
- Orders;
- Wishlist;
- Notifications;
- Security;
- Logout.

Use:

- `Tabs` where the number of sections, or sub-sections within one section, is small;
- otherwise custom sidebar navigation;
- `Table` for desktop order history, with an `Order status` column rendered as `Badge` and an `Order date` column always shown in full date format;
- Cards for mobile order history, each showing the same `Badge` status and full date.

---

## 17. Order Details and Tracking

Include:

- order number;
- date — shown in full (e.g. "Placed on September 26, 2026"), never abbreviated;
- payment status — `Badge`;
- fulfillment status — `Badge`;
- delivery address;
- items;
- summary;
- tracking timeline — each timeline entry paired with its own `Badge` status and full date/time;
- support action;
- invoice download.

Use:

- `Badge` for every status (payment, fulfillment, and each timeline step);
- `Separator`;
- custom vertical timeline;
- `Button`;
- `Alert` where action is needed.

Status color rules (all statuses use `Badge`, and every badge is shown next to its full date, e.g. "Shipped — September 24, 2026"):

- Processing: neutral;
- Confirmed: blue or primary;
- Shipped: olive;
- Delivered: green;
- Cancelled: destructive;
- Refunded: muted purple or neutral.

Do not rely only on color; always include text and icon alongside the `Badge`.

---

## 18. Content and Editorial Pages

Pages such as About, Mission, Fair Trade, and Sustainability should use an editorial layout.

Recommended structure:

- large title;
- lead paragraph;
- full-width image;
- narrow text columns;
- quotes or statistics;
- alternating image/text sections;
- CTA to shop or meet artisans.

Use custom layouts with limited Cards. Preserve a premium magazine feel.

---

## 19. shadcn/ui Component Map

| Use case | shadcn/ui component |
|---|---|
| Navigation | `NavigationMenu`, `Sheet`, `DropdownMenu` |
| Language selection | `Combobox` (searchable, flag + native name) |
| Country / province / wilaya selection | `Combobox` (dependent pair: province list scoped to chosen country) |
| Search | `Command`, `Dialog`, `Sheet` |
| Product filters | `Accordion`, `Checkbox`, `Slider`, `RadioGroup`, `Select`, `Combobox` (for Artisan/Region filters) |
| Forms | `Form`, `Input`, `Textarea`, `Checkbox`, `Select`, `Combobox` |
| Product options | `RadioGroup`, custom swatches, `Combobox` (long variant lists) |
| Cart | `Sheet`, `Separator`, `Progress` |
| Checkout | `Form`, `RadioGroup`, `Combobox` (country/province), `Alert`, `Button` |
| Product details | `Tabs` (desktop), `Accordion` (mobile), `Dialog` |
| Grouped content panels | `Tabs` (product details desktop, artisan profile desktop, account sub-sections) |
| Product gallery | `Carousel`, `Dialog` |
| Feedback | `Sonner`, `Alert`, `AlertDialog` |
| Loading | `Skeleton` |
| Status (order, stock, sale, verification) | `Badge` — always paired with a full date when time-relevant |
| User identity | `Avatar` |
| Tooltips | `Tooltip` |
| Pagination | `Pagination` |
| Data display | `Table` for desktop, custom cards for mobile |
| Date selection | `Calendar`, `Popover` |
| Theme switching | `ThemeToggle` (custom, `Button` + `DropdownMenu`: Light / Dark / System) |

### Component rule

Do not use shadcn/ui components without visual adaptation. Apply AISHA spacing, colors, typography, and radius tokens to all components.

### Select vs. Combobox rule

Use `Select` only for short, fixed-length option lists (typically fewer than ~10 items, e.g. sort order, currency, shipping method). Use `Combobox` whenever the list is long, searchable, or grows over time — language, country, state/province/wilaya, artisan filter, region filter, and long variant lists.

---

## 20. Reusable Application Components

Create the following project-level components:

```text
SiteHeader
AnnouncementBar
DesktopNavigation
MobileNavigation
MegaMenu
LanguageSwitcher        (Combobox-based)
ThemeToggle             (Light / Dark / System)
CurrencySwitcher
CountrySelect            (Combobox)
ProvinceSelect           (Combobox, dependent on CountrySelect)
GlobalSearch
SiteFooter
Breadcrumbs
SectionHeading
ProductCard
ProductCardSkeleton
ProductGrid
ProductGallery
ProductPrice
ProductBadge             (Badge, with optional full-date slot)
WishlistButton
QuantitySelector
VariantSelector
AddToCartButton
BuyNowButton
CartDrawer
CartLineItem
OrderSummary
FilterSidebar
MobileFilterSheet
SortSelect
ArtisanCard
ArtisanPreview
ArtisanTabs               (Tabs — profile page sections)
ProductDetailTabs         (Tabs — desktop details panel)
RegionCard
CategoryCard
EditorialBanner
TrustPillars
ReviewSummary
ReviewCard
EmptyState
ErrorState
NewsletterForm
AccountSidebar
AccountSectionTabs        (Tabs — sub-views within an account section)
OrderStatusBadge          (Badge + full date)
OrderTimeline
LocalizedLink
ResponsiveImage
```

---

## 21. Responsive Design

### Breakpoints

Use Tailwind defaults unless a project requirement demands otherwise.

```text
sm: 640px
md: 768px
lg: 1024px
xl: 1280px
2xl: 1536px
```

### Mobile-first rules

- Design all interactions for touch first.
- Minimum touch target: `44 × 44px`.
- Never place essential actions only on hover.
- Use Sheets for filters, menus, and cart.
- Avoid horizontal overflow in page layouts; long modal and drawer content may scroll horizontally when its content cannot safely collapse.
- For artisan private-file rows, replace wide tables with stacked, labelled cards below `sm`; keep filename, type, date, and actions readable without page-level horizontal scrolling.
- Keep primary checkout and add-to-cart actions visible.
- Use one-column forms on mobile.
- On mobile, collapse `Tabs` panels (Product details, Artisan profile) into `Accordion` where more than 3–4 tabs would otherwise crowd the screen.

### Desktop rules

- Use larger product grids;
- sticky sidebars when useful;
- hover previews;
- mega menu navigation;
- wider editorial compositions;
- prefer `Tabs` over `Accordion` for grouped content panels where horizontal space allows.

---

## 22. RTL and Localization

### Direction

Set direction at the locale layout level:

```tsx
<html lang={locale} dir={locale === "ar" ? "rtl" : "ltr"}>
```

### RTL requirements

- Use logical CSS properties where possible.
- Avoid hardcoded `left` and `right` for layout.
- Mirror navigation and directional icons.
- Keep product media order natural.
- Make breadcrumb separators direction-aware.
- Test `Sheet`, `DropdownMenu`, `Combobox`, `Tabs`, and `Carousel` in RTL.
- `Combobox` popovers must open aligned to the correct logical edge in RTL, and the search input must accept Arabic input with correct caret direction.
- `Tabs` order must follow reading direction (rightmost tab first in RTL).

### Localized content

Translate:

- navigation;
- buttons;
- validation messages;
- status labels (every `Badge` label);
- filters;
- metadata labels;
- checkout content, including country and province/wilaya names shown in the `Combobox`;
- empty states;
- system messages;
- date formatting — full dates must use the correct locale calendar/format conventions (e.g. Gregorian, day-month-year order per locale) while remaining unabbreviated.

Product and artisan content should support localized fields where available.

---

## 23. Accessibility

The interface should target WCAG 2.2 AA.

Requirements:

- semantic HTML;
- keyboard navigation;
- visible focus indicators;
- sufficient color contrast;
- alt text for every meaningful image;
- labels for every form field, including `Combobox` triggers and `Tabs` triggers;
- accessible names for icon buttons;
- no information communicated only by color — every status uses a `Badge` with text, and every `Combobox`/`Tabs` state is announced;
- reduced-motion support;
- correct dialog focus trapping (including the `Combobox` popover);
- correct heading order;
- live announcements for cart updates, form errors, and `Tabs` panel changes.

Use shadcn/ui accessibility behavior as a base, but verify every custom composition.

---

## 24. Motion and Interaction

Motion should be subtle and functional.

Recommended transitions:

- button and card hover: `150–200ms`;
- drawer, dialog, and Combobox popover: `200–300ms`;
- image zoom: `300ms`;
- accordion and tab-panel switch: default Radix animation;
- page content reveal: limited and optional.

Avoid large parallax effects, excessive scroll animation, and animation that delays shopping actions.

Respect:

```css
@media (prefers-reduced-motion: reduce)
```

---

## 25. Loading, Error, and Empty States

Every data-driven screen must include:

- initial loading state;
- partial loading state;
- empty state;
- recoverable error state;
- unavailable product state.

### Loading

Use `Skeleton` matching the final layout dimensions, including skeleton rows for `Combobox` option lists and `Tabs` panels while data loads.

### Error

Explain what failed and show a retry action.

### Empty

Provide a useful next step, such as clear filters or continue shopping. A `Combobox` with no matching results shows an inline "No results" row rather than an empty popover.

### Product unavailable

Allow:

- back-in-stock notification;
- similar product discovery;
- wishlist saving.

---

## 26. Trust and Commerce Signals

Use trust messages contextually rather than in one overloaded section.

Examples:

- `Quality checked by AISHA` on the product page;
- `Verified artisan` on artisan profiles, as a `Badge` with full verification date on hover;
- `Made in Algeria` on product metadata;
- `Secure payment` in cart and checkout;
- delivery estimates near purchase actions, shown as a full date;
- return policy near checkout;
- origin certification in product details.

Use `Badge`, `Tooltip`, and small icon-text rows.

---

## 27. Performance Guidelines

- Use Next.js Image for product and editorial images.
- Provide responsive sizes.
- Lazy-load below-the-fold images.
- Preload the hero image only.
- Avoid shipping large client components; `Combobox` and `Tabs` interactivity should be isolated Client Components.
- Use Server Components by default.
- Load filters and interactive controls as focused Client Components.
- Use image placeholders to prevent layout shift.
- Keep icon libraries tree-shakeable.
- Avoid autoplay video on mobile.
- For the country/province `Combobox`, load the country list eagerly (small dataset) and lazy-load each country's province/wilaya list only after that country is selected.

---

## 28. Recommended Page Wireframes

## 28.1 Homepage

```text
Announcement bar
Header
Hero
Featured categories
Featured products
Artisan story
Brand values
Shop by region
Editorial collection
More products
Testimonials / press
Newsletter
Footer
```

## 28.2 Product listing

```text
Header
Breadcrumb
Title + description
Toolbar: product count + filters + sort
Sidebar filters | Product grid
Pagination
Recently viewed
Footer
```

## 28.3 Product detail

```text
Header
Breadcrumb
Product gallery | Purchase panel
Product story and specifications (Tabs desktop / Accordion mobile)
Artisan preview
Shipping and trust information
More from this artisan
Similar products
Recently viewed
Footer
```

## 28.4 Checkout

```text
Minimal checkout header
Contact and address form (Country/Province as Combobox) | Order summary
Shipping method
Payment method
Review
Place order
Legal and security note
```

---

## 29. MVP Design Priorities

The MVP must fully design and implement:

1. Global header and footer, including the `Combobox`-based language switcher;
2. Homepage;
3. Product listing and filters;
4. Product card, with `Badge` status + full-date pattern;
5. Product detail page, with `Tabs`/`Accordion` details panel;
6. Artisan listing and profile;
7. Search;
8. Cart drawer and cart page;
9. Checkout, with `Combobox` country/province fields;
10. Login and registration;
11. Customer orders, with `Badge` status + full date throughout;
12. Multilingual and RTL behavior;
13. Responsive states;
14. Loading, empty, and error states.

The following can be added after the MVP:

- live artisan chat;
- advanced personalization preview;
- interactive regional map;
- augmented reality preview;
- loyalty program;
- gift registry;
- complex editorial animations.

---

## 30. Design Acceptance Criteria

The UI is accepted when:

- it clearly looks like a polished international e-commerce application;
- product photography dominates the browsing experience;
- artisan identity and product origin are visible throughout the journey;
- users can complete checkout without confusion;
- Arabic RTL works correctly on every route, including `Combobox` and `Tabs`;
- all pages are usable on mobile screens from 320px width;
- all forms provide accessible validation;
- filters and search work without visual instability;
- loading and error states match final layouts;
- every status anywhere in the app is a `Badge`, and every status tied to a point in time shows the full, unabbreviated, localized date;
- language, country, and province/wilaya selection all use a searchable `Combobox`, with province scoped correctly to the chosen country;
- grouped content panels (product details, artisan profile, account sub-sections) use `Tabs` on desktop and collapse sensibly on mobile;
- shadcn/ui components follow a consistent AISHA theme;
- the application fully supports **Light**, **Dark**, and **System** modes, with no flash of the wrong theme on load, correct AA contrast in both, and product photography left untouched by the theme;
- the application maintains a premium, warm, trustworthy, and culturally authentic appearance in both light and dark mode.

---

## 31. Implementation Notes for Codex

When generating the UI:

- use Server Components by default;
- add `"use client"` only where interaction requires it — `Combobox` and `Tabs` triggers are client components;
- create shared components instead of duplicating page markup;
- keep mock content in typed data files, including a full country list and a per-country province/wilaya dataset for the `Combobox` fields;
- use strict TypeScript;
- use `cn()` for conditional Tailwind classes;
- use translation keys for all visible labels, including every `Badge` status label and every `Tabs` trigger label;
- use logical spacing utilities compatible with RTL;
- do not place business-critical pricing, stock, or payment rules only in the frontend;
- preserve accessibility attributes provided by Radix and shadcn/ui, especially for `Combobox` (Command + Popover) and `Tabs`;
- do not use placeholder gradients when suitable product imagery is available;
- avoid generic dashboard styling for customer-facing pages;
- keep commerce actions visually stronger than secondary editorial actions;
- centralize date formatting in one utility (`formatFullDate(date, locale)`) so every status/date pairing across the app stays consistent;
- centralize theming in one provider (e.g. a `ThemeProvider` wrapping the locale layout) that reads/writes the persisted preference and toggles the `.dark` class, applied before hydration to avoid a flash of the wrong theme; never hardcode a color — always consume the CSS variable tokens so both themes render correctly everywhere, including inside `Combobox` popovers, `Sheet`/`Dialog` overlays, and `Tabs` panels.

### Brand entry and navigation behavior

- On the first browser visit, show a short AISHA splash with the wordmark, an
  animated story line, and a progress bar. Store a completion timestamp locally
  so normal navigation does not replay it; allow the opening splash again only
  after a two-minute cooldown.
- Language switching replaces only the locale segment of the current URL and
  preserves the current page, query string, and hash. The header control uses
  a flag and two-letter abbreviation while its options retain full language
  names.
- The theme menu supports Light, Dark, and System modes, closes with Escape or
  an outside click, and uses the same persisted preference as the pre-hydration
  theme bootstrap.
