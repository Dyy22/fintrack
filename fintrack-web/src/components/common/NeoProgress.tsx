type NeoProgressProps = {
  value: number;
  className?: string;
};

export function NeoProgress({ value, className = "" }: NeoProgressProps) {
  const normalizedValue = Math.max(0, Math.min(100, value));

  return (
    <div
      className={`h-4 flex-1 rounded-full border-2 border-slate-950 bg-[#fffdf7] dark:border-slate-100 dark:bg-slate-900 ${className}`}
    >
      <div
        className="h-full rounded-full bg-blue-300"
        style={{ width: `${normalizedValue}%` }}
      />
    </div>
  );
}
