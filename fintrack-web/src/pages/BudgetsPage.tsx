import { useEffect, useMemo, useState } from "react";
import { Button } from "../components/common/Button";
import { Card } from "../components/common/Card";
import { NeoAlert } from "../components/common/NeoAlert";
import { NeoEmptyState } from "../components/common/NeoEmptyState";
import { NeoInput } from "../components/common/NeoInput";
import { NeoPageHeader } from "../components/common/NeoPageHeader";
import { NeoSelect } from "../components/common/NeoSelect";
import { SkeletonCard } from "../components/common/Skeleton";
import { useBudgetStore } from "../stores/budgetStore";
import { formatIDR } from "../utils/format";
import { usePageTitle } from "../utils/usePageTitle";

const MONTHS = [
  { value: "1", label: "January" },
  { value: "2", label: "February" },
  { value: "3", label: "March" },
  { value: "4", label: "April" },
  { value: "5", label: "May" },
  { value: "6", label: "June" },
  { value: "7", label: "July" },
  { value: "8", label: "August" },
  { value: "9", label: "September" },
  { value: "10", label: "October" },
  { value: "11", label: "November" },
  { value: "12", label: "December" },
];

const now = new Date();
const currentMonth = now.getMonth() + 1;
const currentYear = now.getFullYear();

