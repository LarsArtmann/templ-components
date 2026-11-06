module github.com/larsartmann/templ-components/errorpage

go 1.27

require (
	github.com/a-h/templ v0.3.1100
	github.com/larsartmann/go-error-family v0.12.0
	github.com/larsartmann/templ-components/icons v1.27.0
	github.com/larsartmann/templ-components/utils v1.27.0
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/stretchr/testify v1.12.1 // indirect
	go.yaml.in/yaml/v3 v3.0.6 // indirect
)

replace github.com/larsartmann/templ-components/utils => ../utils

replace github.com/larsartmann/templ-components/icons => ../icons
