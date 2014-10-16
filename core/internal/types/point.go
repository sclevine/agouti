package types

type Point interface {
	X() (x int, present bool)
	Y() (y int, present bool)
}

type XYPoint struct {
	XPos int
	YPos int
}

func (p XYPoint) X() (x int, present bool) {
	return p.XPos, true
}

func (p XYPoint) Y() (y int, present bool) {
	return p.YPos, true
}

type XPoint int

func (p XPoint) X() (x int, present bool) {
	return int(p), true
}

func (p XPoint) Y() (y int, present bool) {
	return 0, false
}

type YPoint int

func (p YPoint) X() (x int, present bool) {
	return 0, false
}

func (p YPoint) Y() (y int, present bool) {
	return int(p), true
}
