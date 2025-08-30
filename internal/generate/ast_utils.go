package generate

import (
	"go/ast"
)

// AST utility functions for common traversal patterns

// findFunctionDeclarations finds all function declarations in an AST node
// Returns a slice of all function declarations found during AST traversal
func findFunctionDeclarations(node ast.Node) []*ast.FuncDecl {
	var funcs []*ast.FuncDecl

	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			funcs = append(funcs, fn)
		}
		return true
	})

	return funcs
}

// findTypeDeclarations finds all type declarations in an AST node
// Returns a slice of all type specifications found during AST traversal
func findTypeDeclarations(node ast.Node) []*ast.TypeSpec {
	var types []*ast.TypeSpec

	ast.Inspect(node, func(n ast.Node) bool {
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			types = append(types, typeDecl)
		}
		return true
	})

	return types
}

// findStructTypes finds all struct type declarations in an AST node
// Returns a slice of all struct type declarations found during AST traversal
func findStructTypes(node ast.Node) []*ast.StructType {
	var structs []*ast.StructType

	ast.Inspect(node, func(n ast.Node) bool {
		if structType, ok := n.(*ast.StructType); ok {
			structs = append(structs, structType)
		}
		return true
	})

	return structs
}

// findCallExpressions finds all call expressions in an AST node
// Returns a slice of all function/method calls found during AST traversal
func findCallExpressions(node ast.Node) []*ast.CallExpr {
	var calls []*ast.CallExpr

	ast.Inspect(node, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			calls = append(calls, call)
		}
		return true
	})

	return calls
}

// findSelectorExpressions finds all selector expressions in an AST node
// Returns a slice of all field/method selections found during AST traversal
func findSelectorExpressions(node ast.Node) []*ast.SelectorExpr {
	var selectors []*ast.SelectorExpr

	ast.Inspect(node, func(n ast.Node) bool {
		if sel, ok := n.(*ast.SelectorExpr); ok {
			selectors = append(selectors, sel)
		}
		return true
	})

	return selectors
}

// isElsaCall checks if a call expression is an elsa.Generate or elsa.Set call
// Validates that the call is to the "elsa" package with the specified method name
func isElsaCall(call *ast.CallExpr, method string) bool {
	if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
		if x, ok := sel.X.(*ast.Ident); ok && x.Name == "elsa" && sel.Sel.Name == method {
			return true
		}
	}
	return false
}

// extractIdentifiers extracts all identifiers from a call expression's arguments
// Returns a slice of identifier names from the call arguments
func extractIdentifiers(call *ast.CallExpr) []string {
	var ids []string

	for _, arg := range call.Args {
		if ident, ok := arg.(*ast.Ident); ok {
			ids = append(ids, ident.Name)
		}
	}

	return ids
}

// findStructByName finds a struct type declaration by name in an AST node
// Searches through the AST to locate a struct with the specified name
// Returns the struct type if found, nil otherwise
func findStructByName(node ast.Node, structName string) *ast.StructType {
	var found *ast.StructType

	ast.Inspect(node, func(n ast.Node) bool {
		typeDecl, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if typeDecl.Name.Name != structName {
			return true
		}

		if structType, ok := typeDecl.Type.(*ast.StructType); ok {
			found = structType
			return false // Stop traversal after finding
		}

		return true
	})

	return found
}
