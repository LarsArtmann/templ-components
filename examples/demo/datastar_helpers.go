package main

import "github.com/larsartmann/templ-components/utils"

// demoBaseProps returns BaseProps pre-configured with the demo CSP nonce.
func demoBaseProps() utils.BaseProps {
	return utils.BaseProps{Nonce: "demo-nonce"}
}

// datastarPolledStatsBaseProps returns the PolledRegion demo's BaseProps: the
// fixed id is what both the tick endpoint's fragment root and the patch
// selector target, so the region self-patches; the AriaLabel gives the
// region landmark its accessible name.
func datastarPolledStatsBaseProps() utils.BaseProps {
	props := demoBaseProps()
	props.ID = "datastar-polled-stats"
	props.AriaLabel = "Live demo statistics"

	return props
}
