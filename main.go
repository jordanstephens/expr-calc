package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

type tokenKind int

const (
	op tokenKind = iota
	val
)

type operatorKind int

const (
	mul operatorKind = iota
	div
	add
	sub
)

type operator struct {
	kind       operatorKind
	precedence int
}

var operators = map[string]operator{
	"*": {kind: mul, precedence: 1},
	"/": {kind: div, precedence: 1},
	"+": {kind: add, precedence: 0},
	"-": {kind: sub, precedence: 0},
}

type token struct {
	kind     tokenKind
	operator *operator
	value    *float64
}

func makeToken(input string) (*token, error) {
	operator, known := operators[input]
	if known {
		return &token{
			kind:     op,
			operator: &operator,
			value:    nil,
		}, nil
	}

	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return nil, err
	}

	return &token{
		kind:     val,
		operator: nil,
		value:    &value,
	}, nil
}

func pop(stack []*token) (*token, []*token) {
	i := len(stack) - 1
	return stack[i], stack[0:i]
}

func parse(input []string) ([]*token, error) {
	output := make([]*token, 0)
	operators := make([]*token, 0)

	for i := 0; i < len(input); i++ {
		token, err := makeToken(input[i])
		if err != nil {
			return nil, err
		}
		if token.kind == val {
			output = append(output, token)
		} else if token.kind == op {
			for {
				opslen := len(operators)
				if opslen == 0 {
					break
				}

				lastop := operators[opslen-1]
				if token.operator.precedence > lastop.operator.precedence {
					break
				}

				op, ops := pop(operators)
				operators = ops
				output = append(output, op)
			}
			operators = append(operators, token)
		} else {
			return nil, errors.New("Unknown token kind")
		}
	}

	for i := len(operators) - 1; i >= 0; i-- {
		output = append(output, operators[i])
	}

	return output, nil
}

func pop2(stack []float64) (float64, float64, []float64) {
	i := len(stack) - 1
	return stack[i], stack[i-1], stack[0 : i-1]
}

func compute(op operator, left float64, right float64) (float64, error) {
	if op.kind == mul {
		return left * right, nil
	}
	if op.kind == div {
		return left / right, nil
	}
	if op.kind == add {
		return left + right, nil
	}
	if op.kind == sub {
		return left - right, nil
	}

	return math.NaN(), errors.New("Unknown operator")
}

func eval(tokens []*token) (float64, error) {
	stack := make([]float64, 0)
	for i := 0; i < len(tokens); i++ {
		token := *tokens[i]
		if token.kind == op {
			if len(stack) < 2 {
				return math.NaN(), errors.New("invalid eval stack")
			}
			left, right, rest := pop2(stack)
			stack = rest
			result, err := compute(*token.operator, left, right)
			if err != nil {
				return math.NaN(), err
			}
			stack = append(stack, result)
		} else {
			stack = append(stack, *token.value)
		}
	}

	if len(stack) > 1 {
		return math.NaN(), errors.New("invalid result stack")
	}

	return stack[0], nil
}

func run(input []string) (float64, error) {
	tokens, err := parse(input)

	if err != nil {
		return math.NaN(), err
	}

	result, err := eval(tokens)
	return result, err
}

func main() {
	terms := os.Args[1:]
	result, err := run(terms)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
