package elsa

// Set configures dependency data that will be used for building.
// It accepts variadic arguments following the pattern: pkg.funcname
// This function sets up the necessary data structures and dependencies
// that will be parsed and used by the Generate function.
func Set(node ...any) bool {
	return true
}

// Generate performs the building process based on the data set by Set function.
// It parses the configured dependencies and executes the building logic
// according to the previously set configuration.
func Generate(node ...any) bool {
	return true
}
