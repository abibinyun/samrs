import { test, expect, describe } from "bun:test";
import { API_ENDPOINTS } from "@/constants/endpoints";

describe("COMPLAINTS API endpoints", () => {
  test("BASE endpoint is defined", () => {
    expect(API_ENDPOINTS.COMPLAINTS.BASE).toBeDefined();
    expect(typeof API_ENDPOINTS.COMPLAINTS.BASE).toBe("string");
  });

  test("DETAIL endpoint returns correct path", () => {
    const id = "test-123";
    const result = API_ENDPOINTS.COMPLAINTS.DETAIL(id);
    expect(result).toContain(id);
  });
});
