module github.com/larsartmann/templ-components/website

go 1.26

ignore (
	dist
	node_modules
)

require (
	github.com/a-h/templ v0.3.1020
	github.com/alecthomas/chroma/v2 v2.27.0
	github.com/larsartmann/templ-components v1.17.0
	github.com/larsartmann/templ-components/errorpage v1.17.0
	github.com/larsartmann/templ-components/icons v1.17.0
	github.com/larsartmann/templ-components/utils v1.17.0
	github.com/yuin/goldmark v1.8.6
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/dlclark/regexp2/v2 v2.8.0 // indirect
	github.com/larsartmann/go-error-family v0.10.1 // indirect
	github.com/larsartmann/templ-components/htmx v1.17.0 // indirect
)

replace github.com/larsartmann/templ-components => ../

replace github.com/larsartmann/templ-components/charts/echarts => ../charts/echarts

replace github.com/larsartmann/templ-components/datastar => ../datastar

replace github.com/larsartmann/templ-components/errorpage => ../errorpage

replace github.com/larsartmann/templ-components/htmx => ../htmx

replace github.com/larsartmann/templ-components/icons => ../icons

replace github.com/larsartmann/templ-components/utils => ../utils
