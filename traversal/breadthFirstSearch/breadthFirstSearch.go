// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package breadthFirstSearch

import (
	GraphBLAS "github.com/RossMerr/Caudex.GraphBLAS"
	"golang.org/x/net/context"
)

// BreadthFirstSearch a breadth-first search v is the source
func BreadthFirstSearch(ctx context.Context, a GraphBLAS.Matrix, s int, c func(GraphBLAS.Vector) bool) GraphBLAS.Matrix {
	n := a.Rows()
	// vertices visited in each level
	var frontier GraphBLAS.Vector = GraphBLAS.NewDenseVector(n)
	frontier.SetVec(s, 1)

	//visited := frontier.Copy().(GraphBLAS.Vector)

	// result
	v := GraphBLAS.NewDenseVector(n)

	GraphBLAS.MatrixVectorMultiply(ctx, a, frontier, nil, v)

	// // level in BFS traversal
	// d := 0

	// // when some successor found
	// for {
	// 	d++

	// 	GraphBLAS.MatrixVectorMultiply(ctx, a, frontier, v)

	// 	if c(v) {
	// 		break
	// 	}

	// 	GraphBLAS.ElementWiseVectorAdd(ctx, visited, v, visited)
	// 	frontier = v
	// 	//todo need to skip the visited edges
	// }

	return v
}
