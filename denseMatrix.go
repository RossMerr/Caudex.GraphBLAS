// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package GraphBLAS

import (
	"log"
)

// DenseMatrix a dense matrix
type DenseMatrix struct {
	c    int // number of rows in the sparse matrix
	r    int // number of columns in the sparse matrix
	data [][]float64
}

// NewDenseMatrix returns a GraphBLAS.DenseMatrix
func NewDenseMatrix(r, c int) *DenseMatrix {
	return newMatrix(r, c, nil)
}

// NewDenseMatrixFromArray returns a GraphBLAS.DenseMatrix
func NewDenseMatrixFromArray(data [][]float64) *DenseMatrix {
	r := len(data)
	c := len(data[0])
	s := &DenseMatrix{data: data, r: r, c: c}

	return s
}

func newMatrix(r, c int, initialise func([]float64, int)) *DenseMatrix {
	s := &DenseMatrix{data: make([][]float64, r), r: r, c: c}

	for i := 0; i < r; i++ {
		s.data[i] = make([]float64, c)

		if initialise != nil {
			initialise(s.data[i], i)
		}
	}

	return s
}

// Columns the number of columns of the matrix
func (s *DenseMatrix) Columns() int {
	return s.c
}

// Rows the number of rows of the matrix
func (s *DenseMatrix) Rows() int {
	return s.r
}

// Update does a At and Set on the matrix element at r-th, c-th
func (s *DenseMatrix) Update(r, c int, f func(float64) float64) {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	s.data[r][c] = f(s.data[r][c])

	return
}

// At returns the value of a matrix element at r-th, c-th
func (s *DenseMatrix) At(r, c int) float64 {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	return s.data[r][c]
}

// Set sets the value at r-th, c-th of the matrix
func (s *DenseMatrix) Set(r, c int, value float64) {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	s.data[r][c] = value
}

// ColumnsAt return the columns at c-th
func (s *DenseMatrix) ColumnsAt(c int) Vector {
	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	columns := NewDenseVector(s.r)

	for r := 0; r < s.r; r++ {
		columns.SetVec(r, s.data[r][c])
	}

	return columns
}

// RowsAt return the rows at r-th
func (s *DenseMatrix) RowsAt(r int) Vector {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	rows := NewDenseVector(s.c)
	for i := 0; i < s.c; i++ {
		rows.SetVec(i, s.data[r][i])
	}

	return rows
}

// Copy copies the matrix
func (s *DenseMatrix) Copy() Matrix {
	v := 0.0
	matrix := newMatrix(s.Rows(), s.Columns(), func(row []float64, r int) {
		for c := 0; c < s.Columns(); c++ {
			v = s.data[r][c]
			if v != 0.0 {
				row[c] = v
			} else {
				row[c] = v
			}
		}
	})

	return matrix
}

// Scalar multiplication of a matrix by alpha
func (s *DenseMatrix) Scalar(alpha float64) Matrix {
	return Scalar(s, alpha)
}

// Multiply multiplies a matrix by another matrix
func (s *DenseMatrix) Multiply(m Matrix) Matrix {
	matrix := newMatrix(s.Rows(), m.Columns(), nil)

	Multiply(s, m, matrix)
	return matrix
}

// Add addition of a matrix by another matrix
func (s *DenseMatrix) Add(m Matrix) Matrix {
	matrix := s.Copy()
	Add(s, m, matrix)
	return matrix
}

// Subtract subtracts one matrix from another matrix
func (s *DenseMatrix) Subtract(m Matrix) Matrix {
	matrix := m.Copy()
	Subtract(s, m, matrix)
	return matrix
}

// Negative the negative of a matrix
func (s *DenseMatrix) Negative() Matrix {
	matrix := s.Copy()
	Negative(s, matrix)
	return matrix
}

// Transpose swaps the rows and columns
func (s *DenseMatrix) Transpose() Matrix {
	matrix := newMatrix(s.Columns(), s.Rows(), nil)

	Transpose(s, matrix)
	return matrix
}

// Equal the two matrices are equal
func (s *DenseMatrix) Equal(m Matrix) bool {
	return Equal(s, m)
}

// NotEqual the two matrices are not equal
func (s *DenseMatrix) NotEqual(m Matrix) bool {
	return NotEqual(s, m)
}

// Size of the matrix
func (s *DenseMatrix) Size() int {
	return s.r * s.c
}

// Values the number of elements in the matrix
func (s *DenseMatrix) Values() int {
	return s.r * s.c
}

