# PRD-03 · Design System, Localization & Responsive UI

| | |
| --- | --- |
| **Epic** | 3 — Design System, Localization & Responsive UI |
| **Stories** | US-VIS-001 … US-VIS-003; REQ-I18N-001 |
| **Priority** | High |
| **Status** | ✅ Implemented |
| **Surfaces** | `apps/web/app/[locale]`, `apps/web/core/components`, `apps/web/messages` |
| **API** | Localized public catalogue contracts |

## 1. Summary

Deliver a premium, accessible, mobile-first storefront that communicates Algerian craft heritage while keeping discovery and purchase actions clear.

## 2. Problem

A single-language or generic marketplace interface excludes users and hides the cultural context that makes AISHA valuable. Inconsistent responsive and RTL behavior also creates accessibility and conversion failures.

## 3. Goals

- Support Arabic, French, English, and Spanish with translated visible text.
- Make Arabic RTL, responsive layouts, keyboard access, and focus behavior first-class.
- Establish reusable Tailwind and shadcn/ui primitives for storefront and operational screens.

### Non-goals

- A separate mobile application or decorative animation system.

## 4. Users

- Visitor browsing public content.
- Customer completing a purchase.
- Artisan or operator using authenticated workspaces.
- Arabic, French, English, and Spanish users on mobile or desktop.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Use locale-prefixed routes and translation keys for visible text. | REQ-I18N-001 |
| R2 | Render Arabic with RTL and mirror directional controls. | REQ-I18N-001, US-VIS-001 |
| R3 | Provide reusable layout, commerce, editorial, form, loading, empty, and error components. | US-VIS-001..003 |
| R4 | Cover desktop, tablet, mobile, keyboard, focus, contrast, and reduced-motion behavior. | Frontend accessibility stories |
| R5 | Keep product photography and artisan storytelling central to hierarchy. | US-VIS-002, US-VIS-003 |

## 6. Flow

```text
Detect locale → render translated route → apply dir → fetch localized content → handle loading/empty/error
```

## 7. Technical notes

- Use Server Components by default and Client Components only for interaction.
- Use logical CSS properties and locale-aware dates, numbers, and currencies.
- The backend remains authoritative for permissions, price, stock, and publication.

## 8. Success metrics

- Public pages render correctly in all four locales.
- Arabic RTL defects are found before release.
- Core public journeys pass keyboard and mobile smoke tests.

## 9. Risks & open questions

- Long translations may expose layout defects without dedicated visual regression coverage.
- Final accessibility audit and operational-screen parity remain to be completed.

## Source traceability

`docs/ui_design.md`, `docs/frontend.md`, `REQ-I18N-001`, `US-VIS-001..003`.
