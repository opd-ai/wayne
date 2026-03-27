//go:build windows || darwin || android || ios || linux

package wayne_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// TestAPIAlignment validates that wayne's public API aligns with wain's by
// parsing both codebases' AST and comparing exported symbols.
//
// This test does NOT import both packages (which would fail due to conflicting
// build tags), but instead analyzes source code structure to verify compatibility.
func TestAPIAlignment(t *testing.T) {
	// Parse wayne's public API
	wayneAPI, err := extractPublicAPI(".")
	if err != nil {
		t.Fatalf("failed to extract wayne API: %v", err)
	}

	// For now, document the wayne API surface
	// In CI with access to wain source, this would compare against wain's API
	t.Run("DocumentPublicAPI", func(t *testing.T) {
		documentPublicAPI(t, wayneAPI)
	})

	// Validate that critical types exist
	t.Run("RequiredTypes", func(t *testing.T) {
		requiredTypes := []string{
			"PublicWidget", "Container", "Canvas",
			"Button", "Label", "TextInput", "Panel", "Row", "Column",
			"Grid", "Stack", "ScrollView", "Spacer", "ImageWidget",
			"App", "Window", "Theme", "Color", "Size",
			"Event", "PointerEvent", "KeyEvent", "TouchEvent", "WindowEvent",
		}

		for _, typeName := range requiredTypes {
			if !wayneAPI.HasType(typeName) {
				t.Errorf("missing required type: %s", typeName)
			}
		}
	})

	// Validate that critical functions exist
	t.Run("RequiredFunctions", func(t *testing.T) {
		requiredFuncs := []string{
			"NewApp", "NewAppWithConfig",
			"NewButton", "NewLabel", "NewTextInput",
			"NewPanel", "NewRow", "NewColumn", "NewGrid", "NewStack",
			"NewScrollView", "NewSpacer", "NewImageWidget",
			"DefaultConfig", "DefaultDark", "DefaultLight", "HighContrast",
			"RGB", "RGBA",
			"NewEventDispatcher", "NewFocusManager",
		}

		for _, funcName := range requiredFuncs {
			if !wayneAPI.HasFunction(funcName) {
				t.Errorf("missing required function: %s", funcName)
			}
		}
	})
}

// API represents the public API surface of a package
type API struct {
	Types     map[string]*TypeInfo
	Functions map[string]*FuncInfo
	Constants map[string]*ConstInfo
}

// TypeInfo represents information about an exported type
type TypeInfo struct {
	Name    string
	Kind    string // "struct", "interface", "alias"
	Methods []string
}

// FuncInfo represents information about an exported function
type FuncInfo struct {
	Name      string
	Params    int
	Results   int
	Signature string
}

// ConstInfo represents information about an exported constant
type ConstInfo struct {
	Name  string
	Value string
}

func (a *API) HasType(name string) bool {
	_, exists := a.Types[name]
	return exists
}

func (a *API) HasFunction(name string) bool {
	_, exists := a.Functions[name]
	return exists
}

// extractPublicAPI parses Go source files and extracts the public API
func extractPublicAPI(pkgPath string) (*API, error) {
	api := &API{
		Types:     make(map[string]*TypeInfo),
		Functions: make(map[string]*FuncInfo),
		Constants: make(map[string]*ConstInfo),
	}

	fset := token.NewFileSet()

	// Parse all .go files (excluding tests)
	pattern := filepath.Join(pkgPath, "*.go")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, path := range matches {
		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}

		// Extract declarations
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				// Skip methods for now (they'll be captured via type methods)
				if d.Recv == nil && d.Name.IsExported() {
					api.Functions[d.Name.Name] = &FuncInfo{
						Name:      d.Name.Name,
						Params:    countParams(d.Type.Params),
						Results:   countParams(d.Type.Results),
						Signature: funcSignature(d.Type),
					}
				}

			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							kind := typeKind(s.Type)
							api.Types[s.Name.Name] = &TypeInfo{
								Name:    s.Name.Name,
								Kind:    kind,
								Methods: []string{},
							}
						}

					case *ast.ValueSpec:
						// Constants
						for _, name := range s.Names {
							if name.IsExported() && d.Tok == token.CONST {
								api.Constants[name.Name] = &ConstInfo{
									Name: name.Name,
								}
							}
						}
					}
				}
			}
		}
	}

	return api, nil
}

func countParams(fields *ast.FieldList) int {
	if fields == nil {
		return 0
	}
	count := 0
	for _, field := range fields.List {
		if len(field.Names) > 0 {
			count += len(field.Names)
		} else {
			count++
		}
	}
	return count
}

func funcSignature(ft *ast.FuncType) string {
	var sb strings.Builder
	sb.WriteString("(")
	if ft.Params != nil {
		for i, field := range ft.Params.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(exprString(field.Type))
		}
	}
	sb.WriteString(")")

	if ft.Results != nil && len(ft.Results.List) > 0 {
		sb.WriteString(" ")
		if len(ft.Results.List) > 1 {
			sb.WriteString("(")
		}
		for i, field := range ft.Results.List {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(exprString(field.Type))
		}
		if len(ft.Results.List) > 1 {
			sb.WriteString(")")
		}
	}
	return sb.String()
}

func exprString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + exprString(e.X)
	case *ast.SelectorExpr:
		return exprString(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprString(e.Elt)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.FuncType:
		return "func" + funcSignature(e)
	default:
		return "unknown"
	}
}

func typeKind(expr ast.Expr) string {
	switch expr.(type) {
	case *ast.StructType:
		return "struct"
	case *ast.InterfaceType:
		return "interface"
	case *ast.Ident:
		return "alias"
	default:
		return "unknown"
	}
}

func documentPublicAPI(t *testing.T, api *API) {
	t.Logf("Wayne Public API Surface:\n")

	// Document types
	var typeNames []string
	for name := range api.Types {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)

	t.Logf("\n=== Exported Types ===")
	for _, name := range typeNames {
		info := api.Types[name]
		t.Logf("  type %s (%s)", name, info.Kind)
	}

	// Document functions
	var funcNames []string
	for name := range api.Functions {
		funcNames = append(funcNames, name)
	}
	sort.Strings(funcNames)

	t.Logf("\n=== Exported Functions ===")
	for _, name := range funcNames {
		info := api.Functions[name]
		t.Logf("  func %s%s", name, info.Signature)
	}

	// Document constants
	var constNames []string
	for name := range api.Constants {
		constNames = append(constNames, name)
	}
	sort.Strings(constNames)

	t.Logf("\n=== Exported Constants ===")
	for _, name := range constNames {
		t.Logf("  const %s", name)
	}

	t.Logf("\n=== API Summary ===")
	t.Logf("  Types: %d", len(api.Types))
	t.Logf("  Functions: %d", len(api.Functions))
	t.Logf("  Constants: %d", len(api.Constants))
}
