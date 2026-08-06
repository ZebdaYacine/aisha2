export function FormErrorSummary({ message, id }: { message?: string; id: string }) {
  if (!message) return null;
  return <div id={id} tabIndex={-1} role="alert" className="border-s-2 border-destructive bg-destructive/5 px-4 py-3 text-sm text-destructive">{message}</div>;
}
