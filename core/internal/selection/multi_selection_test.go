package selection_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
)

var _ = Describe("MultiSelection", func() {
	var (
		client         *mocks.Client
		multiSelection *MultiSelection
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		selection := &Selection{Client: client}
		multiSelection = &MultiSelection{selection}

	})

	Describe("#At", func() {
		Context("when called on a multi-selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(multiSelection.At(1).String()).To(Equal(""))
			})
		})

		Context("when called on a multi-selection with selectors", func() {
			It("should select an index of the current selection", func() {
				Expect(multiSelection.All("#selector").At(1).String()).To(Equal("CSS: #selector [1]"))
			})
		})
	})

	Describe("#Single", func() {
		Context("when called on a multi-selection with no selectors", func() {
			It("should return an empty selection", func() {
				Expect(multiSelection.Single().String()).To(Equal(""))
			})
		})

		Context("when called on a multi-selection with selectors", func() {
			It("should select a single element of the current selection", func() {
				Expect(multiSelection.All("#selector").Single().String()).To(Equal("CSS: #selector [single]"))
			})
		})
	})
})
