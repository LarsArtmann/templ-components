module github.com/larsartmann/templ-components/website

go 1.26.7

require (
	github.com/a-h/templ v0.3.1020
	github.com/alecthomas/chroma/v2 v2.20.0
	github.com/larsartmann/templ-components v1.16.0
	github.com/larsartmann/templ-components/icons v1.16.0
	github.com/larsartmann/templ-components/utils v1.16.0
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/dlclark/regexp2 v1.11.5 // indirect
)

replace github.com/larsartmann/templ-components => ../

replace github.com/larsartmann/templ-components/charts/echarts => ../charts/echarts

replace github.com/larsartmann/templ-components/datastar => ../datastar

replace github.com/larsartmann/templ-components/errorpage => ../errorpage

replace github.com/larsartmann/templ-components/htmx => ../htmx

replace github.com/larsartmann/templ-components/icons => ../icons

replace github.com/larsartmann/templ-components/utils => ../utils
