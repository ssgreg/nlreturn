package nlreturn

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const (
	linterName = "nlreturn"
	linterDoc  = `Linter requires a new line before return and branch statements except when the return is alone inside a statement group (such as an if statement) to increase code clarity.`
)

// NewAnalyzer returns a new nlreturn analyzer.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: linterName,
		Doc:  linterDoc,
		Run:  run,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			switch c := node.(type) {
			case *ast.CaseClause:
				inspectBlock(pass, c.Body)
			case *ast.CommClause:
				inspectBlock(pass, c.Body)
			case *ast.BlockStmt:
				inspectBlock(pass, c.List)
			}

			return true
		})
	}

	return nil, nil
}

func inspectBlock(pass *analysis.Pass, block []ast.Stmt) {
	for i, stmt := range block {
		switch stmt.(type) {
		case *ast.BranchStmt, *ast.ReturnStmt:
			if i == 0 {
				return
			}

			if !strings.HasSuffix(pass.Fset.Position(stmt.Pos()).Filename, ".go") {
				return
			}

			if isFalsePositive(stmt) {
				return
			}

			if line(pass, stmt.Pos())-line(pass, block[i-1].End()) <= 1 {
				pass.Report(analysis.Diagnostic{
					Pos:     stmt.Pos(),
					Message: fmt.Sprintf("%s with no blank line before", name(stmt)),
					SuggestedFixes: []analysis.SuggestedFix{
						{
							TextEdits: []analysis.TextEdit{
								{
									Pos:     stmt.Pos(),
									NewText: []byte("\n"),
									End:     stmt.Pos(),
								},
							},
						},
					},
				})
			}
		}
	}
}

func name(stmt ast.Stmt) string {
	switch c := stmt.(type) {
	case *ast.BranchStmt:
		return c.Tok.String()
	case *ast.ReturnStmt:
		return "return"
	default:
		return "unknown"
	}
}

func line(pass *analysis.Pass, pos token.Pos) int {
	return pass.Fset.Position(pos).Line
}

func isFalsePositive(node ast.Node) bool {
	r, ok := node.(*ast.ReturnStmt)
	if !ok {
		return false
	}

	if len(r.Results) != 1 {
		return false
	}

	call, ok := r.Results[0].(*ast.CallExpr)
	if !ok {
		return false
	}

	fun, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false
	}

	if fun.Name == "_Cfunc_GoBytes" {
		return true
	}

	return false
}
