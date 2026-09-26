# PRD-14 · Reviews & Customer Engagement

| | |
| --- | --- |
| **Epic** | 14 — Reviews & Customer Engagement |
| **Stories** | Review and wishlist capabilities from the architecture roadmap |
| **Priority** | Medium |
| **Status** | ⬜ Not Started |
| **Surfaces** | Product detail, account reviews, artisan/product profiles |
| **API** | Review and rating endpoints to be defined after MVP inclusion is confirmed |

## 1. Summary

Allow eligible customers to share useful product feedback and help future customers evaluate products and artisans without compromising moderation or privacy.

## 2. Problem

Customers lack post-purchase social proof, while unmoderated reviews can contain abuse, private information, or claims that the platform cannot verify.

## 3. Goals

- Confirm whether reviews are in MVP scope.
- If included, restrict reviews to eligible purchasers and preserve ownership.
- Moderate, aggregate, and display ratings safely.

### Non-goals

- Public anonymous reviews, social feeds, or automated sentiment claims.

## 4. Users

- Customer who completed an eligible purchase.
- Visitor reading published review summaries.
- Moderator managing reported or invalid content.
- Artisan viewing feedback about owned products.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Allow a customer to review only an eligible purchased product. | Review capability |
| R2 | Prevent duplicate or cross-account review ownership. | Review capability |
| R3 | Support moderation, reporting, editing, and removal rules without deleting audit history. | Review capability |
| R4 | Aggregate ratings from accepted reviews only. | Review capability |
| R5 | Keep private customer data out of public review content. | Security rules |

## 6. Flow

```text
Delivered order → eligible review prompt → submit → moderate → publish aggregate/detail → report or remove
```

## 7. Technical notes

- Review status and verified-purchase state belong in the backend.
- Public responses must use DTOs and redacted author identity.
- Review media, if later added, follows private upload rules.

## 8. Success metrics

- Eligible-purchase review rate.
- Published-review moderation turnaround.
- Reported-review resolution rate.

## 9. Risks & open questions

- Confirm MVP inclusion, rating scale, edit window, and moderation policy.
- Do not expose ratings or verification badges until the data is authoritative.

## Source traceability

Architecture review feature, `docs/frontend.md`, `docs/jira_project.md`.
