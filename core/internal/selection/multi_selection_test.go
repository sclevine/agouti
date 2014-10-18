package selection_test

import (
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Selection", func() {
	var (
		selection      types.Selection
		multiSelection types.MultiSelection
		driver         *mocks.Driver
		firstElement   *mocks.Element
		secondElement  *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		firstElement = &mocks.Element{}
		secondElement = &mocks.Element{}
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
			driver.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		Context("when we fail to retrieve the list of elements", func() {
			It("returns an error", func() {
				driver.GetElementsCall.Err = errors.New("some error")
				_, err := multiSelection.Visible()
				Expect(err).To(MatchError("failed to retrieve elements with 'CSS: #selector - All': some error"))
			})
		})

		Context("when no elements are returned", func() {
			It("returns an error", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{}
				_, err := multiSelection.Visible()
				Expect(err).To(MatchError("no elements found for 'CSS: #selector - All'"))
			})
		})

		Context("when the driver fails to retrieve any element's visible status", func() {
			It("returns an error", func() {
				firstElement.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.Err = errors.New("some error")
				_, err := multiSelection.Visible()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector - All' is visible: some error"))
			})
		})

		Context("when the driver succeeds in retrieving all elements' visible status", func() {
			It("returns true when all elements are visible", func() {
				firstElement.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.ReturnDisplayed = true
				value, _ := multiSelection.Visible()
				Expect(value).To(BeTrue())
			})

			It("returns false when not all elements are visible", func() {
				firstElement.IsDisplayedCall.ReturnDisplayed = true
				secondElement.IsDisplayedCall.ReturnDisplayed = false
				value, _ := multiSelection.Visible()
				Expect(value).To(BeFalse())
			})

			It("returns false when no elements are visible", func() {
				firstElement.IsDisplayedCall.ReturnDisplayed = false
				secondElement.IsDisplayedCall.ReturnDisplayed = false
				value, _ := multiSelection.Visible()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := multiSelection.Visible()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
