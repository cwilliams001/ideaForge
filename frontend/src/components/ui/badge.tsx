"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface BadgeProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "secondary" | "outline" | "homelab" | "coding" | "personal" | "learning" | "creative";
}

const Badge = React.forwardRef<HTMLDivElement, BadgeProps>(
  ({ className, variant = "default", ...props }, ref) => {
    const variants = {
      default: "bg-accent/20 text-accent border-accent/50",
      secondary: "bg-accent-secondary/20 text-accent-secondary border-accent-secondary/50",
      outline: "bg-transparent text-muted-foreground border-border",
      // Category-specific variants
      homelab: "bg-accent-tertiary/20 text-accent-tertiary border-accent-tertiary/50",
      coding: "bg-accent/20 text-accent border-accent/50",
      personal: "bg-accent-secondary/20 text-accent-secondary border-accent-secondary/50",
      learning: "bg-yellow-500/20 text-yellow-400 border-yellow-500/50",
      creative: "bg-purple-500/20 text-purple-400 border-purple-500/50",
    };

    return (
      <div
        ref={ref}
        className={cn(
          "inline-flex items-center",
          "px-2.5 py-0.5",
          "text-xs font-mono uppercase tracking-wider",
          "border",
          "cyber-chamfer-sm",
          variants[variant],
          className
        )}
        {...props}
      />
    );
  }
);

Badge.displayName = "Badge";

export { Badge };
