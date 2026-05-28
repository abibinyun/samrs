import { useNavigate } from "@tanstack/react-router";
import { mockFormTemplatesApi } from "../actions/mockApi";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Plus, Edit, Trash2, Eye } from "lucide-react";
import { format } from "date-fns";
import * as React from "react";

export default function FormTemplatesPage() {
  const navigate = useNavigate();
  const [templates, setTemplates] = React.useState(mockFormTemplatesApi.listTemplates());
  const handleDelete = (id: string) => {
    if (confirm("Are you sure you want to delete this form template?")) {
      mockFormTemplatesApi.deleteTemplate(id);
      setTemplates(mockFormTemplatesApi.listTemplates());
    }
  };

  return (
    <div className="p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-3xl font-bold">Form Templates</h1>
          <p className="text-muted-foreground">Create and manage form templates</p>
        </div>
        <Button onClick={() => navigate({ to: "create" })}>
          <Plus className="h-4 w-4 mr-2" />
          Create Form
        </Button>
      </div>

      {templates.length === 0 ? (
        <Card className="p-12 text-center">
          <p className="text-muted-foreground mb-4">No form templates yet</p>
          <Button onClick={() => navigate({ to: "create" })}>
            <Plus className="h-4 w-4 mr-2" />
            Create Your First Form
          </Button>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {templates.map((template) => (
            <Card key={template.id} className="p-4">
              <div className="flex items-start justify-between mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold">{template.name}</h3>
                  {template.description && (
                    <p className="text-sm text-muted-foreground mt-1">
                      {template.description}
                    </p>
                  )}
                </div>
                <Badge variant="secondary">
                  {template.schema?.fields?.length || 0} fields
                </Badge>
              </div>

              <div className="text-xs text-muted-foreground mb-4">
                Updated {format(new Date(template.updatedAt), "MMM d, yyyy")}
              </div>

              <div className="flex items-center gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => navigate({ to: `fill/${template.id}` })}
                >
                  <Eye className="h-3 w-3 mr-1" />
                  Fill
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => navigate({ to: `edit/${template.id}`  })}
                >
                  <Edit className="h-3 w-3 mr-1" />
                  Edit
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => handleDelete(template.id)}
                >
                  <Trash2 className="h-3 w-3" />
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}