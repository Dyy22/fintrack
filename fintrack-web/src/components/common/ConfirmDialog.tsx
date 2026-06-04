import { useEffect, useId, useRef, type KeyboardEvent } from "react";

export function ConfirmDialog({
  open,
  title,
  message,
  buttonLabel = "Deactivate",
  onConfirm,
  onCancel,
}: {
  open: boolean;
  title: string;
  message: string;
  buttonLabel?: string;
  onConfirm: () => void;
  onCancel: () => void;
}) {
  const cancelButtonRef = useRef<HTMLButtonElement>(null);
  const titleID = useId();
  const descriptionID = useId();

  useEffect(() => {
    if (open) {
      cancelButtonRef.current?.focus();
    }
  }, [open]);

  if (!open) return null;

  function handleKeyDown(event: KeyboardEvent<HTMLDivElement>) {
    if (event.key === "Escape") {
      onCancel();
    }
  }

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4 backdrop-blur-sm"
      onKeyDown={handleKeyDown}
    >
      <div
        className="neo-card w-full max-w-sm bg-blue-100 p-4 dark:bg-slate-800 sm:p-6"
        role="dialog"
        aria-modal="true"
        aria-labelledby={titleID}
        aria-describedby={descriptionID}
      >
        <h2
          id={titleID}
          className="text-xl font-black uppercase text-slate-950 dark:text-slate-100"
        >
          {title}
        </h2>
        <p
          id={descriptionID}
          className="mt-2 text-sm font-semibold text-slate-700 dark:text-slate-200"
        >
          {message}
        </p>
        <div className="mt-6 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <button
            ref={cancelButtonRef}
            className="neo-button bg-[#fffdf7] dark:bg-slate-800 dark:text-slate-100"
            onClick={onCancel}
          >
            Cancel
          </button>
          <button
            className="neo-button bg-red-300 dark:text-slate-950"
            onClick={onConfirm}
          >
            {buttonLabel}
          </button>
        </div>
      </div>
    </div>
  );
}
