import { useEffect, useRef, useState } from "react";

type NeoSelectOption = {
  value: string;
  label: string;
  disabled?: boolean;
};

type NeoSelectProps = {
  value: string;
  options: NeoSelectOption[];
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
  ariaLabel?: string;
};

export function NeoSelect({
  value,
  options,
  onChange,
  placeholder = "Select option",
  className = "",
  ariaLabel,
}: NeoSelectProps) {
  const rootRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const selectedOption = options.find((option) => option.value === value);

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

  return (
    <div ref={rootRef} className={`relative ${className}`}>
      <button
        type="button"
        className="neo-button min-h-11 w-full justify-between bg-[#fffdf7] normal-case tracking-normal dark:bg-slate-800 dark:text-slate-100"
        onClick={() => setOpen((current) => !current)}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label={ariaLabel ?? placeholder}
      >
        <span
          className={`min-w-0 truncate text-left ${
            selectedOption ? "" : "text-slate-500 dark:text-slate-200"
          }`}
        >
          {selectedOption?.label ?? placeholder}
        </span>
        <span className="shrink-0" aria-hidden="true">
          ⌄
        </span>
      </button>

      {open ? (
        <div
          className="neo-card absolute left-0 top-[calc(100%+0.75rem)] z-40 max-h-72 w-full overflow-y-auto bg-[#fffdf7] p-2 dark:bg-slate-800"
          role="listbox"
        >
          {options.map((option) => {
            const selected = option.value === value;
            return (
              <button
                key={option.value || "empty-option"}
                type="button"
                className={`mb-2 flex w-full items-center justify-between rounded-lg border-2 border-slate-950 px-3 py-2 text-left text-sm font-black shadow-[2px_2px_0_0_#101828] transition last:mb-0 hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none disabled:cursor-not-allowed disabled:opacity-50 dark:border-slate-100 dark:shadow-[2px_2px_0_0_#f8fafc] ${
                  selected
                    ? "bg-blue-300 text-slate-950"
                    : "bg-[#fffdf7] text-slate-950 dark:bg-slate-800 dark:text-slate-100"
                }`}
                onClick={() => {
                  if (option.disabled) return;
                  onChange(option.value);
                  setOpen(false);
                }}
                disabled={option.disabled}
                role="option"
                aria-selected={selected}
              >
                <span className="min-w-0 break-words">{option.label}</span>
                {selected ? <span aria-hidden="true">✓</span> : null}
              </button>
            );
          })}
        </div>
      ) : null}
    </div>
  );
}
