# HTTP Status Codes — Practical Reference (Go)

A practical reference for choosing the right HTTP status code, grouped by category, with the matching `net/http` constants and guidance on when to use each.

## 2xx — Success

| Status | Constant | When to use |
| --- | --- | --- |
| 200 | `http.StatusOK` | Default for GET responses. Implicit when you call `w.Write()` without `WriteHeader`. |
| 201 | `http.StatusCreated` | POST that creates a new resource. |
| 202 | `http.StatusAccepted` | Request accepted but processing hasn't completed yet (async jobs). |
| 204 | `http.StatusNoContent` | Success but no body to return. Typical for DELETE and PUT/PATCH. |

## 3xx — Redirection

| Status | Constant | When to use |
| --- | --- | --- |
| 301 | `http.StatusMovedPermanently` | Resource permanently moved to a new URL. |
| 302 | `http.StatusFound` | Temporary redirect. |
| 304 | `http.StatusNotModified` | Client's cached version is still valid (used with ETags). |

## 4xx — Client errors

| Status | Constant | When to use |
| --- | --- | --- |
| 400 | `http.StatusBadRequest` | Malformed request, invalid JSON, missing required fields. |
| 401 | `http.StatusUnauthorized` | Not authenticated — no valid token/session. |
| 403 | `http.StatusForbidden` | Authenticated but not allowed to access this resource. |
| 404 | `http.StatusNotFound` | Resource doesn't exist. |
| 405 | `http.StatusMethodNotAllowed` | Wrong HTTP method (e.g. POST on a GET-only route). |
| 409 | `http.StatusConflict` | Conflict with current state — e.g. duplicate entry, optimistic lock failure. |
| 422 | `http.StatusUnprocessableEntity` | Request is well-formed but semantically invalid — e.g. validation errors. |
| 429 | `http.StatusTooManyRequests` | Rate limit exceeded. |

## 5xx — Server errors

| Status | Constant | When to use |
| --- | --- | --- |
| 500 | `http.StatusInternalServerError` | Unexpected server-side failure. Catch-all for unhandled errors. |
| 501 | `http.StatusNotImplemented` | Endpoint exists but functionality isn't built yet. |
| 502 | `http.StatusBadGateway` | Upstream service returned an invalid response. |
| 503 | `http.StatusServiceUnavailable` | Server is overloaded or down for maintenance. |

## Applied to the controller

| Handler | Success status | Reason |
| --- | --- | --- |
| Create | 201 Created | New resource was made. |
| Update | 204 No Content | Updated but nothing to return. |
| Delete | 204 No Content | Deleted but nothing to return. |
| GetAll | 200 OK | Returning data (implicit). |
| GetById | 200 OK | Returning data (implicit). |
