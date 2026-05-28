import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { Plus, Trash2 } from "lucide-react";
import type { FormField } from "../types";

type Props = {
  field: FormField | null;
  onUpdate: (updates: Partial<FormField>) => void;
};

export function FieldProperties({ field, onUpdate }: Props) {
  if (!field) {
    return (
      <div className="max-w-[22%] border-l bg-muted/10 p-4">
        <p className="text-sm text-muted-foreground">
          Select a field to edit properties
        </p>
      </div>
    );
  }

  return (
    <div className="max-w-[22%] border-l bg-muted/10">
      <div className="p-4 border-b">
        <h3 className="font-semibold">Field Properties</h3>
      </div>
      
      <div className="p-4 space-y-4">
        <div className="space-y-2">
          <Label>Label</Label>
          <Input
            value={field.label}
            onChange={(e) => onUpdate({ label: e.target.value })}
          />
        </div>

        <div className="space-y-2">
          <Label>Placeholder</Label>
          <Input
            value={field.placeholder || ""}
            onChange={(e) => onUpdate({ placeholder: e.target.value })}
          />
        </div>

        <div className="space-y-2">
          <Label>Column Span</Label>
          <Select
            value={String(field.gridColumn || 1)}
            onValueChange={(v) => onUpdate({ gridColumn: Number(v) as 1 | 2 | 3 })}
          >
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="1">1 Column</SelectItem>
              <SelectItem value="2">2 Columns</SelectItem>
              <SelectItem value="3">3 Columns</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="flex items-center space-x-2">
          <Checkbox
            checked={field.required}
            onCheckedChange={(checked) => onUpdate({ required: !!checked })}
          />
          <Label>Required</Label>
        </div>

        {(field.type === "select" || field.type === "radio") && (
          <>
            <Separator />
            <div className="space-y-2">
              <Label>Options</Label>
              {field.options?.map((option, index) => (
                <div key={index} className="flex gap-2">
                  <Input
                    placeholder="Label"
                    value={option.label}
                    onChange={(e) => {
                      const newOptions = [...(field.options || [])];
                      newOptions[index] = { ...option, label: e.target.value };
                      onUpdate({ options: newOptions });
                    }}
                  />
                  <Input
                    placeholder="Value"
                    value={option.value}
                    onChange={(e) => {
                      const newOptions = [...(field.options || [])];
                      newOptions[index] = { ...option, value: e.target.value };
                      onUpdate({ options: newOptions });
                    }}
                  />
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => {
                      const newOptions = field.options?.filter((_, i) => i !== index);
                      onUpdate({ options: newOptions });
                    }}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
              <Button
                variant="outline"
                size="sm"
                className="w-full"
                onClick={() => {
                  const newOptions = [
                    ...(field.options || []),
                    { label: `Option ${(field.options?.length || 0) + 1}`, value: `opt${(field.options?.length || 0) + 1}` }
                  ];
                  onUpdate({ options: newOptions });
                }}
              >
                <Plus className="h-4 w-4 mr-2" />
                Add Option
              </Button>
            </div>
          </>
        )}
      </div>
    </div>
  );
}
