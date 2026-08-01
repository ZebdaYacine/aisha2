import type { HTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Badge({ className, ...props }: HTMLAttributes<HTMLSpanElement>) {
  return <span className={cn("inline-flex min-h-6 items-center bg-background px-2 text-[0.6875rem] font-medium uppercase tracking-[0.12em] text-foreground", className)} {...props} />;
}
