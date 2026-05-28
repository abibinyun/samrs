import * as React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Eye, Save } from "lucide-react";
import { ComponentPalette } from "../components/ComponentPalette";
import { FormCanvas } from "../components/FormCanvas";
import { FieldProperties } from "../components/FieldProperties";
import { FormPreview } from "../components/FormPreview";
import type { FormField, FormSchema, FieldType } from "../types";

type Props = {
  initialName?: string;
  initialDescription?: string;
  initialSchema?: FormSchema;
  onSave: (data: { name: string; description?: string; schema: FormSchema }) => Promise<void>;
};

export function FormBuilder({ 
  initialName = "", 
  initialDescription = "",
  initialSchema,
  onSave 
}: Props) {
  const [name, setName] = React.useState(initialName);
  const [description, setDescription] = React.useState(initialDescription);
  const [fields, setFields] = React.useState<FormField[]>(initialSchema?.fields || []);
  const [selectedFieldId, setSelectedFieldId] = React.useState<string | null>(null);
  const [showPreview, setShowPreview] = React.useState(false);
  const [isSaving, setIsSaving] = React.useState(false);

  const selectedField = fields.find(f => f.id === selectedFieldId) || null;

  const handleAddField = (type: FieldType) => {
    const newField: FormField = {
      id: `field_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type,
      label: `${type.charAt(0).toUpperCase() + type.slice(1)} Field`,
      placeholder: "",
      required: false,
      gridColumn: 1,
      ...(type === "select" || type === "radio" ? { options: [
        { label: "Option 1", value: "opt1" },
        { label: "Option 2", value: "opt2" }
      ]} : {})
    };

    setFields([...fields, newField]);
    setSelectedFieldId(newField.id);
  };

  const handleUpdateField = (updates: Partial<FormField>) => {
    if (!selectedFieldId) return;
    
    setFields(fields.map(f => 
      f.id === selectedFieldId ? { ...f, ...updates } : f
    ));
  };

  const handleDeleteField = (id: string) => {
    setFields(fields.filter(f => f.id !== id));
    if (selectedFieldId === id) {
      setSelectedFieldId(null);
    }
  };

  const handleMoveField = (id: string, direction: "up" | "down") => {
    const index = fields.findIndex(f => f.id === id);
    if (index === -1) return;
    
    const newIndex = direction === "up" ? index - 1 : index + 1;
    if (newIndex < 0 || newIndex >= fields.length) return;

    const newFields = [...fields];
    [newFields[index], newFields[newIndex]] = [newFields[newIndex], newFields[index]];
    setFields(newFields);
  };

  const handleSave = async () => {
    if (!name.trim()) {
      alert("Please enter a form name");
      return;
    }

    setIsSaving(true);
    try {
      await onSave({
        name,
        description,
        schema: { fields }
      });
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <>
      <div className="h-screen flex flex-col">
        {/* Header */}
        <div className="border-b p-4">
          <div className="flex items-center justify-between">
            <div className="flex-1 max-w-md space-y-2">
              <Input
                placeholder="Form Name"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="font-semibold"
              />
              <Input
                placeholder="Description (optional)"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="text-sm"
              />
            </div>

            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                onClick={() => setShowPreview(true)}
                disabled={fields.length === 0}
              >
                <Eye className="h-4 w-4 mr-2" />
                Preview
              </Button>
              <Button onClick={handleSave} disabled={isSaving}>
                <Save className="h-4 w-4 mr-2" />
                {isSaving ? "Saving..." : "Save"}
              </Button>
            </div>
          </div>
        </div>

        {/* Main Content */}
        <div className="flex-1 flex overflow-hidden">
          <ComponentPalette onAddField={handleAddField} />
          
          <FormCanvas
            fields={fields}
            selectedFieldId={selectedFieldId}
            onSelectField={setSelectedFieldId}
            onDeleteField={handleDeleteField}
            onMoveField={handleMoveField}
          />
          
          <FieldProperties
            field={selectedField}
            onUpdate={handleUpdateField}
          />
        </div>
      </div>

      {showPreview && (
        <FormPreview
          schema={{ fields }}
          onClose={() => setShowPreview(false)}
        />
      )}
    </>
  );
}
