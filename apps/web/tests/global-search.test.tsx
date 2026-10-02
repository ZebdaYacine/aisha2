import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { GlobalSearch } from "@/core/components/layout/global-search";
import { storeCopy } from "@/core/lib/store-copy";

describe("GlobalSearch", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads catalogue data only after the search is opened", async () => {
    const fetchMock = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        categories: [],
        artisans: [],
        workshops: [],
        products: [{
          slug: "vase",
          name: { en: "Vase", fr: "Vase", ar: "Vase", es: "Vase" },
          summary: { en: "", fr: "", ar: "", es: "" },
          story: { en: "", fr: "", ar: "", es: "" },
          artisanSlug: "artisan",
          categorySlug: "pottery",
          region: { en: "", fr: "", ar: "", es: "" },
          materials: { en: "", fr: "", ar: "", es: "" },
          method: { en: "", fr: "", ar: "", es: "" },
          priceMinor: 1000,
          currency: "DZD",
          availability: "out_of_stock",
          images: [],
          rating: 0,
          reviewCount: 0,
        }],
      }),
    } as Response);
    globalThis.fetch = fetchMock;

    render(<GlobalSearch locale="en" copy={storeCopy("en")} />);
    expect(fetchMock).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Search" }));

    await waitFor(() => expect(screen.getByRole("link", { name: "Vase · Products" })).toBeInTheDocument());
    expect(fetchMock).toHaveBeenCalledWith("/api/catalogue?locale=en");
  });

  it("lists products, artisans, workshops, and categories related to the query", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        products: [{ slug: "bowl", name: { en: "Painted bowl", fr: "Bol peint", ar: "وعاء مرسوم", es: "Cuenco pintado" }, summary: { en: "Handmade clay", fr: "Argile faite main", ar: "طين مصنوع يدوياً", es: "Arcilla artesanal" }, story: { en: "A traditional craft", fr: "Un artisanat traditionnel", ar: "حرفة تقليدية", es: "Artesanía tradicional" }, artisanName: "Amina", workshop: "Amina Studio", region: { en: "Algiers", fr: "Alger", ar: "الجزائر", es: "Argel" }, materials: { en: "Clay", fr: "Argile", ar: "طين", es: "Arcilla" }, method: { en: "Hand painted", fr: "Peint à la main", ar: "مرسوم يدوياً", es: "Pintado a mano" }, categorySlug: "pottery", images: [] }],
        artisans: [{ slug: "amina", name: "Amina", workshop: "Amina Studio", biography: { en: "Ceramic maker", fr: "Céramiste", ar: "صانعة خزف", es: "Ceramista" }, craft: { en: "Pottery", fr: "Poterie", ar: "فخار", es: "Cerámica" }, region: { en: "Algiers", fr: "Alger", ar: "الجزائر", es: "Argel" } }],
        workshops: [{ slug: "amina-studio", name: "Amina Studio", description: "Ceramic workshop", wilaya: "Algiers", location: "Kasbah", craft: "Pottery", artisanId: "amina", artisanName: "Amina", image: "", productCount: 1 }],
        categories: [{ slug: "pottery", name: { en: "Pottery", fr: "Poterie", ar: "فخار", es: "Cerámica" }, description: { en: "Clay craft", fr: "Art de l'argile", ar: "حرفة الطين", es: "Arte de arcilla" }, image: "" }],
      }),
    }) as unknown as typeof fetch;

    const copy = storeCopy("en");
    render(<GlobalSearch locale="en" copy={copy} />);
    fireEvent.click(screen.getByRole("button", { name: copy.search }));
    fireEvent.change(screen.getByRole("textbox", { name: copy.search }), { target: { value: "craft workshop" } });

    await waitFor(() => expect(screen.getByRole("link", { name: /Amina Studio · Workshop/ })).toBeInTheDocument());
    expect(screen.getByRole("link", { name: /Painted bowl · Products/ })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Pottery · Categories/ })).toBeInTheDocument();
  });
});
