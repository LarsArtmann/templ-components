module github.com/larsartmann/templ-components/errorpage

go 1.26.5

require (
	github.com/a-h/templ v0.3.1020
	github.com/larsartmann/go-error-family v0.10.0
	github.com/larsartmann/templ-components/icons v0.0.0
	github.com/larsartmann/templ-components/utils v0.0.0
)

require github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect

replace github.com/larsartmann/templ-components/icons => ../icons

replace github.com/larsartmann/templ-components/utils => ../utils
