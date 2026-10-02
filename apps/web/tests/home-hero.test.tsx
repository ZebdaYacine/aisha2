import { render, screen } from "@testing-library/react";

import { HomeHero } from "@/core/components/editorial/home-hero";
import { dictionary } from "@/core/lib/i18n";

describe("HomeHero", () => {
  it("makes the down arrow scroll to the next landing section", () => {
    const messages = dictionary("en");

    render(<HomeHero locale="en" messages={messages} />);

    expect(screen.getByRole("link", { name: messages.navigation.story })).toHaveAttribute("href", "#story");
  });
});
