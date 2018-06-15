package lint

import (
	"go/ast"
)

func init() {
	addChecker(flagDerefChecker{}, &ruleInfo{})
}

type flagDerefChecker struct {
	baseExprChecker

	flagPtrFuncs map[string]bool
}

func (c flagDerefChecker) New(ctx *context) func(*ast.File) {
	return wrapExprChecker(&flagDerefChecker{
		baseExprChecker: baseExprChecker{ctx: ctx},

		flagPtrFuncs: map[string]bool{
			"flag.Bool":     true,
			"flag.Duration": true,
			"flag.Float64":  true,
			"flag.Int":      true,
			"flag.Int64":    true,
			"flag.String":   true,
			"flag.Uint":     true,
			"flag.Uint64":   true,
		},
	})
}

func (c *flagDerefChecker) CheckExpr(expr ast.Expr) {
	if expr, ok := expr.(*ast.StarExpr); ok {
		call, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return
		}
		called := qualifiedName(call.Fun)
		if c.flagPtrFuncs[called] {
			c.warn(expr, called+"Var")
		}
	}
}

func (c *flagDerefChecker) warn(x ast.Node, suggestion string) {
	c.ctx.Warn(x, "immediate deref in %s is most likely an error; consider using %s",
		x, suggestion)
}
