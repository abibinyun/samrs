import { test, expect, describe } from "bun:test";
import { complaintStatus } from "../components/ComplaintStatus";
import type { ComplaintStatus } from "../types";

describe("complaintStatus", () => {
  test("returns correct label and class for open status", () => {
    const result = complaintStatus("open");
    expect(result.label).toBe("Open");
    expect(result.className).toContain("blue");
  });

  test("returns correct label and class for in_progress status", () => {
    const result = complaintStatus("in_progress");
    expect(result.label).toBe("In Progress");
    expect(result.className).toContain("amber");
  });

  test("returns correct label and class for done status", () => {
    const result = complaintStatus("done");
    expect(result.label).toBe("Done");
    expect(result.className).toContain("green");
  });

  test("returns fallback for unknown status", () => {
    const result = complaintStatus("unknown" as ComplaintStatus);
    expect(result.label).toBe("unknown");
    expect(result.className).toContain("gray");
  });

  test("returns fallback for null status", () => {
    const result = complaintStatus(null as unknown as ComplaintStatus);
    expect(result.label).toBe("unknown");
  });
});
