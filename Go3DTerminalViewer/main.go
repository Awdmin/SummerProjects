package main

import(
	"fmt"
	"time"
	"math"
	//"log"
)

var FACE_CHARS []byte = []byte{'@', '#', '$', '&', '÷'}
const ANGLE float64 = math.Pi/2.5
const CL_RATIO int = 2
const SCENE_DURATION int = 0
const FPS int = 12
const F_DURATION time.Duration = time.Duration(1000/FPS)
const PLANE_SIZE = 80

type Scene struct {
	Objects	[]Cube
}

type Square struct {
	Side	int
	X		int
	Y		int
	Angle	float64
	Top		bool
}

type Cube struct {
	Squares	[]Square
	Side	int
	X		int
	Y		int
	Angle	float64
}


func (c *Cube) Init() {
	

	s1 := Square{
		Side:	int(math.Abs(float64(c.Side) * math.Cos(c.Angle))),
		X:		c.X,
		Y:		c.Y, //+ int(math.Abs(float64(c.Side) * math.Cos(math.Pi - 2*math.Pi/3.0 - c.Angle)) * math.Tan(math.Pi - 2*math.Pi/3 - c.Angle)), 
		Angle:	-c.Angle,
		Top:	false,
	}
	//fmt.Println(s1.Side, s1.Angle)

	s2 := Square{
		Side: 	int(math.Abs(float64(c.Side) * math.Cos(math.Pi - 2*math.Pi/3.0 - c.Angle))),
		X:		c.X+s1.Side,
		Y:		c.Y-int(math.Abs(float64(c.Side) * math.Cos(math.Pi - 2*math.Pi/3.0 - c.Angle))),
		Angle:	math.Pi - 2*math.Pi/3 - c.Angle,
		Top:	false,
	}
	//fmt.Print(s2.Side, s2.Angle)

	s3 := Square{
		Side:	int(math.Abs(float64(c.Side) * math.Cos(c.Angle))),
		X:		c.X,
		Y:		c.Y-c.Side,
		Angle:	-c.Angle,
		Top:	true,
	}

	s4 := Square{
		Side: 	int(math.Abs(float64(c.Side) * math.Cos(math.Pi - 2*math.Pi/3.0 - c.Angle))),
		X:		c.X+s1.Side,
		Y:		c.Y-c.Side-int(math.Abs(float64(c.Side) * math.Cos(math.Pi - 2*math.Pi/3.0 - c.Angle))),
		Angle:	math.Pi - 2*math.Pi/3 - c.Angle,
		Top:	true,
	}

	c.Squares = []Square{s1, s2, s3, s4}

}


func (s *Square) VerticalOffsetForX(x int) int {
	var x1 int
		x1 = x - s.X
		//x1 = s.Side - (x - s.X)
	if x1 < 0 {
		return -1
	}
	offset := float64(s.Side) * math.Tan(s.Angle) - float64(x1) * math.Tan(s.Angle)
	//fmt.Print(offset)

	return int(offset)
}


func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

/*
func translateCoordinates(x int, y int) (int, int) {
	return PLANE_SIZE/2 + x, (PLANE_SIZE/CL_RATIO)/2 + y/CL_RATIO
}
*/

func (c *Cube) DrawPixle(x int, y int) bool {
	for i, s := range c.Squares {
		offsetBottom := s.VerticalOffsetForX(x)
		offsetTop := offsetBottom
		ch := 4
		
		if s.Top {
			offsetTop = (int(float64(s.Side) * math.Tan(s.Angle)) - offsetBottom)
			ch = 3
		}

		if i == 1 {
			ch = 1
		}


		topBorder, bottomBorder := s.Y/2 + offsetTop, s.Y/2 + c.Side/2 + offsetBottom
		leftBorder, rightBorder := s.X, s.X + s.Side

		if x > leftBorder && x <= rightBorder && y > topBorder && y <= bottomBorder {
			fmt.Print(string(FACE_CHARS[ch]))
			return true
		}
	}

	return false
}


func (s *Scene) Draw(size int, char int) {
	for i := 0; i < size/CL_RATIO; i++ {
		for j := 0; j < size; j++ {
			if j == 0 {
				fmt.Printf("%02d", i)
			}
			pixle := false
			for _, e := range s.Objects {
				if e.DrawPixle(j, i) {
					pixle = true
				}
			}
			if !pixle {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n");
	}
}


func main() {
	clearScreen()

	c1 := Cube{
		Side:	20,
		X: 10,
		Y: 45,
		Angle:	math.Pi/6,
	}
	c1.Init()

	var scene Scene
	scene.Objects = []Cube{c1}

	for i := 0; i < SCENE_DURATION*FPS; i++ {
		scene.Draw(PLANE_SIZE, i % 5)
		time.Sleep(F_DURATION * time.Millisecond)
		clearScreen()
	}

	scene.Draw(PLANE_SIZE, 1)
}

