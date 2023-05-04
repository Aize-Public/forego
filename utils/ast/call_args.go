package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/Aize-Public/forego/ctx"
	util "github.com/Aize-Public/forego/utils"
)

type Call struct {
	Name string
	Args []Arg
	All  string
}

type Arg struct {
	//Expr       ast.Expr
	Src        string
	Assignment string
}

func (a Arg) String() string {
	return a.Src
}

// expand the source code of the caller (use above==0 for the current function)
func Caller(above int) (call *Call, err error) {
	stack := util.Caller(above + 2) // caller of the caller
	fname := util.Caller(above + 1).FuncName()
	src, err := os.ReadFile(stack.AbsFile())
	if err != nil {
		return nil, ctx.Errorf(nil, "opening %q: %w", stack.AbsFile(), err)
	}
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, stack.AbsFile(), src, 0)
	if err != nil {
		return nil, ctx.Errorf(nil, "parsing %q: %w", stack.AbsFile(), err)
	}

	ast.Walk(visitor{
		fset:  fset,
		src:   src,
		fname: fname,
		stack: stack,
		found: func(found Call) {
			call = &found
		},
	}, f)

	if call == nil {
		return nil, ctx.Errorf(nil, "can't find function %q in %s:%d", fname, stack.File, stack.Line)
	}
	return call, nil
}

type visitor struct {
	scope map[string]string
	fset  *token.FileSet
	src   []byte
	fname string
	stack util.StackFrame
	found func(Call)
}

func (this visitor) Line(p token.Pos) int {
	return this.fset.PositionFor(p, false).Line
}

func (this visitor) Src(a, b token.Pos) string {
	p0 := this.fset.Position(a).Offset
	p1 := this.fset.Position(b).Offset
	return string(this.src[p0:p1])
}

func (this visitor) Visit(n ast.Node) ast.Visitor {
	//log.Printf("visit(%+v)", n)
	switch n := n.(type) {

	case *ast.BlockStmt:
		//fmt.Printf("block: %d, %d\n", n.Lbrace, n.Rbrace)
		m := map[string]string{}
		for k, v := range this.scope {
			m[k] = v
		}
		this.scope = m

	case *ast.AssignStmt:
		src := this.Src(n.Rhs[0].Pos(), n.End())
		for _, lh := range n.Lhs {
			switch lh := lh.(type) {
			case *ast.Ident:
				this.scope[lh.Name] = src
			}
		}

	case *ast.CallExpr:
		line := this.Line(n.Pos())
		//fmt.Printf("%d: %+v\n", line, n)
		if line == this.stack.Line {
			//log.Printf("at :%d found %T %+v", line, n.Fun, n.Fun)
			switch f := n.Fun.(type) {
			case *ast.Ident:
				//log.Printf("Ident name: %q == %q", f.Name, this.fname)
				if f.Name == this.fname {
					call := Call{
						Name: f.Name,
						All:  this.Src(f.Pos(), f.End()),
					}
					for _, a := range n.Args {
						arg := Arg{
							Src: this.Src(a.Pos(), a.End()),
						}
						switch a := a.(type) {
						case *ast.FuncLit:
							start := a.Body.List[0].Pos()
							end := a.Body.List[len(a.Body.List)-1].End()
							arg.Src = this.Src(start, end)
						default:
							//log.Printf("a: %#v", a)
						}
						arg.Assignment = this.scope[arg.Src]
						call.Args = append(call.Args, arg)
					}
					this.found(call)
				}
			case *ast.SelectorExpr:
				//log.Printf("searching for %q: %q", this.fname, f.Sel.Name)
				if f.Sel.Name == this.fname {
					call := Call{
						Name: f.Sel.Name,
						All:  this.Src(f.Pos(), f.End()),
					}
					switch fn := f.X.(type) {
					case *ast.Ident:
						call.Name = fn.Name + "." + f.Sel.Name
					}
					for _, a := range n.Args {
						arg := Arg{
							Src: this.Src(a.Pos(), a.End()),
						}
						switch a := a.(type) {
						case *ast.FuncLit:
							start := a.Body.List[0].Pos()
							end := a.Body.List[len(a.Body.List)-1].End()
							arg.Src = this.Src(start, end) // `func() { BODY }` => `BODY`
						}
						arg.Assignment = this.scope[arg.Src]
						call.Args = append(call.Args, arg)
					}
					this.found(call)
				}

			case *ast.ArrayType: // ignore
			case *ast.ParenExpr: // ignore

			default:
				src := this.Src(f.Pos(), f.End())
				log.Printf("call_args: unexpected %T %s", f, src)
			}
		}
	}
	return this
}
