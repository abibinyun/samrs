import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { 
  GripVertical, 
  Trash2,
  ChevronUp,
  ChevronDown
} from "lucide-react";
import type { FormField } from "../types";

type Props = {
  fields: FormField[];
  selectedFieldId: string | null;
  onSelectField: (id: string) => void;
  onDeleteField: (id: string) => void;
  onMoveField: (id: string, direction: "up" | "down") => void;
};

export function FormCanvas({ 
  fields, 
  selectedFieldId, 
  onSelectField, 
  onDeleteField,
  onMoveField 
}: Props) {
  if (fields.length === 0) {
    return (
      <div className="flex-1 flex items-center justify-center bg-muted/5">
        <div className="text-center text-muted-foreground">
          <p className="text-lg font-medium">No fields yet</p>
          <p className="text-sm">Add components from the left panel</p>
        </div>
      </div>
    );
  }

  const groupedFields: FormField[][] = [];
  let currentRow: FormField[] = [];
  
  fields.forEach((field) => {
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
    <div className="flex-1 p-6 overflow-auto bg-muted/5">
      <div className="max-w-4xl mx-auto space-y-4">
        {groupedFields.map((row, rowIndex) => (
          <div 
            key={rowIndex} 
            className="grid grid-cols-3 gap-4"
          >
            {row.map((field) => {
              const fieldIndex = fields.findIndex(f => f.id === field.id);
              const isSelected = selectedFieldId === field.id;
              const colSpanClass = field.gridColumn === 3 ? "col-span-3" : field.gridColumn === 2 ? "col-span-2" : "col-span-1";
              
              return (
                <Card
                  key={field.id}
                  className={`p-4 cursor-pointer transition-all ${
                    isSelected ? "ring-2 ring-primary" : ""
                  } ${colSpanClass}`}
                  onClick={() => onSelectField(field.id)}
                >
                  <div className="flex items-start justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <GripVertical className="h-4 w-4 text-muted-foreground" />
                      <Badge variant="secondary" className="text-xs">
                        {field.type}
                      </Badge>
                    </div>
                    
                    <div className="flex items-center gap-1">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={(e) => {
                          e.stopPropagation();
                          onMoveField(field.id, "up");
                        }}
                        disabled={fieldIndex === 0}
                      >
                        <ChevronUp className="h-3 w-3" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={(e) => {
                          e.stopPropagation();
                          onMoveField(field.id, "down");
                        }}
                        disabled={fieldIndex === fields.length - 1}
                      >
                        <ChevronDown className="h-3 w-3" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 w-6 p-0"
                        onClick={(e) => {
                          e.stopPropagation();
                          onDeleteField(field.id);
                        }}
                      >
                        <Trash2 className="h-3 w-3" />
                      </Button>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <Label>
                      {field.label}
                      {field.required && <span className="text-destructive ml-1">*</span>}
                    </Label>
                    
                    {field.type === "input" && (
                      <Input placeholder={field.placeholder} disabled />
                    )}
                    
                    {field.type === "textarea" && (
                      <textarea 
                        className="flex min-h-20 w-full rounded-md border border-input bg-background px-3 py-2 text-sm" 
                        placeholder={field.placeholder}
                        disabled
                      />
                    )}
                    
                    {field.type === "select" && (
                      <div className="text-sm text-muted-foreground">
                        {field.options?.length || 0} options
                      </div>
                    )}
                    
                    {(field.type === "checkbox" || field.type === "switch") && (
                      <div className="text-sm text-muted-foreground">
                        Toggle field
                      </div>
                    )}
                  </div>
                </Card>
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}
