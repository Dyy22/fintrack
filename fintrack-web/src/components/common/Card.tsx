import type { PropsWithChildren } from "react";

type CardProps = PropsWithChildren<{
  className?: string;
}>;

export function Card({ children, className = "" }: CardProps) {
  return <div className={`neo-card p-4 sm:p-6 ${className}`}>{children}</div>;
}
