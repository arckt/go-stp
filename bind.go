package bind

// #cgo linux CFLAGS: -I/usr/local/lib
// #cgo linux LDFLAGS: /usr/local/lib/libstp.so
// #include <stp/c_interface.h>
// #include <stdio.h>
// #include <stdlib.h>
/*
void errorHandlerC(const char* err_msg) {
    printf("Error: %s\n", err_msg);
    exit(1);
}
typedef void (*closure)();
typedef void * VC;
typedef void * Expr;
typedef void * Type;
typedef void * WholeCounterExample;
*/
import "C"

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math"
	"strconv"
	"strings"
)

type VC C.VC
type Expr C.Expr
type Type C.Type
type WholeCounterExample C.WholeCounterExample

type Solver struct {
	handle C.VC
	signed bool
	width  uint
	vecs   map[string]C.Expr
}

//H2i turns a string 0x00000 to its base 10 counterpart in uint64
func (s *Solver) H2i(in string) uint64 {
	tmp := strings.ReplaceAll(in, "0x", "")
	tmp = strings.TrimLeft(tmp, "0")
	tmp = strings.TrimRight(tmp, " ")

	if tmp == "" {
		log.Fatal("Unsat")
	}

	res, err := strconv.ParseUint(tmp, 16, int(s.width))
	if err != nil {
		log.Fatal(err)
	}

	return res
}

//SH2i turns a string 0x00000 to its base 10 counterpart in int64
func (s *Solver) SH2i(in string) int64 {
	tmp := strings.ReplaceAll(in, "0x", "")
	tmp = strings.TrimLeft(tmp, "0")
	tmp = strings.TrimRight(tmp, " ")

	if tmp == "" {
		log.Fatal("Unsat")
	}

	res, err := strconv.ParseUint(tmp, 16, int(s.width))
	if err != nil {
		log.Fatal(err)
	}

	if res&(1<<s.width-1) > 0 {
		return -((^(int64(res)) & int64(math.Pow(2, float64(s.width))-1)) + 1)
	}

	return int64(res)
}

//SOLVER RELATED FUNCTIONS
//RegisterErrorHandler registers an error handler which triggers when there is an error in the C api
func RegisterErrorHandler() {
	C.vc_registerErrorHandler((*[0]byte)(C.errorHandlerC))
}

//CreateValidityChecker initializes the stp checker
func Init() Solver {
	RegisterErrorHandler()
	return Solver{C.VC(C.vc_createValidityChecker()), false, 32, map[string]C.Expr{}}
}

//AssertFormula adds constraint e to the solver
func (s *Solver) AssertFormula(e Expr) {
	C.vc_assertFormula(s.handle, C.Expr(e))
}

///Destroy tells the solver to sleep
func (s *Solver) Destroy() {
	C.vc_Destroy(s.handle)
}

//PrintAsserts prints to stdout the assertions made to the checker
func (s *Solver) PrintAsserts(simplify_print int) {
	C.vc_printAsserts(s.handle, C.int(simplify_print))
}

//Query tells the solver to validate an Expr e
func (s *Solver) Query(e Expr) int {
	return int(C.vc_query(s.handle, C.Expr(e)))
}

//PrintQuery prints to stdout the most recent Query
func (s *Solver) PrintQuery() {
	C.vc_printQuery(s.handle)
}

//PrintCounterExample prints to stdout the CounterExample which invalidates the Queried Expr
func (s *Solver) PrintCounterExample() {
	C.vc_printCounterExample(s.handle)
}

//GetWholeCounterExample returns the full CounterExample that invalidates the Queried Expr
func (s *Solver) GetWholeCounterExample() WholeCounterExample {
	return WholeCounterExample(C.vc_getWholeCounterExample(s.handle))
}

//GetTermFromCounterExample returns an Expr that is associated with Expr e from the WholeCounterExample c
func (s *Solver) GetTermFromCounterExample(e Expr, c WholeCounterExample) Expr {
	return Expr(C.vc_getTermFromCounterExample(s.handle, C.Expr(e), C.WholeCounterExample(c)))
}

//BITVECTOR RELATED FUNCTIONS
//BvType returns a type of bitwidth no_bits that could be used in VarExpr
func (s *Solver) BvType(no_bits int) Type {
	return Type(C.vc_bvType(s.handle, C.int(no_bits)))
}

//VarExpr returns an Expr that is initialized from a Type and a symbollic handle name
func (s *Solver) VarExpr(name string, typ Type) Expr {
	res := C.vc_varExpr(s.handle, C.CString(name), C.Type(typ))
	s.vecs[name] = res
	return Expr(res)
}

//BVPlusExpr returns an Expr of bit length bitWidth between two Exprs left and right which are added
func (s *Solver) BvPlusExpr(bitWidth int, left Expr, right Expr) Expr {
	return Expr(C.vc_bvPlusExpr(s.handle, C.int(bitWidth), C.Expr(left), C.Expr(right)))
}

