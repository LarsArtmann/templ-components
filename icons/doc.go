// Package icons provides a type-safe icon system with named constants and SVG rendering
// for use across all component packages.
//
// # Animated Icons
//
// AnimatedIcon renders any icon with a hover-triggered CSS animation, inspired by
// https://www.heroicons-animated.com/. Animations are pure CSS (zero JavaScript),
// respect prefers-reduced-motion, and trigger on both :hover and :focus-within.
//
//	@icons.AnimatedIcon(icons.Heart, "h-6 w-6 text-red-500")
//	@icons.AnimatedIconWithAnimation(icons.Bell, icons.AnimWiggle, "h-6 w-6")
//
// Each icon has a sensible default animation (Heart→pulse, Bell→wiggle, Settings→spin, etc.).
// Override with AnimatedIconWithAnimation and any Animation constant.
// Requires the .tc-anim-* CSS classes from templates/custom.css.
package icons
