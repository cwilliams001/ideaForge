"use client";

import * as React from "react";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface NoteInputProps {
  onSubmit: (content: string) => void;
  isLoading?: boolean;
}

export function NoteInput({ onSubmit, isLoading = false }: NoteInputProps) {
  const [content, setContent] = React.useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (content.trim() && !isLoading) {
      onSubmit(content.trim());
      setContent("");
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter" && (e.metaKey || e.ctrlKey)) {
      handleSubmit(e);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="relative">
        <Textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          onKeyDown={handleKeyDown}
          placeholder="enter your idea or note..."
          disabled={isLoading}
          className={cn(
            "min-h-[140px]",
            isLoading && "opacity-50"
          )}
        />
        <div className="absolute bottom-3 right-3 text-xs text-muted-foreground font-mono">
          {content.length > 0 && `${content.length} chars`}
        </div>
      </div>

      <div className="flex items-center justify-between">
        <p className="text-xs text-muted-foreground font-mono">
          <span className="text-accent">ctrl+enter</span> to submit
        </p>
        <Button
          type="submit"
          variant="glitch"
          disabled={!content.trim() || isLoading}
        >
          {isLoading ? (
            <>
              <span className="cyber-cursor">processing</span>
            </>
          ) : (
            "forge idea"
          )}
        </Button>
      </div>
    </form>
  );
}
