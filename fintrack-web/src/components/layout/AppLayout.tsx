import { useEffect, useState } from "react";
import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { useAuthStore } from "../../stores/authStore";
import { useTheme } from "../../stores/themeStore";
import { usePWAInstall } from "../../hooks/usePWAInstall";

const navItems = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/accounts", label: "Accounts" },
  { to: "/markets", label: "Markets" },
  { to: "/transactions", label: "Transactions" },
  { to: "/budgets", label: "Budgets" },
  { to: "/reports", label: "Reports" },
];

const mobileNavItems = [
  { to: "/dashboard", label: "Home", icon: "🏠" },
  { to: "/accounts", label: "Account", icon: "💳" },
  { to: "/transactions", label: "Tx", icon: "➕" },
  { to: "/reports", label: "Reports", icon: "📊" },
];

const mobileMoreItems = [
  { to: "/markets", label: "Markets", icon: "📈" },
  { to: "/budgets", label: "Budgets", icon: "🎯" },
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
  const { pathname } = useLocation();
  const { isDark, toggle } = useTheme();
  const { canInstall, install } = usePWAInstall();
  const [isMoreOpen, setIsMoreOpen] = useState(false);
  const isMoreActive = mobileMoreItems.some((item) =>
    pathname.startsWith(item.to),
  );

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname]);

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
          {canInstall && (
            <button
              className="neo-button bg-emerald-300 dark:text-slate-950"
              onClick={install}
            >
              Install App
            </button>
          )}
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
            {canInstall && (
              <button
                className="neo-button bg-emerald-300 px-3 py-1 dark:text-slate-950"
                onClick={install}
                aria-label="Install app"
              >
                ⬇
              </button>
            )}
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

      {isMoreOpen ? (
        <button
          type="button"
          className="fixed inset-0 z-20 bg-transparent lg:hidden"
          aria-label="Close more menu"
          onClick={() => setIsMoreOpen(false)}
        />
      ) : null}

      {isMoreOpen ? (
        <div
          id="mobile-more-menu"
          className="fixed inset-x-3 bottom-[calc(6rem+env(safe-area-inset-bottom))] z-30 rounded-2xl border-2 border-slate-950 bg-[#f7f1e3] p-3 shadow-[6px_6px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[6px_6px_0_0_#f8fafc] lg:hidden"
        >
          <p className="mb-2 px-2 text-xs font-black uppercase tracking-[0.16em] text-slate-500 dark:text-slate-300">
            More
          </p>
          <div className="grid gap-2">
            {mobileMoreItems.map((item) => (
              <NavLink
                key={item.to}
                to={item.to}
                onClick={() => setIsMoreOpen(false)}
                className={({ isActive }) =>
                  `flex items-center justify-between rounded-xl border-2 border-slate-950 px-3 py-2 text-sm font-black uppercase dark:border-slate-100 ${
                    isActive
                      ? "bg-blue-300 text-slate-950 shadow-[3px_3px_0_0_#101828] dark:shadow-[3px_3px_0_0_#f8fafc]"
                      : "bg-[#fffdf7] text-slate-700 dark:bg-slate-800 dark:text-slate-100"
                  }`
                }
              >
                <span>{item.label}</span>
                <span aria-hidden="true">{item.icon}</span>
              </NavLink>
            ))}
          </div>
        </div>
      ) : null}

      <nav className="fixed inset-x-2 bottom-[calc(0.75rem+env(safe-area-inset-bottom))] z-40 grid grid-cols-5 gap-1 rounded-2xl border-2 border-slate-950 bg-[#f7f1e3] p-1.5 shadow-[6px_6px_0_0_#101828] dark:border-slate-100 dark:bg-slate-900 dark:shadow-[6px_6px_0_0_#f8fafc] lg:hidden">
        {mobileNavItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            onClick={() => setIsMoreOpen(false)}
            className={({ isActive }) =>
              `flex min-w-0 flex-col items-center gap-1 rounded-xl px-0.5 py-2 text-center text-[0.58rem] font-black uppercase tracking-[-0.02em] transition min-[380px]:text-[0.62rem] sm:text-xs ${
                isActive
                  ? "bg-blue-300 text-slate-950"
                  : "text-slate-700 dark:text-slate-100"
              }`
            }
          >
            <span className="text-base leading-none" aria-hidden="true">
              {item.icon}
            </span>
            <span className="w-full truncate">{item.label}</span>
          </NavLink>
        ))}
        <button
          type="button"
          className={`flex min-w-0 flex-col items-center gap-1 rounded-xl px-0.5 py-2 text-center text-[0.58rem] font-black uppercase tracking-[-0.02em] transition min-[380px]:text-[0.62rem] sm:text-xs ${
            isMoreOpen || isMoreActive
              ? "bg-blue-300 text-slate-950"
              : "text-slate-700 dark:text-slate-100"
          }`}
          aria-controls="mobile-more-menu"
          aria-expanded={isMoreOpen}
          onClick={() => setIsMoreOpen((open) => !open)}
        >
          <span className="text-base leading-none" aria-hidden="true">
            ☰
          </span>
          <span className="w-full truncate">More</span>
        </button>
      </nav>
    </div>
  );
}