// Apply modifies edge weights by the UnaryOperator
// C ⊕= f(A)
func (s *DenseMatrix) Apply(u UnaryOperator) {
	for iterator := s.Map(); iterator.HasNext(); {
		iterator.Map(func(r, c int, v float64) (result float64) {
			u(v, result)
			return
		})
	}
}

// Reduced row echelon form of matrix (Gauss-Jordan elimination)
// rref
func (s *DenseMatrix) Reduced() Matrix {
	m := s.Copy()
	lead := 0
	rowCount := m.Rows()
	columnCount := m.Columns()
	for r := 0; r < rowCount; r++ {
		if lead >= columnCount {
			return m
		}
		i := r
		for m.At(i, lead) == 0 {
			i++
			if rowCount == i {
				i = r
				lead++
				if columnCount == lead {
					return m
				}
			}
		}

		v1 := m.RowsAt(i)
		v2 := m.RowsAt(r)

		for iterator := v1.Map(); iterator.HasNext(); {
			iterator.Map(func(c, _ int, value float64) float64 {
				m.Set(r, c, value)
				return value
			})
		}

		for iterator := v2.Map(); iterator.HasNext(); {
			iterator.Map(func(c, _ int, value float64) float64 {
				m.Set(i, c, value)
				return value
			})
		}

		f := 1 / m.At(r, lead)

		vector := m.RowsAt(r)
		for iterator := vector.Map(); iterator.HasNext(); {
			iterator.Map(func(c, _ int, value float64) float64 {
				value *= f
				m.Set(r, c, value)
				return value
			})
		}

		for i = 0; i < rowCount; i++ {
			if i != r {
				f = m.At(i, lead)

				vector := m.RowsAt(r)
				for iterator := vector.Map(); iterator.HasNext(); {
					iterator.Map(func(c, _ int, value float64) float64 {
						v := m.At(i, c)
						v -= value * f
						m.Set(i, c, v)
						return value
					})
				}
			}
		}
		lead++
	}

	return m
}

// ReduceToScalar perform's a reduction on the Matrix
func (s *DenseMatrix) ReduceToScalar() int {
	// https://people.eecs.berkeley.edu/~aydin/GraphBLAS_API_C.pdf
	// TODO need to reduce computes the result of performing a reduction
	// across each of the elements of an input matrix

	// monoid := binaryOp.NewMonoIDBool(true, binaryOp.LXOR)

	// monoid.Reduce()
	// for i := 0; i < s.Values(); i++ {
	// 	xor.Apply()
	// }
	return 0
}

// RawMatrix returns the raw matrix
func (s *DenseMatrix) RawMatrix() [][]float64 {
	return s.data
}

// Enumerate iterates through all non-zero elements, order is not guaranteed
func (s *DenseMatrix) Enumerate() Enumerate {
	return s.iterator()
}

func (s *DenseMatrix) iterator() *denseMatrixIterator {
	i := &denseMatrixIterator{
		matrix: s,
		size:   s.Values(),
		last:   0,
		c:      0,
		r:      0,
	}
	return i
}

type denseMatrixIterator struct {
	matrix *DenseMatrix
	size   int
	last   int
	c      int
	r      int
	cOld   int
	rOld   int
}

// HasNext checks the iterator has any more values
func (s *denseMatrixIterator) HasNext() bool {
	if s.last >= s.size {
		return false
	}
	return true
}

func (s *denseMatrixIterator) next() {
	if s.c == s.matrix.Columns() {
		s.c = 0
		s.r++
	}
	s.cOld = s.c
	s.c++
	s.last++
}

// Next moves the iterator and returns the row, column and value
func (s *denseMatrixIterator) Next() (int, int, float64) {
	s.next()
	return s.r, s.cOld, s.matrix.At(s.r, s.cOld)
}

// Map replace each element with the result of applying a function to its value
func (s *DenseMatrix) Map() Map {
	t := s.iterator()
	i := &denseMatrixMap{t}
	return i
}

type denseMatrixMap struct {
	*denseMatrixIterator
}

// HasNext checks the iterator has any more values
func (s *denseMatrixMap) HasNext() bool {
	return s.denseMatrixIterator.HasNext()
}

// Map move the iterator and uses a higher order function to changes the elements current value
func (s *denseMatrixMap) Map(f func(int, int, float64) float64) {
	s.next()
	s.matrix.Set(s.r, s.cOld, f(s.r, s.cOld, s.matrix.At(s.r, s.cOld)))
}
