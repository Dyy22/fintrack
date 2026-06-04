export function getTodayDateInputValue() {
  const today = new Date();
  const year = today.getFullYear();
  const month = String(today.getMonth() + 1).padStart(2, "0");
  const day = String(today.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

export function clampFutureDateInput(value: string, max = getTodayDateInputValue()) {
  if (!value) return value;
  return value > max ? max : value;
}
