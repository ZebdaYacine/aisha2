import "@testing-library/jest-dom";

jest.mock("next/navigation", () => ({
  usePathname: () => "/en",
  useRouter: () => ({ replace: jest.fn(), push: jest.fn(), refresh: jest.fn() }),
}));
