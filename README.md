# Tablq Video Menu Plugin

Plugin for displaying video-based menus, sliders, and random hero videos with tag filtering.

## Structure

The plugin is organized into two separate API endpoint groups:

### `/api/websites/v1` - Website API
- **Location**: `websites/` folder
- **Resolvers**: 
  - `websites/get_menu/` - Menu positions and slider queries
  - `websites/random_video_main/` - Random video queries
- **GraphQL Queries**:
  - `positionsVideoMenu(tagsFilter: [String!]!)` - Fetch menu positions with tag filtering
  - `positionsVideoMenuSlider(tagsFilter: [String!]!, currentID: ID!)` - Fetch menu items for slider
  - `randomVideoMain` - Fetch random videos for landing page
- **Access**: Requires `@requireRole(role: "website")`

### `/api/platform/v1` - Platform API
- **Location**: `platform/menu_modules/` folder
- **Resolvers**:
  - `platform/menu_modules/` - Menu positions and slider queries (platform version)
  - `platform/menu_modules/menu_edit/` - Menu editing functionality (mutations)
- **GraphQL Queries**:
  - `platformPositionsVideoMenu(tagsFilter: [String!]!)` - Fetch menu positions with tag filtering
  - `platformPositionsVideoMenuSlider(tagsFilter: [String!]!, currentID: ID!)` - Fetch menu items for slider
- **GraphQL Mutations** (future):
  - Menu editing mutations will be added in `menu_edit/` module
- **Access**: 
  - Requires `@requireRole(role: "platform")`
  - **OAuth 2.0 Authentication Required**: All platform endpoints require Bearer token authentication
  - **Required Scopes**: 
    - Queries: `read` scope
    - Mutations: `write` scope
- **Authentication**: Uses core OAuth middleware (`core/middleware/security/oauth_token.go`)

## Implementation Details

Both API versions share the same business logic and MongoDB queries, but are separated for:
1. Clear separation of concerns
2. Independent versioning
3. Different access control (website vs platform roles)
4. Potential future customization per API type

## Package Structure

```
tq_video_menu/
├── websites/
│   ├── get_menu/
│   │   ├── main.go           # Website resolvers (PositionsVideoMenu, PositionsVideoMenuSlider)
│   │   ├── pipeline.go       # MongoDB aggregation pipelines
│   │   ├── positions.go      # Position fetching logic
│   │   └── slider.go         # Slider-specific logic
│   └── random_video_main/
│       └── main.go           # Website random video resolver
├── platform/
│   └── menu_modules/
│       ├── main.go           # Platform resolvers (PlatformPositionsVideoMenu, etc.)
│       ├── pipeline.go       # MongoDB aggregation pipelines
│       ├── positions.go      # Position fetching logic
│       ├── slider.go         # Slider-specific logic
│       └── menu_edit/
│           └── main.go       # Menu editing mutations (Create, Update, Delete)
├── schema.graphqls           # GraphQL schema (defines both website and platform queries)
└── plugin.json              # Plugin metadata
```

## Usage

### Website API Example

```graphql
query {
  positionsVideoMenu(tagsFilter: ["breakfast", "lunch"]) {
    groupName
    items {
      id
      name
      videoUrlHevc
      urlPosterPrevVideo
    }
  }
}
```

### Platform API Example

```graphql
query {
  platformPositionsVideoMenu(tagsFilter: ["breakfast", "lunch"]) {
    groupName
    items {
      id
      name
      videoUrlHevc
      urlPosterPrevVideo
    }
  }
}

# Future mutations for menu editing:
# mutation {
#   createPosition(input: PositionInput!) { ... }
#   updatePosition(id: ID!, input: PositionInput!) { ... }
#   deletePosition(id: ID!) { ... }
# }
```

## MongoDB Collections

- `video_menu` - Main menu positions collection
- `plugins` - Plugin configuration (for groupOrder settings)

## OAuth 2.0 Integration

Platform endpoints are integrated with OAuth 2.0 authentication through the core middleware system.

**See [OAUTH_INTEGRATION.md](./OAUTH_INTEGRATION.md) for detailed documentation.**
