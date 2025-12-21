# OAuth 2.0 Integration for tq_video_menu Plugin

## Overview

The `tq_video_menu` plugin is integrated with OAuth 2.0 authentication through the core middleware system. Platform endpoints require OAuth Bearer token authentication.

## Authentication Flow

### Middleware Chain

The OAuth authentication is handled by `core/middleware/security/oauth_token.go`:

1. **OAuthTokenMiddleware** - Validates Bearer tokens from `Authorization: Bearer <token>` header
2. **AuthMiddleware** - Falls back to API key authentication if no Bearer token is present
3. **Resolver Level** - Platform resolvers check for OAuth authentication and required scopes

### Platform Endpoints

All platform endpoints (`/api/platform/v1`) require:

1. **OAuth Bearer Token** in `Authorization` header:
   ```
   Authorization: Bearer <access_token>
   ```

2. **Required Scopes**:
   - **Queries** (`platformPositionsVideoMenu`, `platformPositionsVideoMenuSlider`): `read` scope
   - **Mutations** (`CreatePosition`, `UpdatePosition`, `DeletePosition`, `UpdateMenuGroup`): `write` scope

## Implementation Details

### Resolver Integration

Platform resolvers use core security functions:

```go
import "web100now-clients-platform/core/middleware/security"

// Check OAuth authentication
if err := security.RequireOAuthAuth(ctx); err != nil {
    return nil, errors.New("OAuth authentication required")
}

// Check required scope
if err := security.RequireScope(ctx, "read"); err != nil {
    return nil, fmt.Errorf("insufficient permissions: %w", err)
}

// Get OAuth user data
oauthUser := security.GetOAuthUserData(ctx)
// oauthUser.UserID, oauthUser.ClientID, oauthUser.Scope, etc.
```

### Available Functions

- `security.RequireOAuthAuth(ctx)` - Returns error if OAuth auth is missing
- `security.IsOAuthAuthenticated(ctx)` - Returns bool if OAuth is present
- `security.GetOAuthUserData(ctx)` - Returns `*OAuthUserData` or `nil`
- `security.RequireScope(ctx, scope)` - Returns error if scope is missing

### OAuth User Data

```go
type OAuthUserData struct {
    UserID    string    // Authenticated user ID
    ClientID  string    // OAuth client ID
    Scope     string    // Token scopes (space-separated)
    DeviceID  string    // Device identifier
    JTI       string    // JWT ID (for revocation)
    ExpiresAt time.Time // Token expiration
}
```

## Usage Examples

### Query with OAuth

```graphql
query {
  platformPositionsVideoMenu(tagsFilter: ["breakfast"]) {
    groupName
    items {
      id
      name
    }
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
X-API-KEY: <api_key>  # Still required for endpoint routing
```

### Mutation with OAuth

```graphql
mutation {
  createPosition(input: {
    name: "New Item"
    groupName: "Breakfast"
    # ... other fields
  }) {
    id
    name
  }
}
```

**Headers:**
```
Authorization: Bearer <access_token>
X-API-KEY: <api_key>
```

**Required Scope:** `write`

## Fallback Behavior

If no Bearer token is provided:
- OAuthTokenMiddleware skips validation
- Falls back to API key authentication (AuthMiddleware)
- Platform resolvers will still require OAuth (returns error)

## Security Features

1. **Token Validation**: JWT RS256 signature validation
2. **Blacklist Check**: Revoked tokens are rejected
3. **Scope Validation**: Endpoints check for required scopes
4. **Device Fingerprinting**: Device validation (when enabled)
5. **Token Expiration**: Automatic expiration check

## Error Responses

### Missing OAuth Token
```json
{
  "error": "OAuth authentication required for platform endpoints"
}
```

### Insufficient Scope
```json
{
  "error": "insufficient permissions: write scope required"
}
```

### Invalid/Expired Token
```json
{
  "error": "Unauthorized",
  "error_description": "Token has been revoked"
}
```

## Testing

To test OAuth integration:

1. Obtain access token via OAuth flow:
   ```bash
   curl -X POST https://your-backend.com/oauth/token \
     -d "grant_type=authorization_code" \
     -d "client_id=your_client_id" \
     -d "client_secret=your_secret" \
     -d "code=authorization_code" \
     -d "redirect_uri=your_redirect_uri"
   ```

2. Use token in requests:
   ```bash
   curl -X POST https://your-backend.com/api/platform/v1 \
     -H "Authorization: Bearer <access_token>" \
     -H "X-API-KEY: <api_key>" \
     -H "Content-Type: application/json" \
     -d '{"query": "{ platformPositionsVideoMenu(tagsFilter: [\"breakfast\"]) { groupName } }"}'
   ```
