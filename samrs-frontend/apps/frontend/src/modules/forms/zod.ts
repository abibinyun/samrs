import { z } from "zod";
import type { FormField } from "./types";

export function buildZodSchema(fields: FormField[]) {
  const shape: Record<string, z.ZodTypeAny> = {};

  for (const f of fields) {
    switch (f.type) {
      case "input":
      case "textarea": {
        shape[f.id] = f.required
          ? z.string().min(1, `${f.label} is required`)
          : z.string().optional();
        break;
      }

      case "date": {
        shape[f.id] = f.required
          ? z.date({ message: `${f.label} is required` })
          : z.date().optional();
        break;
      }

      case "select":
      case "radio": {
        shape[f.id] = f.required
          ? z.string().min(1, `${f.label} is required`)
          : z.string().optional();
        break;
      }

      case "checkbox":
      case "switch": {
        shape[f.id] = z.boolean().optional();
        break;
      }

      default:
        shape[f.id] = z.any().optional();
    }
  }

  return z.object(shape);
}
