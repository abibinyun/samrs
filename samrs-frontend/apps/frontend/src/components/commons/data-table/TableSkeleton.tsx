import { cn } from "@/lib/utils";

interface TableSkeletonProps {
  rows?: number;
  className?: string;
  columns?: number;
  isFilterBoxVisible?: boolean;
}

export function FilterBoxSkeleton() {
  return (
    <div className="w-full h-32 my-5 p-4 rounded-lg border border-border bg-card animate-pulse">
      {/* Title / Label */}
      <div className="bg-muted/60 rounded h-4 w-1/3 mb-4" />

      {/* Input / Controls */}
      <div className="bg-muted/60 rounded h-10 w-full" />
    </div>
  );
}

export function TableSkeleton({ 
  rows = 10, 
  className,
  isFilterBoxVisible = false, 
  // columns = 3 
}: TableSkeletonProps) {
  return (
    <div className={cn("w-full space-y-3", className)}>
      {isFilterBoxVisible && <FilterBoxSkeleton />}
      {Array.from({ length: rows }).map((_, i) => (
        <div
          key={i}
          className="flex items-center gap-4 p-4 rounded-lg border bg-card animate-pulse"
        >
          {/* Column 1: Small (ID/Icon) */}
          <div className="w-10 h-4 rounded bg-muted"></div>
          
          {/* Column 2: Large (Name/Info) */}
          <div className="flex-1 h-4 rounded bg-muted"></div>
          
          {/* Column 3: Medium (Status) */}
          <div className="w-24 h-4 rounded bg-muted"></div>
          
          {/* Optional Column 4: Action icon */}
          <div className="w-8 h-4 rounded bg-muted"></div>
        </div>
      ))}
    </div>
  );
}