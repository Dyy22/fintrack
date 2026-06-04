export default {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        border: "#101828",
        main: "#2563eb",
        "main-foreground": "#ffffff",
        "secondary-background": "#fffdf7",
        foreground: "#101828",
        "neo-cream": "#f7f1e3",
        "neo-blue": "#60a5fa",
        "neo-mint": "#86efac",
        "neo-yellow": "#fde68a",
        "neo-coral": "#fca5a5",
      },
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
        heading: ["Inter", "system-ui", "sans-serif"],
      },
      fontWeight: {
        base: "600",
        heading: "900",
      },
      borderRadius: {
        base: "0.75rem",
      },
      boxShadow: {
        shadow: "5px 5px 0 0 #101828",
      },
      translate: {
        boxShadowX: "5px",
        boxShadowY: "5px",
        reverseBoxShadowX: "-5px",
        reverseBoxShadowY: "-5px",
      },
    },
  },
  plugins: [],
};
