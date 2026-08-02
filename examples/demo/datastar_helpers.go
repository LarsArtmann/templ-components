package main

import "github.com/larsartmann/templ-components/utils"

// demoBaseProps returns BaseProps pre-configured with the demo CSP nonce.
func demoBaseProps() utils.BaseProps {
	return utils.BaseProps{Nonce: "demo-nonce"}
}
