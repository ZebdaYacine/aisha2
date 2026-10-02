"use client";

import { ArrowUpRight, Bot, Hand, Send, Sparkles, X } from "lucide-react";
import Link from "next/link";
import Image from "next/image";
import { useEffect, useRef, useState } from "react";

import type { Locale } from "@/core/lib/i18n";
import type { StoreCopy } from "@/core/lib/store-copy";
import { localized } from "@/features/catalogue/format";
import type { Product } from "@/features/catalogue/types";

type ChatMessage = {
  id: number;
  role: "assistant" | "user";
  text: string;
  link?: { href: string; label: string };
  products?: Product[];
};

type AssistantChatProps = {
  locale: Locale;
  copy: StoreCopy;
  products: Product[];
};

const initialMessage = (copy: StoreCopy): ChatMessage => ({
  id: 1,
  role: "assistant",
  text: copy.assistantGreeting,
});

function responseFor(query: string, copy: StoreCopy, locale: Locale, products: Product[]): Omit<ChatMessage, "id" | "role"> {
  const value = query.toLocaleLowerCase(locale);
  if (/best|top|five|5|recommend|meilleur|meilleure|cinq|أفضل|خمس|خمسة|mejores|cinco/.test(value)) {
    return { text: copy.assistantBestFiveReply, products: products.slice(0, 5) };
  }
  if (/product|object|collection|produit|objet|produit|منتج|قطع|colecci|producto/.test(value)) {
    return {
      text: copy.assistantProductsReply,
      link: { href: `/${locale}/products`, label: copy.assistantViewProducts },
    };
  }
  if (/artisan|maker|workshop|atelier|حرف|ورشة|artesano|taller/.test(value)) {
    return {
      text: copy.assistantArtisansReply,
      link: { href: `/${locale}/artisans`, label: copy.assistantViewArtisans },
    };
  }
  if (/deliver|shipping|ship|delivery|livraison|توصيل|شحن|entrega|envío/.test(value)) {
    return { text: copy.assistantDeliveryReply };
  }
  if (/story|history|histoire|قصة|قصتنا|historia/.test(value)) {
    return {
      text: copy.brandStoryBody,
      link: { href: `/${locale}/#story`, label: copy.assistantViewStory },
    };
  }
  return { text: copy.assistantUnderDevelopment };
}

