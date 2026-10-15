package tone

// ColorClass is a complete, scanner-safe color utility string: the light
// class and its dark: half in ONE literal (e.g. "text-blue-600
// dark:text-blue-400"). Chart series, spinners, and progress fills accept
// it so a typo cannot silently render nothing — pass a named constant
// (display.ChartColorBlue, …) or a complete literal of your own; a bare
// "text-blue-600" without the dark half compiles but trips the dark-mode
// guards.
//
// The string kind keeps every existing untyped literal assigning (TODO
// #426 non-breaking); a typed value converts with string(c) where a raw
// string is needed.
type ColorClass string
