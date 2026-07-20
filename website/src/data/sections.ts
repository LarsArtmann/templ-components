import type { StepCard, ComparisonItem, UseCase, ComparisonMatrix } from "./types";

export const steps: StepCard[] = [
  {
    step: "1",
    stepColor: "accent",
    title: "Import",
    desc: "Pick only the packages you need. No monolithic bundle — pay for what you use.",
    code: `import "github.com/larsartmann/templ-components/display"`,
  },
  {
    step: "2",
    stepColor: "accent",
    title: "Render",
    desc: "Pass typed props structs. Every invalid state is a compile-time error, not a runtime check.",
    code: `@display.Card(display.CardProps{Title: "Hello"}) { ... }`,
  },
  {
    step: "3",
    stepColor: "amber",
    title: "Generate",
    desc: "templ generate compiles .templ files to Go. The generated code is committed for library consumers.",
    code: "templ generate ./...",
  },
  {
    step: "4",
    stepColor: "amber",
    title: "Ship",
    desc: "Server renders HTML. HTMX enhances where needed. No client-side framework required.",
    code: "// HTML over the wire — done",
  },
];

export const comparisons: ComparisonItem[] = [
  {
    variant: "templUI",
    accent: false,
    pros: ["40+ components", "Active community", "Alpine.js interactivity"],
    cons: ["Alpine.js dependency", "CSS custom properties (not Tailwind v4)"],
  },
  {
    variant: "templ-components",
    accent: true,
    pros: [
      "98 components across 9 packages",
      "43 typed string enums",
      "Tailwind v4 CSS-first config",
      "CSP nonce on every inline script",
      "Built-in HTMX integration package",
      "102 SVG icons (no icon library dep)",
    ],
    cons: [],
  },
  {
    variant: "goshipit",
    accent: false,
    pros: ["Full starter kit", "DaisyUI component set"],
    cons: ["Requires Node.js", "DaisyUI JS dependency", "Not a standalone library"],
  },
];

export const comparisonMatrix: ComparisonMatrix = {
  columns: ["templUI", "goshipit", "templ-components"],
  rows: [
    {
      feature: "CSS approach",
      values: ["Tailwind + vars", "Tailwind + DaisyUI", "Tailwind v4 CSS-first"],
    },
    { feature: "JavaScript", values: ["Alpine.js", "DaisyUI JS", "HATEOAS (enhances HTML)"] },
    { feature: "Requires Node.js", values: ["no", "yes", "no"] },
    { feature: "Typed props enums", values: ["no", "no", "34"] },
    { feature: "CSP nonce support", values: ["yes", "no", "yes"] },
    { feature: "Dark mode", values: ["CSS vars", "DaisyUI", "Tailwind dark: (tested)"] },
    { feature: "HTMX integration", values: ["no", "no", "yes"] },
    { feature: "Standalone library", values: ["no", "no", "yes"] },
  ],
};

export const useCases: UseCase[] = [
  {
    title: "Admin Dashboards",
    desc: "Tables, stat cards, charts, and sidebars. Build data-dense panels in minutes, not days.",
    icon: "grid",
  },
  {
    title: "Forms & CRUD",
    desc: "21 form components with validation, comboboxes, date pickers, and accessible error handling.",
    icon: "form",
  },
  {
    title: "SaaS Interfaces",
    desc: "Navigation, modals, drawers, tabs, and breadcrumbs. Full app chrome without a frontend framework.",
    icon: "nav",
  },
];
