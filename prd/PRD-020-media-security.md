# PRD-20 · File & Media Security

| | |
| --- | --- |
| **Epic** | 20 — File & Media Security |
| **Stories** | US-ART-005; US-PROD-001; US-CUSTOM-001; media security requirements |
| **Priority** | Critical |
| **Status** | 🟡 Partial |
| **Surfaces** | Product media, artisan verification, custom-order attachments, moderation evidence |
| **API** | Feature-specific upload, authorized-read, signed-URL, and delete endpoints |

## 1. Summary

Securely accept, store, publish, and read product media, private artisan documents, custom-order references, and moderation evidence.

## 2. Problem

Uploaded files are attacker-controlled. Trusting filenames, extensions, public buckets, or permanent URLs can expose private documents or execute unsafe content.

## 3. Goals

- Validate content and size before storage.
- Generate safe object keys and keep private buckets private.
- Issue signed URLs only to authorized readers.
- Publish product media only after moderation approval.

### Non-goals

- User-controlled executable uploads or permanent public access to private files.

## 4. Users

- Artisan uploading product or verification media.
- Customer uploading a custom-order reference.
- Moderator or admin reviewing authorized evidence.
- Public visitor reading approved product media.

## 5. Requirements

| # | Requirement | Story |
| --- | --- | --- |
| R1 | Enforce request/file size and MIME/signature validation. | US-ART-002, US-ART-005 |
| R2 | Generate server-side object keys and never use the original filename as a path. | REQ-MEDIA-001 |
| R3 | Store drafts, verification documents, references, and moderation files privately. | US-ART-005, US-CUSTOM-001 |
| R4 | Provide authorized signed reads with expiry and audit where required. | US-ART-005..006 |
| R5 | Publish product media only after moderation/publication approval. | US-PROD-005 |
| R6 | Clean up failed, deleted, or orphaned objects safely. | Security upload rules |

## 6. Flow

```text
Upload request → size/signature validation → generated key → private MinIO object → metadata → scan/read authorization
Product approval → controlled public publication
```

## 7. Technical notes

- MinIO buckets are private by default; PostgreSQL stores metadata.
- Re-encode images or strip metadata when practical and scan when configured.
- Never return storage credentials or permanent private URLs.

## 8. Success metrics

- Invalid/spoofed/oversized uploads rejected.
- Private-file access denials and signed-read audit coverage.
- Orphaned-object cleanup rate.

## 9. Risks & open questions

- Maximum sizes, malware scanner, retention, and legal-document policy are unconfirmed.
- Public media bucket policy must remain aligned with moderation state.

## Source traceability

`REQ-MEDIA-001`, `US-ART-005..006`, `US-PROD-001`, `US-CUSTOM-001`, `docs/security.md`.
