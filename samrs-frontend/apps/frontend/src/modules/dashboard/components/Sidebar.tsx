import type { MenuItem, MenuSection, SubMenuItem } from "@/constants/navigations";
import { MENU_ITEMS } from "@/constants/navigations";
import { cn } from "@/lib/utils";
import { Link, useLocation, useParams } from "@tanstack/react-router";
import { ChevronDown, type LucideIcon } from "lucide-react";
import { useMemo, useState } from "react";
import { useAppSelector } from "@/store/hooks";

/* ============================================================
 * Permission Utilities (Deterministic & Safe)
 * ============================================================
 */

export const hasPermission = (
  permissions: readonly string[],
  permission?: string,
  isSystemUser?: boolean
): boolean => {
  // Super Admin bypass (check is_system flag or role.is_system)
  if (isSystemUser) return true;
  
  return !permission || permissions.includes(permission);
};

/* ============================================================
 * Menu Filtering (No null, No side-effects)
 * ============================================================
 */

export const filterMenuByPermission = (
  sections: readonly MenuSection[],
  permissions: readonly string[],
  isSystemUser?: boolean
): MenuSection[] => {
  return sections.flatMap(section => {
    const items = section.items.flatMap<MenuItem>(item => {
      // Item with children
      if (item.children) {
        const children = item.children.filter(child =>
          hasPermission(permissions, child.permission, isSystemUser)
        );

        return children.length
          ? [{ ...item, children }]
          : [];
      }

      // Single item
      return hasPermission(permissions, item.permission, isSystemUser)
        ? [item]
        : [];
    });

    return items.length
      ? [{ ...section, items }]
      : [];
  });
};

/* ============================================================
 * Route Matching (Precise & Hardened)
 * ============================================================
 */

const isRouteActive = (pathname: string, target: string): boolean => {
  return pathname === target || pathname.startsWith(`${target}/`);
};

/* ============================================================
 * Sidebar Link Component
 * ============================================================
 */

interface SidebarLinkProps {
  label: string;
  icon?: LucideIcon;
  to?: string;
  children?: SubMenuItem[];
}

export const SidebarLink = ({
  label,
  icon: Icon,
  to,
  children,
}: SidebarLinkProps) => {
  const { tenantId } = useParams({ from: "/$tenantId" });
  const { pathname } = useLocation();
  const [manualOpen, setManualOpen] = useState(false);

  const isChildActive = useMemo(() => {
    if (!children?.length) return false;

    return children.some(child =>
      isRouteActive(
        pathname,
        `/${tenantId}${child.to}`
      )
    );
  }, [children, pathname, tenantId]);

  const isOpen = isChildActive || manualOpen;

  /* ---------- Parent with children ---------- */
  if (children?.length) {
    return (
      <div className="space-y-1">
        <button
          type="button"
          aria-expanded={isOpen}
          onClick={() => setManualOpen(o => !o)}
          className={cn(
            "w-full flex items-center justify-between px-3 py-2 rounded-md text-sm transition-colors",
            isChildActive
              ? "text-primary font-medium bg-primary/5"
              : "text-muted-foreground hover:bg-muted"
          )}
        >
          <div className="flex items-center gap-3">
            {Icon && <Icon className="w-4 h-4" />}
            <span>{label}</span>
          </div>
          <ChevronDown
            className={cn(
              "w-3 h-3 transition-transform",
              isOpen && "rotate-180"
            )}
          />
        </button>

        {isOpen && (
          <div className="ml-5 border-l pl-4 space-y-1">
            {children.map(child => (
              <Link
                key={child.to}
                to={`/$tenantId${child.to}`}
                params={{ tenantId }}
                activeProps={{
                  className: "text-primary font-semibold",
                }}
                className="block px-3 py-2 text-xs text-muted-foreground rounded-md hover:bg-muted/50 hover:text-primary"
              >
                {child.label}
              </Link>
            ))}
          </div>
        )}
      </div>
    );
  }

  /* ---------- Single link ---------- */
  return (
    <Link
      to={`/$tenantId${to === "/" ? "" : to}`}
      params={{ tenantId }}
      activeProps={{
        className: "bg-primary/10 text-primary font-medium",
      }}
      className="flex items-center gap-3 px-3 py-2 rounded-md text-sm text-muted-foreground hover:bg-muted hover:text-primary"
    >
      {Icon && <Icon className="w-4 h-4" />}
      <span>{label}</span>
    </Link>
  );
};

/* ============================================================
 * Sidebar Component
 * ============================================================
 */

export const Sidebar = () => {
  const { permissions, user } = useAppSelector(
    (state) => state?.auth ?? { permissions: [], user: null }
  );

  const sections = filterMenuByPermission(MENU_ITEMS, permissions, user?.is_system);

  return (
    <aside className="w-64 border-r bg-card h-screen sticky top-0 hidden md:flex flex-col dark:scrollbar-thumb-cyan-800 dark:scrollbar-track-sidebar">
      {/* Header */}
      <div className="p-6">
        <div className="flex items-center gap-2 font-bold text-xl text-primary">
          <div className="w-8 h-8 bg-primary rounded-lg flex items-center justify-center text-white text-sm">
            S
          </div>
          SAMRS
          <span className="text-[10px] font-mono border px-1 rounded">
            V3
          </span>
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex-1 px-3 space-y-6 overflow-y-auto pb-10 scrollbar-thin">
        {sections.map(section => (
          <div key={section.group} className="space-y-2">
            <h3 className="px-3 text-[10px] font-bold text-muted-foreground/60 uppercase tracking-widest">
              {section.group}
            </h3>

            <div className="space-y-1">
              {section.items.map(item => (
                <SidebarLink
                  key={item.to ?? item.label}
                  label={item.label}
                  icon={item.icon}
                  to={item.to}
                  children={item.children}
                />
              ))}
            </div>
          </div>
        ))}
      </nav>

      {/* Footer */}
      <div className="p-4 border-t bg-muted/20">
        <div className="flex items-center justify-center gap-2">
          <div className="w-2 h-2 rounded-full bg-green-500 animate-pulse" />
          <span className="text-[9px] uppercase tracking-widest text-muted-foreground">
            System Active
          </span>
        </div>
      </div>
    </aside>
  );
};
