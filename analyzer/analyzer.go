package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "gorecover",
	Doc:      "Checks that goroutine has recover in defer function",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GoStmt)(nil),
	}
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var call *ast.CallExpr
		if g, ok := n.(*ast.GoStmt); !ok {
			return
		} else {
			call = g.Call
		}

		if call == nil {
			return
		}

		var body ast.Node
		switch f := call.Fun.(type) {
		case *ast.Ident:
			fd, ok := f.Obj.Decl.(*ast.FuncDecl)
			if !ok {
				return
			}
			body = fd
		case *ast.FuncLit:
			body = f
		default:
			return
		}

		ast.Inspect(body, func(n ast.Node) bool {
			block, ok := n.(*ast.BlockStmt)
			if !ok {
				return true
			}
			if len(block.List) == 0 {
				pass.Reportf(call.Pos(), "goroutine should have recover in defer func")
				return false
			}

			if d, ok := block.List[0].(*ast.DeferStmt); !ok || d.Call == nil {
				pass.Reportf(call.Pos(), "goroutine should have recover in defer func")
				return false
			} else {
				hasRecoverStmt := false
				// 看看defer里面有没有调recover()
				inspect.Nodes([]ast.Node{d.Call}, func(node ast.Node, push bool) (proceed bool) {
					if !push {
						return false
					}
					switch x := node.(type) {
					case *ast.CallExpr:
						if ex, ok := x.Fun.(*ast.Ident); !ok {
							break
						} else if ex.Name == "recover" {
							hasRecoverStmt = true
							return false
						}
					}
					return true
				})

				if !hasRecoverStmt {
					pass.Reportf(call.Pos(), "goroutine should have recover in defer func")
				}
				return false
			}
			return false
		})
	})
	return nil, nil
}

func hasRecover(bs *ast.BlockStmt) bool {
	for _, blockStmt := range bs.List {
		deferStmt, ok := blockStmt.(*ast.DeferStmt) // 是否包含defer 语句
		if !ok {
			return false
		}
		switch deferStmt.Call.Fun.(type) {
		case *ast.SelectorExpr:
			// 判断是否defer中包含  helper.Recover()
			selectorExpr := deferStmt.Call.Fun.(*ast.SelectorExpr)
			if "Recover" == selectorExpr.Sel.Name {
				return true
			}
		case *ast.FuncLit:
			// 判断是否有 defer func(){ }()
			fl := deferStmt.Call.Fun.(*ast.FuncLit)
			for i := range fl.Body.List {

				stmt := fl.Body.List[i]
				switch stmt.(type) {
				case *ast.ExprStmt:
					exprStmt := stmt.(*ast.ExprStmt)
					if isRecoverExpr(exprStmt.X) { // recover()
						return true
					}
				case *ast.IfStmt:
					is := stmt.(*ast.IfStmt) // if r:=recover();r!=nil{}
					as, ok := is.Init.(*ast.AssignStmt)
					if !ok {
						continue
					}
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}
				case *ast.AssignStmt:
					as := stmt.(*ast.AssignStmt) // r=:recover
					if isRecoverExpr(as.Rhs[0]) {
						return true
					}

				}
			}
		}
	}
	return false
}

func isRecoverExpr(expr ast.Expr) bool {
	ac, ok := expr.(*ast.CallExpr) // r:=recover()
	if !ok {
		return false
	}
	id, ok := ac.Fun.(*ast.Ident)
	if !ok {
		return false
	}
	if "recover" == id.Name {
		return true
	}
	return false
}
