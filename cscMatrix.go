package GraphBLAS

import (
	"fmt"
)

// CSCMatrix compressed storage by columns (CSC)
type CSCMatrix struct {
	r        int // number of rows in the sparse matrix
	c        int // number of columns in the sparse matrix
	values   []float64
	rows     []int
	colStart []int
}

// NewCSCMatrix returns a GraphBLAS.CSCMatrix.
func NewCSCMatrix(r, c int) *CSCMatrix {
	return newCSCMatrix(r, c, 0)
}

func newCSCMatrix(r, c, l int) *CSCMatrix {
	s := &CSCMatrix{
		r:        r,
		c:        c,
		values:   make([]float64, l),
		rows:     make([]int, l),
		colStart: make([]int, c+1),
	}

	return s
}

func (s *CSCMatrix) Columns() int {
	return s.c
}

func (s *CSCMatrix) Rows() int {
	return s.r
}

// At returns the value of a matrix element at r-th, c-th.
func (s *CSCMatrix) At(r, c int) (float64, error) {
	if r < 0 || r >= s.r {
		return 0, fmt.Errorf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.c {
		return 0, fmt.Errorf("Column '%+v' is invalid", c)
	}

	pointerStart, pointerEnd := s.rowIndex(r, c)

	if pointerStart < pointerEnd && s.rows[pointerStart] == r {
		return s.values[pointerStart], nil
	}

	return 0, nil
}

func (s *CSCMatrix) Set(r, c int, value float64) error {
	if r < 0 || r >= s.r {
		return fmt.Errorf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.c {
		return fmt.Errorf("Column '%+v' is invalid", c)
	}

	pointerStart, pointerEnd := s.rowIndex(r, c)

	if pointerStart < pointerEnd && s.rows[pointerStart] == r {
		if value == 0 {
			s.remove(pointerStart, c)
		} else {
			s.values[pointerStart] = value
		}
	} else {
		s.insert(pointerStart, r, c, value)
	}

	return nil
}

func (s *CSCMatrix) ColumnsAt(c int) ([]float64, error) {
	if c < 0 || c >= s.c {
		return nil, fmt.Errorf("Column '%+v' is invalid", c)
	}

	start := s.colStart[c]
	end := s.colStart[c+1]

	columns := make([]float64, s.r)
	for i := start; i < end; i++ {
		columns[s.rows[i]] = s.values[i]
	}

	return columns, nil
}

func (s *CSCMatrix) RowsAt(r int) ([]float64, error) {
	if r < 0 || r >= s.r {
		return nil, fmt.Errorf("Row '%+v' is invalid", r)
	}

	rows := make([]float64, s.c)

	for c := range s.colStart[:s.c] {
		pointerStart, _ := s.rowIndex(r, c)
		rows[c] = s.values[pointerStart]
	}

	return rows, nil
}

func (s *CSCMatrix) insert(pointer, r, c int, value float64) {
	if value == 0 {
		return
	}

	s.rows = append(s.rows[:pointer], append([]int{r}, s.rows[pointer:]...)...)
	s.values = append(s.values[:pointer], append([]float64{value}, s.values[pointer:]...)...)

	for i := c + 1; i <= s.c; i++ {
		s.colStart[i]++
	}
}

func (s *CSCMatrix) remove(pointer, c int) {
	s.rows = append(s.rows[:pointer], s.rows[pointer+1:]...)
	s.values = append(s.values[:pointer], s.values[pointer+1:]...)

	for i := c + 1; i <= s.c; i++ {
		s.colStart[i]--
	}
}

func (s *CSCMatrix) rowIndex(r, c int) (int, int) {

	start := s.colStart[c]
	end := s.colStart[c+1]

	if start-end == 0 {
		return start, end
	}

	if r > s.rows[end-1] {
		return end, end
	}

	for start < end {
		p := (start + end) / 2
		if s.rows[p] > r {
			end = p
		} else if s.rows[p] < r {
			start = p + 1
		} else {
			return p, end
		}
	}

	return start, end
}

func (s *CSCMatrix) Copy() Matrix {
	return s.copy(func(value float64) float64 {
		return value
	})
}

func (s *CSCMatrix) copy(action func(float64) float64) *CSCMatrix {
	matrix := newCSCMatrix(s.r, s.c, len(s.values))

	for i := range s.values {
		matrix.values[i] = action(s.values[i])
		matrix.rows[i] = s.rows[i]
	}

	for i := range s.colStart {
		matrix.colStart[i] = s.colStart[i]
	}

	return matrix
}

// Scalar multiplication
func (s *CSCMatrix) Scalar(alpha float64) Matrix {
	return s.copy(func(value float64) float64 {
		return alpha * value
	})
}

// Multiply multiplies a Matrix structure by another Matrix structure.
func (s *CSCMatrix) Multiply(m Matrix) (Matrix, error) {
	if s.Rows() != m.Columns() {
		return nil, fmt.Errorf("Can not multiply matrices found length miss match %+v, %+v", s.Rows(), m.Columns())
	}

	matrix := newCSCMatrix(s.Rows(), m.Columns(), 0)

	for r := 0; r < s.Rows(); r++ {
		rows, _ := s.RowsAt(r)

		for c := 0; c < m.Columns(); c++ {
			column, _ := m.ColumnsAt(c)

			sum := 0.0
			for l := 0; l < len(rows); l++ {
				sum += rows[l] * column[l]
			}

			matrix.Set(r, c, sum)
		}

	}

	return matrix, nil
}

// Add addition of a Matrix structure by another Matrix structure.
func (s *CSCMatrix) Add(m Matrix) (Matrix, error) {
	if s.Columns() != m.Columns() {
		return nil, fmt.Errorf("Column miss match %+v, %+v", s.Columns(), m.Columns())
	}

	if s.Rows() != m.Rows() {
		return nil, fmt.Errorf("Row miss match %+v, %+v", s.Rows(), m.Rows())
	}

	matrix := newCSCMatrix(s.Rows(), m.Columns(), 0)

	for c := 0; c < s.Columns(); c++ {
		sColumn, _ := s.ColumnsAt(c)

		mColumn, _ := m.ColumnsAt(c)

		for r := 0; r < s.Rows(); r++ {
			matrix.Set(r, c, sColumn[r]+mColumn[r])
		}
	}

	return matrix, nil
}

// Subtract subtracts one matrix from another.
func (s *CSCMatrix) Subtract(m Matrix) (Matrix, error) {
	if s.Columns() != m.Columns() {
		return nil, fmt.Errorf("Column miss match %+v, %+v", s.Columns(), m.Columns())
	}

	if s.Rows() != m.Rows() {
		return nil, fmt.Errorf("Row miss match %+v, %+v", s.Rows(), m.Rows())
	}

	matrix := newCSCMatrix(s.Rows(), m.Columns(), 0)

	for c := 0; c < s.Columns(); c++ {
		sColumn, _ := s.ColumnsAt(c)

		mColumn, _ := m.ColumnsAt(c)

		for r := 0; r < s.Rows(); r++ {
			matrix.Set(r, c, sColumn[r]-mColumn[r])
		}
	}

	return matrix, nil
}

// Negative the negative of a matrix.
func (s *CSCMatrix) Negative() Matrix {
	return s.copy(func(value float64) float64 {
		return -value
	})
}
