type SkeletonProps = {
  className?: string;
};

export function Skeleton({ className = "" }: SkeletonProps) {
  return (
    <div
      className={`animate-pulse rounded-xl border-2 border-slate-950 bg-slate-200 shadow-[3px_3px_0_0_#101828] dark:border-slate-100 dark:bg-slate-700 dark:shadow-[3px_3px_0_0_#f8fafc] ${className}`}
    />
  );
}

export function SkeletonCard({ className = "" }: SkeletonProps) {
  return (
    <div className={`neo-card p-6 ${className}`}>
      <Skeleton className="mb-3 h-4 w-24" />
      <Skeleton className="h-8 w-48" />
      <div className="mt-4 space-y-2">
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-3/4" />
      </div>
    </div>
  );
}

export function SkeletonTable({ rows = 5 }: { rows?: number }) {
  return (
    <div className="space-y-3">
      <Skeleton className="h-4 w-full" />
      {Array.from({ length: rows }).map((_, i) => (
        <Skeleton key={i} className="h-12 w-full" />
      ))}
    </div>
  );
}
