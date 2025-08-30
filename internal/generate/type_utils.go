package generate

import (
	"fmt"
	"go/ast"
)

// Type utility functions for handling Go types and expressions

// TypeResolver handles type resolution with import context
// Provides methods to resolve Go types with full package path information
type TypeResolver struct {
	imports map[string]string
}

// NewTypeResolver creates a new TypeResolver with the given imports
// The imports map should contain alias-to-package-path mappings
func NewTypeResolver(imports map[string]string) *TypeResolver {
	return &TypeResolver{imports: imports}
}

// ResolveType resolves a type expression with package path context
// Handles various Go type expressions including pointers, arrays, maps, channels, etc.
// Returns the fully qualified type name with package path when applicable
func (tr *TypeResolver) ResolveType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		// Check if this is a type from an imported package
		if pkgPath, exists := tr.imports[t.Name]; exists {
			// This is a package alias, use the full package path
			return pkgPath + "." + t.Name
		}
		// This is a local or built-in type
		return t.Name
	case *ast.StarExpr:
		return "*" + tr.ResolveType(t.X)
	case *ast.SelectorExpr:
		// Handle selector expressions like package.Type
		if x, ok := t.X.(*ast.Ident); ok {
			if pkgPath, exists := tr.imports[x.Name]; exists {
				// This is a package alias, use the full package path
				return pkgPath + "." + t.Sel.Name
			}
			// This is a selector from a local type
			return x.Name + "." + t.Sel.Name
		}
		return exprToString(t)
	case *ast.ArrayType:
		return "[]" + tr.ResolveType(t.Elt)
	case *ast.MapType:
		return "map[" + tr.ResolveType(t.Key) + "]" + tr.ResolveType(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		if t.Dir == ast.SEND {
			return "chan<- " + tr.ResolveType(t.Value)
		} else if t.Dir == ast.RECV {
			return "<-chan " + tr.ResolveType(t.Value)
		}
		return "chan " + tr.ResolveType(t.Value)
	default:
		return exprToString(t)
	}
}

// exprToString converts an AST expression to a string representation
// Provides a fallback string representation for complex expressions
// Used when detailed type resolution is not available
func exprToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + exprToString(t.X)
	case *ast.SelectorExpr:
		return exprToString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(t.Elt)
	case *ast.MapType:
		return "map[" + exprToString(t.Key) + "]" + exprToString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func(...)"
	case *ast.ChanType:
		if t.Dir == ast.SEND {
			return "chan<- " + exprToString(t.Value)
		} else if t.Dir == ast.RECV {
			return "<-chan " + exprToString(t.Value)
		}
		return "chan " + exprToString(t.Value)
	default:
		return fmt.Sprintf("%T", t)
	}
}

// extractFieldTag extracts the tag from a field
// Returns the raw tag string or empty string if no tag is present
func extractFieldTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return tag.Value
}

// isStructType checks if a type is a struct type
// Provides basic struct type detection without full context
// Note: This is a simplified check and may not catch all struct types
func isStructType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return isStructType(t.X)
	case *ast.Ident:
		// We can't determine if this is a struct without more context
		// This would require looking up the actual type definition
		return false
	case *ast.SelectorExpr:
		// This could be a struct from another package
		return true
	default:
		return false
	}
}