//BVPlusExpr returns an Expr of bit length bitWidth between two Exprs left and right which are subtracted
func (s *Solver) BvMinusExpr(bitWidth int, left Expr, right Expr) Expr {
	return Expr(C.vc_bvMinusExpr(s.handle, C.int(bitWidth), C.Expr(left), C.Expr(right)))
}

//BvMultExpr returns an Expr of bit length bitWidth between two Exprs left and right which are multiplied
func (s *Solver) BvMultExpr(bitWidth int, left Expr, right Expr) Expr {
	return Expr(C.vc_bvMultExpr(s.handle, C.int(bitWidth), C.Expr(left), C.Expr(right)))
}

//BvDivExpr returns an Expr of bit length bitWidth between two Exprs dividend and divisor which are divided
func (s *Solver) BvDivExpr(bitWidth int, dividend Expr, divisor Expr) Expr {
	return Expr(C.vc_bvDivExpr(s.handle, C.int(bitWidth), C.Expr(dividend), C.Expr(divisor)))
}

//BvDivExpr returns an Expr of bit length bitWidth between two Exprs dividend and divisor which are in modulo
func (s *Solver) BvModExpr(bitWidth int, dividend Expr, divisor Expr) Expr {
	return Expr(C.vc_bvModExpr(s.handle, C.int(bitWidth), C.Expr(dividend), C.Expr(divisor)))
}

