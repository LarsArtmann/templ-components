package main

import "github.com/larsartmann/templ-components/display"

// demoStatusBadgeMapper is the StatusBadgeWith demo mapper: it extends the
// built-in status mapping with a demo-only provisioning state and delegates
// everything else to the built-in map (unknown statuses render neutral).
func demoStatusBadgeMapper(status string) display.BadgeType {
	if status == "provisioning" {
		return display.BadgeInfo
	}

	return display.MapStatusToBadgeType(status)
}
