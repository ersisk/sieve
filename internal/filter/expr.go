package filter

import (
	"fmt"
	"strings"
	"unicode"
)

// Operator represents a comparison or logical operator.
type Operator int

const (
	OpUnknown Operator = iota
	OpEqual
	OpNotEqual
	OpGreater
	OpLess
	OpGreaterEqual
	OpLessEqual
	OpContains
	OpMatches
	OpAnd
	OpOr
	OpNot
)

func (op Operator) String() string {
	switch op {
	case OpEqual:
		return "=="
	case OpNotEqual:
		return "!="
	case OpGreater:
		return ">"
	case OpLess:
		return "<"
	case OpGreaterEqual:
		return ">="
	case OpLessEqual:
		return "<="
	case OpContains:
		return "contains"
	case OpMatches:
		return "matches"
	case OpAnd:
		return "and"
	case OpOr:
		return "or"
	case OpNot:
		return "not"
	default:
		return "unknown"
	}
}

// Expr is the interface for all expression types.
type Expr interface {
	Eval(evalFunc Evaluator) (bool, error)
	String() string
}

// FieldAccess represents accessing a field from an entry (.field).
type FieldAccess struct {
	Field string
}

func (f FieldAccess) Eval(evalFunc Evaluator) (bool, error) {
	val, err := evalFunc(FieldValue{Field: f.Field}) //nolint:gosimple
	if err != nil {
		return false, err
	}
	if val == nil {
		return false, nil
	}
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	return true, nil
}

func (f FieldAccess) String() string {
	return "." + f.Field
}

// Literal represents a constant value.
type Literal struct {
	Value any
}

func (l Literal) Eval(_ Evaluator) (bool, error) {
	return true, nil
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

// BinaryOp represents a binary operation (left OP right).
type BinaryOp struct {
	Left  Expr
	Op    Operator
	Right Expr
}

func (b BinaryOp) Eval(evalFunc Evaluator) (bool, error) {
	leftResult, err := b.Left.Eval(evalFunc)
	if err != nil {
		return false, err
	}

	if b.Op == OpAnd || b.Op == OpOr {
		return b.evalLogical(leftResult, evalFunc)
	}

	return b.evalComparison(evalFunc)
}

func (b BinaryOp) evalLogical(leftResult bool, evalFunc Evaluator) (bool, error) {
	switch b.Op {
	case OpAnd:
		if !leftResult {
			return false, nil
		}
		return b.Right.Eval(evalFunc)
	case OpOr:
		if leftResult {
			return true, nil
		}
		return b.Right.Eval(evalFunc)
	default:
		return false, fmt.Errorf("invalid logical operator: %s", b.Op)
	}
}

func (b BinaryOp) evalComparison(evalFunc Evaluator) (bool, error) {
	leftVal, err := b.resolveValue(b.Left, evalFunc)
	if err != nil {
		return false, err
	}

	rightVal, err := b.resolveValue(b.Right, evalFunc)
	if err != nil {
		return false, err
	}

	return compareValues(leftVal, rightVal, b.Op)
}

func (b BinaryOp) resolveValue(expr Expr, evalFunc Evaluator) (any, error) {
	switch e := expr.(type) {
	case FieldAccess:
		return evalFunc(FieldValue{Field: e.Field}) //nolint:gosimple
	case Literal:
		return e.Value, nil
	default:
		return nil, fmt.Errorf("cannot resolve value from expression type %T", expr)
	}
}

func (b BinaryOp) String() string {
	return fmt.Sprintf("%s %s %s", b.Left, b.Op, b.Right)
}

// UnaryOp represents a unary operation (not expr).
type UnaryOp struct {
	Op   Operator
	Expr Expr
}

func (u UnaryOp) Eval(evalFunc Evaluator) (bool, error) {
	if u.Op == OpNot {
		result, err := u.Expr.Eval(evalFunc)
		return !result, err
	}
	return false, fmt.Errorf("unsupported unary operator: %s", u.Op)
}

func (u UnaryOp) String() string {
	return fmt.Sprintf("%s %s", u.Op, u.Expr)
}

// CompoundExpr represents a compound expression (expr and expr, expr or expr).
type CompoundExpr struct {
	Left  Expr
	Op    Operator
	Right Expr
}

func (c CompoundExpr) Eval(evalFunc Evaluator) (bool, error) {
	leftResult, err := c.Left.Eval(evalFunc)
	if err != nil {
		return false, err
	}

	switch c.Op {
	case OpAnd:
		if !leftResult {
			return false, nil
		}
		return c.Right.Eval(evalFunc)
	case OpOr:
		if leftResult {
			return true, nil
		}
		return c.Right.Eval(evalFunc)
	default:
		return false, fmt.Errorf("invalid compound operator: %s", c.Op)
	}
}

func (c CompoundExpr) String() string {
	return fmt.Sprintf("%s %s %s", c.Left, c.Op, c.Right)
}

// parser represents a filter expression parser.
type parser struct {
	input string
	pos   int
}

// Parse parses a filter expression string into an AST.
func Parse(input string) (Expr, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty expression")
	}

	p := &parser{input: input}
	expr, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	if p.pos < len(p.input) {
		return nil, fmt.Errorf("unexpected character at position %d: %c", p.pos, p.current())
	}

	return expr, nil
}

func (p *parser) parseExpression() (Expr, error) {
	return p.parseOr()
}

