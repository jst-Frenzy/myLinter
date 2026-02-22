package analyzer

import (
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

type LoggerConfig struct {
	PkgPath    string
	ExtractPkg func(n *ast.CallExpr, pass *analysis.Pass) *types.PkgName
}

func getSlogPkgName(n *ast.CallExpr, pass *analysis.Pass) *types.PkgName {
	selExpr, fun := n.Fun.(*ast.SelectorExpr)
	if !fun {
		return nil
	}

	switch calls := selExpr.X.(type) {
	case *ast.Ident:
		return getPkgNameFromIdent(calls, pass)
	case *ast.SelectorExpr:
		if leftPart, ok := calls.X.(*ast.Ident); ok {
			return getPkgNameFromIdent(leftPart, pass)
		}
		return getSlogPkgName(&ast.CallExpr{Fun: calls}, pass)
	case *ast.CallExpr:
		return getSlogPkgName(calls, pass)
	}
	return nil
}

func getZapPkgName(n *ast.CallExpr, pass *analysis.Pass) *types.PkgName {
	selExpr, fun := n.Fun.(*ast.SelectorExpr)
	if !fun {
		return nil
	}

	switch calls := selExpr.X.(type) {
	case *ast.Ident:
		return getPkgNameFromIdent(calls, pass)
	case *ast.SelectorExpr:
		if leftPart, ok := calls.X.(*ast.Ident); ok {
			return getPkgNameFromIdent(leftPart, pass)
		}
		return getZapPkgName(&ast.CallExpr{Fun: calls}, pass)
	case *ast.CallExpr:
		return getZapPkgName(calls, pass)
	}
	return nil
}

func getPkgNameFromIdent(ident *ast.Ident, pass *analysis.Pass) *types.PkgName {
	obj := pass.TypesInfo.ObjectOf(ident)
	if obj == nil {
		return nil
	}

	if pkgName, ok := obj.(*types.PkgName); ok {
		return pkgName
	}

	return getPkgNameFromVar(obj)
}

func getPkgNameFromVar(obj types.Object) *types.PkgName {
	if v, ok := obj.(*types.Var); ok {
		if ptr, ok := v.Type().(*types.Pointer); ok {
			if named, ok := ptr.Elem().(*types.Named); ok {
				if pkg := named.Obj().Pkg(); pkg != nil {
					return types.NewPkgName(0, nil, pkg.Name(), pkg)
				}
			}
		}
	}
	return nil
}
