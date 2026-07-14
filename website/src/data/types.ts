export const featureIconKeys = ["shield", "lightning", "code", "moon", "check", "bolt"] as const;
export type FeatureIcon = (typeof featureIconKeys)[number];

export interface Feature {
  icon: FeatureIcon;
  title: string;
  desc: string;
}

export interface StepCard {
  step: string;
  stepColor: "accent" | "amber";
  title: string;
  desc: string;
  code?: string;
}

export type ComparisonVariant = "templUI" | "goshipit" | "templ-components";

export interface ComparisonItem {
  variant: ComparisonVariant;
  pros: string[];
  cons: string[];
  accent: boolean;
}

export type MatrixValue = "yes" | "no" | string;

export interface MatrixRow {
  feature: string;
  values: [MatrixValue, MatrixValue, MatrixValue];
}

export interface ComparisonMatrix {
  columns: [ComparisonVariant, ComparisonVariant, ComparisonVariant];
  rows: MatrixRow[];
}

export const useCaseIconKeys = ["grid", "form", "nav", "chart", "database"] as const;
export type UseCaseIcon = (typeof useCaseIconKeys)[number];

export interface UseCase {
  title: string;
  desc: string;
  icon: UseCaseIcon;
}

export const uiIconKeys = [
  "arrow-external",
  "arrow-right",
  "github",
  "menu",
  "close",
  "sun",
  "moon",
  "star",
  "eye",
] as const;
export type UIIcon = (typeof uiIconKeys)[number];

export type IconName = FeatureIcon | UseCaseIcon | UIIcon;
