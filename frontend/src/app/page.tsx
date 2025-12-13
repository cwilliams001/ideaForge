"use client";

import * as React from "react";
import { NoteInput } from "@/components/note-input";
import { NoteCard } from "@/components/note-card";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api, ProcessedNote } from "@/lib/api";

export default function Home() {
  const [notes, setNotes] = React.useState<ProcessedNote[]>([]);
  const [isLoading, setIsLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [selectedCategory, setSelectedCategory] = React.useState<string | null>(null);

  const categories = ["homelab", "coding", "personal", "learning", "creative"];

  const loadNotes = React.useCallback(async () => {
    try {
      const response = await api.getNotes(selectedCategory || undefined);
      setNotes(response.notes || []);
    } catch {
      // API not available yet - that's fine for dev
      console.log("API not available");
    }
  }, [selectedCategory]);

  // Load notes on mount and when category changes
  React.useEffect(() => {
    loadNotes();
  }, [loadNotes]);

  const handleSubmit = async (content: string) => {
    setIsLoading(true);
    setError(null);

    try {
      const note = await api.createNote(content);
      setNotes((prev) => [note, ...prev]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create note");
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    await api.deleteNote(id);
    setNotes((prev) => prev.filter((note) => note.id !== id));
  };

  return (
    <main className="min-h-screen">
      {/* Header */}
      <header className="border-b border-border bg-background/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1
              className="text-2xl md:text-3xl font-black uppercase tracking-widest text-accent neon-text cyber-glitch"
              data-text="IDEA FORGE"
              style={{ fontFamily: "var(--font-heading)" }}
            >
              IDEA FORGE
            </h1>
            <div className="flex items-center gap-2 text-xs font-mono text-muted-foreground">
              <span className="hidden sm:inline">system</span>
              <span className="text-accent">online</span>
              <span className="w-2 h-2 rounded-full bg-accent animate-pulse" />
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-4xl mx-auto px-4 py-8 space-y-8">
        {/* Input Section */}
        <section>
          <Card variant="holographic" className="p-6">
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <span className="text-accent font-mono text-sm">$</span>
                <h2
                  className="text-lg font-semibold uppercase tracking-wide"
                  style={{ fontFamily: "var(--font-heading)" }}
                >
                  new_entry
                </h2>
              </div>
              <NoteInput onSubmit={handleSubmit} isLoading={isLoading} />
              {error && (
                <p className="text-sm text-destructive font-mono">
                  [ ERROR ] {error}
                </p>
              )}
            </div>
          </Card>
        </section>

        {/* Category Filter */}
        <section className="flex flex-wrap gap-2">
          <button
            onClick={() => setSelectedCategory(null)}
            className={`transition-all ${
              selectedCategory === null ? "opacity-100" : "opacity-50 hover:opacity-75"
            }`}
          >
            <Badge variant={selectedCategory === null ? "default" : "outline"}>
              all
            </Badge>
          </button>
          {categories.map((cat) => (
            <button
              key={cat}
              onClick={() => setSelectedCategory(cat)}
              className={`transition-all ${
                selectedCategory === cat ? "opacity-100" : "opacity-50 hover:opacity-75"
              }`}
            >
              <Badge
                variant={selectedCategory === cat ? (cat as "homelab" | "coding" | "personal" | "learning" | "creative") : "outline"}
              >
                {cat}
              </Badge>
            </button>
          ))}
        </section>

        {/* Notes List */}
        <section className="space-y-4">
          {notes.length === 0 ? (
            <Card variant="terminal" className="pt-8">
              <CardContent className="text-center py-12">
                <p className="text-muted-foreground font-mono text-sm">
                  <span className="text-accent">&gt;</span> no entries found
                </p>
                <p className="text-muted-foreground font-mono text-xs mt-2">
                  enter an idea above to get started
                </p>
              </CardContent>
            </Card>
          ) : (
            notes.map((note) => (
              <NoteCard
                key={note.id}
                note={note}
                onDelete={handleDelete}
              />
            ))
          )}
        </section>
      </div>

      {/* Footer */}
      <footer className="border-t border-border mt-auto py-6">
        <div className="max-w-4xl mx-auto px-4 text-center">
          <p className="text-xs font-mono text-muted-foreground">
            [ IDEA FORGE v0.1.0 ] // transform notes into action
          </p>
        </div>
      </footer>
    </main>
  );
}
