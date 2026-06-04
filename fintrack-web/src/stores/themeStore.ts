import { useSyncExternalStore } from "react";

const THEME_KEY = "fintrack_theme";

function getIsDark(): boolean {
  return document.documentElement.classList.contains("dark");
}

function subscribeToTheme(callback: () => void): () => void {
  const observer = new MutationObserver(() => {
    callback();
  });
  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ["class"],
  });
  return () => observer.disconnect();
}

function getSnapshot(): boolean {
  return getIsDark();
}

export function useTheme() {
  const isDark = useSyncExternalStore(subscribeToTheme, getSnapshot);

  function toggle() {
    const root = document.documentElement;
    const hasDark = root.classList.contains("dark");
    if (hasDark) {
      root.classList.remove("dark");
      localStorage.setItem(THEME_KEY, "light");
    } else {
      root.classList.add("dark");
      localStorage.setItem(THEME_KEY, "dark");
    }
  }

  return { isDark, toggle };
}
