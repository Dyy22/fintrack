import type { PropsWithChildren } from "react";

type NeoBadgeProps = PropsWithChildren<{
  variant?: "neutral" | "success" | "danger" | "info" | "warning";
  className?: string;
}>;

export function NeoBadge({
  children,
  variant = "neutral",
  className = "",
}: NeoBadgeProps) {
  const variants = {
    neutral: "bg-slate-200 text-slate-800",
    success: "bg-emerald-200 text-slate-950",
    danger: "bg-red-200 text-slate-950",
    info: "bg-blue-100 text-slate-950 dark:bg-blue-300",
    warning: "bg-amber-200 text-slate-950",
  };

  return (
    <span
      className={`inline-flex items-center rounded-full border-2 border-slate-950 px-2 py-0.5 text-xs font-black uppercase tracking-wide dark:border-slate-100 ${variants[variant]} ${className}`}
    >
      {children}
    </span>
  );
}
