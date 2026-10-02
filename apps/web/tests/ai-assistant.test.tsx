import { fireEvent, render, screen } from "@testing-library/react";

import { AIAssistant } from "@/features/home/components/ai-assistant";
import { storeCopy } from "@/core/lib/store-copy";
import type { Product } from "@/features/catalogue/types";

const products: Product[] = Array.from({ length: 5 }, (_, index) => ({
  slug: `piece-${index + 1}`,
  name: { en: `Piece ${index + 1}`, fr: `Pièce ${index + 1}`, ar: `قطعة ${index + 1}`, es: `Pieza ${index + 1}` },
  summary: { en: "", fr: "", ar: "", es: "" },
  story: { en: "", fr: "", ar: "", es: "" },
  artisanSlug: "artisan",
  artisanName: "Atelier AISHA",
  categorySlug: "pottery",
  region: { en: "Algiers", fr: "Alger", ar: "الجزائر", es: "Argel" },
  materials: { en: "Clay", fr: "Argile", ar: "طين", es: "Arcilla" },
  method: { en: "Handmade", fr: "Fait main", ar: "مصنوع يدوياً", es: "Hecho a mano" },
  priceMinor: 1000,
  currency: "EUR",
  availability: "in_stock",
  images: ["/images/aisha/A1.jpg"],
  rating: 0,
  reviewCount: 0,
}));

describe("AIAssistant", () => {
  it("opens inline and answers a product question without leaving the page", () => {
    render(<AIAssistant locale="en" copy={storeCopy("en")} products={products} />);

    const openButton = screen.getByRole("button", { name: "Open AISHA assistant" });
    expect(openButton.querySelector(".assistant-bot-wave")).toBeInTheDocument();
    fireEvent.click(openButton);
    expect(screen.getByRole("region", { name: "AISHA assistant" })).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Show me products" }));

    expect(screen.getByText(storeCopy("en").assistantProductsReply)).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Browse products/i })).toHaveAttribute("href", "/en/products");
  });

  it("shows five collection images for the best-pieces prompt", () => {
    render(<AIAssistant locale="en" copy={storeCopy("en")} products={products} />);

    fireEvent.click(screen.getByRole("button", { name: "Open AISHA assistant" }));
    fireEvent.click(screen.getByRole("button", { name: "Show me the five best pieces" }));

    expect(screen.getByText(storeCopy("en").assistantBestFiveReply)).toBeInTheDocument();
    expect(screen.getAllByRole("link").filter((link) => link.getAttribute("href")?.startsWith("/en/products/piece-")).length).toBe(5);
  });

  it("submits free-form questions and exposes an accessible close action", () => {
    render(<AIAssistant locale="fr" copy={storeCopy("fr")} products={products} />);

    fireEvent.click(screen.getByRole("button", { name: "Ouvrir l’assistant AISHA" }));
    const input = screen.getByRole("textbox", { name: /Parlez-moi/i });
    fireEvent.change(input, { target: { value: "Comment fonctionne la livraison ?" } });
    fireEvent.submit(input.closest("form")!);

    expect(screen.getByText(storeCopy("fr").assistantDeliveryReply)).toBeInTheDocument();
    fireEvent.click(screen.getAllByRole("button", { name: "Fermer l’assistant AISHA" }).at(-1)!);
    expect(screen.queryByRole("region", { name: "Assistant AISHA" })).not.toBeInTheDocument();
  });

  it("explains that unsupported questions are still under development", () => {
    const copy = storeCopy("en");
    render(<AIAssistant locale="en" copy={copy} products={products} />);

    fireEvent.click(screen.getByRole("button", { name: copy.assistantOpen }));
    const input = screen.getByRole("textbox", { name: copy.assistantPlaceholder });
    fireEvent.change(input, { target: { value: "Can you compare these pieces for me?" } });
    fireEvent.submit(input.closest("form")!);

    expect(screen.getByText(copy.assistantUnderDevelopment)).toBeInTheDocument();
  });
});
