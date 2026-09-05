module github.com/larsartmann/templ-components/visualtest

go 1.26.7

require (
	github.com/a-h/templ v0.3.1020
	github.com/chromedp/cdproto v0.0.0-20260804232424-e85f50dbfd32
	github.com/chromedp/chromedp v0.16.0
	github.com/larsartmann/go-datastar/static v0.5.0
	github.com/larsartmann/templ-components v0.0.0-00010101000000-000000000000
	github.com/larsartmann/templ-components/datastar v1.13.0
	github.com/larsartmann/templ-components/errorpage v1.13.0
	github.com/larsartmann/templ-components/htmx v1.13.0
	github.com/larsartmann/templ-components/icons v1.13.0
	github.com/larsartmann/templ-components/utils v1.13.0
	github.com/orisano/pixelmatch v0.0.0-20230914042517-fa304d1dc785
)

require (
	github.com/Oudwins/tailwind-merge-go v0.2.3 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/go-json-experiment/json v0.0.0-20260820222146-c27c302e5fc3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/larsartmann/go-error-family v0.10.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)

replace github.com/larsartmann/templ-components => ..

replace github.com/larsartmann/templ-components/datastar => ../datastar

replace github.com/larsartmann/templ-components/errorpage => ../errorpage

replace github.com/larsartmann/templ-components/htmx => ../htmx

replace github.com/larsartmann/templ-components/icons => ../icons

replace github.com/larsartmann/templ-components/utils => ../utils
