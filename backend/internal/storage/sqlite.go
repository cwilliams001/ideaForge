package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kilo40/idea-forge/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// Database handles SQLite operations
type Database struct {
	db *sql.DB
}

// NewDatabase creates a new database connection
func NewDatabase() (*Database, error) {
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/ideaforge.db"
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	database := &Database{db: db}
	if err := database.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// migrate runs database migrations
func (d *Database) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		original TEXT NOT NULL,
		title TEXT NOT NULL,
		category TEXT NOT NULL,
		markdown TEXT NOT NULL,
		links TEXT NOT NULL DEFAULT '[]',
		created_at DATETIME NOT NULL,
		synced_at DATETIME
	);

	CREATE INDEX IF NOT EXISTS idx_notes_category ON notes(category);
	CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at DESC);
	`

	_, err := d.db.Exec(schema)
	return err
}

// CreateNote inserts a new note into the database
func (d *Database) CreateNote(note *models.ProcessedNote) error {
	if note.ID == "" {
		note.ID = "note_" + uuid.New().String()[:8]
	}

	linksJSON, err := json.Marshal(note.Links)
	if err != nil {
		return fmt.Errorf("failed to marshal links: %w", err)
	}

	_, err = d.db.Exec(`
		INSERT INTO notes (id, original, title, category, markdown, links, created_at, synced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, note.ID, note.Original, note.Title, note.Category, note.Markdown, string(linksJSON), note.CreatedAt, note.SyncedAt)

	if err != nil {
		return fmt.Errorf("failed to insert note: %w", err)
	}

	return nil
}

// GetNote retrieves a note by ID
func (d *Database) GetNote(id string) (*models.ProcessedNote, error) {
	var note models.ProcessedNote
	var linksJSON string
	var syncedAt sql.NullTime

	err := d.db.QueryRow(`
		SELECT id, original, title, category, markdown, links, created_at, synced_at
		FROM notes WHERE id = ?
	`, id).Scan(&note.ID, &note.Original, &note.Title, &note.Category, &note.Markdown, &linksJSON, &note.CreatedAt, &syncedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get note: %w", err)
	}

	if err := json.Unmarshal([]byte(linksJSON), &note.Links); err != nil {
		return nil, fmt.Errorf("failed to unmarshal links: %w", err)
	}

	if syncedAt.Valid {
		note.SyncedAt = &syncedAt.Time
	}

	return &note, nil
}

// ListNotes retrieves notes with optional filtering
func (d *Database) ListNotes(category string, limit, offset int) ([]models.ProcessedNote, int, error) {
	var args []interface{}
	whereClause := ""

	if category != "" {
		whereClause = "WHERE category = ?"
		args = append(args, category)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM notes " + whereClause
	if err := d.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count notes: %w", err)
	}

	// Get notes
	query := fmt.Sprintf(`
		SELECT id, original, title, category, markdown, links, created_at, synced_at
		FROM notes %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	args = append(args, limit, offset)
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	var notes []models.ProcessedNote
	for rows.Next() {
		var note models.ProcessedNote
		var linksJSON string
		var syncedAt sql.NullTime

		if err := rows.Scan(&note.ID, &note.Original, &note.Title, &note.Category, &note.Markdown, &linksJSON, &note.CreatedAt, &syncedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan note: %w", err)
		}

		if err := json.Unmarshal([]byte(linksJSON), &note.Links); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal links: %w", err)
		}

		if syncedAt.Valid {
			note.SyncedAt = &syncedAt.Time
		}

		notes = append(notes, note)
	}

	return notes, total, nil
}

// DeleteNote removes a note by ID
func (d *Database) DeleteNote(id string) error {
	result, err := d.db.Exec("DELETE FROM notes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

// UpdateSyncedAt updates the synced_at timestamp for a note
func (d *Database) UpdateSyncedAt(id string, syncedAt time.Time) error {
	_, err := d.db.Exec("UPDATE notes SET synced_at = ? WHERE id = ?", syncedAt, id)
	return err
}

// GetCategoryCounts returns the count of notes per category
func (d *Database) GetCategoryCounts() (map[string]int, error) {
	rows, err := d.db.Query("SELECT category, COUNT(*) FROM notes GROUP BY category")
	if err != nil {
		return nil, fmt.Errorf("failed to query category counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		counts[category] = count
	}

	return counts, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}
