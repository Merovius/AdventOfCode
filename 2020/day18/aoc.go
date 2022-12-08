package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(0)

	var (
		total1 int
		total2 int
	)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		e, err := ParseExpression(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		total1 += e.Eval()
		e, err = ParseExpressionCorrect(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		total2 += e.Eval()
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Total part 1:", total1)
	fmt.Println("Total part 2:", total2)
}

type Expr struct {
	V  int
	L  *Expr
	R  *Expr
	Op byte
}

func ParseExpression(str string) (*Expr, error) {
	var stack []*Expr
	push := func(e *Expr) {
		stack = append(stack, e)
	}
	pop := func() (e *Expr) {
		stack, e = stack[:len(stack)-1], stack[len(stack)-1]
		return e
	}
	peek := func() *Expr {
		return stack[len(stack)-1]
	}
	str = strings.TrimSpace(str)
	for len(str) > 0 {
		var t string
		t, str = split(str)
		switch t {
		case "*", "+":
			if len(stack) == 0 {
				return nil, fmt.Errorf("unexpected token %q at start of line", t)
			}
			top := pop()
			if (top.L == nil) != (top.R == nil) {
				return nil, fmt.Errorf("unexpected token %q", t)
			}
			e := &Expr{L: top, Op: t[0]}
			push(e)
		case "(":
			push(nil)
		case ")":
			if len(stack) < 2 {
				return nil, fmt.Errorf("unbalanced parenthesis")
			}
			e := pop()
			if (e.L == nil) != (e.R == nil) {
				return nil, fmt.Errorf("unexpected token %q", t)
			}
			if pop() != nil {
				return nil, fmt.Errorf("unbalanced parenthesis")
			}
			if len(stack) == 0 || peek() == nil {
				push(e)
				continue
			}
			top := pop()
			if top.L == nil || top.R != nil {
				return nil, errors.New("unbalanced expression")
			}
			top.R = e
			push(top)
		default:
			v, err := strconv.Atoi(t)
			if err != nil {
				return nil, err
			}
			e := &Expr{V: v}
			if len(stack) == 0 || peek() == nil {
				push(e)
				continue
			}
			top := pop()
			if top.L == nil || top.R != nil {
				return nil, errors.New("unbalanced expression")
			}
			top.R = e
			push(top)
		}
	}
	if len(stack) != 1 {
		return nil, errors.New("unbalanced expression")
	}
	return stack[0], nil
}

type stack []string

func (s *stack) push(e string) {
	*s = append(*s, e)
}

func (s *stack) pop() (e string) {
	*s, e = (*s)[:len(*s)-1], (*s)[len(*s)-1]
	return e
}

func (s *stack) peek() string {
	return (*s)[len(*s)-1]
}

func ParseExpressionCorrect(s string) (*Expr, error) {
	var (
		opStack  stack
		outStack stack
	)
	s = strings.TrimSpace(s)
	for len(s) > 0 {
		var tok string
		tok, s = split(s)
		switch tok {
		case "+", "*":
			for {
				if len(opStack) == 0 || !isOp(opStack.peek()) {
					break
				}
				if tok == "+" && opStack.peek() == "*" {
					break
				}
				outStack.push(opStack.pop())
			}
			opStack.push(tok)
		case "(":
			opStack.push(tok)
		case ")":
			for len(opStack) > 0 && opStack.peek() != "(" {
				outStack.push(opStack.pop())
			}
			if len(opStack) == 0 {
				return nil, errors.New("unbalanced parenthesis")
			}
			opStack.pop()
		default:
			outStack.push(tok)
		}
	}
	for len(opStack) > 0 {
		outStack.push(opStack.pop())
	}
	var build func() (*Expr, error)
	build = func() (e *Expr, err error) {
		if len(outStack) == 0 {
			return nil, errors.New("unbalanced expression")
		}
		tok := outStack.pop()
		switch tok {
		case "+", "*":
			r, err := build()
			if err != nil {
				return nil, err
			}
			l, err := build()
			if err != nil {
				return nil, err
			}
			return &Expr{L: l, R: r, Op: tok[0]}, nil
		default:
			v, err := strconv.Atoi(tok)
			if err != nil {
				return nil, err
			}
			return &Expr{V: v}, nil
		}
	}
	e, err := build()
	if err != nil {
		return nil, err
	}
	if len(outStack) != 0 {
		return nil, errors.New("unbalanced expressions")
	}
	return e, nil
}

func (e *Expr) Eval() int {
	if e.L == nil || e.R == nil {
		return e.V
	}
	switch e.Op {
	case '+':
		return e.L.Eval() + e.R.Eval()
	case '*':
		return e.L.Eval() * e.R.Eval()
	default:
		panic(fmt.Errorf("unknown operator %q", e.Op))
	}
}

func (e *Expr) String() string {
	if e == nil {
		return "<nil>"
	}
	if e.L == nil && e.R == nil {
		return fmt.Sprintf("%d", e.V)
	}
	return fmt.Sprintf("(%v %q %v)", e.L, e.Op, e.R)
}

func split(s string) (tok, rest string) {
	switch s[0] {
	case '(', ')', '+', '*':
		return string(s[0]), strings.TrimSpace(s[1:])
	}
	for i := 0; i < len(s); i++ {
		if !(s[i] >= '0' && s[i] <= '9') {
			return s[:i], strings.TrimSpace(s[i:])
		}
	}
	return s, ""
}

func isOp(s string) bool {
	return s == "+" || s == "*"
}
