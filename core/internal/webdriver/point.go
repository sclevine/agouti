package webdriver

type Point interface {
	x() (x int, present bool)
	y() (y int, present bool)
}

type XYPoint struct {
	X int
	Y int
}

func (p XYPoint) x() (x int, present bool) {
	return p.X, true
}

func (p XYPoint) y() (y int, present bool) {
	return p.Y, true
}

type XPoint int

func (p XPoint) x() (x int, present bool) {
	return int(p), true
}

func (p XPoint) y() (y int, present bool) {
	return 0, false
}

type YPoint int

func (p YPoint) x() (x int, present bool) {
	return 0, false
}

func (p YPoint) y() (y int, present bool) {
	return int(p), true
}
