# Frontend API Reference

This document describes the current HTTP API surface exposed by the backend and is written for frontend integration.

## Base behavior

- **Base URL**: whatever host the backend is running on.
- **Content-Type**: send and expect `application/json` for JSON endpoints.
- **IDs**: all resource IDs are UUID strings.
- **Auth**: protected routes require authentication through either:
  - the `Authorization` cookie set by the login endpoint, or
  - an `Authorization: Bearer <token>` header.
- **Response envelope**:
  - success: `{"data": ...}`
  - error: `{"error": "message"}`

## Authentication notes

The login endpoint sets an `Authorization` cookie.

- **Dev**: `SameSite=Lax`, not secure
- **Prod**: `SameSite=None`, `Secure=true`

For browser clients, use `credentials: "include"` / `withCredentials: true` if you want cookie-based auth.

## Common response codes

- `200 OK` for successful reads and updates
- `201 Created` is supported by the response helper, but not currently used by routes
- `400 Bad Request` for invalid JSON or invalid UUID/path params
- `401 Unauthorized` when auth is missing or invalid
- `500 Internal Server Error` for unexpected backend failures

> Error handling is not fully uniform across handlers yet, so some authorization/domain failures currently surface as `500` instead of `401` or `403`.

## Endpoint summary

| Method | Route | Auth | Purpose |
| --- | --- | --- | --- |
| GET | `/` | No | Simple welcome response |
| POST | `/user/create` | No | Create a user |
| POST | `/user/login` | No | Authenticate user and issue token |
| GET | `/user/me` | Yes | Return current user |
| GET | `/workspace` | Yes | List current user workspaces |
| POST | `/workspace` | Yes | Create a workspace |
| GET | `/workspace/:workspace_id` | Yes | Fetch workspace details |
| PATCH | `/workspace/:workspace_id` | Yes | Rename workspace |
| DELETE | `/workspace/:workspace_id` | Yes | Delete workspace |
| GET | `/workspace/:workspace_id/members` | Yes | List workspace members |
| PATCH | `/workspace/:workspace_id/members/:user_id` | Yes | Update a member role |
| DELETE | `/workspace/:workspace_id/members/:user_id` | Yes | Remove a member |
| GET | `/workspace/:workspace_id/integrations` | Yes | List integrations for a workspace |
| GET | `/workspace/:workspace_id/integrations/github/connect` | Yes | GitHub connect stub |
| GET | `/workspace/:workspace_id/integrations/github/callback` | Yes | GitHub callback stub |
| DELETE | `/workspace/:workspace_id/integrations/:integration_id` | Yes | Delete integration stub |

---

## `GET /`

Simple health/welcome route.

**Response**

```json
"Welcome to resourcify"
```

---

## `POST /user/create`

Create a new user.

### Request body

```json
{
  "name": "Aman",
  "email": "aman@example.com",
  "password": "secret"
}
```

### Field details

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `name` | string | Yes | Display name |
| `email` | string | Yes | Must be unique |
| `password` | string | Yes | Stored as bcrypt hash |

### Success response

```json
{
  "data": {
    "id": "2f4f4a0f-4d5a-4e1d-8a0f-7f0d4f01a2c1"
  }
}
```

### Errors

- `400` if JSON binding fails
- `500` if the email already exists or user creation fails

---

## `POST /user/login`

Authenticate a user and return the session token.

### Request body

```json
{
  "email": "aman@example.com",
  "password": "secret"
}
```

### Field details

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `email` | string | Yes | User email |
| `password` | string | Yes | Password or master password |

### Success response

```json
{
  "data": {
    "user": {
      "token": "jwt-token-here",
      "id": "2f4f4a0f-4d5a-4e1d-8a0f-7f0d4f01a2c1",
      "name": "Aman",
      "email": "aman@example.com",
      "role": "user",
      "profilePicture": ""
    }
  }
}
```

### Side effects

- Sets the `Authorization` cookie with the JWT token.
- The token also comes back in the JSON payload.

### Errors

- `400` if JSON binding fails
- `500` if credentials are invalid or token generation fails

---

## `GET /user/me`

Returns the current authenticated user.

### Auth

Requires auth via cookie or bearer token.

### Success response

```json
{
  "data": {
    "id": "2f4f4a0f-4d5a-4e1d-8a0f-7f0d4f01a2c1",
    "name": "Aman",
    "email": "aman@example.com",
    "profilePicture": "",
    "role": "user"
  }
}
```

---

## Workspace APIs

All workspace routes require auth.

### Workspace object

```json
{
  "id": "workspace-uuid",
  "name": "Acme",
  "slug": "acme",
  "billing_plan": "free"
}
```

### Workspace list item

```json
{
  "id": "workspace-uuid",
  "name": "Acme",
  "slug": "acme",
  "billing_plan": "free",
  "member_count": 3,
  "project_count": 2
}
```

---

## `GET /workspace`

List all workspaces the current user belongs to.

### Success response

