export function formatIDR(value: number | null | undefined): string {
  if (value == null) return "-";
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(value);
}

export function formatDate(value: string): string {
  return new Intl.DateTimeFormat("id-ID", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  }).format(new Date(value));
}

export function formatGoldGrams(value: number | null | undefined): string {
  if (value == null) return "-";
  return `${new Intl.NumberFormat("id-ID", {
    maximumFractionDigits: 4,
  }).format(value)} gr`;
}

export function transactionAmountLabel(type: string, amount: number): string {
  if (type === "expense") return `- ${formatIDR(amount)}`;
  if (type === "income") return `+ ${formatIDR(amount)}`;
  return formatIDR(amount);
}
