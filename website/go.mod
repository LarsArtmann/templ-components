module github.com/larsartmann/templ-components/website

go 1.26.7

ignore (
	node_modules
	dist
)

require (
	github.com/a-h/templ v0.3.1020
	github.com/alecthomas/chroma/v2 v2.27.0
	github.com/larsartmann/templ-components v0.0.0-00010101000000-000000000000
	github.com/larsartmann/templ-components/icons v1.16.0
	github.com/larsartmann/templ-components/utils v1.16.0
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/dlclark/regexp2/v2 v2.2.1 // indirect
	github.com/larsartmann/templ-components/htmx v1.16.0 // indirect
)

replace github.com/larsartmann/templ-components => ../

replace github.com/larsartmann/templ-components/charts/echarts => ../charts/echarts

replace github.com/larsartmann/templ-components/datastar => ../datastar

replace github.com/larsartmann/templ-components/errorpage => ../errorpage

replace github.com/larsartmann/templ-components/htmx => ../htmx

replace github.com/larsartmann/templ-components/icons => ../icons

replace github.com/larsartmann/templ-components/utils => ../utils
