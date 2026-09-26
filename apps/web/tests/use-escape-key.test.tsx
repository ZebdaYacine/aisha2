import { fireEvent, render, screen } from "@testing-library/react";

import { useEscapeKey } from "@/core/hooks/use-escape-key";

function EscapeDismissible({ onClose }: { onClose: () => void }) {
  useEscapeKey(onClose);
  return <div role="dialog">Modal content</div>;
}

describe("useEscapeKey", () => {
  it("dismisses the active modal on Escape", () => {
    const onClose = jest.fn();
    render(<EscapeDismissible onClose={onClose} />);

    fireEvent.keyDown(screen.getByRole("dialog"), { key: "Escape" });

    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
