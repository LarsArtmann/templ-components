module github.com/larsartmann/templ-components

go 1.26.0

require github.com/a-h/templ v0.3.1020

require (
	github.com/larsartmann/templ-components/charts/echarts v1.19.3
	github.com/larsartmann/templ-components/datastar v1.19.3
	github.com/larsartmann/templ-components/errorpage v1.19.3
	github.com/larsartmann/templ-components/htmx v1.19.3
	github.com/larsartmann/templ-components/icons v1.19.3
	github.com/larsartmann/templ-components/utils v1.19.3
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/alecthomas/chroma/v2 v2.27.0 // indirect
	github.com/chromedp/cdproto v0.0.0-20260922220944-a19bff23514f // indirect
	github.com/chromedp/chromedp v0.16.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/larsartmann/go-datastar/static v0.6.0 // indirect
	github.com/larsartmann/go-error-family v0.10.2 // indirect
	github.com/larsartmann/templ-components/visualtest v0.0.0-20260924162719-86d2a43421e6 // indirect
	github.com/larsartmann/templ-components/website v0.0.0-20260924162719-86d2a43421e6 // indirect
	github.com/orisano/pixelmatch v0.0.0-20230914042517-fa304d1dc785 // indirect
	github.com/yuin/goldmark v1.8.6 // indirect
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20230729083705-37449abec8cc // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/larsartmann/templ-components => ./

replace github.com/larsartmann/templ-components/utils => ./utils

replace github.com/larsartmann/templ-components/icons => ./icons

replace github.com/larsartmann/templ-components/errorpage => ./errorpage

replace github.com/larsartmann/templ-components/charts/echarts => ./charts/echarts

replace github.com/larsartmann/templ-components/datastar => ./datastar

replace github.com/larsartmann/templ-components/htmx => ./htmx
