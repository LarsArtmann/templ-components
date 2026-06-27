package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
)

func TestFormMethod(t *testing.T) {
	t.Parallel()

	t.Run("post passes through", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", Method: FormPost}))
		utils.AssertContains(t, output, `method="POST"`)
	})

	t.Run("get passes through", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/search", Method: FormGet}))
		utils.AssertContains(t, output, `method="GET"`)
	})

	t.Run("empty falls back to GET", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x"}))
		utils.AssertContains(t, output, `method="GET"`)
	})

	t.Run("invalid method falls back to GET", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/x", Method: FormMethod("DELETE")}))
		utils.AssertContains(t, output, `method="GET"`)
		utils.AssertNotContains(t, output, `DELETE`)
	})

	t.Run("csrf token renders hidden input", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", Method: FormPost, CSRFToken: "abc"}))
		utils.AssertContains(t, output, `name="csrf_token"`)
		utils.AssertContains(t, output, `value="abc"`)
	})
}
