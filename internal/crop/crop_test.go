package crop

import (
	"image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crop", func() {
	It("crops the image", func() {
		r, err := Crop(getImage(), 512, 400, image.Point{})
		Expect(err).NotTo(HaveOccurred())
		Expect(r.Bounds().Dx()).To(Equal(512))
		Expect(r.Bounds().Dy()).To(Equal(400))
		Expect(r.Bounds().Min.X).To(Equal(0))
		Expect(r.Bounds().Min.Y).To(Equal(0))
	})

	Context("when a different anchor point is used", func() {
		It("crops the image", func() {
			r, err := Crop(getImage(), 512, 400, image.Point{X: 100, Y: 50})
			Expect(err).NotTo(HaveOccurred())
			Expect(r.Bounds().Dx()).To(Equal(512))
			Expect(r.Bounds().Dy()).To(Equal(400))
			Expect(r.Bounds().Min.X).To(Equal(100))
			Expect(r.Bounds().Min.Y).To(Equal(50))
		})
	})
})

func getImage() image.Image {
	return image.NewGray(image.Rect(0, 0, 1600, 1437))
}
