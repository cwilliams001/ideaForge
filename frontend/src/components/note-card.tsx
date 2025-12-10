"use client";

import * as React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface Link {
  title: string;
  url: string;
  type: string;
  description?: string;
}

interface ProcessedNote {
  id: string;
  original: string;
  title: string;
  category: string;
  markdown: string;
  links: Link[];
  created_at: string;
  synced_at?: string;
}

interface NoteCardProps {
  note: ProcessedNote;
  className?: string;
}

const linkTypeIcons: Record<string, string> = {
  github: "[ GH ]",
  docs: "[ DOCS ]",
  youtube: "[ YT ]",
  article: "[ ART ]",
};

export function NoteCard({ note, className }: NoteCardProps) {
  const formattedDate = new Date(note.created_at).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });

  return (
    <Card variant="terminal" hoverEffect className={cn("pt-8", className)}>
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <CardTitle className="text-lg">{note.title}</CardTitle>
          <Badge variant={note.category as "homelab" | "coding" | "personal" | "learning" | "creative"}>
            {note.category}
          </Badge>
        </div>
        <p className="text-xs text-muted-foreground font-mono mt-1">
          {formattedDate}
          {note.synced_at && (
            <span className="text-accent ml-2">[ synced ]</span>
          )}
        </p>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* Original note */}
        <div className="text-sm text-muted-foreground border-l-2 border-border pl-3 italic">
          &quot;{note.original}&quot;
        </div>

        {/* Markdown preview */}
        <div
          className="prose prose-invert prose-sm max-w-none
            prose-headings:font-[var(--font-heading)] prose-headings:text-accent prose-headings:uppercase prose-headings:tracking-wide
            prose-p:text-foreground
            prose-li:text-foreground prose-li:marker:text-accent
            prose-a:text-accent-tertiary prose-a:no-underline hover:prose-a:underline
            prose-code:text-accent-secondary prose-code:bg-muted prose-code:px-1 prose-code:rounded-none"
          style={{ fontFamily: "var(--font-body)" }}
        >
          {/* Simple markdown rendering - in real app use react-markdown */}
          <div className="whitespace-pre-wrap font-mono text-sm">
            {note.markdown}
          </div>
        </div>

        {/* Links */}
        {note.links.length > 0 && (
          <div className="space-y-2 pt-2 border-t border-border">
            <h4 className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
              resources
            </h4>
            <ul className="space-y-1">
              {note.links.map((link, index) => (
                <li key={index} className="text-sm">
                  <a
                    href={link.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex items-center gap-2 text-accent-tertiary hover:text-accent transition-colors group"
                  >
                    <span className="text-xs text-muted-foreground font-mono">
                      {linkTypeIcons[link.type] || "[ LINK ]"}
                    </span>
                    <span className="group-hover:underline">{link.title}</span>
                  </a>
                </li>
              ))}
            </ul>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
