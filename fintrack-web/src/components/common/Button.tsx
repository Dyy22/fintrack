import type { ButtonHTMLAttributes, PropsWithChildren } from "react";

type ButtonProps = PropsWithChildren<
  ButtonHTMLAttributes<HTMLButtonElement>
> & {
  variant?: "primary" | "secondary" | "danger";
};

export function Button({
  children,
  className = "",
  variant = "primary",
  ...props
}: ButtonProps) {
  const variants = {
    primary:
      "bg-blue-300 hover:bg-blue-200 dark:bg-blue-300 dark:text-slate-950",
    secondary:
      "bg-[#fffdf7] hover:bg-emerald-100 dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700",
    danger: "bg-red-300 hover:bg-red-200 dark:bg-red-300 dark:text-slate-950",
  };

  return (
    <button
      className={`neo-button ${variants[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}
