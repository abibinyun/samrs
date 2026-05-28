import type { FormTemplate } from "../types";

const STORAGE_KEY = "form_templates";

// Mock API menggunakan localStorage
export const mockFormTemplatesApi = {
  // Get all templates
  listTemplates: (): FormTemplate[] => {
    const data = localStorage.getItem(STORAGE_KEY);
    return data ? JSON.parse(data) : [];
  },

  // Get single template
  getTemplate: (id: string): FormTemplate | null => {
    const templates = mockFormTemplatesApi.listTemplates();
    return templates.find(t => t.id === id) || null;
  },

  // Create template
  createTemplate: (data: { name: string; description?: string; schema: any }): FormTemplate => {
    const templates = mockFormTemplatesApi.listTemplates();
    const newTemplate: FormTemplate = {
      id: `template_${Date.now()}`,
      name: data.name,
      description: data.description,
      schema: data.schema,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
    };
    
    templates.push(newTemplate);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(templates));
    return newTemplate;
  },

  // Update template
  updateTemplate: (id: string, data: { name?: string; description?: string; schema?: any }): FormTemplate | null => {
    const templates = mockFormTemplatesApi.listTemplates();
    const index = templates.findIndex(t => t.id === id);
    
    if (index === -1) return null;
    
    templates[index] = {
      ...templates[index],
      ...data,
      updatedAt: new Date().toISOString(),
    };
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(templates));
    return templates[index];
  },

  // Delete template
  deleteTemplate: (id: string): boolean => {
    const templates = mockFormTemplatesApi.listTemplates();
    const filtered = templates.filter(t => t.id !== id);
    
    if (filtered.length === templates.length) return false;
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(filtered));
    return true;
  },
};
