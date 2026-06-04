import type { InputHTMLAttributes } from "react";

type NeoInputProps = InputHTMLAttributes<HTMLInputElement>;

export function NeoInput({ className = "", ...props }: NeoInputProps) {
  return (
    <input
      className={`mt-1 w-full rounded-xl border-2 border-slate-950 bg-[#fffdf7] px-3 py-2 font-bold text-slate-950 shadow-[3px_3px_0_0_#101828] outline-none transition focus:-translate-y-0.5 focus:shadow-[5px_5px_0_0_#101828] dark:border-slate-100 dark:bg-slate-800 dark:text-slate-100 dark:shadow-[3px_3px_0_0_#f8fafc] dark:focus:shadow-[5px_5px_0_0_#f8fafc] ${className}`}
      {...props}
    />
  );
}
