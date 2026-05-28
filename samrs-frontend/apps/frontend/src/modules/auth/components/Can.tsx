import { usePermission } from "@/hooks/usePermission";

interface CanProps {
  perform: string | string[];
  children: React.ReactNode;
}

export const Can = ({ perform, children }: CanProps) => {
  const { hasPermission } = usePermission();

  if (!hasPermission(perform)) return null;

  return <>{children}</>;
};