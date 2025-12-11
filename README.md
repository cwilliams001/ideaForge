# Idea Forge

A PWA that transforms quick notes/ideas into structured markdown todos with contextual resources, synced to your Obsidian vault.

## Quick Start

1. Copy `.env.example` to `.env` and configure:
   ```bash
   cp .env.example .env
   # Edit .env with your API keys and paths
   ```

2. Get a Tailscale auth key from your [admin console](https://login.tailscale.com/admin/settings/keys) and add it to `.env`

3. Start with Docker Compose:
   ```bash
   docker compose up -d
   ```

4. Access via Tailscale at `https://ts-ideaforge.<your-tailnet>.ts.net`

## Features

- Quick note input via PWA (installable on mobile/desktop)
- AI-powered expansion into structured markdown checklists
- Automatic link discovery (GitHub, docs, tutorials)
- Auto-categorization (homelab, coding, personal, learning, creative)
- Direct sync to Obsidian vault

## Architecture

- **Frontend**: Next.js 15.5.7 (PWA) with Cyberpunk design system
- **Backend**: Go with Gin, SQLite
- **AI**: Claude API (Anthropic)
- **Search**: SearXNG (self-hosted) or Tavily API

## Configuration

See `.env.example` for all configuration options.

## Security

Designed to run on a private network (Tailscale). No authentication is implemented - access is controlled by network access.

---

Built for homelab enthusiasts who want to capture ideas on the go.
