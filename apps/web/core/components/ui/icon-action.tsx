import type { ButtonHTMLAttributes, ReactNode } from "react";
import Link from "next/link";

import { Button } from "@/core/components/ui/button";
import { cn } from "@/core/lib/utils";

type IconActionProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  label: string;
  icon: ReactNode;
  variant?: "primary" | "secondary" | "outline" | "ghost" | "destructive";
};

/** Compact action used in dense tables. The label remains available to screen readers and as a native hint. */
export function IconAction({ label, icon, className, variant = "outline", ...props }: IconActionProps) {
  return (
    <Button
      {...props}
      type={props.type ?? "button"}
      variant={variant}
      aria-label={label}
      title={label}
      className={cn("size-10 min-h-10 shrink-0 rounded-md p-0", className)}
    >
      {icon}
    </Button>
  );
}

export function IconActionLink({ href, label, icon, className }: { href: string; label: string; icon: ReactNode; className?: string }) {
  return (
    <Link href={href} aria-label={label} title={label} className={cn("inline-flex size-10 shrink-0 items-center justify-center rounded-md border border-current transition-colors hover:bg-foreground hover:text-background", className)}>
      {icon}
    </Link>
  );
}
