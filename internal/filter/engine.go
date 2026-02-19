package filter

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ersanisk/sieve/pkg/logentry"
)

// FieldValue represents a field name for evaluation.
type FieldValue struct {
	Field string
}

// Evaluator is a function that evaluates a field or literal value against an entry.
type Evaluator func(val FieldValue) (any, error)

// CompiledFilter is a compiled expression ready for evaluation.
type CompiledFilter struct {
	expr Expr
}

// Evaluate evaluates the compiled filter against an entry.
func (c *CompiledFilter) Evaluate(entry logentry.Entry) (bool, error) {
	evalFunc := func(fv FieldValue) (any, error) {
		switch fv.Field {
		case "message", "msg":
			return entry.Message, nil
		case "caller", "source":
			return entry.Caller, nil
		case "level":
			if entry.Fields != nil {
				if val, ok := entry.Fields["level"]; ok {
					return val, nil
				}
			}
			return int(entry.Level), nil
		default:
			if val, ok := entry.GetField(fv.Field); ok {
				return val, nil
			}
			return nil, nil
		}
	}
	return c.expr.Eval(evalFunc)
}

// Compile compiles an expression into an executable filter.
func Compile(expr Expr) (*CompiledFilter, error) {
	return &CompiledFilter{
		expr: expr,
	}, nil
}

// ByLevel creates a filter that matches entries at or above the specified level.
func ByLevel(minLevel logentry.Level) *CompiledFilter {
	levelValue := 0
	switch minLevel {
	case logentry.Debug:
		levelValue = 10
	case logentry.Info:
		levelValue = 30
	case logentry.Warn:
		levelValue = 40
	case logentry.Error:
		levelValue = 50
	case logentry.Fatal:
		levelValue = 60
	}

	expr := BinaryOp{
		Left:  FieldAccess{Field: "level"},
		Op:    OpGreaterEqual,
		Right: Literal{Value: levelValue},
	}

	return &CompiledFilter{expr: expr}
}

// ByValue creates a filter that matches entries with a specific field value.
func ByValue(field string, value any, op Operator) *CompiledFilter {
	return &CompiledFilter{
		expr: BinaryOp{
			Left:  FieldAccess{Field: field},
			Op:    op,
			Right: Literal{Value: value},
		},
	}
}

// compareValues compares two values using the specified operator.
func compareValues(left, right any, op Operator) (bool, error) {
	if left == nil && right == nil {
		return op == OpEqual, nil
	}
	if left == nil || right == nil {
		return op == OpNotEqual, nil
	}

	leftType := reflect.TypeOf(left)
	rightType := reflect.TypeOf(right)

	if leftType != rightType {
		var err error
		left, right, err = coerceTypes(left, right)
		if err != nil {
			return false, err
		}
	}

	switch op {
	case OpEqual:
		return reflect.DeepEqual(left, right), nil
	case OpNotEqual:
		return !reflect.DeepEqual(left, right), nil
	case OpGreater, OpLess, OpGreaterEqual, OpLessEqual:
		return compareNumeric(left, right, op)
	case OpContains:
		return compareContains(left, right), nil
	case OpMatches:
		return compareMatches(left, right)
	default:
		return false, fmt.Errorf("unsupported operator: %s", op)
	}
}

func coerceTypes(left, right any) (any, any, error) {
	leftStr, leftIsStr := left.(string)
	rightStr, rightIsStr := right.(string)

	if leftIsStr && rightIsStr {
		return left, right, nil
	}

	leftFloat, leftIsFloat := toFloat(left)
	rightFloat, rightIsFloat := toFloat(right)

	if leftIsFloat && rightIsFloat {
		return leftFloat, rightFloat, nil
	}

	leftInt, leftIsInt := toInt(left)
	rightInt, rightIsInt := toInt(right)

	if leftIsInt && rightIsInt {
		return leftInt, rightInt, nil
	}

	if leftIsStr {
		if rightIsFloat {
			return leftStr, fmt.Sprintf("%v", right), nil
		}
		if rightIsInt {
			return leftStr, fmt.Sprintf("%v", right), nil
		}
	}

	if rightIsStr {
		if leftIsFloat {
			return fmt.Sprintf("%v", left), rightStr, nil
		}
		if leftIsInt {
			return fmt.Sprintf("%v", left), rightStr, nil
		}
	}

	return nil, nil, fmt.Errorf("cannot coerce types: %T and %T", left, right)
}

func toFloat(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		f, err := strconv.ParseFloat(v, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func toInt(val any) (int, bool) {
	switch v := val.(type) {
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case float64:
		if v == float64(int(v)) {
			return int(v), true
		}
		return 0, false
	case string:
		i, err := strconv.Atoi(v)
		return i, err == nil
	default:
		return 0, false
	}
}

func compareNumeric(left, right any, op Operator) (bool, error) {
	leftFloat, leftOk := toFloat(left)
	rightFloat, rightOk := toFloat(right)

	if !leftOk || !rightOk {
		return false, fmt.Errorf("cannot compare non-numeric values")
	}

	switch op {
	case OpGreater:
		return leftFloat > rightFloat, nil
	case OpLess:
		return leftFloat < rightFloat, nil
	case OpGreaterEqual:
		return leftFloat >= rightFloat, nil
	case OpLessEqual:
		return leftFloat <= rightFloat, nil
	default:
		return false, fmt.Errorf("invalid comparison operator: %s", op)
	}
}

func compareContains(left, right any) bool {
	leftStr, leftOk := left.(string)
	rightStr, rightOk := right.(string)

	if !leftOk || !rightOk {
		return false
	}

	return strings.Contains(leftStr, rightStr)
}

func compareMatches(left, right any) (bool, error) {
	leftStr, leftOk := left.(string)
	rightStr, rightOk := right.(string)

	if !leftOk || !rightOk {
		return false, nil
	}

	re, err := regexp.Compile(rightStr)
	if err != nil {
		return false, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return re.MatchString(leftStr), nil
}
