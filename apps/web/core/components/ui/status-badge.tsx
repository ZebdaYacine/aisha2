import { Badge } from "@/core/components/ui/badge";
import { cn } from "@/core/lib/utils";
import type { Locale } from "@/core/lib/i18n";
import { statusLabel } from "@/core/lib/common-copy";

const toneByStatus: Record<string, string> = {
  ACTIVE: "bg-success/15 text-success",
  APPROVED: "bg-success/15 text-success",
  VERIFIED: "bg-success/15 text-success",
  INSPECTED: "bg-success/15 text-success",
  DELIVERED: "bg-success/15 text-success",
  PAID: "bg-success/15 text-success",
  PROCESSING: "bg-accent/25 text-foreground",
  PENDING: "bg-accent/25 text-foreground",
  PENDING_REVIEW: "bg-accent/25 text-foreground",
  RECEIVED_PENDING_INSPECTION: "bg-accent/25 text-foreground",
  CHANGES_REQUESTED: "bg-primary/15 text-primary",
  SUSPENDED: "bg-destructive/10 text-destructive",
  REJECTED: "bg-destructive/10 text-destructive",
  DISABLED: "bg-destructive/10 text-destructive",
  CANCELLED: "bg-destructive/10 text-destructive",
  RETURNED: "bg-primary/15 text-primary",
  DAMAGED: "bg-destructive/10 text-destructive",
  QUARANTINED: "bg-primary/15 text-primary",
  ARCHIVED: "bg-muted text-muted-foreground",
};

export function StatusBadge({ status, locale = "en", className }: { status: string; locale?: Locale; className?: string }) {
  return (
    <Badge className={cn("rounded-full border border-transparent px-3 py-1", toneByStatus[status] ?? "bg-muted text-muted-foreground", className)}>
      {statusLabel(status, locale)}
    </Badge>
  );
}