export function BudgetsPage() {
  usePageTitle("Budgets");

  const {
    budgets,
    categories,
    loading,
    error,
    fetchBudgets,
    fetchCategories,
    createBudget,
    updateBudget,
    deleteBudget,
  } = useBudgetStore();

  const [selectedMonth, setSelectedMonth] = useState(String(currentMonth));
  const [selectedYear, setSelectedYear] = useState(String(currentYear));

  // modal state
  const [showModal, setShowModal] = useState(false);
  const [modalMode, setModalMode] = useState<"create" | "edit">("create");
  const [editBudgetId, setEditBudgetId] = useState<string | null>(null);
  const [formCategoryID, setFormCategoryID] = useState("");
  const [formAmount, setFormAmount] = useState("");
  const [formMonth, setFormMonth] = useState(String(currentMonth));
  const [formYear, setFormYear] = useState(String(currentYear));
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  // delete confirm
  const [deleteBudgetId, setDeleteBudgetId] = useState<string | null>(null);

  const monthNum = Number(selectedMonth);
  const yearNum = Number(selectedYear);

  useEffect(() => {
    fetchBudgets(monthNum, yearNum);
  }, [fetchBudgets, monthNum, yearNum]);

  useEffect(() => {
    fetchCategories();
  }, [fetchCategories]);

  const categoryOptions = useMemo(() => {
    const usedCategoryIDs = new Set(budgets.map((b) => b.category_id));
    return categories
      .filter((c) => c.type === "expense")
      .filter((c) => !usedCategoryIDs.has(c.id) || modalMode === "edit")
      .map((c) => ({ value: c.id, label: c.name }));
  }, [categories, budgets, modalMode]);

  function openCreateModal() {
    setModalMode("create");
    setEditBudgetId(null);
    setFormCategoryID("");
    setFormAmount("");
    setFormMonth(String(currentMonth));
    setFormYear(String(currentYear));
    setFormError(null);
    setShowModal(true);
  }

  function openEditModal(budget: {
    id: string;
    amount: number;
    category_id?: string;
    month?: number;
    year?: number;
  }) {
    setModalMode("edit");
    setEditBudgetId(budget.id);
    setFormCategoryID(budget.category_id ?? "");
    setFormAmount(String(budget.amount));
    setFormMonth(String(currentMonth));
    setFormYear(String(currentYear));
    setFormError(null);
    setShowModal(true);
  }

  function closeModal() {
    setShowModal(false);
    setEditBudgetId(null);
    setFormCategoryID("");
    setFormAmount("");
    setFormError(null);
  }

  async function handleSubmit() {
    setFormError(null);
    const amount = Number(formAmount);
    if (!amount || amount <= 0) {
      setFormError("Amount must be > 0");
      return;
    }
    if (modalMode === "create" && !formCategoryID) {
      setFormError("Select a category");
      return;
    }
    setSubmitting(true);
    try {
      if (modalMode === "create") {
        await createBudget({
          category_id: formCategoryID,
          month: Number(formMonth),
          year: Number(formYear),
          amount,
        });
      } else if (editBudgetId) {
        await updateBudget(editBudgetId, amount);
      }
      closeModal();
      fetchBudgets(monthNum, yearNum);
    } catch {
      setFormError("Failed to save budget");
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id: string) {
    setDeleteBudgetId(null);
    try {
      await deleteBudget(id);
    } catch {
      // error state handled silently
    }
  }

  const years = useMemo(() => {
    const y = currentYear;
    return Array.from({ length: 5 }, (_, i) => String(y - 2 + i));
  }, []);

  return (
    <div className="space-y-6">
      <NeoPageHeader
        title="Budgets"
        description="Set and track monthly spending limits per category."
        eyebrow="Spending limits"
        icon="💰"
        actions={
          <Button onClick={openCreateModal} disabled={categories.length === 0}>
            Add Budget
          </Button>
        }
      />

      {/* Month / Year selector */}
      <Card>
        <div className="flex flex-wrap items-end gap-3">
          <div className="min-w-[140px] flex-1">
            <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
              Month
            </label>
            <NeoSelect
              value={selectedMonth}
              options={MONTHS}
              onChange={setSelectedMonth}
              ariaLabel="Select month"
            />
          </div>
          <div className="min-w-[110px] flex-1">
            <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
              Year
            </label>
            <NeoSelect
              value={selectedYear}
              options={years.map((y) => ({ value: y, label: y }))}
              onChange={setSelectedYear}
              ariaLabel="Select year"
            />
          </div>
        </div>
      </Card>

      {/* Error banner */}
      {error ? <NeoAlert variant="danger">{error}</NeoAlert> : null}

      {/* Loading */}
      {loading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <SkeletonCard key={i} />
          ))}
        </div>
      ) : budgets.length === 0 ? (
        <Card>
          <NeoEmptyState
            title="No budgets set"
            description='Click "Add Budget" to set a spending limit for a category.'
            icon="💰"
            action={
              <Button
                onClick={openCreateModal}
                disabled={categories.length === 0}
              >
                Add Budget
              </Button>
            }
          />
        </Card>
      ) : (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {budgets.map((budget) => (
            <BudgetCard
              key={budget.id}
              budget={budget}
              onEdit={() => openEditModal(budget)}
              onDelete={() => setDeleteBudgetId(budget.id)}
            />
          ))}
        </div>
      )}

      {/* Create / Edit Modal */}
      {showModal ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4 backdrop-blur-sm">
          <div className="neo-card w-full max-w-md bg-blue-100 p-4 dark:bg-slate-800 sm:p-6">
            <h2 className="text-xl font-black uppercase text-slate-950 dark:text-slate-100">
              {modalMode === "create" ? "Set Budget" : "Edit Budget"}
            </h2>

            <div className="mt-4 space-y-4">
              {modalMode === "create" ? (
                <>
                  <div>
                    <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
                      Category
                    </label>
                    <NeoSelect
                      value={formCategoryID}
                      options={categoryOptions}
                      onChange={setFormCategoryID}
                      placeholder="Select category"
                      ariaLabel="Select category"
                    />
                  </div>
                  <div className="flex gap-3">
                    <div className="flex-1">
                      <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
                        Month
                      </label>
                      <NeoSelect
                        value={formMonth}
                        options={MONTHS}
                        onChange={setFormMonth}
                        ariaLabel="Budget month"
                      />
                    </div>
                    <div className="flex-1">
                      <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
                        Year
                      </label>
                      <NeoSelect
                        value={formYear}
                        options={years.map((y) => ({ value: y, label: y }))}
                        onChange={setFormYear}
                        ariaLabel="Budget year"
                      />
                    </div>
                  </div>
                </>
              ) : null}

              <div>
                <label className="mb-1 block text-xs font-black uppercase tracking-wide text-slate-600 dark:text-slate-200">
                  Amount (IDR)
                </label>
                <NeoInput
                  type="number"
                  min="1"
                  value={formAmount}
                  onChange={(e) => setFormAmount(e.target.value)}
                  placeholder="e.g. 1000000"
                />
              </div>

              {formError ? (
                <NeoAlert variant="danger">{formError}</NeoAlert>
              ) : null}

              <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
                <button
                  className="neo-button bg-[#fffdf7] dark:bg-slate-800 dark:text-slate-100"
                  onClick={closeModal}
                  disabled={submitting}
                >
                  Cancel
                </button>
                <button
                  className="neo-button bg-blue-300 dark:text-slate-950"
                  onClick={handleSubmit}
                  disabled={submitting}
                >
                  {submitting
                    ? "Saving..."
                    : modalMode === "create"
                      ? "Create Budget"
                      : "Update Budget"}
                </button>
              </div>
            </div>
          </div>
        </div>
      ) : null}

      {/* Delete confirmation */}
      {deleteBudgetId ? (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-950/60 p-4 backdrop-blur-sm">
          <div className="neo-card w-full max-w-sm bg-blue-100 p-4 dark:bg-slate-800 sm:p-6">
            <h2 className="text-xl font-black uppercase text-slate-950 dark:text-slate-100">
              Delete Budget
            </h2>
            <p className="mt-2 text-sm font-semibold text-slate-700 dark:text-slate-200">
              Remove this budget? This cannot be undone.
            </p>
            <div className="mt-6 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
              <button
                className="neo-button bg-[#fffdf7] dark:bg-slate-800 dark:text-slate-100"
                onClick={() => setDeleteBudgetId(null)}
              >
                Cancel
              </button>
              <button
                className="neo-button bg-red-300 dark:text-slate-950"
                onClick={() => handleDelete(deleteBudgetId)}
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}

/* ─── Budget Card ─────────────────────────────────── */

type BudgetCardProps = {
  budget: {
    id: string;
    category: { name: string; type?: string };
    amount: number;
    spent: number;
    remaining: number;
    percent: number;
  };
  onEdit: () => void;
  onDelete: () => void;
};

function BudgetCard({ budget, onEdit, onDelete }: BudgetCardProps) {
  const pct = Math.max(0, Math.min(100, budget.percent));
  const spentColor =
    pct >= 100
      ? "text-red-600"
      : pct >= 80
        ? "text-yellow-600"
        : "text-emerald-600";

  const barColor =
    pct >= 100 ? "bg-red-400" : pct >= 80 ? "bg-yellow-300" : "bg-emerald-400";

  return (
    <Card>
      <div className="flex items-start justify-between gap-2">
        <div className="min-w-0">
          <p className="truncate text-lg font-black uppercase text-slate-950 dark:text-slate-100">
            {budget.category.name}
          </p>
          {budget.category.type ? (
            <p className="text-xs font-bold uppercase text-slate-500">
              {budget.category.type}
            </p>
          ) : null}
        </div>
        <div className="flex shrink-0 gap-1">
          <button
            className="rounded-lg border-2 border-slate-950 bg-[#fffdf7] px-2 py-1 text-xs font-black uppercase shadow-[2px_2px_0_0_#101828] transition hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none dark:border-slate-100 dark:bg-slate-800 dark:text-slate-100 dark:shadow-[2px_2px_0_0_#f8fafc]"
            onClick={onEdit}
            aria-label="Edit budget"
          >
            ✏️
          </button>
          <button
            className="rounded-lg border-2 border-slate-950 bg-[#fffdf7] px-2 py-1 text-xs font-black uppercase shadow-[2px_2px_0_0_#101828] transition hover:translate-x-0.5 hover:translate-y-0.5 hover:shadow-none dark:border-slate-100 dark:bg-slate-800 dark:text-slate-100 dark:shadow-[2px_2px_0_0_#f8fafc]"
            onClick={onDelete}
            aria-label="Delete budget"
          >
            🗑️
          </button>
        </div>
      </div>

      <div className="mt-4 space-y-2">
        <div className="flex justify-between text-sm">
          <span className="font-semibold text-slate-600 dark:text-slate-200">
            Budget
          </span>
          <span className="font-black text-slate-950 dark:text-slate-100">
            {formatIDR(budget.amount)}
          </span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="font-semibold text-slate-600 dark:text-slate-200">
            Spent
          </span>
          <span className={`font-black ${spentColor}`}>
            {formatIDR(budget.spent)}
          </span>
        </div>
        <div className="flex justify-between text-sm">
          <span className="font-semibold text-slate-600 dark:text-slate-200">
            Remaining
          </span>
          <span
            className={`font-black ${
              budget.remaining < 0
                ? "text-red-600"
                : "text-slate-950 dark:text-slate-100"
            }`}
          >
            {formatIDR(budget.remaining)}
          </span>
        </div>
      </div>

      {/* Progress bar */}
      <div className="mt-4">
        <div className="flex items-center gap-2">
          <div className="h-4 flex-1 rounded-full border-2 border-slate-950 bg-[#fffdf7] dark:border-slate-100 dark:bg-slate-900">
            <div
              className={`h-full rounded-full transition-all ${barColor}`}
              style={{ width: `${pct}%` }}
            />
          </div>
          <span className="text-xs font-black tabular-nums text-slate-600 dark:text-slate-200">
            {Math.round(pct)}%
          </span>
        </div>
      </div>
    </Card>
  );
}
