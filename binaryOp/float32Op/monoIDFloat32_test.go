// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package float32Op_test

import (
	"testing"

	"github.com/RossMerr/Caudex.GraphBLAS/binaryOp/float32Op"
)

func Test_Reduce(t *testing.T) {
	done := make(chan struct{})
	slice := make(chan float32)
	defer close(slice)
	defer close(done)

	monoID := float32Op.NewMonoIDFloat32(1, float32Op.Addition)

	result := monoID.Reduce(done, slice)

	zero := monoID.Zero()

	if zero != 1 {
		t.Errorf("Zero = %+v want %+v", zero, 1)
	}

	slice <- 1
	done <- struct{}{}
	for out := range result {
		if 2 != out {
			t.Errorf("MonoIDBool = %+v, want %+v", out, 2)
		}
	}
}
