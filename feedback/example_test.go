package feedback_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/feedback"
)

func ExampleAlert() {
	props := feedback.DefaultAlertProps()
	props.Type = feedback.FeedbackSuccess
	props.Title = "Success"
	props.Message = "Your changes have been saved."

	var buf bytes.Buffer

	_ = feedback.Alert(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleToast() {
	props := feedback.DefaultToastProps()
	props.Type = feedback.FeedbackInfo
	props.Title = "Notification"
	props.Message = "You have a new message."

	var buf bytes.Buffer

	_ = feedback.Toast(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleSpinner() {
	var buf bytes.Buffer

	_ = feedback.Spinner(feedback.SpinnerProps{Size: feedback.SpinnerMD, Color: "text-blue-600"}).
		Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleSkeletonCardGrid() {
	var buf bytes.Buffer

	_ = feedback.SkeletonCardGrid(feedback.SkeletonCardGridProps{Count: 6}).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain a responsive grid of skeleton loading cards
}

func ExampleCircularProgress() {
	var buf bytes.Buffer

	_ = feedback.CircularProgress(feedback.CircularProgressProps{
		Value:     75,
		ShowLabel: true,
	}).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain a circular SVG progress ring showing 75%
}
