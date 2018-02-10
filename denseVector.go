// Copyright (c) 2018 Ross Merrigan
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package GraphBLAS

import (
	"log"
)

// DenseVector a vector
type DenseVector struct {
	l      int // length of the sparse vector
	values []float64
}

// NewDenseVector returns a GraphBLAS.DenseVector
func NewDenseVector(l int) *DenseVector {
	return &DenseVector{l: l, values: make([]float64, l)}
}

// NewDenseVectorFromArray returns a GraphBLAS.SparseVector
func NewDenseVectorFromArray(data []float64) *DenseVector {
	arr := make([]float64, len(data), len(data))
	arr = append(arr, data...)
	return &DenseVector{l: len(data), values: arr}
}

// AtVec returns the value of a vector element at i-th
func (s *DenseVector) AtVec(i int) float64 {
	if i < 0 || i >= s.Length() {
		log.Panicf("Length '%+v' is invalid", i)
	}

	return s.values[i]
}

// SetVec sets the value at i-th of the vector
func (s *DenseVector) SetVec(i int, value float64) {
	if i < 0 || i >= s.Length() {
		log.Panicf("Length '%+v' is invalid", i)
	}

	s.values[i] = value
}

// Length of the vector
func (s *DenseVector) Length() int {
	return s.l
}

// Columns the number of columns of the vector
func (s *DenseVector) Columns() int {
	return 1
}

// Rows the number of rows of the vector
func (s *DenseVector) Rows() int {
	return s.l
}

// Update does a At and Set on the vector element at r-th, c-th
func (s *DenseVector) Update(r, c int, f func(float64) float64) {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	v := s.AtVec(r)
	s.SetVec(r, f(v))
}

// At returns the value of a vector element at r-th, c-th
func (s *DenseVector) At(r, c int) (value float64) {
	s.Update(r, c, func(v float64) float64 {
		value = v
		return v
	})

	return
}

// Set sets the value at r-th, c-th of the vector
func (s *DenseVector) Set(r, c int, value float64) {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	s.SetVec(r, value)
}

// ColumnsAt return the columns at c-th
func (s *DenseVector) ColumnsAt(c int) Vector {
	if c < 0 || c >= s.Columns() {
		log.Panicf("Column '%+v' is invalid", c)
	}

	return s.copy()
}

// RowsAt return the rows at r-th
func (s *DenseVector) RowsAt(r int) Vector {
	if r < 0 || r >= s.Rows() {
		log.Panicf("Row '%+v' is invalid", r)
	}

	v := s.AtVec(1)
	rows := NewDenseVector(1)
	rows.SetVec(0, v)

	return rows
}

// Iterator iterates through all non-zero elements, order is not guaranteed
func (s *DenseVector) Iterator(i func(r, c int, v float64) bool) bool {
	for c := 0; c < s.Columns(); c++ {
		for r := 0; r < s.Rows(); r++ {
			v := s.At(r, c)
			if v != 0.0 {
				if i(r, c, v) == false {
					return false
				}
			}
		}
	}

	return true
}

func (s *DenseVector) copy() *DenseVector {
	vector := NewDenseVector(s.l)

	for i, v := range s.values {
		vector.SetVec(i, v)
	}

	return vector
}

// Copy copies the vector
func (s *DenseVector) Copy() Matrix {
	return s.CopyArithmetic(func(value float64) float64 {
		return value
	})
}

// CopyArithmetic copies the matrix and applies a arithmetic function through all non-zero elements, order is not guaranteed
func (s *DenseVector) CopyArithmetic(action func(float64) float64) Matrix {
	vector := NewDenseVector(s.l)

	for i, v := range s.values {
		if v != 0.0 {
			vector.SetVec(i, action(v))
		} else {
			vector.SetVec(i, v)
		}
	}

	return vector
}

// Scalar multiplication of a vector by alpha
func (s *DenseVector) Scalar(alpha float64) Matrix {
	return Scalar(s, alpha)
}

// Multiply multiplies a vector by another vector
func (s *DenseVector) Multiply(m Matrix) Matrix {
	matrix := newMatrix(m.Rows(), s.Columns(), nil)

	return multiplyVector(s, m, matrix)
}

// Add addition of a vector by another vector
func (s *DenseVector) Add(m Matrix) Matrix {
	return Add(s, m)
}

// Subtract subtracts one vector from another vector
func (s *DenseVector) Subtract(m Matrix) Matrix {
	return Subtract(s, m)
}

// Negative the negative of a metrix
func (s *DenseVector) Negative() Matrix {
	return Negative(s)
}

// Transpose swaps the rows and columns
func (s *DenseVector) Transpose() Matrix {
	matrix := newMatrix(s.Columns(), s.Rows(), nil)

	return Transpose(s, matrix)
}

// Equal the two vectors are equal
func (s *DenseVector) Equal(m Matrix) bool {
	return Equal(s, m)
}

// NotEqual the two vectors are not equal
func (s *DenseVector) NotEqual(m Matrix) bool {
	return NotEqual(s, m)
}
