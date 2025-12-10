# Idea Forge - Project Plan

> A PWA that transforms quick notes/ideas into structured markdown todos with contextual resources, synced to Obsidian.

## Security Note

**Next.js Version Requirement**: Must use **15.5.7** or **16.0.7** (patched for CVE-2025-66478/CVE-2025-55182 - Critical RCE vulnerability in React Server Components).

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Tailscale Network                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌─────────────┐      ┌─────────────┐      ┌─────────────────┐  │
│  │   Browser   │      │   Go API    │      │  Obsidian Vault │  │
│  │  (PWA/Next) │─────▶│   Server    │─────▶│   (Filesystem)  │  │
│  └─────────────┘      └──────┬──────┘      └─────────────────┘  │
│         │                    │                                   │
│         │                    ▼                                   │
│         │            ┌─────────────┐                            │
│         │            │  LLM API    │                            │
│         │            │  (Claude)   │                            │
│         │            └─────────────┘                            │
│         │                    │                                   │
│         │                    ▼                                   │
│         │            ┌─────────────┐                            │
│         └───────────▶│ Search API  │                            │
│                      │ (SearXNG)   │                            │
│                      └─────────────┘                            │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

---

## Tech Stack

### Frontend (PWA)
- **Framework**: Next.js 15.5.7+ (App Router)
- **Styling**: Tailwind CSS 4.x with Cyberpunk design system
- **UI Components**: Custom components (no shadcn - full cyberpunk aesthetic)
- **State**: React hooks + Context (simple needs)
- **PWA**: next-pwa or manual service worker

### Backend (API)
- **Language**: Go 1.22+
- **Router**: Chi or Gin (Chi recommended for learning - more idiomatic)
- **Database**: SQLite (simple, file-based, easy backup)
- **LLM Client**: Anthropic Go SDK or HTTP client

### External Services
- **LLM**: Claude API (Anthropic) - for note expansion
- **Search**: Self-hosted SearXNG or Tavily API - for finding relevant links

### Infrastructure
- **Hosting**: Your homelab server
- **Access**: Tailscale (already configured)
- **Sync**: Syncthing (already configured for Obsidian)

---

## Project Structure

```
idea-forge/
├── frontend/                 # Next.js PWA
│   ├── app/
│   │   ├── layout.tsx       # Root layout with cyberpunk globals
│   │   ├── page.tsx         # Main input/dashboard
│   │   ├── history/
│   │   │   └── page.tsx     # View past notes
│   │   └── globals.css      # Cyberpunk design tokens
│   ├── components/
│   │   ├── ui/              # Base UI components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── input.tsx
│   │   │   └── ...
│   │   ├── note-input.tsx   # Main input component
│   │   ├── note-card.tsx    # Display processed note
│   │   └── category-badge.tsx
│   ├── lib/
│   │   ├── api.ts           # API client
│   │   └── utils.ts
│   ├── public/
│   │   └── manifest.json    # PWA manifest
│   ├── next.config.js
│   ├── tailwind.config.ts
│   └── package.json
│
├── backend/                  # Go API server
│   ├── cmd/
│   │   └── server/
│   │       └── main.go      # Entry point
│   ├── internal/
│   │   ├── api/
│   │   │   ├── router.go    # HTTP routes
│   │   │   ├── handlers.go  # Request handlers
│   │   │   └── middleware.go
│   │   ├── llm/
│   │   │   └── claude.go    # Claude API client
│   │   ├── search/
│   │   │   └── searxng.go   # Search client
│   │   ├── storage/
│   │   │   ├── sqlite.go    # Database operations
│   │   │   └── obsidian.go  # Write to vault
│   │   └── models/
│   │       └── note.go      # Data structures
│   ├── go.mod
│   └── go.sum
│
├── docker-compose.yml        # Optional: containerized deployment
├── .env.example
├── PLAN.md
└── README.md
```

---

## Data Models

### Note (Input)
```go
type NoteInput struct {
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}
```

### Note (Processed)
```go
type ProcessedNote struct {
    ID          string    `json:"id"`
    Original    string    `json:"original"`
    Title       string    `json:"title"`
    Category    string    `json:"category"`    // homelab, coding, personal, etc.
    Markdown    string    `json:"markdown"`    // Formatted todo list
    Links       []Link    `json:"links"`
    CreatedAt   time.Time `json:"created_at"`
    SyncedAt    *time.Time `json:"synced_at"`  // When written to Obsidian
}

type Link struct {
    Title       string `json:"title"`
    URL         string `json:"url"`
    Type        string `json:"type"`          // github, docs, youtube, article
    Description string `json:"description"`
}
```

### Categories
- `homelab` - Infrastructure, servers, networking
- `coding` - Programming projects, libraries
- `personal` - Life admin, tasks
- `learning` - Courses, tutorials, skills
- `creative` - Ideas, projects, writing

---

## API Endpoints

### POST /api/notes
Create and process a new note.

**Request:**
```json
{
  "content": "setup and install bezel for monitoring homelab"
}
```

**Response:**
```json
{
  "id": "note_abc123",
  "original": "setup and install bezel for monitoring homelab",
  "title": "Set Up Bezel Monitoring",
  "category": "homelab",
  "markdown": "# Set Up Bezel Monitoring\n\n## Tasks\n- [ ] Review Bezel documentation...",
  "links": [
    {
      "title": "Bezel GitHub Repository",
      "url": "https://github.com/bezel-hq/bezel",
      "type": "github",
      "description": "Official source code and installation instructions"
    }
  ],
  "created_at": "2025-12-10T10:30:00Z",
  "synced_at": "2025-12-10T10:30:05Z"
}
```