func (p *parser) parseOr() (Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.match("or") {
		p.skip(2)
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		return CompoundExpr{Left: left, Op: OpOr, Right: right}, nil
	}

	return left, nil
}

func (p *parser) parseAnd() (Expr, error) {
	left, err := p.parseComparison()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()

	if p.match("and") {
		p.skip(3)
		right, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		return CompoundExpr{Left: left, Op: OpAnd, Right: right}, nil
	}

	return left, nil
}

func (p *parser) parseComparison() (Expr, error) {
	p.skipWhitespace()

	if p.match("not") {
		p.skip(3)
		expr, err := p.parseComparison()
		if err != nil {
			return nil, err
		}
		return UnaryOp{Op: OpNot, Expr: expr}, nil
	}

	left, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()

	op, err := p.parseOperator()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()

	right, err := p.parseOperand()
	if err != nil {
		return nil, err
	}

	return BinaryOp{Left: left, Op: op, Right: right}, nil
}

func (p *parser) parseOperand() (Expr, error) {
	p.skipWhitespace()

	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	ch := p.input[p.pos]

	if ch == '.' {
		return p.parseFieldAccess()
	}

	if ch == '"' || ch == '\'' {
		return p.parseStringLiteral()
	}

	if unicode.IsDigit(rune(ch)) || ch == '-' {
		return p.parseNumericLiteral()
	}

	if ch == 't' || ch == 'f' {
		return p.parseBooleanLiteral()
	}

	return nil, fmt.Errorf("unexpected character: %c", ch)
}

func (p *parser) parseFieldAccess() (Expr, error) {
	if p.pos >= len(p.input) || p.input[p.pos] != '.' {
		return nil, fmt.Errorf("expected '.' for field access")
	}
	p.pos++

	start := p.pos
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_' {
			p.pos++
		} else {
			break
		}
	}

	field := p.input[start:p.pos]
	if field == "" {
		return nil, fmt.Errorf("empty field name")
	}

	return FieldAccess{Field: field}, nil
}

func (p *parser) parseStringLiteral() (Expr, error) {
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input")
	}

	quote := p.input[p.pos]
	if quote != '"' && quote != '\'' {
		return nil, fmt.Errorf("expected string literal")
	}
	p.pos++

	var sb strings.Builder
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == quote {
			p.pos++
			return Literal{Value: sb.String()}, nil
		}
		if ch == '\\' && p.pos+1 < len(p.input) {
			p.pos++
			ch = p.input[p.pos]
		}
		sb.WriteByte(ch)
		p.pos++
	}

	return nil, fmt.Errorf("unterminated string literal")
}

func (p *parser) parseNumericLiteral() (Expr, error) {
	start := p.pos
	if p.input[p.pos] == '-' {
		p.pos++
	}

	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if unicode.IsDigit(rune(ch)) || ch == '.' {
			p.pos++
		} else {
			break
		}
	}

	numStr := p.input[start:p.pos]
	var value any
	if strings.Contains(numStr, ".") {
		f, err := parseFloat(numStr)
		if err != nil {
			return nil, err
		}
		value = f
	} else {
		i, err := parseInt(numStr)
		if err != nil {
			return nil, err
		}
		value = i
	}

	return Literal{Value: value}, nil
}

func (p *parser) parseBooleanLiteral() (Expr, error) {
	if p.match("true") {
		p.skip(4)
		return Literal{Value: true}, nil
	}
	if p.match("false") {
		p.skip(5)
		return Literal{Value: false}, nil
	}

	return nil, fmt.Errorf("expected boolean literal")
}

func (p *parser) parseOperator() (Operator, error) {
	p.skipWhitespace()

	if p.match("==") {
		p.skip(2)
		return OpEqual, nil
	}
	if p.match("!=") {
		p.skip(2)
		return OpNotEqual, nil
	}
	if p.match(">=") {
		p.skip(2)
		return OpGreaterEqual, nil
	}
	if p.match("<=") {
		p.skip(2)
		return OpLessEqual, nil
	}
	if p.match(">") {
		p.skip(1)
		return OpGreater, nil
	}
	if p.match("<") {
		p.skip(1)
		return OpLess, nil
	}
	if p.match("contains") {
		p.skip(8)
		return OpContains, nil
	}
	if p.match("matches") {
		p.skip(7)
		return OpMatches, nil
	}

	return OpUnknown, fmt.Errorf("unknown operator at position %d", p.pos)
}

func (p *parser) skipWhitespace() {
	for p.pos < len(p.input) && unicode.IsSpace(rune(p.input[p.pos])) {
		p.pos++
	}
}

func (p *parser) skip(n int) {
	p.pos += n
}

func (p *parser) match(s string) bool {
	if p.pos+len(s) > len(p.input) {
		return false
	}
	return p.input[p.pos:p.pos+len(s)] == s
}

func (p *parser) current() byte {
	if p.pos >= len(p.input) {
		return 0
	}
	return p.input[p.pos]
}

func parseInt(s string) (int, error) {
	var result int
	n, err := fmt.Sscanf(s, "%d", &result)
	if n != 1 || err != nil {
		return 0, fmt.Errorf("invalid integer: %s", s)
	}
	return result, nil
}

func parseFloat(s string) (float64, error) {
	var result float64
	n, err := fmt.Sscanf(s, "%f", &result)
	if n != 1 || err != nil {
		return 0, fmt.Errorf("invalid float: %s", s)
	}
	return result, nil
}
