package _test

import(
	"testing"
	"fmt"

	"github.com/Awdmin/SummerProjects/Go3DTerminalViewer/mt"
)

func TestMatrixMultiplication2x2(t *testing.T) {
	m1 := mt.Matrix{
		Values:		[][]float64{{2, 1},{3, 4}},
	}

	m2 := mt.Matrix{
		Values:		[][]float64{{5, 0},{2, 4}},
	}

	res, err := m1.Multiply(m2)

	if err != nil {
		t.Fatal("Incorect dimentions of the input!")
	}

	expected := [][]float64{{12, 4},{23, 16}}

	if len(res.Values) != len(expected) || len(res.Values[0]) != len(expected[0]) {
		t.Fatal("Incorect dimentions of the result!")
	} else {
		for i, _ := range expected {
			for j, _ := range expected[0] {
				if res.Values[i][j] != expected[i][j] {
					t.Fatal(fmt.Sprintf("Incorrect value inside matrix, i: %d j: %d x: %f x1: %f", i, j, res.Values[i][j], expected[i][j]))
				}
			}
		}
	}
}


func TestMatrixMultiplication3x3(t *testing.T) {
	m1 := mt.Matrix{
		Values:		[][]float64{{2, 1, 0},{1, 3, 4},{1, 1, 1}},
	}

	m2 := mt.Matrix{
		Values:		[][]float64{{5, 0, 1},{2, 2, 4},{1, 3, 2}},
	}

	res, err := m1.Multiply(m2)

	if err != nil {
		t.Fatal("Incorect dimentions of the input!")
	}

	expected := [][]float64{{12, 2, 6},{15, 18, 21},{8, 5, 7}}

	if len(res.Values) != len(expected) || len(res.Values[0]) != len(expected[0]) {
		t.Fatal("Incorect dimentions of the result!")
	} else {
		for i, _ := range expected {
			for j, _ := range expected[0] {
				if res.Values[i][j] != expected[i][j] {
					t.Fatal("Incorrect value inside matrix")
				}
			}
		}
	}
}


