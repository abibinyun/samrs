import { test, expect, describe } from "bun:test";

test("basic test", () => {
  expect(1 + 1).toBe(2);
});

describe("complaint module", () => {
  test("placeholder", () => {
    expect(true).toBe(true);
  });
});