### GET /api/notes
List all notes with optional filters.

**Query Params:**
- `category` - Filter by category
- `limit` - Number of results (default 50)
- `offset` - Pagination offset

### GET /api/notes/:id
Get a specific note by ID.

### DELETE /api/notes/:id
Delete a note (also removes from Obsidian vault).

### GET /api/categories
List available categories with counts.

---

## LLM Prompt Strategy

### System Prompt
```
You are a productivity assistant that transforms quick notes into structured, actionable markdown todo lists.

Given a brief note or idea, you will:
1. Determine the most appropriate category (homelab, coding, personal, learning, creative)
2. Create a clear, descriptive title
3. Expand the note into a markdown checklist with logical steps
4. Keep steps actionable and specific
5. Add brief context where helpful

Output format (JSON):
{
  "title": "Clear Title Here",
  "category": "category_name",
  "markdown": "# Title\n\n## Tasks\n- [ ] First step\n- [ ] Second step\n..."
}

Keep the markdown concise but comprehensive. Each task should be completable in one sitting.
```

### Example Expansion
**Input:** "setup and install bezel for monitoring homelab"

**Output:**
```markdown
# Set Up Bezel Monitoring

## Prerequisites
- [ ] Ensure Docker is installed and running
- [ ] Verify server meets minimum requirements

## Installation
- [ ] Clone/download Bezel from GitHub
- [ ] Configure environment variables
- [ ] Run initial setup/deployment
- [ ] Access web interface and complete setup wizard

## Configuration
- [ ] Add homelab hosts/services to monitor
- [ ] Configure alerting (Discord/email/etc.)
- [ ] Set up dashboards for key metrics

## Resources
- [Bezel GitHub](link)
- [Official Documentation](link)
- [Setup Tutorial (YouTube)](link)
```

---

## Search Strategy

For finding relevant links, query SearXNG/Tavily with:
1. `"{project_name}" github` - Find official repo
2. `"{project_name}" documentation official` - Find docs
3. `"{project_name}" tutorial setup 2024 2025` - Find recent guides

Parse results and categorize:
- GitHub URLs → `github` type
- docs.*, *.io/docs → `docs` type
- youtube.com → `youtube` type
- Others → `article` type

---

## Obsidian Integration

### Vault Structure
```
Obsidian Vault/
└── IdeaForge/
    ├── homelab/
    │   └── 2025-12-10-set-up-bezel-monitoring.md
    ├── coding/
    ├── personal/
    ├── learning/
    └── creative/
```

### File Naming Convention
`YYYY-MM-DD-slugified-title.md`

### File Template
```markdown
---
created: 2025-12-10T10:30:00Z
category: homelab
source: idea-forge
id: note_abc123
---

# Set Up Bezel Monitoring

## Tasks
- [ ] Review Bezel documentation
- [ ] ...

## Resources
- [Bezel GitHub Repository](https://github.com/...)
- [Official Documentation](https://...)
```

---

## Implementation Phases

### Phase 1: Foundation (MVP)
1. Set up Next.js with cyberpunk design system
2. Create Go API with basic structure
3. Implement note input → LLM expansion → markdown output
4. Basic search integration (3 links per note)
5. Write to Obsidian vault (filesystem)
6. Simple history view

### Phase 2: Polish
1. PWA manifest and service worker
2. Offline viewing of history
3. Category filtering and search
4. Improved link quality/relevance
5. Edit note before saving

### Phase 3: Enhancements
1. Voice input (Web Speech API)
2. Templates per category
3. Bulk operations
4. Note archiving
5. Analytics/stats dashboard

---

## Configuration

### Environment Variables

**Backend (.env):**
```bash
# Server
PORT=8080
HOST=0.0.0.0

# Database
DATABASE_PATH=./data/ideaforge.db

# Obsidian
OBSIDIAN_VAULT_PATH=/path/to/obsidian/vault
OBSIDIAN_FOLDER=IdeaForge

# LLM
ANTHROPIC_API_KEY=sk-ant-...
LLM_MODEL=claude-sonnet-4-20250514

# Search (choose one)
SEARXNG_URL=http://localhost:8888
# OR
TAVILY_API_KEY=tvly-...
```

**Frontend (.env.local):**
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## Security Considerations

1. **Tailscale Only**: API not exposed to public internet
2. **No Auth Required**: Trust Tailscale network authentication
3. **Input Sanitization**: Validate note content length/format
4. **Path Traversal**: Sanitize filenames for Obsidian writes
5. **API Keys**: Store securely, never commit to git

---

## Development Commands

```bash
# Frontend
cd frontend
npm install
npm run dev          # Dev server on :3000

# Backend
cd backend
go mod tidy
go run cmd/server/main.go  # Dev server on :8080

# Both (with air for Go hot reload)
# Terminal 1: cd frontend && npm run dev
# Terminal 2: cd backend && air
```

---

## Questions to Address Later

1. **Search API**: Self-host SearXNG or use Tavily? (SearXNG = free, Tavily = easier)
2. **Link Caching**: Cache search results to reduce API calls for similar topics?
3. **Duplicate Detection**: Warn if similar note already exists?
4. **Mobile Input**: Any special handling for mobile keyboards?

---

## Next Steps

1. [ ] Initialize Next.js frontend with patched version
2. [ ] Set up Tailwind with cyberpunk design tokens
3. [ ] Create base UI components (Button, Card, Input)
4. [ ] Initialize Go module and basic server
5. [ ] Create note input → API → response flow
6. [ ] Integrate Claude API for expansion
7. [ ] Add search for links
8. [ ] Implement Obsidian filesystem write
