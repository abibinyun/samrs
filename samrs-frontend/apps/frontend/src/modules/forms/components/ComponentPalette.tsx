import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import { 
  Type, 
  TextCursorInput, 
  List, 
  CheckSquare, 
  ToggleLeft, 
  Calendar,
  Circle
} from "lucide-react";
import type { FieldType } from "../types";

const FIELD_COMPONENTS = [
  { type: "input", label: "Input", icon: Type },
  { type: "textarea", label: "Textarea", icon: TextCursorInput },
  { type: "select", label: "Select", icon: List },
  { type: "checkbox", label: "Checkbox", icon: CheckSquare },
  { type: "radio", label: "Radio", icon: Circle },
  { type: "switch", label: "Switch", icon: ToggleLeft },
  { type: "date", label: "Date", icon: Calendar },
] as const;

type Props = {
  onAddField: (type: FieldType) => void;
};

export function ComponentPalette({ onAddField }: Props) {
  return (
    <div className="max-w-[15%] border-r bg-muted/10">
      <div className="p-4 border-b">
        <h3 className="font-semibold">Components</h3>
        <p className="text-xs text-muted-foreground mt-1">Click to add to form</p>
      </div>
      
      <ScrollArea className="h-[calc(100vh-180px)]">
        <div className="p-4 space-y-2">
          {FIELD_COMPONENTS.map(({ type, label, icon: Icon }) => (
            <Button
              key={type}
              variant="outline"
              className="w-full justify-start"
              onClick={() => onAddField(type as FieldType)}
            >
              <Icon className="h-4 w-4 mr-2" />
              {label}
            </Button>
          ))}
        </div>
      </ScrollArea>
    </div>
  );
}
