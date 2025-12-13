#!/bin/bash
# Migration script to fix ownership of existing files
# Run this ONCE before upgrading to the new version with non-root containers

set -e

# Source .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
else
    echo "Error: .env file not found"
    echo "Please create a .env file with PUID, PGID, and other required variables"
    exit 1
fi

# Default values
PUID=${PUID:-1000}
PGID=${PGID:-1000}

echo "================================================================"
echo "IdeaForge Ownership Migration Script"
echo "================================================================"
echo "This script will fix ownership of existing files to match your"
echo "configured user (UID:GID = ${PUID}:${PGID})"
echo ""
echo "This is required for Syncthing to properly sync IdeaForge notes."
echo "================================================================"
echo ""

# Fix backend data volume (requires docker volume inspect to find mount point)
echo "1. Fixing backend data volume..."
VOLUME_PATH=$(docker volume inspect ideaforge_backend-data --format '{{ .Mountpoint }}' 2>/dev/null || echo "")
if [ -n "$VOLUME_PATH" ] && [ -d "$VOLUME_PATH" ]; then
    echo "   Found volume at: $VOLUME_PATH"
    sudo chown -R ${PUID}:${PGID} "$VOLUME_PATH"
    echo "   ✓ Fixed ownership: $VOLUME_PATH"
else
    echo "   ! Volume not found (this is OK if you haven't run IdeaForge yet)"
fi

echo ""

# Fix Obsidian vault IdeaForge folder
echo "2. Fixing Obsidian vault folder..."
if [ -n "$OBSIDIAN_VAULT_PATH" ] && [ -d "$OBSIDIAN_VAULT_PATH" ]; then
    IDEAFORGE_PATH="${OBSIDIAN_VAULT_PATH}/${OBSIDIAN_FOLDER:-IdeaForge}"
    if [ -d "$IDEAFORGE_PATH" ]; then
        echo "   Found IdeaForge folder at: $IDEAFORGE_PATH"
        sudo chown -R ${PUID}:${PGID} "$IDEAFORGE_PATH"
        echo "   ✓ Fixed ownership: $IDEAFORGE_PATH"
    else
        echo "   ! IdeaForge folder not found in vault (this is OK if no notes exist yet)"
    fi
else
    echo "   ! OBSIDIAN_VAULT_PATH not set or doesn't exist"
    echo "   Please check your .env file"
fi

echo ""
echo "================================================================"
echo "Migration complete!"
echo "================================================================"
echo ""
echo "Next steps:"
echo "  1. Stop containers:     docker compose down"
echo "  2. Rebuild (no cache):  docker compose build --no-cache"
echo "  3. Start containers:    docker compose up -d"
echo "  4. Verify user:         docker exec ideaforge-backend id"
echo ""
echo "Expected output from step 4: uid=${PUID}(${DOCKER_USER:-ct}) gid=${PGID}(${DOCKER_USER:-ct})"
echo "================================================================"
