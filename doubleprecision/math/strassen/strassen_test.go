// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package strassen_test

import (
	"testing"

	"golang.org/x/net/context"

	"github.com/rossmerr/graphblas/doubleprecision"
	"github.com/rossmerr/graphblas/doubleprecision/math/strassen"
)

func TestMatrix_Multiply(t *testing.T) {
	setup := func(m doubleprecision.Matrix) {
		m.Set(0, 0, 1)
		m.Set(0, 1, 2)
		m.Set(0, 2, 3)
		m.Set(0, 3, 4)
		m.Set(1, 0, 1)
		m.Set(1, 1, 2)
		m.Set(1, 2, 3)
		m.Set(1, 3, 4)
		m.Set(2, 0, 1)
		m.Set(2, 1, 2)
		m.Set(2, 2, 3)
		m.Set(2, 3, 4)
		m.Set(3, 0, 1)
		m.Set(3, 1, 2)
		m.Set(3, 2, 3)
		m.Set(3, 3, 4)
	}

	matrix := doubleprecision.NewDenseMatrix(4, 4)
	matrix.Set(0, 0, 1)
	matrix.Set(0, 1, 2)
	matrix.Set(0, 2, 3)
	matrix.Set(0, 3, 4)
	matrix.Set(1, 0, 1)
	matrix.Set(1, 1, 2)
	matrix.Set(1, 2, 3)
	matrix.Set(1, 3, 4)
	matrix.Set(2, 0, 1)
	matrix.Set(2, 1, 2)
	matrix.Set(2, 2, 3)
	matrix.Set(2, 3, 4)
	matrix.Set(3, 0, 1)
	matrix.Set(3, 1, 2)
	matrix.Set(3, 2, 3)
	matrix.Set(3, 3, 4)

	want := doubleprecision.NewDenseMatrix(4, 4)
	want.Set(0, 0, 10)
	want.Set(0, 1, 20)
	want.Set(0, 2, 30)
	want.Set(0, 3, 40)
	want.Set(1, 0, 10)
	want.Set(1, 1, 20)
	want.Set(1, 2, 30)
	want.Set(1, 3, 40)
	want.Set(2, 0, 10)
	want.Set(2, 1, 20)
	want.Set(2, 2, 30)
	want.Set(2, 3, 40)
	want.Set(3, 0, 10)
	want.Set(3, 1, 20)
	want.Set(3, 2, 30)
	want.Set(3, 3, 40)

	tests := []struct {
		name string
		s    doubleprecision.Matrix
	}{
		{
			name: "DenseMatrix",
			s:    doubleprecision.NewDenseMatrix(4, 4),
		},
		{
			name: "CSCMatrix",
			s:    doubleprecision.NewCSCMatrix(4, 4),
		},
		{
			name: "CSRMatrix",
			s:    doubleprecision.NewCSRMatrix(4, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(tt.s)
			if got := strassen.MultiplyCrossoverPoint(context.Background(), tt.s, matrix, 2); !got.Equal(want) {
				t.Errorf("%+v Multiply = got %+v, want %+v", tt.name, got, want)
			}
		})
	}
}