export function AIAssistant({ locale, copy, products }: AssistantChatProps) {
  const [open, setOpen] = useState(false);
  const [input, setInput] = useState("");
  const [messages, setMessages] = useState<ChatMessage[]>(() => [initialMessage(copy)]);
  const endRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!open) return;
    inputRef.current?.focus();
    endRef.current?.scrollIntoView?.({ behavior: "smooth", block: "nearest" });
  }, [open, messages.length]);

  useEffect(() => {
    if (!open) return;
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };
    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [open]);

  const ask = (question: string) => {
    const trimmed = question.trim();
    if (!trimmed) return;
    const response = responseFor(trimmed, copy, locale, products);
    setMessages((current) => [
      ...current,
      { id: Date.now(), role: "user", text: trimmed },
      { id: Date.now() + 1, role: "assistant", ...response },
    ]);
    setInput("");
  };

  return (
    <div className="fixed inset-x-4 bottom-4 z-[80] flex justify-end sm:inset-x-auto sm:end-6">
      {open && (
        <section
          id="aisha-assistant-panel"
          aria-label={copy.assistantName}
          className="absolute bottom-16 end-0 flex h-[min(34rem,calc(100dvh-7rem))] w-[min(25rem,calc(100vw-2rem))] flex-col overflow-hidden rounded-2xl border border-border bg-background shadow-2xl"
        >
          <header className="flex items-center justify-between border-b border-border bg-foreground px-4 py-3 text-background">
            <div className="flex min-w-0 items-center gap-3">
              <span className="grid size-9 shrink-0 place-items-center rounded-full bg-primary text-primary-foreground">
                <Sparkles aria-hidden="true" size={17} />
              </span>
              <div className="min-w-0">
                <p className="truncate text-sm font-medium">{copy.assistantName}</p>
                <p className="text-xs text-background/65">AISHA · {copy.verified}</p>
              </div>
            </div>
            <button
              type="button"
              onClick={() => setOpen(false)}
              className="grid size-9 shrink-0 place-items-center rounded-full transition hover:bg-background/10"
              aria-label={copy.assistantClose}
            >
              <X aria-hidden="true" size={18} />
            </button>
          </header>

          <div className="min-h-0 flex-1 space-y-3 overflow-y-auto p-4" aria-live="polite">
            {messages.map((message) => (
              <div key={message.id} className={`flex ${message.role === "user" ? "justify-end" : "justify-start"}`}>
                <div
                  className={`max-w-[88%] rounded-2xl px-3.5 py-3 text-sm leading-6 ${message.role === "user" ? "rounded-ee-sm bg-foreground text-background" : "rounded-es-sm bg-muted text-foreground"}`}
                >
                  {message.role === "assistant" && (
                    <div className="mb-1 flex items-center gap-1.5 text-[0.65rem] font-medium uppercase tracking-[0.16em] text-primary">
                      <Bot aria-hidden="true" size={12} /> AISHA
                    </div>
                  )}
                  <p>{message.text}</p>
                  {message.link && (
                    <Link
                      href={message.link.href}
                      onClick={() => setOpen(false)}
                      className="mt-2 inline-flex items-center gap-1 text-xs font-medium text-primary underline underline-offset-4"
                    >
                      {message.link.label}
                      <ArrowUpRight aria-hidden="true" size={13} />
                    </Link>
                  )}
                  {message.products && (
                    <div className="mt-3 grid gap-2">
                      {message.products.map((product) => (
                        <Link
                          key={product.slug}
                          href={`/${locale}/products/${product.slug}`}
                          onClick={() => setOpen(false)}
                          className="group flex min-w-0 items-center gap-2 rounded-lg border border-border/80 bg-background p-1.5 text-foreground transition hover:border-primary"
                        >
                          <span className="relative size-12 shrink-0 overflow-hidden rounded-md bg-muted">
                            <Image
                              src={product.images[0]}
                              alt={localized(product.name, locale)}
                              fill
                              sizes="48px"
                              className="object-cover transition group-hover:scale-105"
                            />
                          </span>
                          <span className="min-w-0 flex-1">
                            <span className="block truncate text-xs font-medium">{localized(product.name, locale)}</span>
                            <span className="block text-[0.68rem] text-muted-foreground">{product.artisanName ?? copy.artisan}</span>
                          </span>
                          <ArrowUpRight aria-hidden="true" size={14} className="me-1 shrink-0 text-primary" />
                        </Link>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            ))}
            <div ref={endRef} />
          </div>

          <div className="border-t border-border p-3">
            <div className="mb-3 flex gap-2 overflow-x-auto pb-1 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
              {[copy.assistantBestFive, copy.assistantProducts, copy.assistantArtisans, copy.assistantDelivery].map((question) => (
                <button
                  key={question}
                  type="button"
                  onClick={() => ask(question)}
                  className="min-h-9 shrink-0 rounded-full border border-border px-3 text-xs transition hover:border-primary hover:text-primary"
                >
                  {question}
                </button>
              ))}
            </div>
            <form
              className="flex items-center gap-2 rounded-xl border border-input bg-background p-1.5 focus-within:border-primary"
              onSubmit={(event) => {
                event.preventDefault();
                ask(input);
              }}
            >
              <input
                ref={inputRef}
                value={input}
                onChange={(event) => setInput(event.target.value)}
                placeholder={copy.assistantPlaceholder}
                aria-label={copy.assistantPlaceholder}
                className="min-w-0 flex-1 bg-transparent px-2 text-sm outline-none placeholder:text-muted-foreground"
              />
              <button
                type="submit"
                disabled={!input.trim()}
                className="grid size-9 shrink-0 place-items-center rounded-lg bg-foreground text-background transition hover:bg-primary disabled:cursor-not-allowed disabled:opacity-40"
                aria-label={copy.assistantSend}
              >
                <Send aria-hidden="true" size={16} />
              </button>
            </form>
          </div>
        </section>
      )}
      <button
        type="button"
        onClick={() => setOpen((current) => !current)}
        className="relative grid size-14 place-items-center rounded-full bg-foreground text-background shadow-lg transition hover:-translate-y-0.5 hover:bg-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
        aria-label={open ? copy.assistantClose : copy.assistantOpen}
        aria-expanded={open}
        aria-controls="aisha-assistant-panel"
      >
        {open ? (
          <X aria-hidden="true" size={21} />
        ) : (
          <span className="relative grid size-8 place-items-center" aria-hidden="true">
            <Bot size={25} strokeWidth={1.7} />
            <Hand className="assistant-bot-wave absolute -end-2 -top-2 origin-bottom" size={16} strokeWidth={1.8} />
          </span>
        )}
        {!open && <span className="absolute -end-0.5 -top-0.5 size-3 rounded-full border-2 border-background bg-primary" aria-hidden="true" />}
      </button>
    </div>
  );
}
