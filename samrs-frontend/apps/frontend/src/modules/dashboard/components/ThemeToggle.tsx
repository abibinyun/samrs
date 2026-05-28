import * as React from "react";
import { Moon, Sun } from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { cn } from "@/lib/utils";

type Theme = "light" | "dark";

function setTheme(theme: Theme) {
  const html = document.documentElement;

  if (theme === "light") {
    html.removeAttribute("class");
  } else {
    html.setAttribute("class", theme);
  }

  localStorage.setItem("theme", theme);
}

export function ThemeToggle() {
  const [theme, setThemeState] = React.useState<Theme>("light");

  React.useEffect(() => {
    const stored = localStorage.getItem("theme") as Theme | null;
    if (stored) {
      setTheme(stored);
      setThemeState(stored);
    }
  }, []);

  const isDark = theme === "dark";

  return (
    <div className="flex items-center gap-2">
      <Sun className="h-4 w-4 text-muted-foreground" />
      <Switch
        checked={isDark}
        onCheckedChange={(checked) => {
          const nextTheme = checked ? "dark" : "light";
          setTheme(nextTheme);
          setThemeState(nextTheme);
        }}
        aria-label="Toggle dark mode"
      />
      <Moon className={cn("h-4 w-4", isDark ? "text-foreground" : "text-muted-foreground")} />
    </div>
  );
}
