// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package reduced_test

import (
	"testing"

	"github.com/RossMerr/Caudex.GraphBLAS/doublePrecision"
	"github.com/RossMerr/Caudex.GraphBLAS/doublePrecision/math/reduced"
)

func TestMatrix_Reduced(t *testing.T) {

	setup := func(m doublePrecision.Matrix) {
		m.Set(0, 0, 1)
		m.Set(0, 1, 2)
		m.Set(0, 2, -1)
		m.Set(0, 3, -4)

		m.Set(1, 0, 2)
		m.Set(1, 1, 3)
		m.Set(1, 2, -1)
		m.Set(1, 3, -11)

		m.Set(2, 0, -2)
		m.Set(2, 1, 0)
		m.Set(2, 2, -3)
		m.Set(2, 3, 22)
	}

	want := doublePrecision.NewDenseMatrix(3, 4)
	want.Set(0, 0, 1)
	want.Set(0, 1, 0)
	want.Set(0, 2, 0)
	want.Set(0, 3, -8)

	want.Set(1, 0, -0)
	want.Set(1, 1, 1)
	want.Set(1, 2, 0)
	want.Set(1, 3, 1)

	want.Set(2, 0, -0)
	want.Set(2, 1, -0)
	want.Set(2, 2, 1)
	want.Set(2, 3, -2)

	tests := []struct {
		name string
		s    doublePrecision.Matrix
	}{
		{
			name: "DenseMatrix",
			s:    doublePrecision.NewDenseMatrix(3, 4),
		},
		{
			name: "CSCMatrix",
			s:    doublePrecision.NewCSCMatrix(3, 4),
		},
		{
			name: "CSRMatrix",
			s:    doublePrecision.NewCSRMatrix(3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup(tt.s)
			got := reduced.Reduced(tt.s)
			if got.NotEqual(want) {
				t.Errorf("%+v Reduced = %+v, want %+v", tt.name, got, want)
			}
		})
	}
}
