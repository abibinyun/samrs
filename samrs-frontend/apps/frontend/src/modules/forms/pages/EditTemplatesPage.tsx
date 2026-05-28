import { useNavigate, useParams } from "@tanstack/react-router";
import { mockFormTemplatesApi } from "../actions/mockApi";
import { FormBuilder } from "../components/FormBuilder";

export default function EditTemplatePage() {
  const navigate = useNavigate();
  const search = useParams({ strict: false }) as any;
  const id = search?.id;
  
  // Mock API call - ganti dengan real API nanti
  const template = id ? mockFormTemplatesApi.getTemplate(id) : null;

  if (!id) {
    return <div className="p-6">No template ID provided</div>;
  }

  if (!template) {
    return <div className="p-6">Template not found</div>;
  }

  return (
    <FormBuilder
      initialName={template.name}
      initialDescription={template.description}
      initialSchema={template.schema}
      onSave={async ({ name, description, schema }) => {
        // Mock API call - ganti dengan real API nanti
        mockFormTemplatesApi.updateTemplate(id, { name, description, schema });
        navigate({ to: "/forms" });
      }}
    />
  );
}