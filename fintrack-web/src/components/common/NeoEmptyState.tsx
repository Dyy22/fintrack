import type { PropsWithChildren, ReactNode } from "react";

type NeoEmptyStateProps = PropsWithChildren<{
  title: string;
  description?: string;
  icon?: ReactNode;
  action?: ReactNode;
  className?: string;
}>;

export function NeoEmptyState({
  title,
  description,
  icon = "✦",
  action,
  children,
  className = "",
}: NeoEmptyStateProps) {
  return (
    <div className={`neo-empty ${className}`}>
      <div
        className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-xl border-2 border-slate-950 bg-blue-300 text-2xl font-black text-slate-950 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:shadow-[3px_3px_0_0_#f8fafc]"
        aria-hidden="true"
      >
        {icon}
      </div>
      <p className="text-lg font-black uppercase text-slate-950 dark:text-slate-100">
        {title}
      </p>
      {description ? (
        <p className="mx-auto mt-2 max-w-md text-sm font-semibold text-slate-600 dark:text-slate-200">
          {description}
        </p>
      ) : null}
      {children}
      {action ? <div className="mt-4">{action}</div> : null}
    </div>
  );
}
