import { fireEvent, render, screen, waitFor } from "@testing-library/react";

import { WarehouseOperations } from "@/features/admin/components/warehouse-operations";

function jsonResponse(body: unknown) {
  return { ok: true, status: 200, json: async () => body } as Response;
}

describe("warehouse operations", () => {
  afterEach(() => {
    delete (globalThis as { fetch?: typeof fetch }).fetch;
  });

  it("loads receptions and opens inspection details", async () => {
    globalThis.fetch = jest.fn().mockResolvedValue(
      jsonResponse({
        items: [
          {
            id: "reception-1",
            productId: "product-1",
            productCode: "AISHA-ABC1234567",
            productName: "Copper bowl",
            artisanName: "Artisan",
            workshopName: "Workshop",
            receivedQuantity: 5,
            referenceKey: "parcel-1",
            status: "RECEIVED_PENDING_INSPECTION",
            evidence: [],
          },
        ],
        page: 1,
        pageSize: 10,
        total: 1,
      }),
    ) as unknown as typeof fetch;

    render(<WarehouseOperations />);
    await waitFor(() => expect(screen.getByText("Copper bowl")).toBeInTheDocument());
    fireEvent.click(screen.getByRole("button", { name: "Details" }));
    expect(screen.getByRole("dialog")).toBeInTheDocument();
    expect(screen.getByText("Inspect batch")).toBeInTheDocument();
    expect(screen.getByText(/must total 5/)).toBeInTheDocument();
  });

  it("finds validated products by artisan phone and workshop", async () => {
    const fetchMock = jest
      .fn()
      .mockResolvedValueOnce(jsonResponse({ items: [], page: 1, pageSize: 10, total: 0 }))
      .mockResolvedValueOnce(
        jsonResponse({
          items: [
            {
              productId: "product-2",
              productCode: "AISHA-ABC1234567",
              productName: "Copper bowl",
              productStatus: "APPROVED",
              artisanName: "Amina",
              artisanPhone: "0550123456",
              workshopId: "workshop-1",
              workshopName: "Amina Atelier",
              priceMinor: 2500,
              currency: "EUR",
              availableQuantity: 0,
            },
          ],
        }),
      );
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<WarehouseOperations />);
    fireEvent.change(screen.getByLabelText("Artisan phone"), { target: { value: "0550123456" } });
    fireEvent.click(screen.getByRole("button", { name: "Find products" }));
    await waitFor(() => expect(screen.getByRole("combobox", { name: "Workshop" })).not.toBeDisabled());
    fireEvent.focus(screen.getByRole("combobox", { name: "Workshop" }));
    fireEvent.click(screen.getByText("Amina Atelier"));
    fireEvent.focus(screen.getByRole("combobox", { name: "Validated product" }));
    fireEvent.click(screen.getByText(/AISHA-ABC1234567 · Copper bowl/));
    expect(screen.getByText((_, element) => element?.tagName === "P" && element.textContent?.includes("0550123456") === true)).toBeInTheDocument();
    expect(screen.getByRole("combobox", { name: "Validated product" })).toHaveValue("AISHA-ABC1234567 · Copper bowl · 0 available");
    expect(screen.getByLabelText("Supplier or artisan")).toHaveValue("Amina");
  });

  it("opens the inspection step after recording a reception", async () => {
    const reception = {
      id: "reception-2",
      productId: "product-2",
      productCode: "AISHA-ABC1234567",
      productName: "Copper bowl",
      artisanName: "Amina",
      workshopName: "Amina Atelier",
      receivedQuantity: 5,
      referenceKey: "parcel-2",
      status: "RECEIVED_PENDING_INSPECTION",
      evidence: [],
    };
    const validatedProduct = {
      productId: "product-2",
      productCode: "AISHA-ABC1234567",
      productName: "Copper bowl",
      productStatus: "APPROVED",
      artisanName: "Amina",
      artisanPhone: "0550123456",
      workshopId: "workshop-2",
      workshopName: "Amina Atelier",
      priceMinor: 2500,
      currency: "EUR",
      availableQuantity: 0,
    };
    const fetchMock = jest
      .fn()
      .mockResolvedValueOnce(jsonResponse({ items: [], page: 1, pageSize: 10, total: 0 }))
      .mockResolvedValueOnce(jsonResponse({ items: [validatedProduct], page: 1, pageSize: 10, total: 1 }))
      .mockResolvedValueOnce(jsonResponse(reception))
      .mockResolvedValueOnce(jsonResponse({ items: [reception], page: 1, pageSize: 10, total: 1 }));
    globalThis.fetch = fetchMock as unknown as typeof fetch;

    render(<WarehouseOperations />);
    fireEvent.change(screen.getByLabelText("Artisan phone"), { target: { value: "0550123456" } });
    fireEvent.click(screen.getByRole("button", { name: "Find products" }));
    await waitFor(() => expect(screen.getByRole("combobox", { name: "Workshop" })).not.toBeDisabled());
    fireEvent.focus(screen.getByRole("combobox", { name: "Workshop" }));
    fireEvent.click(screen.getByText("Amina Atelier"));
    fireEvent.focus(screen.getByRole("combobox", { name: "Validated product" }));
    fireEvent.click(screen.getByText(/AISHA-ABC1234567 · Copper bowl/));
    fireEvent.change(screen.getByLabelText("Reception reference"), { target: { value: "parcel-2" } });
    fireEvent.change(screen.getByLabelText("Received quantity"), { target: { value: "5" } });

    const submit = screen.getByRole("button", { name: "Record reception" });
    fireEvent.click(submit);

    await waitFor(() => expect(screen.getByRole("dialog")).toBeInTheDocument());
    expect(screen.getByText("Inspect batch")).toBeInTheDocument();
    expect(fetchMock).toHaveBeenCalledWith("/api/warehouse/receptions", expect.objectContaining({ method: "POST" }));
  });
});
