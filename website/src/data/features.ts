import type { Feature } from "./types";

export const features: Feature[] = [
  {
    icon: "shield",
    title: "Type-Safe Props",
    desc: "37 typed string enums make invalid states unrepresentable. Every props struct embeds BaseProps for consistent ID, class, ARIA, and CSP nonce propagation.",
  },
  {
    icon: "lightning",
    title: "Zero Node.js",
    desc: "Pure Go + templ + Tailwind CSS v4. No build pipeline beyond templ generate. No npm, no bundlers, no SPA framework. CSS-first config, class-based dark mode.",
  },
  {
    icon: "code",
    title: "Server-Rendered",
    desc: "Every component renders HTML on the server. HATEOAS-aligned: JavaScript enhances rather than replaces HTML. Interactive features use minimal vanilla JS with CSP nonces.",
  },
  {
    icon: "moon",
    title: "Built-in Dark Mode",
    desc: "Every component has dark: variants for all neutral and semantic colors, enforced by regression tests. ThemeScript prevents FOUC. color-scheme for native form controls.",
  },
  {
    icon: "check",
    title: "CSP-Ready",
    desc: "All inline scripts use nonce attributes. No eval(), no inline event handlers. Integration test suite verifies nonce compliance on every inline script across all components.",
  },
  {
    icon: "bolt",
    title: "HTMX Integration",
    desc: "Dedicated htmx package with loading indicators, error handling, CSRF protection, out-of-band swaps, and View Transitions. First-class server-rendered HTMX patterns.",
  },
];
