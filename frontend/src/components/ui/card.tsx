"use client";

import * as React from "react";
import { cn } from "@/lib/utils";

export interface CardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "terminal" | "holographic";
  hoverEffect?: boolean;
}

const Card = React.forwardRef<HTMLDivElement, CardProps>(
  ({ className, variant = "default", hoverEffect = false, children, ...props }, ref) => {
    const baseStyles = cn(
      "cyber-chamfer",
      "transition-all duration-300"
    );

    const variants = {
      default: cn(
        "bg-card border border-border",
        hoverEffect && "hover:-translate-y-0.5 hover:border-accent hover:shadow-[var(--shadow-neon)]"
      ),
      terminal: cn(
        "bg-background border border-border relative",
        "before:absolute before:top-0 before:left-0 before:right-0 before:h-6",
        "before:bg-muted before:border-b before:border-border",
        hoverEffect && "hover:border-accent hover:shadow-[var(--shadow-neon-sm)]"
      ),
      holographic: cn(
        "bg-muted/30 border border-accent/30",
        "shadow-[var(--shadow-neon-sm)]",
        "backdrop-blur-sm",
        hoverEffect && "hover:border-accent hover:shadow-[var(--shadow-neon)]"
      ),
    };

    return (
      <div
        className={cn(baseStyles, variants[variant], className)}
        ref={ref}
        {...props}
      >
        {variant === "terminal" && (
          <div className="absolute top-1.5 left-3 flex gap-1.5 z-10">
            <div className="w-2.5 h-2.5 rounded-full bg-destructive" />
            <div className="w-2.5 h-2.5 rounded-full bg-yellow-500" />
            <div className="w-2.5 h-2.5 rounded-full bg-accent" />
          </div>
        )}
        {children}
      </div>
    );
  }
);

Card.displayName = "Card";

const CardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("flex flex-col space-y-1.5 p-6", className)}
    {...props}
  />
));
CardHeader.displayName = "CardHeader";

const CardTitle = React.forwardRef<
  HTMLHeadingElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h3
    ref={ref}
    className={cn(
      "font-[var(--font-heading)] text-xl font-semibold uppercase tracking-wide text-accent",
      className
    )}
    style={{ fontFamily: "var(--font-heading)" }}
    {...props}
  />
));
CardTitle.displayName = "CardTitle";

const CardDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p
    ref={ref}
    className={cn("text-sm text-muted-foreground", className)}
    {...props}
  />
));
CardDescription.displayName = "CardDescription";

const CardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div ref={ref} className={cn("p-6 pt-0", className)} {...props} />
));
CardContent.displayName = "CardContent";

const CardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("flex items-center p-6 pt-0", className)}
    {...props}
  />
));
CardFooter.displayName = "CardFooter";

export { Card, CardHeader, CardFooter, CardTitle, CardDescription, CardContent };
