package feedback_test

import (
	"bytes"
	"context"
	"fmt"

	"github.com/larsartmann/templ-components/feedback"
)

func ExampleAlert() {
	props := feedback.DefaultAlertProps()
	props.Type = feedback.AlertSuccess
	props.Title = "Success"
	props.Message = "Your changes have been saved."

	var buf bytes.Buffer
	_ = feedback.Alert(props).Render(context.Background(), &buf)
	fmt.Println(buf.String())
}

func ExampleToast() {
	props := feedback.DefaultToastProps()
	props.Type = feedback.ToastInfo
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
	_ = feedback.SkeletonCardGrid(6).Render(context.Background(), &buf)
	fmt.Println(buf.String())
	// Output will contain a responsive grid of skeleton loading cards
}
