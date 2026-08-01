import { fireEvent, render, screen } from "@testing-library/react";

import { Modal } from "@/components/ui/modal";

describe("Modal", () => {
  it("traps focus, closes with Escape, and restores the trigger", () => {
    const onClose = jest.fn();
    const { rerender } = render(<><button>Open filters</button></>);
    const trigger = screen.getByRole("button", { name: "Open filters" });
    trigger.focus();
    rerender(<><button>Open filters</button><Modal label="Filters" closeLabel="Close" onClose={onClose}><button>Apply</button></Modal></>);
    expect(screen.getByRole("dialog", { name: "Filters" })).toHaveFocus();
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalledTimes(1);
    rerender(<><button>Open filters</button></>);
    expect(screen.getByRole("button", { name: "Open filters" })).toHaveFocus();
  });
});
