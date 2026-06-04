import type { PropsWithChildren } from "react";

type NeoAlertProps = PropsWithChildren<{
  variant?: "success" | "danger" | "info";
  className?: string;
}>;

export function NeoAlert({
  children,
  variant = "info",
  className = "",
}: NeoAlertProps) {
  const variants = {
    success: "bg-emerald-200",
    danger: "bg-red-200",
    info: "bg-blue-100",
  };

  return (
    <div
      className={`rounded-xl border-2 border-slate-950 px-4 py-3 text-sm font-bold text-slate-950 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:shadow-[3px_3px_0_0_#f8fafc] ${variants[variant]} ${className}`}
    >
      {children}
    </div>
  );
}
