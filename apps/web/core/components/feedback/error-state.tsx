import { AlertTriangle } from "lucide-react";
import type { ReactNode } from "react";

export function ErrorState({ action, body, title }: { action?: ReactNode; body: string; title: string }) {
  return <section role="alert" className="mx-auto flex max-w-xl flex-col items-center py-24 text-center">
    <AlertTriangle aria-hidden="true" size={36} strokeWidth={1.2} />
    <h2 className="mt-6 font-serif text-3xl">{title}</h2>
    <p className="mt-3 text-muted-foreground">{body}</p>
    {action && <div className="mt-7">{action}</div>}
  </section>;
}
