// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package boolop_test

import (
	"testing"

	"github.com/rossmerr/graphblas/binaryop"
	"github.com/rossmerr/graphblas/binaryop/boolop"
)

func Test_LOR(t *testing.T) {
	tests := []struct {
		name   string
		s      binaryop.BinaryOp
		in1    bool
		in2    bool
		result bool
	}{
		{
			name:   "1",
			s:      boolop.LOR,
			in1:    true,
			in2:    true,
			result: true,
		},
		{
			name:   "2",
			s:      boolop.LOR,
			in1:    false,
			in2:    true,
			result: true,
		},
		{
			name:   "3",
			s:      boolop.LOR,
			in1:    true,
			in2:    false,
			result: true,
		},
		{
			name:   "4",
			s:      boolop.LOR,
			in1:    false,
			in2:    false,
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if op, ok := tt.s.(boolop.BinaryOpBool); ok {
				if tt.result != op.Apply(tt.in1, tt.in2) {
					t.Errorf("%+v LOR = %+v, want %+v", tt.name, !tt.result, tt.result)
				}
			} else {
				t.Errorf("%+v not a BinaryOpBool", tt.name)
			}
		})
	}
}

func Test_LAND(t *testing.T) {
	tests := []struct {
		name   string
		s      binaryop.BinaryOp
		in1    bool
		in2    bool
		result bool
	}{
		{
			name:   "1",
			s:      boolop.LAND,
			in1:    true,
			in2:    true,
			result: true,
		},
		{
			name:   "2",
			s:      boolop.LAND,
			in1:    false,
			in2:    true,
			result: false,
		},
		{
			name:   "3",
			s:      boolop.LAND,
			in1:    true,
			in2:    false,
			result: false,
		},
		{
			name:   "4",
			s:      boolop.LAND,
			in1:    false,
			in2:    false,
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if op, ok := tt.s.(boolop.BinaryOpBool); ok {
				if tt.result != op.Apply(tt.in1, tt.in2) {
					t.Errorf("%+v LAND = %+v, want %+v", tt.name, !tt.result, tt.result)
				}
			} else {
				t.Errorf("%+v not a BinaryOpBool", tt.name)
			}
		})
	}
}

func Test_LXOR(t *testing.T) {
	tests := []struct {
		name   string
		s      binaryop.BinaryOp
		in1    bool
		in2    bool
		result bool
	}{
		{
			name:   "1",
			s:      boolop.LXOR,
			in1:    true,
			in2:    true,
			result: false,
		},
		{
			name:   "2",
			s:      boolop.LXOR,
			in1:    false,
			in2:    true,
			result: true,
		},
		{
			name:   "3",
			s:      boolop.LXOR,
			in1:    true,
			in2:    false,
			result: true,
		},
		{
			name:   "4",
			s:      boolop.LXOR,
			in1:    false,
			in2:    false,
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if op, ok := tt.s.(boolop.BinaryOpBool); ok {
				if tt.result != op.Apply(tt.in1, tt.in2) {
					t.Errorf("%+v LXOR = %+v, want %+v", tt.name, !tt.result, tt.result)
				}
			} else {
				t.Errorf("%+v not a BinaryOpBool", tt.name)
			}
		})
	}
}

func Test_Associative(t *testing.T) {
	tests := []struct {
		name   string
		s      binaryop.BinaryOp
		a      bool
		b      bool
		c      bool
		result bool
	}{
		{
			name:   "1",
			s:      boolop.LOR,
			a:      true,
			b:      true,
			c:      true,
			result: true,
		},
		{
			name:   "2",
			s:      boolop.LOR,
			a:      false,
			b:      true,
			c:      true,
			result: true,
		},
		{
			name:   "3",
			s:      boolop.LOR,
			a:      true,
			b:      false,
			c:      true,
			result: true,
		},
		{
			name:   "4",
			s:      boolop.LOR,
			a:      true,
			b:      true,
			c:      false,
			result: true,
		},
		{
			name:   "5",
			s:      boolop.LOR,
			a:      false,
			b:      false,
			c:      false,
			result: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if op, ok := tt.s.(boolop.BinaryOpBool); ok {
				result := op.Apply(op.Apply(tt.a, tt.b), tt.c) == op.Apply(tt.a, op.Apply(tt.b, tt.c))
				if tt.result != result {
					t.Errorf("%+v Associative = %+v, want %+v", tt.name, !tt.result, tt.result)
				}
			} else {
				t.Errorf("%+v not a Associative", tt.name)
			}
		})
	}
}

func Test_Commutative(t *testing.T) {
	tests := []struct {
		name   string
		s      binaryop.BinaryOp
		a      bool
		b      bool
		c      bool
		result bool
	}{
		{
			name:   "1",
			s:      boolop.LOR,
			a:      true,
			b:      true,
			result: false,
		},
		{
			name:   "2",
			s:      boolop.LOR,
			a:      false,
			b:      true,
			result: false,
		},
		{
			name:   "3",
			s:      boolop.LOR,
			a:      true,
			b:      false,
			result: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if op, ok := tt.s.(boolop.BinaryOpBool); ok {
				result := op.Apply(tt.a, tt.b) != op.Apply(tt.b, tt.a)
				if tt.result != result {
					t.Errorf("%+v Commutative = %+v, want %+v", tt.name, !tt.result, tt.result)
				}
			} else {
				t.Errorf("%+v not a Commutative", tt.name)
			}
		})
	}
}

func Test_Operator(t *testing.T) {
	boolop.LOR.Operator()
}

func Test_BinaryOp(t *testing.T) {
	boolop.LOR.BinaryOp()
}

func Test_Semigroup(t *testing.T) {
	boolop.LOR.Semigroup()
}
