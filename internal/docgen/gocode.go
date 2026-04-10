package docgen

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

// ExtractGoCode extracts a function or struct declaration from a Go source file.
// Transforms supported:
//   - "stripImports": Remove import declarations
//   - "rename:Old→New": Replace type/variable names
//   - "package:Name": Change package name
func ExtractGoCode(filePath, name string, transforms []string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return "", fmt.Errorf("failed to parse file: %w", err)
	}

	// Parse transforms
	transformOpts := parseTransforms(transforms)

	// Find the target declaration
	var targetDecl ast.Decl
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if d.Name.Name == name {
				targetDecl = d
				break
			}
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ts.Name.Name == name {
						targetDecl = d
						break
					}
				}
			}
		}
		if targetDecl != nil {
			break
		}
	}

	if targetDecl == nil {
		return "", fmt.Errorf("declaration %q not found", name)
	}

	// Build a new file with just the target declaration
	newFile := &ast.File{
		Name:    ast.NewIdent(file.Name.Name),
		Decls:   []ast.Decl{targetDecl},
		Scope:   ast.NewScope(nil),
		Imports: file.Imports,
	}

	// Apply package rename if requested
	if transformOpts.packageName != "" {
		newFile.Name.Name = transformOpts.packageName
	}

	// Strip imports if requested
	if transformOpts.stripImports {
		newFile.Imports = nil
		// Remove import declarations from Decls
		var filteredDecls []ast.Decl
		for _, decl := range newFile.Decls {
			if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.IMPORT {
				continue
			}
			filteredDecls = append(filteredDecls, decl)
		}
		newFile.Decls = filteredDecls
	}

	// Apply renames if requested
	if len(transformOpts.renames) > 0 {
		ast.Inspect(newFile, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.Ident:
				for oldName, newName := range transformOpts.renames {
					if node.Name == oldName {
						node.Name = newName
					}
				}
			}
			return true
		})
	}

	// Format the result
	var result strings.Builder
	if err := format.Node(&result, fset, newFile); err != nil {
		return "", fmt.Errorf("failed to format code: %w", err)
	}

	return result.String(), nil
}

// transformOptions holds parsed transform options
type transformOptions struct {
	stripImports bool
	renames      map[string]string
	packageName  string
}

// parseTransforms parses transform strings into options
func parseTransforms(transforms []string) transformOptions {
	opts := transformOptions{
		renames: make(map[string]string),
	}

	for _, t := range transforms {
		switch {
		case t == "stripImports":
			opts.stripImports = true
		case strings.HasPrefix(t, "rename:"):
			// Parse "rename:Old→New" or "rename:Old->New"
			rest := strings.TrimPrefix(t, "rename:")
			// Support both arrow types
			var parts []string
			if strings.Contains(rest, "→") {
				parts = strings.SplitN(rest, "→", 2)
			} else {
				parts = strings.SplitN(rest, "->", 2)
			}
			if len(parts) == 2 {
				opts.renames[parts[0]] = parts[1]
			}
		case strings.HasPrefix(t, "package:"):
			opts.packageName = strings.TrimPrefix(t, "package:")
		}
	}

	return opts
}

// extractStringLiteral safely extracts string value from a basic lit
func extractStringLiteral(lit *ast.BasicLit) string {
	if lit == nil || lit.Kind != token.STRING {
		return ""
	}
	s, err := strconv.Unquote(lit.Value)
	if err != nil {
		return lit.Value
	}
	return s
}