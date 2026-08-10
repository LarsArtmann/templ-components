module github.com/larsartmann/templ-components

go 1.26.5

require github.com/a-h/templ v0.3.1020

require (
	github.com/larsartmann/templ-components/charts/echarts v0.0.0
	github.com/larsartmann/templ-components/errorpage v0.0.0
	github.com/larsartmann/templ-components/icons v0.0.0
	github.com/larsartmann/templ-components/utils v0.0.0
	github.com/stretchr/testify v1.10.0
)

replace github.com/larsartmann/templ-components/charts/echarts => ./charts/echarts

replace github.com/larsartmann/templ-components/errorpage => ./errorpage

replace github.com/larsartmann/templ-components/icons => ./icons

replace github.com/larsartmann/templ-components/utils => ./utils

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/larsartmann/go-error-family v0.10.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
