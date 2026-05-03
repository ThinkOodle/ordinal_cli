# Ordinal API Coverage

Last updated: May 3, 2026

This file tracks the Ordinal API v1 resource coverage currently implemented in this repository.

Sources:
- https://docs.tryordinal.com/api/introduction
- https://docs.tryordinal.com/api/openapi.json

Notes:
- Coverage here reflects the current repository state.
- "Implemented endpoints" means there is a typed API method and CLI surface for that operation.
- The `agency-api/*` resource groups are intentionally out of scope for this CLI (they require a company-scoped API key). MCP integration is also out of scope.
- Utility commands such as `auth`, `skill`, and `completion` are not included in the API coverage totals because they do not map to Ordinal API resource groups.
- Post channel JSON bodies support `linkedIn`, `x`, `instagram`, `tikTok`, and `youTubeShorts`.
- Idea channel JSON bodies support `linkedIn`, `x`, `tikTok`, and `youTubeShorts`.
- Channel attachments use `assets` arrays of objects (`{"assetId":"<uuid>"}`), not the old `assetIds` array shape. Instagram asset objects may include `tags`; TikTok and YouTube Shorts each accept exactly one video asset object.

## Summary

- Resource groups covered: `19 / 19` (of the in-scope Core REST API)
- Implemented API operations: `63`

## Covered Resources

| Resource | Implemented endpoints |
| --- | --- |
| `analytics` | `GET /analytics/cpm`, `PUT /analytics/cpm`, `GET /analytics/linkedin/{profileId}/followers`, `GET /analytics/linkedin/{profileId}/posts`, `GET /analytics/x/{profileId}/followers`, `GET /analytics/x/{profileId}/posts` |
| `approval` | `POST /approvals`, `DELETE /approvals/{id}`, `GET /posts/{postId}/approvals` |
| `comment` | `GET /posts/{postId}/comments`, `POST /posts/{postId}/comments`, `DELETE /comments/{commentId}` |
| `engagement` | `GET /posts/{postId}/engagements`, `POST /posts/{postId}/engagements`, `PATCH /engagements/{id}`, `DELETE /engagements/{id}` |
| `upload` | `POST /uploads`, `GET /uploads/{id}` |
| `idea` | `GET /ideas`, `POST /ideas`, `GET /ideas/{id}`, `PATCH /ideas/{id}`, `POST /ideas/{id}/archive`, `POST /ideas/{id}/unarchive`, `POST /ideas/{id}/add-to-calendar` |
| `inline-comment` | `GET /posts/{postId}/inline-comments` |
| `instagram` | `GET /instagram/locations/search` |
| `invite` | `GET /invites`, `POST /invites`, `DELETE /invites/{id}` |
| `label` | `GET /labels`, `POST /labels`, `DELETE /labels/{id}` |
| `linkedin` | `GET /linkedin/profile/{urn}`, `GET /linkedin/{username}/mentions` |
| `linkedin-leads` | `GET /linkedin/leads/{profileId}/posts`, `GET /linkedin/leads/{profileId}/posts/{postId}` |
| `post` | `GET /posts`, `POST /posts`, `GET /posts/{id}`, `PATCH /posts/{id}`, `POST /posts/{id}/archive`, `POST /posts/{id}/unarchive`, `POST /posts/{id}/schedule`, `POST /posts/{id}/unschedule` |
| `profile` | `GET /profiles/engagement`, `GET /profiles/scheduling` |
| `slack-boost` | `GET /posts/{postId}/slack-boosts`, `POST /slack-boosts`, `GET /slack-boosts/{id}`, `PATCH /slack-boosts/{id}`, `DELETE /slack-boosts/{id}` |
| `slack-webhook` | `GET /slack-webhooks` |
| `subscriber` | `GET /posts/{postId}/subscribers`, `POST /subscribers`, `DELETE /subscribers/{id}` |
| `user` | `GET /users` |
| `webhook` | `GET /webhooks`, `POST /webhooks`, `GET /webhooks/{id}`, `PATCH /webhooks/{id}`, `DELETE /webhooks/{id}` |
| `workspace` | `GET /workspace` |

## Missing Resources

- None for the in-scope API surface

## Status

All currently in-scope resource groups are implemented. The root CLI also includes non-resource utility commands such as `auth`, `skill`, and `completion`. Remaining uncovered API areas are the ones intentionally excluded from this CLI's scope: the Agency API (multi-workspace / company-scoped) and the MCP integration.
