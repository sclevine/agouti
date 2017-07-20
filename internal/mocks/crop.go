package mocks

import "image"

type Cropper struct {
	Image       image.Image
	Width       int
	Height      int
	Anchor      image.Point
	ReturnImage image.Image
	Err         error
}

func (c *Cropper) Crop(img image.Image, width, height int, anchor image.Point) (image.Image, error) {
	c.Image = img
	c.Width = width
	c.Height = height
	c.Anchor = anchor
	return c.ReturnImage, c.Err
}
