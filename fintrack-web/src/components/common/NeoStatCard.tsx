import type { ReactNode } from "react";

type NeoStatCardProps = {
  label: string;
  value: ReactNode;
  icon?: ReactNode;
  tone?: "blue" | "emerald" | "red" | "neutral";
  helper?: ReactNode;
  className?: string;
};

export function NeoStatCard({
  label,
  value,
  icon,
  tone = "blue",
  helper,
  className = "",
}: NeoStatCardProps) {
  const tones = {
    blue: "bg-blue-100",
    emerald: "bg-emerald-100",
    red: "bg-red-100",
    neutral: "bg-[#fffdf7]",
  };
  const iconTones = {
    blue: "bg-blue-300",
    emerald: "bg-emerald-200",
    red: "bg-red-200",
    neutral: "bg-slate-200",
  };

  return (
    <div className={`neo-card ${tones[tone]} p-4 sm:p-6 ${className}`}>
      <div className="flex items-start justify-between gap-3 sm:gap-4">
        <div className="min-w-0">
          <p className="text-xs font-black uppercase tracking-[0.16em] text-slate-600 dark:text-slate-200">
            {label}
          </p>
          <div className="mt-2 break-words text-2xl font-black text-slate-950 dark:text-slate-100 sm:text-3xl">
            {value}
          </div>
        </div>
        {icon ? (
          <div
            className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl border-2 border-slate-950 text-lg font-black text-slate-950 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:shadow-[3px_3px_0_0_#f8fafc] sm:h-12 sm:w-12 sm:text-xl ${iconTones[tone]}`}
            aria-hidden="true"
          >
            {icon}
          </div>
        ) : null}
      </div>
      {helper ? (
        <div className="mt-4 text-sm font-semibold text-slate-600 dark:text-slate-200">
          {helper}
        </div>
      ) : null}
    </div>
  );
}