//BvAndExpr returns an Expr of bit length bitWidth between two Exprs left and right which are bitwise and
func (s *Solver) BvAndExpr(left, right Expr) Expr {
	return Expr(C.vc_bvAndExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvAndExpr returns an Expr of bit length bitWidth between two Exprs left and right which are bitwise or
func (s *Solver) BvOrExpr(left, right Expr) Expr {
	return Expr(C.vc_bvOrExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvAndExpr returns an Expr of bit length bitWidth between two Exprs left and right which are bitwise xor
func (s *Solver) BvXorExpr(left, right Expr) Expr {
	return Expr(C.vc_bvXorExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvAndExpr returns an Expr which is the logical not of the Expr child
func (s *Solver) BvNotExpr(child Expr) Expr {
	return Expr(C.vc_bvNotExpr(s.handle, C.Expr(child)))
}

//BvLShiftExpr returns an Expr of bit length bitWidth where left << right
func (s *Solver) BvLShiftExpr(bitWidth int, left Expr, right Expr) Expr {
	return Expr(C.vc_bvLeftShiftExprExpr(s.handle, C.int(bitWidth), C.Expr(left), C.Expr(right)))
}

//BvRShiftExpr returns an Expr of bit length bitWidth where left >> right
func (s *Solver) BvRShiftExpr(bitWidth int, left Expr, right Expr) Expr {
	return Expr(C.vc_bvRightShiftExprExpr(s.handle, C.int(bitWidth), C.Expr(left), C.Expr(right)))
}

//BvUMinusExpr returns the arithmetic negation of Expr child
func (s *Solver) BvUMinusExpr(child Expr) Expr {
	return Expr(C.vc_bvUMinusExpr(s.handle, C.Expr(child)))
}

//BvConstExprFromInt returns an Expr of bit length bitWidth which is a constant of value
func (s *Solver) BvConstExprFromInt(bitWidth int, value uint) Expr {
	return Expr(C.vc_bvConstExprFromInt(s.handle, C.int(bitWidth), C.uint(value)))
}

//LOGIC RELATED FUNCTIONS
//EqExpr returns an Expr where child0 == child1
func (s *Solver) EqExpr(child0 Expr, child1 Expr) Expr {
	return Expr(C.vc_eqExpr(s.handle, C.Expr(child0), C.Expr(child1)))
}

//BvLtExpr returns an Expr where left < right
func (s *Solver) BvLtExpr(left, right Expr) Expr {
	return Expr(C.vc_bvLtExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvLeExpr returns an Expr where left <= right
func (s *Solver) BvLeExpr(left, right Expr) Expr {
	return Expr(C.vc_bvLeExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvGtExpr returns an Expr where left > right
func (s *Solver) BvGtExpr(left, right Expr) Expr {
	return Expr(C.vc_bvGtExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BvGeExpr returns an Expr where left >= right
func (s *Solver) BvGeExpr(left, right Expr) Expr {
	return Expr(C.vc_bvGeExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//BOOLEAN RELATED FUNCTIONS
//AndExpr returns an Expr where left && right
func (s *Solver) AndExpr(left, right Expr) Expr {
	return Expr(C.vc_andExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//OrExpr returns an Expr where left || right
func (s *Solver) OrExpr(left, right Expr) Expr {
	return Expr(C.vc_orExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//XorExpr returns an Expr where left ^ right
func (s *Solver) XorExpr(left, right Expr) Expr {
	return Expr(C.vc_xorExpr(s.handle, C.Expr(left), C.Expr(right)))
}

//NotExpr flips the boolean expression of Expr and returns the flipped Expr
func (s *Solver) NotExpr(child Expr) Expr {
	return Expr(C.vc_notExpr(s.handle, C.Expr(child)))
}

//ImpliesExpr returns an Expr where hyp implies conc
func (s *Solver) ImpliesExpr(hyp Expr, conc Expr) Expr {
	return Expr(C.vc_eqExpr(s.handle, C.Expr(hyp), C.Expr(conc)))
}

//TrueExpr returns an expression which has a value of True
func (s *Solver) TrueExpr() Expr {
	return Expr(C.vc_trueExpr(s.handle))
}

//FalseExpr returns an expression which has a value of False
func (s *Solver) FalseExpr() Expr {
	return Expr(C.vc_falseExpr(s.handle))
}

//HIGH LEVEL API FUNCTIONS
//ExprString converts an Expr e into its string counterpart
func ExprString(e Expr) string {
	return C.GoString(C.exprString(C.Expr(e)))
}

//GetBVLen returns the bitvector length of an Expr e
func GetBVLen(e Expr) int {
	return int(C.getBVLength(C.Expr(e)))
}

//Eval recursively traverses the ast and returns the associated Expr form of the ast
func (s *Solver) Eval(exp ast.Expr) Expr {
	switch exp := exp.(type) {
	case *ast.ParenExpr:
		return s.Eval(exp.X)
	case *ast.UnaryExpr:
		return s.EvalUnaryExpr(exp)
	case *ast.BinaryExpr:
		return s.EvalBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.ParseUint(exp.Value, 10, int(s.width))
			return s.BvConstExprFromInt(int(s.width), uint(i))
		}
	case *ast.Ident:
		return Expr(s.vecs[exp.Name])
	}
	return Expr(nil)
}

//EvalUnaryExpr returns an associated Expr which is representive of the UnaryExpr within an ast
func (s *Solver) EvalUnaryExpr(exp *ast.UnaryExpr) Expr {
	switch exp.Op {
	case token.NOT:
		switch exp.X.(type) {
		case *ast.Ident:
			return s.BvNotExpr(s.Eval(exp.X))
		default:
			return s.NotExpr(s.Eval(exp.X))
		}
	case token.SUB:
		return s.BvUMinusExpr(s.Eval(exp.X))
	}

	return Expr(nil)
}

//EvalBinaryExpr returns an associated Expr which is representive of the BinaryExpr within an ast
func (s *Solver) EvalBinaryExpr(exp *ast.BinaryExpr) Expr {
	left := s.Eval(exp.X)
	right := s.Eval(exp.Y)
	width := GetBVLen(left)

	switch exp.Op {
	case token.ADD:
		return s.BvPlusExpr(width, left, right)
	case token.SUB:
		return s.BvMinusExpr(width, left, right)
	case token.MUL:
		return s.BvMultExpr(width, left, right)
	case token.QUO:
		return s.BvDivExpr(width, left, right)
	case token.REM:
		return s.BvModExpr(width, left, right)
	case token.SHL:
		return s.BvLShiftExpr(width, left, right)
	case token.SHR:
		return s.BvRShiftExpr(width, left, right)
	case token.AND:
		return s.BvAndExpr(left, right)
	case token.OR:
		return s.BvOrExpr(left, right)
	case token.XOR:
		return s.BvXorExpr(left, right)
	case token.LAND:
		return s.AndExpr(left, right)
	case token.LOR:
		return s.OrExpr(left, right)
	case token.LSS:
		return s.BvLtExpr(left, right)
	case token.GTR:
		return s.BvGtExpr(left, right)
	case token.NEQ:
		return s.NotExpr(s.EqExpr(left, right))
	case token.LEQ:
		return s.BvLeExpr(left, right)
	case token.GEQ:
		return s.BvGeExpr(left, right)
	case token.EQL:
		return s.EqExpr(left, right)
	}

	return Expr(nil)
}

//Add adds constraints to a solver through a string in
func (s *Solver) Add(in string) {
	tr, _ := parser.ParseExpr(in)
	s.AssertFormula(s.Eval(tr))
}

//BitVec returns an Expr which is symbollic of string in with bit length bithWidth
func (s *Solver) BitVec(in string, bitWidth int) Expr {
	if uint(bitWidth) > s.width {
		s.width = uint(bitWidth)
	}

	return Expr(s.VarExpr(in, s.BvType(bitWidth)))
}

//Solve evaluates the assertions which returns values of symbollic inputs vars in int form
func (s *Solver) Solve(vars ...Expr) []uint64 {
	res := []uint64{}
	s.AssertFormula(s.TrueExpr())
	s.Query(s.NotExpr(s.TrueExpr()))
	whole := s.GetWholeCounterExample()

	for _, x := range vars {
		tmp := ExprString(s.GetTermFromCounterExample(x, whole))
		res = append(res, s.H2i(tmp))
	}

	return res
}

//Signed solve evaluates the assertions which returns values of symbollic inputs vars in int form
func (s *Solver) SSolve(vars ...Expr) []int64 {
	res := []int64{}
	s.AssertFormula(s.TrueExpr())
	s.Query(s.NotExpr(s.TrueExpr()))
	whole := s.GetWholeCounterExample()

	for _, x := range vars {
		tmp := ExprString(s.GetTermFromCounterExample(x, whole))
		res = append(res, s.SH2i(tmp))
	}

	return res
}
