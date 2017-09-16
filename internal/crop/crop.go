package crop

import (
	"image"
	"image/draw"
)

// An interface that is
// image.Image + SubImage method.
type subImageSupported interface {
	SubImage(r image.Rectangle) image.Image
}

// Cropper is the interface used to crop images
type Cropper interface {
	Crop(img image.Image, width, height int, anchor image.Point) (image.Image, error)
}

// CropperFunc exposes a Crop function that calls itself.
// It implements Cropper.
type CropperFunc func(img image.Image, width, height int, anchor image.Point) (image.Image, error)

// Crop calls the CropperFunc
func (c CropperFunc) Crop(img image.Image, width, height int, anchor image.Point) (image.Image, error) {
	return c(img, width, height, anchor)
}

// Crop retrieves an image that is a
// cropped copy of the original img.
func Crop(img image.Image, width, height int, anchor image.Point) (image.Image, error) {
	maxBounds := maxBounds(anchor, img.Bounds())
	size := computeSize(maxBounds, image.Point{width, height})
	cr := computedCropArea(anchor, img.Bounds(), size)
	cr = img.Bounds().Intersect(cr)

	if dImg, ok := img.(subImageSupported); ok {
		return dImg.SubImage(cr), nil
	}
	return cropWithCopy(img, cr)
}

func cropWithCopy(img image.Image, cr image.Rectangle) (image.Image, error) {
	result := image.NewRGBA(cr)
	draw.Draw(result, cr, img, cr.Min, draw.Src)
	return result, nil
}

func maxBounds(anchor image.Point, bounds image.Rectangle) image.Rectangle {
	return image.Rect(anchor.X, anchor.Y, bounds.Max.X, bounds.Max.Y)
}

// computeSize retrieve the effective size of the cropped image.
func computeSize(bounds image.Rectangle, ratio image.Point) image.Point {
	return image.Point{ratio.X, ratio.Y}
}

// computedCropArea retrieve the theorical crop area.
func computedCropArea(anchor image.Point, bounds image.Rectangle, size image.Point) (r image.Rectangle) {
	min := bounds.Min
	rMin := image.Point{min.X + anchor.X, min.Y + anchor.Y}
	return image.Rect(rMin.X, rMin.Y, rMin.X+size.X, rMin.Y+size.Y)
}
