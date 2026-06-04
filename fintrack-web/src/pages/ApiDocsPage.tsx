import { ApiReferenceReact } from "@scalar/api-reference-react";
import "@scalar/api-reference-react/style.css";

export function ApiDocsPage() {
  return (
    <div className="min-h-screen bg-white text-slate-950 dark:bg-slate-950 dark:text-slate-100">
      <ApiReferenceReact
        configuration={{
          url: "/openapi.yaml",
          theme: "default",
          layout: "modern",
          hideDarkModeToggle: true,
        }}
      />
    </div>
  );
}
