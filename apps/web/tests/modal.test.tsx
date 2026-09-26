import { fireEvent, render, screen } from "@testing-library/react";
import { useState } from "react";

import { Modal } from "@/core/components/ui/modal";

describe("Modal", () => {
  it("traps focus, closes with Escape, and restores the trigger", () => {
    const onClose = jest.fn();
    const { rerender } = render(<><button>Open filters</button></>);
    const trigger = screen.getByRole("button", { name: "Open filters" });
    trigger.focus();
    rerender(<><button>Open filters</button><Modal label="Filters" closeLabel="Close" onClose={onClose}><button>Apply</button></Modal></>);
    const dialog = screen.getByRole("dialog", { name: "Filters" });
    expect(dialog).toHaveFocus();
    expect(dialog).toHaveClass("overflow-x-auto");
    fireEvent.keyDown(document, { key: "Escape" });
    expect(onClose).toHaveBeenCalledTimes(1);
    rerender(<><button>Open filters</button></>);
    expect(screen.getByRole("button", { name: "Open filters" })).toHaveFocus();
  });

  it("keeps an input focused while its value changes", () => {
    function EditableModal() {
      const [value, setValue] = useState("");
      return <Modal label="Edit" closeLabel="Close" onClose={jest.fn()}><input aria-label="Name" value={value} onChange={(event) => setValue(event.target.value)} /></Modal>;
    }

    render(<EditableModal />);
    const input = screen.getByRole("textbox", { name: "Name" });
    input.focus();
    fireEvent.change(input, { target: { value: "A" } });
    fireEvent.change(input, { target: { value: "AI" } });

    expect(input).toHaveFocus();
    expect(input).toHaveValue("AI");
  });
});
