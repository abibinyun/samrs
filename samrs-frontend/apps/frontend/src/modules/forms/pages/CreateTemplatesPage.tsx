import { useNavigate } from "@tanstack/react-router";
import { mockFormTemplatesApi } from "../actions/mockApi";
import { FormBuilder } from "../components/FormBuilder";

export default function CreateTemplatePage() {
  const navigate = useNavigate();

  return (
    <FormBuilder
      onSave={async ({ name, description, schema }) => {
        // Mock API call - ganti dengan real API nanti
        mockFormTemplatesApi.createTemplate({ name, description, schema });
        navigate({ to: "/forms" });
      }}
    />
  );
}