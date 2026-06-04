import { NavLink, Outlet, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../stores/authStore";
import { useTheme } from "../../stores/themeStore";

const navItems = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/accounts", label: "Accounts" },
  { to: "/transactions", label: "Transactions" },
  { to: "/reports", label: "Reports" },
];

function navClass({ isActive }: { isActive: boolean }) {
  return `rounded-xl border-2 border-slate-950 px-3 py-2 text-sm font-black uppercase tracking-wide transition dark:border-slate-100 ${
    isActive
      ? "bg-blue-300 text-slate-950 shadow-[4px_4px_0_0_#101828] dark:shadow-[4px_4px_0_0_#f8fafc]"
      : "bg-[#fffdf7] text-slate-700 hover:-translate-y-0.5 hover:bg-emerald-100 hover:shadow-[3px_3px_0_0_#101828] dark:bg-slate-800 dark:text-slate-100 dark:hover:bg-slate-700 dark:hover:shadow-[3px_3px_0_0_#f8fafc]"
  }`;
}

export function AppLayout() {
  const logout = useAuthStore((state) => state.logout);
  const navigate = useNavigate();
  const { isDark, toggle } = useTheme();

  function handleLogout() {
    logout();
    navigate("/login");
  }

  return (
    <div className="min-h-screen text-slate-950 dark:text-slate-100">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r-2 border-slate-950 bg-[#f7f1e3] p-5 dark:border-slate-100 dark:bg-slate-900 lg:block">
        <div className="mb-8 rounded-xl border-2 border-slate-950 bg-emerald-200 p-4 text-slate-950 shadow-[5px_5px_0_0_#101828] dark:border-slate-100 dark:bg-blue-300 dark:text-slate-950 dark:shadow-[5px_5px_0_0_#f8fafc]">
          <p className="text-2xl font-black uppercase">Fintrack</p>
          <p className="text-sm font-bold">Private finance tracker</p>
        </div>
        <nav className="flex flex-col gap-3">
          {navItems.map((item) => (
            <NavLink key={item.to} to={item.to} className={navClass}>
              {item.label}
            </NavLink>
          ))}
        </nav>
        <div className="absolute bottom-5 left-5 right-5 flex flex-col gap-3">
          <button
            className="neo-button bg-blue-300 dark:text-slate-950"
            onClick={toggle}
            aria-label={isDark ? "Switch to light mode" : "Switch to dark mode"}
            aria-pressed={isDark}
          >
            {isDark ? "☀️ Light" : "🌙 Dark"}
          </button>
          <button
            className="neo-button bg-red-300 dark:text-slate-950"
            onClick={handleLogout}
          >
            Logout
          </button>
        </div>
      </aside>

      <header className="sticky top-0 z-10 border-b-2 border-slate-950 bg-[#f7f1e3] px-3 py-3 dark:border-slate-100 dark:bg-slate-900 sm:px-4 lg:hidden">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-lg font-black uppercase dark:text-slate-100">
              Fintrack
            </p>
            <p className="text-xs font-bold text-slate-600 dark:text-slate-200">
              Personal finance
            </p>
          </div>
          <div className="flex shrink-0 gap-2">
            <button
              className="neo-button bg-blue-300 px-3 py-1 dark:text-slate-950"
              onClick={toggle}
              aria-label={
                isDark ? "Switch to light mode" : "Switch to dark mode"
              }
              aria-pressed={isDark}
            >
              {isDark ? "☀️" : "🌙"}
            </button>
            <button
              className="neo-button bg-red-300 px-3 py-1 dark:text-slate-950"
              onClick={handleLogout}
              aria-label="Logout"
            >
              <span className="sm:hidden" aria-hidden="true">
                ⏻
              </span>
              <span className="hidden sm:inline">Logout</span>
            </button>
          </div>
        </div>
      </header>

      <main className="pb-[calc(5rem+env(safe-area-inset-bottom))] lg:ml-64 lg:pb-0">
        <div className="mx-auto max-w-7xl px-3 py-4 sm:px-6 sm:py-6 lg:px-8">
          <Outlet />
        </div>
      </main>

      <nav className="fixed inset-x-0 bottom-0 z-20 grid grid-cols-4 border-t-2 border-slate-950 bg-[#f7f1e3] pb-[env(safe-area-inset-bottom)] dark:border-slate-100 dark:bg-slate-900 lg:hidden">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `min-w-0 truncate px-1 py-3 text-center text-[0.65rem] font-black uppercase tracking-tight sm:text-xs sm:tracking-wide ${
                isActive
                  ? "bg-blue-300 text-slate-950"
                  : "text-slate-700 dark:text-slate-100"
              }`
            }
          >
            {item.label}
          </NavLink>
        ))}
      </nav>
    </div>
  );
}
