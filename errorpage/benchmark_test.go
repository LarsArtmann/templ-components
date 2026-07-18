package errorpage

import (
	"bytes"
	"context"
	"testing"
)

func BenchmarkErrorpageRenders(b *testing.B) {
	benchmarks := []struct {
		name   string
		render func() string
	}{
		{"ErrorPage", func() string {
			var buf bytes.Buffer

			_ = ErrorPage(NotFound()).Render(context.Background(), &buf)

			return buf.String()
		}},
		{"NotFound404", func() string {
			props := DefaultNotFound404Props()
			props.Links = DefaultNotFoundLinks()

			var buf bytes.Buffer

			_ = NotFound404(props).Render(context.Background(), &buf)

			return buf.String()
		}},
		{"ErrorDetail", func() string {
			var buf bytes.Buffer

			_ = ErrorDetail(ErrorDetailProps{
				Family: FamilyCorruption, Code: "test", Title: "Test", Message: "msg",
			}).Render(context.Background(), &buf)

			return buf.String()
		}},
		{"ErrorAlert", func() string {
			var buf bytes.Buffer

			_ = ErrorAlert(ErrorAlertProps{
				Family: FamilyRejection, Title: "Test", Message: "msg",
			}).Render(context.Background(), &buf)

			return buf.String()
		}},
	}
	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for range b.N {
				_ = bm.render()
			}
		})
	}
}
