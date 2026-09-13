module github.com/larsartmann/templ-components

go 1.26.7

require github.com/a-h/templ v0.3.1020

require (
	github.com/larsartmann/templ-components/charts/echarts v1.16.0
	github.com/larsartmann/templ-components/datastar v1.16.0
	github.com/larsartmann/templ-components/errorpage v1.16.0
	github.com/larsartmann/templ-components/htmx v1.16.0
	github.com/larsartmann/templ-components/icons v1.16.0
	github.com/larsartmann/templ-components/utils v1.16.0
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/larsartmann/go-datastar/static v0.5.0 // indirect
	github.com/larsartmann/go-error-family v0.10.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/larsartmann/templ-components => ./

replace github.com/larsartmann/templ-components/utils => ./utils

replace github.com/larsartmann/templ-components/icons => ./icons

replace github.com/larsartmann/templ-components/errorpage => ./errorpage

replace github.com/larsartmann/templ-components/charts/echarts => ./charts/echarts

replace github.com/larsartmann/templ-components/datastar => ./datastar

replace github.com/larsartmann/templ-components/htmx => ./htmx
