import * as React from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { CalendarIcon } from "lucide-react";
import { format } from "date-fns";
import { cn } from "@/lib/utils";
import type { FormSchema, FormField } from "../types";

type Props = {
  schema: FormSchema;
  onClose: () => void;
};

export function FormPreview({ schema, onClose }: Props) {
  const zodSchema = React.useMemo(() => {
    const shape: Record<string, z.ZodTypeAny> = {};
    
    schema.fields.forEach((field) => {
      let fieldSchema: z.ZodTypeAny = z.any();
      
      if (field.type === "input") {
        fieldSchema = z.string();
      } else if (field.type === "textarea") {
        fieldSchema = z.string();
      } else if (field.type === "checkbox") {
        fieldSchema = z.boolean();
      } else if (field.type === "switch") {
        fieldSchema = z.boolean();
      } else if (field.type === "date") {
        fieldSchema = z.date();
      }
      
      if (field.required && field.type !== "checkbox" && field.type !== "switch") {
        fieldSchema = (fieldSchema as z.ZodString).min(1, `${field.label} is required`);
      }
      
      shape[field.id] = fieldSchema;
    });
    
    return z.object(shape);
  }, [schema]);

  const form = useForm({
    resolver: zodResolver(zodSchema),
  });

  const onSubmit = (data: any) => {
    console.log("Form submitted:", data);
    alert("Form submitted! Check console for data.");
  };

  const groupedFields: FormField[][] = [];
  let currentRow: FormField[] = [];
  
  schema.fields.forEach((field) => {
    if (currentRow.length === 0) {
      currentRow.push(field);
    } else {
      const currentRowSpan = currentRow.reduce((sum, f) => sum + (f.gridColumn || 1), 0);
      const newSpan = field.gridColumn || 1;
      
      if (currentRowSpan + newSpan <= 3) {
        currentRow.push(field);
      } else {
        groupedFields.push(currentRow);
        currentRow = [field];
      }
    }
  });
  
  if (currentRow.length > 0) {
    groupedFields.push(currentRow);
  }

  return (
    <div className="fixed inset-0 bg-background z-50 overflow-auto">
      <div className="container max-w-4xl mx-auto py-8">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold">Form Preview</h2>
          <Button variant="outline" onClick={onClose}>
            Close Preview
          </Button>
        </div>

        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          {groupedFields.map((row, rowIndex) => (
            <div key={rowIndex} className="grid grid-cols-3 gap-4">
              {row.map((field) => {
                const colSpanClass = field.gridColumn === 3 ? "col-span-3" : field.gridColumn === 2 ? "col-span-2" : "col-span-1";
                
                return (
                  <div key={field.id} className={`${colSpanClass} space-y-2`}>
                    <Label>
                      {field.label}
                      {field.required && <span className="text-destructive ml-1">*</span>}
                    </Label>

                    {field.type === "input" && (
                      <Input
                        {...form.register(field.id)}
                        placeholder={field.placeholder}
                      />
                  )}

                  {field.type === "textarea" && (
                    <Textarea
                      {...form.register(field.id)}
                      placeholder={field.placeholder}
                    />
                  )}

                  {field.type === "select" && (
                    <Select onValueChange={(value) => form.setValue(field.id, value)}>
                      <SelectTrigger>
                        <SelectValue placeholder={field.placeholder || "Select..."} />
                      </SelectTrigger>
                      <SelectContent>
                        {field.options?.map((option: any) => (
                          <SelectItem key={option.value} value={option.value}>
                            {option.label}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  )}

                  {field.type === "radio" && (
                    <RadioGroup onValueChange={(value) => form.setValue(field.id, value)}>
                      {field.options?.map((option: any) => (
                        <div key={option.value} className="flex items-center space-x-2">
                          <RadioGroupItem value={option.value} />
                          <Label>{option.label}</Label>
                        </div>
                      ))}
                    </RadioGroup>
                  )}

                  {field.type === "checkbox" && (
                    <div className="flex items-center space-x-2">
                      <Checkbox
                        onCheckedChange={(checked) => form.setValue(field.id, checked)}
                      />
                      <Label>{field.placeholder || "Check this"}</Label>
                    </div>
                  )}

                  {field.type === "switch" && (
                    <div className="flex items-center space-x-2">
                      <Switch
                        onCheckedChange={(checked) => form.setValue(field.id, checked)}
                      />
                      <Label>{field.placeholder || "Toggle this"}</Label>
                    </div>
                  )}

                  {field.type === "date" && (
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          variant="outline"
                          className={cn(
                            "w-full justify-start text-left font-normal",
                            !form.watch(field.id) && "text-muted-foreground"
                          )}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                          {form.watch(field.id) ? (
                            format(form.watch(field.id) as Date, "PPP")
                          ) : (
                            <span>{field.placeholder || "Pick a date"}</span>
                          )}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0">
                        <Calendar
                          mode="single"
                          selected={form.watch(field.id) as Date | undefined}
                          onSelect={(date) => form.setValue(field.id, date)}
                        />
                      </PopoverContent>
                    </Popover>
                  )}

                  {form.formState.errors[field.id] && (
                    <p className="text-sm text-destructive">
                      {form.formState.errors[field.id]?.message as string}
                    </p>
                  )}
                </div>
              );
            })}
            </div>
          ))}

          <Button type="submit" className="w-full">
            Submit Form
          </Button>
        </form>
      </div>
    </div>
  );
}
