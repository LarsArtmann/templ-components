package forms

import (
	"testing"

	"github.com/larsartmann/templ-components/utils"
	"github.com/larsartmann/templ-components/utils/wire"
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

func TestFormNoValidate(t *testing.T) {
	t.Parallel()

	t.Run("renders novalidate when set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", NoValidate: true}))
		utils.AssertContains(t, output, "novalidate")
	})

	t.Run("absent by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save"}))
		utils.AssertNotContains(t, output, "novalidate")
	})

	t.Run("composes with hx-validate", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", Validate: true, NoValidate: true}))
		utils.AssertContains(t, output, `hx-validate="true"`)
		utils.AssertContains(t, output, "novalidate")
	})

	t.Run("composes with datastar wire", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			Wire:       &wire.Action{Transport: wire.TransportDatastar, Method: wire.MethodPost, URL: "/api/save"},
			NoValidate: true,
		}))
		utils.AssertContains(t, output, "novalidate")
		utils.AssertContains(t, output, `data-on:submit`)
	})
}

func TestFormContainerAware(t *testing.T) {
	t.Parallel()

	t.Run("viewport breakpoints by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", Layout: FormLayoutGrid}))
		utils.AssertNotContains(t, output, "@container")
		utils.AssertContains(t, output, "sm:grid-cols-")
		utils.AssertNotContains(t, output, "@sm:")
	})

	t.Run("container breakpoints when flag set", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			Action:         "/save",
			Layout:         FormLayoutGrid,
			ContainerAware: true,
		}))
		utils.AssertContains(t, output, "@container")
		utils.AssertContains(t, output, "@sm:grid-cols-")
		utils.AssertNotContains(t, output, " sm:")
	})
}

func TestFormEnctype(t *testing.T) {
	t.Parallel()

	t.Run("multipart renders enctype", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/upload", Enctype: FormEnctypeMultipart}))
		utils.AssertContains(t, output, `enctype="multipart/form-data"`)
	})

	t.Run("default omits enctype (HTML default applies)", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save"}))
		utils.AssertNotContains(t, output, "enctype=")
	})

	t.Run("IsValid accepts zero value and known values", func(t *testing.T) {
		t.Parallel()
		if !FormEnctypeIsValid("") {
			t.Error(`FormEnctypeIsValid("") = false, want true`)
		}
		if !FormEnctypeIsValid(FormEnctypeUrlencoded) || !FormEnctypeIsValid(FormEnctypeMultipart) {
			t.Error("FormEnctypeIsValid rejected known values")
		}
		if FormEnctypeIsValid(FormEnctype("text/plain")) {
			t.Error("FormEnctypeIsValid accepted unknown value")
		}
	})

	t.Run("multipart composes with wire for uploads", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{
			Method:  FormPost,
			Enctype: FormEnctypeMultipart,
			Wire: &wire.Action{
				Transport: wire.TransportDatastar,
				Method:    wire.MethodPost,
				URL:       "/api/upload",
			},
		}))
		utils.AssertContains(t, output, `enctype="multipart/form-data"`)
		utils.AssertContains(t, output, `data-on:submit=`)
	})
}

func TestFormDirtyGuard(t *testing.T) {
	t.Parallel()

	t.Run("renders the opt-in attribute", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save", DirtyGuard: true}))
		utils.AssertContains(t, output, `data-tc-dirty-guard="true"`)
	})

	t.Run("absent by default", func(t *testing.T) {
		t.Parallel()
		output := utils.Render(t, Form(FormProps{Action: "/save"}))
		utils.AssertNotContains(t, output, "data-tc-dirty-guard")
	})
}
