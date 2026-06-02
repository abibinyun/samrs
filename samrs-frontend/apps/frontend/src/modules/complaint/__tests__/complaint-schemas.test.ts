import { test, expect, describe } from "bun:test";
import { complaintSchema, complaintUpdateSchema } from "../schemas";
import type { ComplaintStatus } from "../types";

describe("complaintSchema", () => {
  test("rejects empty asset_id", () => {
    const result = complaintSchema.safeParse({
      asset_id: "",
      title: "Test Title",
      description: "Test Description that is long enough",
    });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("Aset wajib dipilih");
    }
  });

  test("rejects title shorter than 3 chars", () => {
    const result = complaintSchema.safeParse({
      asset_id: "asset-1",
      title: "ab",
      description: "Test Description that is long enough",
    });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("Judul minimal 3 karakter");
    }
  });

  test("rejects description shorter than 10 chars", () => {
    const result = complaintSchema.safeParse({
      asset_id: "asset-1",
      title: "Test Title",
      description: "short",
    });
    expect(result.success).toBe(false);
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("Deskripsi minimal 10 karakter");
    }
  });

  test("accepts valid input", () => {
    const result = complaintSchema.safeParse({
      asset_id: "asset-1",
      title: "Test Title",
      description: "Test Description that is long enough",
    });
    expect(result.success).toBe(true);
  });
});

describe("complaintUpdateSchema", () => {
  test("accepts valid status values", () => {
    const statuses: ComplaintStatus[] = ["open", "in_progress", "done"];
    for (const status of statuses) {
      const result = complaintUpdateSchema.safeParse({ status });
      expect(result.success).toBe(true);
    }
  });

  test("rejects invalid status", () => {
    const result = complaintUpdateSchema.safeParse({ status: "invalid" });
    expect(result.success).toBe(false);
  });

  test("accepts optional assigned_to", () => {
    const result = complaintUpdateSchema.safeParse({
      status: "open",
      assigned_to: null,
    });
    expect(result.success).toBe(true);
  });

  test("accepts optional resolution_note", () => {
    const result = complaintUpdateSchema.safeParse({
      status: "done",
      resolution_note: "Issue resolved",
    });
    expect(result.success).toBe(true);
  });
});
