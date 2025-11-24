package mt

import(
	//"fmt"
	//"math"
	"errors"
)


type Matrix struct {
	Values	[][]float64
}

var DIMENTION_MISMATCH_ERR = errors.New("Matrix dimetions are not correct for multiplication")

func (m1 *Matrix) Multiply(m2 Matrix) (Matrix, error) {
	var result Matrix

	if len(m1.Values[0]) != len(m2.Values) {
		return result, DIMENTION_MISMATCH_ERR
	}
	
	rows := len(m1.Values)
	cols := len(m2.Values[0])
	commonDim := len(m1.Values[0])
	result.Values = make([][]float64, rows)

	for i := 0; i < rows; i++ {
		result.Values[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			for k := 0; k < commonDim; k++ {
				result.Values[i][j] += m1.Values[i][k] * m2.Values[k][j] 
			}
		}
	}

	return result, nil
}

