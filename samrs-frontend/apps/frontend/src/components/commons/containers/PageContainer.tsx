import * as React from "react";
import { cn } from "@/lib/utils";

type PageHeaderProps = {
  title: React.ReactNode;
  description?: React.ReactNode;
  actions?: React.ReactNode;
  top?: React.ReactNode;
  bottom?: React.ReactNode;
};

export type PageContainerProps = {
  className?: string;
  header?: PageHeaderProps;
  children: React.ReactNode;
  isSticky?: boolean;
};

export default function PageContainer({ className, header, children, isSticky }: PageContainerProps) {
  return (
    <div className={cn("space-y-6", className)}>
      {header ? (
        <div 
          className={cn(
            "space-y-3 transition-all",
            isSticky && [
              "sticky z-20 bg-background backdrop-blur-sm",
              "-top-4 md:top-0 -mt-20", 
              "-mx-4 md:-mx-8 px-4 md:px-8 pt-20 pb-5 mb-6 "
            ]
          )}
        >
          {header.top ? <div>{header.top}</div> : null}

          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div>
              <h1 className="text-2xl font-bold tracking-tight">{header.title}</h1>
              {header.description ? (
                <p className="text-muted-foreground text-sm">{header.description}</p>
              ) : null}
            </div>

            {header.actions ? <div className="flex gap-2">{header.actions}</div> : null}
          </div>

          {header.bottom ? <div>{header.bottom}</div> : null}
        </div>
      ) : null}

      <div className="relative">
        {children}
      </div>
    </div>
  );
}