"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface InputProps
  extends React.InputHTMLAttributes<HTMLInputElement> {
  showPrefix?: boolean;
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, showPrefix = true, ...props }, ref) => {
    return (
      <div className="relative">
        {showPrefix && (
          <span className="absolute left-3 top-1/2 -translate-y-1/2 text-accent font-mono">
            {">"}
          </span>
        )}
        <input
          type={type}
          className={cn(
            "flex h-10 w-full",
            "bg-input border border-border",
            "cyber-chamfer-sm",
            showPrefix ? "pl-8 pr-4" : "px-4",
            "py-2",
            "font-mono text-sm text-accent",
            "placeholder:text-muted-foreground placeholder:italic",
            "transition-all duration-200",
            "focus:border-accent focus:shadow-[var(--shadow-neon-sm)]",
            "focus:outline-none",
            "disabled:cursor-not-allowed disabled:opacity-50",
            className
          )}
          ref={ref}
          {...props}
        />
      </div>
    );
  }
);

Input.displayName = "Input";

export { Input };
