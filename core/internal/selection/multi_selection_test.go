package selection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
	"errors"
)

var _ = Describe("Selection", func() {
	var (
		selection types.Selection
		multiSelection types.MultiSelection
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		element = &mocks.Element{}
		selection = &Selection{Driver: driver}
		multiSelection = selection.Find("#selector").All()
	})

	Describe("#String", func() {
		It("returns selection#String with '- All' appended", func() {
			Expect(multiSelection.String()).To(Equal("CSS: #selector - All"))
		})
	})

	Describe("#Visible", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		Context("when we fail to retrieve the list of elements", func() {
			It("returns an error", func() {
				driver.GetElementsCall.Err = errors.New("some error")
				_, err := multiSelection.Visible()
				Expect(err).To(MatchError("failed to retrieve elements with 'CSS: #selector - All': some error"))
			})
		})

//		Context("if the the driver fails to retrieve the element's visible status", func() {
//			It("returns an error", func() {
//				element.IsDisplayedCall.Err = errors.New("some error")
//				_, err := selection.Visible()
//				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is visible: some error"))
//			})
//		})
//
//		Context("if the driver succeeds in retrieving the element's visible status", func() {
//			It("returns the visible status when visible", func() {
//				element.IsDisplayedCall.ReturnDisplayed = true
//				value, _ := selection.Visible()
//				Expect(value).To(BeTrue())
//			})
//
//			It("returns the visible status when not visible", func() {
//				element.IsDisplayedCall.ReturnDisplayed = false
//				value, _ := selection.Visible()
//				Expect(value).To(BeFalse())
//			})
//
//			It("does not return an error", func() {
//				_, err := selection.Visible()
//				Expect(err).NotTo(HaveOccurred())
//			})
//		})
	})
})
