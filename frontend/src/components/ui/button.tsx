"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: "default" | "secondary" | "outline" | "ghost" | "glitch" | "destructive";
  size?: "default" | "sm" | "lg" | "icon";
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "default", size = "default", children, ...props }, ref) => {
    const baseStyles = cn(
      "inline-flex items-center justify-center gap-2",
      "font-mono uppercase tracking-wider",
      "transition-all duration-150",
      "disabled:pointer-events-none disabled:opacity-50",
      "cyber-chamfer-sm"
    );

    const variants = {
      default: cn(
        "bg-transparent border-2 border-accent text-accent",
        "hover:bg-accent hover:text-background",
        "hover:shadow-[var(--shadow-neon)]"
      ),
      secondary: cn(
        "bg-transparent border-2 border-accent-secondary text-accent-secondary",
        "hover:bg-accent-secondary hover:text-background",
        "hover:shadow-[var(--shadow-neon-secondary)]"
      ),
      outline: cn(
        "bg-transparent border border-border text-foreground",
        "hover:border-accent hover:text-accent",
        "hover:shadow-[var(--shadow-neon-sm)]"
      ),
      ghost: cn(
        "bg-transparent border-none text-muted-foreground",
        "hover:bg-accent/10 hover:text-accent"
      ),
      glitch: cn(
        "bg-accent text-background font-bold",
        "hover:brightness-110",
        "shadow-[var(--shadow-neon)]"
      ),
      destructive: cn(
        "bg-transparent border-2 border-destructive text-destructive",
        "hover:bg-destructive hover:text-background"
      ),
    };

    const sizes = {
      default: "h-10 px-6 py-2 text-sm",
      sm: "h-8 px-4 py-1 text-xs",
      lg: "h-12 px-8 py-3 text-base",
      icon: "h-10 w-10 p-2",
    };

    return (
      <button
        className={cn(baseStyles, variants[variant], sizes[size], className)}
        ref={ref}
        {...props}
      >
        {children}
      </button>
    );
  }
);

Button.displayName = "Button";

export { Button };
