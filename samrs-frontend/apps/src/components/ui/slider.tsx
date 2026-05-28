import * as React from "react";
import { cn } from "@/lib/utils";

export function Slider({ value, onValueChange, min = 0, max = 100, step = 1, disabled, className }: {
  value: number[];
  onValueChange: (value: number[]) => void;
  min?: number;
  max?: number;
  step?: number;
  disabled?: boolean;
  className?: string;
}) {
  return (
    <input
      type="range"
      min={min}
      max={max}
      step={step}
      value={value[0]}
      onChange={(e) => onValueChange([Number(e.target.value)])}
      disabled={disabled}
      className={cn("w-full", className)}
    />
  );
}