```json
{
  "data": {
    "workspaces": [
      {
        "id": "workspace-uuid",
        "name": "Acme",
        "slug": "acme",
        "billing_plan": "free",
        "member_count": 3,
        "project_count": 2
      }
    ]
  }
}
```

---

## `POST /workspace`

Create a workspace and add the current user as owner.

### Request body

```json
{
  "name": "Acme"
}
```

### Field details

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `name` | string | Yes* | Used to generate the workspace slug |

> `name` is logically required for the frontend, but the current handler does not enforce a non-empty value before reaching the service.

### Success response

```json
{
  "data": {
    "id": "workspace-uuid",
    "name": "Acme",
    "slug": "acme",
    "billing_plan": "free",
    "member_count": 0,
    "project_count": 0
  }
}
```

---

## `GET /workspace/:workspace_id`

Fetch a workspace with members and projects.

### Path params

| Param | Type | Required | Notes |
| --- | --- | --- | --- |
| `workspace_id` | UUID | Yes | Workspace ID |

### Success response

```json
{
  "data": {
    "id": "workspace-uuid",
    "name": "Acme",
    "slug": "acme",
    "billing_plan": "free",
    "members": [
      {
        "id": "user-uuid",
        "name": "Aman",
        "email": "aman@example.com",
        "profilePicture": "",
        "role": "owner"
      }
    ],
    "projects": [
      {
        "id": "project-uuid",
        "repository_id": "repo-uuid",
        "workspace_id": "workspace-uuid",
        "name": "Marketing site",
        "framework": "nextjs",
        "build_command": "npm run build",
        "output_dir": "dist",
        "root_dir": "/"
      }
    ]
  }
}
```

### Notes

- `members` entries include user info plus role.
- `projects` entries are read-only summaries from the current backend.

---

## `PATCH /workspace/:workspace_id`

Rename a workspace.

### Request body

```json
{
  "name": "New workspace name"
}
```

### Success response

```json
{
  "data": {
    "id": "workspace-uuid",
    "name": "New workspace name",
    "slug": "new-workspace-name",
    "billing_plan": "free"
  }
}
```

### Permission rules

- owner and admin can update
- other roles receive an auth-style error

---

## `DELETE /workspace/:workspace_id`

Delete a workspace and its workspace members/projects.

### Success response

```json
{
  "data": {
    "message": "workspace deleted successfully"
  }
}
```

### Permission rules

- only owner can delete

---

## Workspace members

### Member object

```json
{
  "id": "user-uuid",
  "name": "Aman",
  "email": "aman@example.com",
  "profilePicture": "",
  "role": "member"
}
```

### Allowed roles

- `owner`
- `admin`
- `member`
- `viewer`

---

## `GET /workspace/:workspace_id/members`

List workspace members.

### Success response

```json
{
  "data": [
    {
      "id": "user-uuid",
      "name": "Aman",
      "email": "aman@example.com",
      "profilePicture": "",
      "role": "owner"
    }
  ]
}
```

---

## `PATCH /workspace/:workspace_id/members/:user_id`

Update a member role.

### Request body

```json
{
  "role": "admin"
}
```

### Path params

| Param | Type | Required | Notes |
| --- | --- | --- | --- |
| `workspace_id` | UUID | Yes | Workspace ID |
| `user_id` | UUID | Yes | Member user ID |

### Role rules

- owner/admin can update roles
- admins cannot modify owners or promote to owner
- users cannot modify their own role

### Success response

```json
{
  "data": {
    "message": "member role updated successfully"
  }
}
```

---

## `DELETE /workspace/:workspace_id/members/:user_id`

Remove a member from a workspace.

### Success response

```json
{
  "data": {
    "message": "member removed successfully"
  }
}
```

### Permission rules

- owner/admin can remove members
- admins cannot remove owners
- users cannot remove themselves

---

## Integrations

### Integration object returned by the workspace integration list

```json
{
  "id": "integration-uuid",
  "workspace_id": "workspace-uuid",
  "provider": "github",
  "provider_id": "github-account-id",
  "expires_at": "2026-04-24T07:36:03Z"
}
```

---

## `GET /workspace/:workspace_id/integrations`

List integrations for a workspace.

### Success response

```json
{
  "data": [
    {
      "id": "integration-uuid",
      "workspace_id": "workspace-uuid",
      "provider": "github",
      "provider_id": "github-account-id",
      "expires_at": "2026-04-24T07:36:03Z"
    }
  ]
}
```

---

## GitHub integration routes

These routes currently exist in the router, but their handlers are empty in the backend.

### `GET /workspace/:workspace_id/integrations/github/connect`

- currently stubbed
- no implemented response body

### `GET /workspace/:workspace_id/integrations/github/callback`

- currently stubbed
- no implemented response body

### `DELETE /workspace/:workspace_id/integrations/:integration_id`

- currently stubbed
- no implemented response body

---

## Frontend integration checklist

1. Send JSON bodies with the exact field names above.
2. Store or forward the login token if you do not rely on cookies.
3. Use `credentials: "include"` when calling authenticated routes from the browser.
4. Treat any `401` as a session/auth failure.
5. Expect `data` wrapping on success responses, except the root route.

