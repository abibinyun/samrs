export type FieldType = 
  | "input" 
  | "textarea" 
  | "select" 
  | "checkbox" 
  | "radio" 
  | "switch" 
  | "date";

export type FormField = {
  id: string;
  type: FieldType;
  label: string;
  placeholder?: string;
  required?: boolean;
  options?: { label: string; value: string }[];
  defaultValue?: any;
  gridColumn?: 1 | 2 | 3;
};

export type FormSchema = {
  fields: FormField[];
};

export type FormTemplate = {
  id: string;
  name: string;
  description?: string;
  schema: FormSchema;
  createdAt: string;
  updatedAt: string;
};
