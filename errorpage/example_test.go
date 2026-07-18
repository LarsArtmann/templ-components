package errorpage_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/errorpage"
)

func ExampleErrorPage() {
	var buf bytes.Buffer

	_ = errorpage.ErrorPage(errorpage.ErrorPageProps{
		Family:     errorpage.FamilyRejection,
		Code:       "page.not_found",
		Title:      "Page not found",
		Message:    "The page you requested does not exist.",
		Fix:        "Check the URL or navigate back to the homepage.",
		WayOut:     "Go home",
		WayOutHref: "/",
	}).Render(context.Background(), &buf)

	fmt.Println("renders full-page error view")
	// Output: renders full-page error view
}

func ExampleErrorDetail() {
	var buf bytes.Buffer

	_ = errorpage.ErrorDetail(errorpage.ErrorDetailProps{
		Family:  errorpage.FamilyCorruption,
		Code:    "data.parse_failed",
		Title:   "Parse Failed",
		Message: "config.yaml has invalid syntax.",
		Fix:     "Check YAML indentation.",
		Context: []errorpage.ContextPair{
			{Key: "file", Value: "config.yaml"},
			{Key: "line", Value: "42"},
		},
	}).Render(context.Background(), &buf)

	fmt.Println("renders inline error detail card")
	// Output: renders inline error detail card
}

func ExampleErrorAlert() {
	var buf bytes.Buffer

	_ = errorpage.ErrorAlert(errorpage.ErrorAlertProps{
		Family:  errorpage.FamilyTransient,
		Title:   "Temporary Error",
		Message: "Please try again shortly.",
		Fix:     "Wait a moment and retry.",
	}).Render(context.Background(), &buf)

	fmt.Println("renders family-aware alert banner")
	// Output: renders family-aware alert banner
}

func ExampleFamilyStatusCode() {
	fmt.Println(errorpage.FamilyStatusCode(errorpage.FamilyRejection))
	fmt.Println(errorpage.FamilyStatusCode(errorpage.FamilyConflict))
	fmt.Println(errorpage.FamilyStatusCode(errorpage.FamilyTransient))
	fmt.Println(errorpage.FamilyStatusCode(errorpage.FamilyCorruption))
	fmt.Println(errorpage.FamilyStatusCode(errorpage.FamilyInfrastructure))
	// Output:
	// 400
	// 409
	// 503
	// 500
	// 503
}

func ExampleContextMap() {
	ctx := errorpage.ContextMap(map[string]string{
		"host": "db.internal",
		"port": "5432",
	})
	fmt.Println(len(ctx))
	// Output: 2
}

func ExampleExtractCauseChain() {
	chain := errorpage.ExtractCauseChain(&outerError{}, 10)
	fmt.Println(len(chain))
	// Output: 2
}

type innerError struct{}

func (e *innerError) Error() string { return "connection refused" }

type outerError struct{}

func (e *outerError) Error() string { return "database unavailable" }
func (e *outerError) Unwrap() error { return &middleError{} }

type middleError struct{}

func (e *middleError) Error() string { return "connection pool exhausted" }
func (e *middleError) Unwrap() error { return &innerError{} }
