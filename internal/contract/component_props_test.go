// Package contract holds cross-package test-only contracts that cannot live
// in any single package due to import-cycle restrictions. Tests in this
// package import the entire library and assert cross-cutting invariants.
package contract

import (
	"reflect"
	"testing"

	"github.com/larsartmann/templ-components/display"
	"github.com/larsartmann/templ-components/errorpage"
	"github.com/larsartmann/templ-components/feedback"
	"github.com/larsartmann/templ-components/forms"
	"github.com/larsartmann/templ-components/htmx"
	"github.com/larsartmann/templ-components/navigation"
	"github.com/larsartmann/templ-components/utils"
)

// componentTypes is the canonical inventory of every props struct in the
// library that should satisfy utils.ComponentProps. Maintained as a single
// source of truth so this contract test guards against accidental removal of
// the BaseProps embed across the codebase.
//
// If you add a new component, add its Props type here. The test will fail at
// CI time if the new type does not embed BaseProps, preventing silent
// interface contract breakage for consumers using generic wrappers.
func componentTypes() []any {
	return []any{
		// display (28)
		display.BadgeProps{},
		display.AvatarProps{},
		display.TooltipProps{},
		display.AccordionProps{},
		display.ButtonProps{},
		display.CardProps{},
		display.SimpleCardProps{},
		display.StatCardProps{},
		display.DropdownProps{},
		display.TabsProps{},
		display.TableProps{},
		display.PageHeaderProps{},
		display.ListNoteProps{},
		display.EmptyStateProps{},
		display.DefinitionListProps{},
		display.ModalProps{},
		display.DrawerProps{},
		display.GridProps{},
		display.CopyButtonProps{},
		display.RelativeTimeProps{},
		display.CountBadgeProps{},
		display.DefinitionGridProps{},
		display.ImageProps{},
		display.PopoverProps{},
		display.DataTableProps{},
		display.HoverCardProps{},
		display.ContextMenuProps{},
		display.CarouselProps{},

		// forms (18)
		forms.InputProps{},
		forms.CheckboxProps{},
		forms.SelectProps{},
		forms.TextareaProps{},
		forms.ToggleProps{},
		forms.FormProps{},
		forms.InputGroupProps{},
		forms.FileInputProps{},
		forms.DatePickerProps{},
		forms.ComboboxProps{},
		forms.ValidationSummaryProps{},
		forms.RadioProps{},
		forms.RadioGroupProps{},
		forms.FilterDropdownProps{},
		forms.SliderProps{},
		forms.RatingProps{},
		forms.TagsInputProps{},
		forms.CalendarProps{},

		// feedback (6)
		feedback.AlertProps{},
		feedback.SpinnerProps{},
		feedback.LoadingOverlayProps{},
		feedback.ToastProps{},
		feedback.ProgressBarProps{},
		feedback.StepIndicatorProps{},

		// navigation (8)
		navigation.NavProps{},
		navigation.SimpleNavProps{},
		navigation.NavLinkProps{},
		navigation.BreadcrumbsProps{},
		navigation.PaginationProps{},
		navigation.SidebarNavProps{},
		navigation.LoadMoreProps{},
		navigation.EndOfListProps{},

		// htmx (2)
		htmx.ConfirmDeleteProps{},
		htmx.SwapOOBProps{},

		// errorpage (4)
		errorpage.ErrorPageProps{},
		errorpage.NotFound404Props{},
		errorpage.ErrorDetailProps{},
		errorpage.ErrorAlertProps{},
	}
}

// TestAllComponentPropsSatisfyInterface walks the canonical inventory of
// component props structs and asserts each one satisfies the ComponentProps
// interface (which requires embedding BaseProps).
//
// This is a regression guard for the public contract documented in AGENTS.md:
// "All component props embed utils.BaseProps (exception: layout.PageProps) —
// all auto-satisfy utils.ComponentProps interface".
func TestAllComponentPropsSatisfyInterface(t *testing.T) {
	t.Parallel()

	componentInterface := reflect.TypeFor[utils.ComponentProps]()

	required := componentTypes()
	if len(required) == 0 {
		t.Fatal("componentTypes() returned an empty list; the inventory is broken")
	}

	for _, sample := range required {
		t.Run(reflect.TypeOf(sample).Name(), func(t *testing.T) {
			t.Parallel()

			ptrType := reflect.PointerTo(reflect.TypeOf(sample))

			if !ptrType.Implements(componentInterface) {
				t.Errorf(
					"*%s does not implement ComponentProps — must embed utils.BaseProps to expose GetBaseProps/SetBaseProps via method promotion",
					ptrType.String(),
				)
			}

			if !embedsBaseProps(reflect.TypeOf(sample)) {
				t.Errorf("%s does not embed utils.BaseProps — see AGENTS.md convention",
					reflect.TypeOf(sample).String())
			}
		})
	}
}

// TestComponentPropsInterfaceContract documents the exact method set the
// interface exposes. If the interface changes, this test forces an explicit
// review.
func TestComponentPropsInterfaceContract(t *testing.T) {
	t.Parallel()

	iface := reflect.TypeFor[utils.ComponentProps]()
	if iface.NumMethod() != 2 {
		t.Errorf("ComponentProps interface has %d methods, want 2 (GetBaseProps, SetBaseProps)", iface.NumMethod())
	}

	for _, name := range []string{"GetBaseProps", "SetBaseProps"} {
		m, ok := iface.MethodByName(name)
		if !ok {
			t.Errorf("ComponentProps interface missing method %q", name)

			continue
		}

		mt := m.Type
		if mt.IsVariadic() {
			t.Errorf("%s must not be variadic", name)
		}
	}
}

// embedsBaseProps returns true if typ contains an embedded field of type
// utils.BaseProps (anonymous field, possibly pointer-indirected).
func embedsBaseProps(typ reflect.Type) bool {
	if typ.Kind() != reflect.Struct {
		return false
	}

	for f := range typ.Fields() {
		if !f.Anonymous {
			continue
		}

		ft := f.Type
		if ft.Kind() == reflect.Pointer {
			ft = ft.Elem()
		}

		if ft.Name() == "BaseProps" && ft.PkgPath() == "github.com/larsartmann/templ-components/utils" {
			return true
		}
	}

	return false
}
