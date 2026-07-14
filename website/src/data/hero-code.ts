import { siteConfig } from "./config";

const importPath = siteConfig.github.replace("https://github.com/", "github.com/");

export const heroCode = `package main

import (
    "${importPath}/display"
    "${importPath}/feedback"
    "${importPath}/layout"
)

templ Page() {
    @layout.Base(layout.DefaultPageProps()) {
        @layout.ThemeScript("")
        @display.PageHeader(display.PageHeaderProps{
            Title: "Dashboard",
            Subtitle: "Welcome back",
        })
        @display.Grid(display.GridProps{
            Cols: display.GridCols3,
        }) {
            @display.StatCard(display.StatCardProps{
                Label: "Revenue",
                Value: "$42,189",
                Trend: display.TrendUp,
            })
        }
        @feedback.Toast(feedback.ToastProps{
            Message: "Data refreshed!",
            Type:   feedback.ToastSuccess,
        })
    }
}`;
