package main

import(
	"fmt"
	"time"
	"math"
	"log"

	"github.com/Awdmin/SummerProjects/Go3DTerminalViewer/mt"
)

var FACE_CHARS []byte = []byte{'@', '#', '$', '&', '÷'}
const ANGLE float64 = math.Pi/6
const CL_RATIO int = 2
const SCENE_DURATION int = 1
const FPS int = 12
const F_DURATION time.Duration = time.Duration(1000/FPS)
const PLANE_SIZE = 80

type Scene struct {
	Objects	[]Cube
}


type Cube struct {
	Side	int
	X		int
	Y		int
}

type Point struct {
	X		int
	Y		int
}

type FillLine struct {
	startX	int
	endX	int
}


func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func fillShape(pts []Point) map[int]FillLine {
	fillLines := make(map[int]FillLine)
	for i, _ := range pts {
		p1 := pts[i]
		p2 := pts[(i+1)%len(pts)]
		fmt.Println(i, (i+1)%len(pts))
		diffY := float64(p2.Y - p1.Y)
		diffX := float64(p2.X - p1.X)
		minY := math.Min(float64(p1.Y), float64(p2.Y))
		minX := math.Min(float64(p1.X), float64(p2.X))

		if diffY == 0 {
			idx := int(minY)
			line := FillLine{
						startX:		int(minX),
						endX:		int(minX + math.Abs(diffX)),
					}
			fillLines[idx] = line
		} else {

			k := -(diffX / diffY)
			if diffY < 0 {
				k = 1*k
				diffY *= -1
			} else {
				k *= -1
			}
			fmt.Print(i, (i+1)%len(pts))
			fmt.Printf("-> %f, %f, %f", diffY, diffX, k)
			for j := 0; j < int(diffY); j++ {
				idx := int(minY) + j
				v1 := float64(fillLines[idx].startX)
				v2 := float64(fillLines[idx].endX) // need to fix the 0 def value to get correct starts
				vt := (float64(j)) * k
				var minV int
				if v1 == 0|| v2 == 0 {
					minV = int(minX + vt) 
				} else {
					minV = int(math.Min(v2, math.Min(v1, minX + vt)))
				}

				maxV := int(math.Max(v2, math.Max(v1, minX + vt)))

				fmt.Print(idx)

				line := FillLine{
							startX:		minV,
							endX:		maxV,
						}

				fmt.Println(line)

				fillLines[idx] = line
								
				
			}
		}
	}

	return fillLines

}


func (c *Cube) DrawPixle(x int, y int, angle float64) bool {
	rtX := mt.Matrix{
		Values:		[][]float64{
						{1, 0, 0},
						{0, math.Cos(angle), math.Sin(angle)},
						{0, -math.Sin(angle), math.Cos(angle)},
					},
		}

	rtY := mt.Matrix{
		Values:		[][]float64{
						{math.Cos(angle), 0, -math.Sin(angle)},
						{0, 1, 0},
						{math.Sin(angle), 0, math.Cos(angle)},
					},
		}

	posVector := mt.Matrix{
		Values:		[][]float64{
						{float64(x)}, 
						{float64(y)}, 
						{0},
					},
		}
	
	m1, err := rtX.Multiply(rtY)
	if err != nil {
		log.Fatal(err)
	}

	newPos, err := m1.Multiply(posVector)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%f %f", newPos.Values[0], newPos.Values[1])

	newX := int(newPos.Values[0][0])
	newY := int(newPos.Values[1][0])
	//fmt.Print(newX, newY)

	if newX < c.X + c.Side && newX > c.X && newY < c.Y + c.Side/2 && newY > c.Y {
		return true
	}

	return false
}


func (s *Scene) Draw(size int, char int, angle float64, fillLines map[int]FillLine) {
	for i := 0; i < size/CL_RATIO; i++ {
		for j := 0; j < fillLines[i].startX; j++ {
			if j == 0 {
				fmt.Printf("%02d", i)
			} else {
				fmt.Print(" ")
			}
		}
		for j := fillLines[i].startX; j < fillLines[i].endX; j++ {
			fmt.Print(string(FACE_CHARS[char]))
		}
		for j := fillLines[i].endX; j < size; j++ {
			if j == 0 {
				fmt.Printf("%02d", i)
			} else {
				fmt.Print(" ")
			}		}

		fmt.Print("\n");
	}
}


//
// make a hash map that stores the important pixle values
// rotate a 3d cube (all 7/8 verticies) and write a function that can color inside a convex shape defined by those points (not the bes solution but progress)
//


func main() {
	clearScreen()

	c1 := Cube{
		Side:	15,
		X: 		10,
		Y: 		15,
	}

	var scene Scene
	scene.Objects = []Cube{c1}

	pts := []Point{
		Point{X: 3, Y: 3},
		Point{X: 20, Y: 5},
		Point{X: 20, Y: 15},
		Point{X: 3, Y: 12},
	}

	fillLines := fillShape(pts)


	for i := 0; i < SCENE_DURATION*FPS; i++ {
		scene.Draw(PLANE_SIZE, i % 5, float64(i%6)*(math.Pi/float64(6)), fillLines)
		time.Sleep(F_DURATION * time.Millisecond)
		clearScreen()
	}

	scene.Draw(PLANE_SIZE, 1, math.Pi/6, fillLines)

	fillLines = fillShape(pts)
}
