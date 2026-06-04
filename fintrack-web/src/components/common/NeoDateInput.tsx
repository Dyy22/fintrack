import { useEffect, useMemo, useRef, useState } from "react";
import {
  clampFutureDateInput,
  getTodayDateInputValue,
} from "../../utils/dateInput";

type NeoDateInputProps = {
  value: string;
  onChange: (value: string) => void;
  max?: string;
  className?: string;
  ariaLabel?: string;
};

function parseDateInput(value: string) {
  const [year, month, day] = value.split("-").map(Number);
  if (!year || !month || !day) return null;
  return new Date(year, month - 1, day);
}

function formatDateInput(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function formatDisplayDate(value: string) {
  const date = parseDateInput(value);
  if (!date) return "Select date";
  return date.toLocaleDateString("en-GB", {
    day: "2-digit",
    month: "short",
    year: "numeric",
  });
}

export function NeoDateInput({
  value,
  onChange,
  max = getTodayDateInputValue(),
  className = "",
  ariaLabel = "Select date",
}: NeoDateInputProps) {
  const rootRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const [viewDate, setViewDate] = useState(
    () => parseDateInput(value) ?? parseDateInput(max) ?? new Date(),
  );
  const maxDate = parseDateInput(max) ?? new Date();

  useEffect(() => {
    if (!open) return;

    function handlePointerDown(event: PointerEvent) {
      if (!rootRef.current?.contains(event.target as Node)) {
        setOpen(false);
      }
    }

    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setOpen(false);
      }
    }

    document.addEventListener("pointerdown", handlePointerDown);
    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [open]);

  function handleToggleOpen() {
    if (!open) {
      setViewDate(parseDateInput(value) ?? parseDateInput(max) ?? new Date());
    }
    setOpen(!open);
  }

  const calendarDays = useMemo(() => {
    const year = viewDate.getFullYear();
    const month = viewDate.getMonth();
    const firstDay = new Date(year, month, 1);
    const daysInMonth = new Date(year, month + 1, 0).getDate();
    const leadingBlanks = firstDay.getDay();

    return [
      ...Array.from({ length: leadingBlanks }, () => null),
      ...Array.from(
        { length: daysInMonth },
        (_, index) => new Date(year, month, index + 1),
      ),
    ];
  }, [viewDate]);

  function moveMonth(direction: -1 | 1) {
    setViewDate(
      (current) =>
        new Date(current.getFullYear(), current.getMonth() + direction, 1),
    );
  }

  function handleSelect(date: Date) {
    const nextValue = clampFutureDateInput(formatDateInput(date), max);
    onChange(nextValue);
    setOpen(false);
  }

  const viewMonthLabel = viewDate.toLocaleDateString("en-GB", {
    month: "long",
    year: "numeric",
  });

  return (
    <div ref={rootRef} className={`relative ${className}`}>
      <button
        type="button"
        className="neo-button min-h-11 w-full justify-between bg-[#fffdf7] normal-case tracking-normal dark:bg-slate-800 dark:text-slate-100"
        onClick={handleToggleOpen}
        aria-haspopup="dialog"
        aria-expanded={open}
        aria-label={`${ariaLabel}: ${formatDisplayDate(value)}`}
      >
        <span className="min-w-0 truncate text-left">
          {formatDisplayDate(value)}
        </span>
        <span className="shrink-0" aria-hidden="true">
          📅
        </span>
      </button>

      {open ? (
        <div
          className="neo-card absolute left-0 top-[calc(100%+0.75rem)] z-40 w-[min(20rem,calc(100vw-2rem))] bg-[#fffdf7] p-3 dark:bg-slate-800 sm:p-4"
          role="dialog"
          aria-label={`${ariaLabel} calendar`}
        >
          <div className="mb-3 flex items-center justify-between gap-2">
            <button
              type="button"
              className="neo-button bg-[#fffdf7] px-3 py-1 dark:bg-slate-800 dark:text-slate-100"
              onClick={() => moveMonth(-1)}
              aria-label="Previous month"
            >
              ←
            </button>
            <p className="text-center text-sm font-black uppercase text-slate-950 dark:text-slate-100">
              {viewMonthLabel}
            </p>
            <button
              type="button"
              className="neo-button bg-[#fffdf7] px-3 py-1 dark:bg-slate-800 dark:text-slate-100"
              onClick={() => moveMonth(1)}
              aria-label="Next month"
              disabled={
                viewDate.getFullYear() === maxDate.getFullYear() &&
                viewDate.getMonth() >= maxDate.getMonth()
              }
            >
              →
            </button>
          </div>

          <div className="grid grid-cols-7 gap-1 text-center text-[0.65rem] font-black uppercase text-slate-600 dark:text-slate-200">
            {["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"].map((day) => (
              <span key={day}>{day}</span>
            ))}
          </div>

          <div className="mt-2 grid grid-cols-7 gap-1">
            {calendarDays.map((date, index) => {
              if (!date) return <span key={`blank-${index}`} />;

              const inputValue = formatDateInput(date);
              const isSelected = inputValue === value;
              const isFuture = inputValue > max;

              return (
                <button
                  key={inputValue}
                  type="button"
                  className={`rounded-lg border-2 border-slate-950 px-2 py-1 text-sm font-black shadow-[2px_2px_0_0_#101828] transition hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none dark:border-slate-100 dark:shadow-[2px_2px_0_0_#f8fafc] ${
                    isSelected
                      ? "bg-blue-300 text-slate-950"
                      : "bg-[#fffdf7] text-slate-950 dark:bg-slate-800 dark:text-slate-100"
                  } ${isFuture ? "cursor-not-allowed opacity-40" : ""}`}
                  onClick={() => handleSelect(date)}
                  disabled={isFuture}
                  aria-label={date.toLocaleDateString("en-GB", {
                    day: "2-digit",
                    month: "long",
                    year: "numeric",
                  })}
                  aria-pressed={isSelected}
                >
                  {date.getDate()}
                </button>
              );
            })}
          </div>
        </div>
      ) : null}
    </div>
  );
}
