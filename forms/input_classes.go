// Shared Tailwind class strings for form input elements.
package forms

// baseInputClass returns the shared Tailwind classes for text inputs, selects, and textareas.
func baseInputClass(hasError bool) string {
	base := "block w-full rounded-md border-0 py-1.5 text-gray-900 dark:text-white shadow-xs ring-1 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-inset sm:text-sm sm:leading-6 transition-colors caret-blue-600 dark:caret-blue-400 dark:bg-gray-800 dark:ring-gray-700 dark:placeholder:text-gray-500"
	if hasError {
		return base + " ring-red-300 focus:ring-red-500 dark:ring-red-700 dark:focus:ring-red-500"
	}
	return base + " ring-gray-300 focus:ring-blue-600 dark:focus:ring-blue-500"
}
