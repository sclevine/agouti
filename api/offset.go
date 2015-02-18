package api

type Offset interface {
	x() (x int, present bool)
	y() (y int, present bool)
}

type XYOffset struct {
	X int
	Y int
}

func (o XYOffset) x() (x int, present bool) {
	return o.X, true
}

func (o XYOffset) y() (y int, present bool) {
	return o.Y, true
}

type XOffset int

func (o XOffset) x() (x int, present bool) {
	return int(o), true
}

func (XOffset) y() (y int, present bool) {
	return 0, false
}

type YOffset int

func (YOffset) x() (x int, present bool) {
	return 0, false
}

func (o YOffset) y() (y int, present bool) {
	return int(o), true
}
