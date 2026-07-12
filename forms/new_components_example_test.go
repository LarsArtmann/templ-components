package forms

func ExampleFilterDropdown() {
	_ = FilterDropdown(FilterDropdownProps{
		Name:     "status",
		Label:    "Status",
		Value:    "active",
		HxGet:    "/api/users",
		HxTarget: "#user-list",
		Options: []SelectOption{
			{Value: "all", Label: "All"},
			{Value: "active", Label: "Active"},
			{Value: "inactive", Label: "Inactive"},
		},
	})
	// Output:
}

func ExampleSlider() {
	_ = Slider(SliderProps{
		Name:      "volume",
		Label:     "Volume",
		Min:       0,
		Max:       100,
		Value:     50,
		Step:      5,
		ShowValue: true,
	})
	// Output:
}

func ExampleRating() {
	_ = Rating(RatingProps{
		Name:  "quality",
		Value: 4,
		Max:   5,
		Label: "Product Quality",
	})
	// Output:
}

func ExampleRating_readOnly() {
	_ = Rating(RatingProps{
		Value:    3,
		Max:      5,
		ReadOnly: true,
	})
	// Output:
}
