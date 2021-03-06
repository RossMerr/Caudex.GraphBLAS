// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package symmetric_test

import (
	"github.com/rossmerr/graphblas/doubleprecision"
	"testing"

	"github.com/rossmerr/graphblas/doubleprecision/math/symmetric"
)

func TestMatrix_Symmetric(t *testing.T) {

	setup := func(m doubleprecision.Matrix) {
		m.Set(0, 0, 1)
		m.Set(0, 1, 7)
		m.Set(0, 2, 3)
		m.Set(1, 0, 7)
		m.Set(1, 1, 4)
		m.Set(1, 2, -5)
		m.Set(2, 0, 3)
		m.Set(2, 1, -5)
		m.Set(2, 2, 6)
	}

	tests := []struct {
		name string
		s    doubleprecision.Matrix
		want bool
	}{
		{
			name: "DenseMatrix",
			s:    doubleprecision.NewDenseMatrix(3, 3),
			want: true,
		},
		{
			name: "CSCMatrix",
			s:    doubleprecision.NewCSCMatrix(3, 3),
			want: true,
		},
		{
			name: "CSRMatrix",
			s:    doubleprecision.NewCSRMatrix(3, 3),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(tt.s)
			if got := symmetric.Symmetric(tt.s); got != tt.want {
				t.Errorf("%+v Symmetric = %+v, want %+v", tt.name, got, tt.want)
			}
		})
	}
}
