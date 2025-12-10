"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface TextareaProps
  extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  showPrefix?: boolean;
}

const Textarea = React.forwardRef<HTMLTextAreaElement, TextareaProps>(
  ({ className, showPrefix = true, ...props }, ref) => {
    return (
      <div className="relative">
        {showPrefix && (
          <span className="absolute left-3 top-3 text-accent font-mono">
            {">"}
          </span>
        )}
        <textarea
          className={cn(
            "flex min-h-[120px] w-full",
            "bg-input border border-border",
            "cyber-chamfer-sm",
            showPrefix ? "pl-8 pr-4" : "px-4",
            "py-3",
            "font-mono text-sm text-accent",
            "placeholder:text-muted-foreground placeholder:italic",
            "transition-all duration-200",
            "focus:border-accent focus:shadow-[var(--shadow-neon-sm)]",
            "focus:outline-none",
            "disabled:cursor-not-allowed disabled:opacity-50",
            "resize-none",
            className
          )}
          ref={ref}
          {...props}
        />
      </div>
    );
  }
);

Textarea.displayName = "Textarea";

export { Textarea };
