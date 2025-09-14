package make

import "strings"

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(s string) string {
	// Convert PascalCase to snake_case
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

// toTitleCase converts first letter to uppercase
func toTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

// toCamelCase converts snake_case to camelCase
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	result := strings.ToLower(parts[0])
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += toTitleCase(parts[i])
		}
	}
	return result
}

// toPascalCase converts snake_case to PascalCase
func toPascalCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}

	result := ""
	for _, part := range parts {
		if len(part) > 0 {
			result += toTitleCase(part)
		}
	}
	return result
}

// toPlural converts singular to plural (basic implementation)
func toPlural(s string) string {
	s = strings.ToLower(s)

	if strings.HasSuffix(s, "y") && len(s) > 1 {
		return strings.TrimSuffix(s, "y") + "ies"
	}
	if strings.HasSuffix(s, "s") || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") {
		return s + "es"
	}
	return s + "s"
}

// toSingular converts plural to singular (basic implementation)
func toSingular(s string) string {
	s = strings.ToLower(s)

	if strings.HasSuffix(s, "ies") && len(s) > 3 {
		return strings.TrimSuffix(s, "ies") + "y"
	}
	if strings.HasSuffix(s, "es") && len(s) > 2 {
		return strings.TrimSuffix(s, "es")
	}
	if strings.HasSuffix(s, "s") && len(s) > 1 {
		return strings.TrimSuffix(s, "s")
	}
	return s
}
