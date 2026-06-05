import type { Account } from "../types";

export function accountTypeLabel(type: string): string {
  switch (type) {
    case "stock_broker":
      return "Stock";
    case "ewallet":
      return "E-Wallet";
    case "gold":
      return "Gold";
    case "bank":
      return "Bank";
    case "cash":
      return "Cash";
    default:
      return type.replace(/_/g, " ").replace(/\b\w/g, (char) => char.toUpperCase());
  }
}

export function accountDisplayBalance(account: Account): number {
  if (
    account.type === "stock_broker" &&
    account.stock_lots != null &&
    account.stock_price_per_share != null &&
    account.stock_price_per_share > 0
  ) {
    return account.stock_lots * 100 * account.stock_price_per_share;
  }
  return account.balance;
}
