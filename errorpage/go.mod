module github.com/larsartmann/templ-components/errorpage

go 1.26

require (
	github.com/a-h/templ v0.3.1020
	github.com/larsartmann/go-error-family v0.10.2
	github.com/larsartmann/templ-components/icons v1.19.4
	github.com/larsartmann/templ-components/utils v1.19.4
)

require github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect

replace github.com/larsartmann/templ-components/utils => ../utils

replace github.com/larsartmann/templ-components/icons => ../icons
