import { fireEvent, render, screen } from "@testing-library/react";

import { ArtisanDirectory } from "@/features/catalogue/components/artisan-directory";
import { storeCopy } from "@/core/lib/store-copy";

const text = (value: string) => ({ en: value, fr: value, ar: value, es: value });
const artisans = [
  { slug: "tala", name: "Amina", workshop: "Atelier Tala", region: text("Algiers"), craft: text("Pottery"), biography: text("Clay"), image: "/images/aisha/A1.jpg", verified: false, productCount: 1 },
  { slug: "noura", name: "Noura", workshop: "Atelier Noura", region: text("Ghardaia"), craft: text("Textiles"), biography: text("Wool"), image: "/images/aisha/M1.jpg", verified: false, productCount: 2 },
];

describe("ArtisanDirectory", () => {
  it("filters public artisans by localized search and region", () => {
    render(<ArtisanDirectory artisans={artisans} locale="en" copy={storeCopy("en")} />);
    expect(screen.getByRole("link", { name: /Amina/ })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Noura/ })).toBeInTheDocument();

    fireEvent.change(screen.getByRole("searchbox"), { target: { value: "wool" } });
    expect(screen.queryByRole("link", { name: /Amina/ })).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: /Noura/ })).toBeInTheDocument();

    fireEvent.change(screen.getByRole("searchbox"), { target: { value: "" } });
    fireEvent.change(screen.getByRole("combobox", { name: "Region" }), { target: { value: "Algiers" } });
    expect(screen.getByRole("link", { name: /Amina/ })).toBeInTheDocument();
    expect(screen.queryByRole("link", { name: /Noura/ })).not.toBeInTheDocument();
  });
});
