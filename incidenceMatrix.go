package GraphBLAS

type Matrix [][]int

type IncidenceMatrix struct {
	matrix Matrix
}

func NewIncidenceMatrix(vertices int, edges int) *IncidenceMatrix {
	m := make(Matrix, vertices)
	for i := 0; i < vertices; i++ {
		m[i] = make([]int, edges)
	}
	return &IncidenceMatrix{matrix: m}

}

func NewIncidenceMatrixFromGraph(g *Graph) *IncidenceMatrix {
	m := NewIncidenceMatrix(len(g.Vertices), len(g.Edges))
	mat := m.matrix
	i := 0
	for _, e := range g.Edges {
		mat[e.Source().Index()][i] = -1
		mat[e.Sink().Index()][i] = 1
	}

	return m
}

func (i *IncidenceMatrix) Element(x int, y int) int {
	return i.matrix[x][y]
}