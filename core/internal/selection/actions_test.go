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
		selection     types.Selection
		client        *mocks.Client
		firstElement  *mocks.Element
		secondElement *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		firstElement = &mocks.Element{}
		secondElement = &mocks.Element{}
		selection = &Selection{Client: client}
		selection = selection.All("#selector")
	})

	ItShouldEnsureAtLeastOneElement := func(matcher func() error) {
		Context("when zero elements are returned", func() {
			It("should return an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{}
				Expect(matcher()).To(MatchError("failed to select 'CSS: #selector': no elements found"))
			})
		})
	}

	Describe("#Click", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Click()
		})

		It("should click on all selected elements", func() {
			selection.Click()
			Expect(firstElement.ClickCall.Called).To(BeTrue())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})

		Context("when any click fails", func() {
			It("should return an error", func() {
				secondElement.ClickCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to click on 'CSS: #selector': some error"))
			})
		})

		Context("when all clicks succeed", func() {
			It("should return nil", func() {
				Expect(selection.Click()).To(BeNil())
			})
		})
	})

	// TODO: extend mock to test multiple calls
	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.DoubleClick()
		})

		It("should move the mouse to the middle of each selected element", func() {
			selection.DoubleClick()
			Expect(client.MoveToCall.Element).To(Equal(secondElement))
			Expect(client.MoveToCall.Point).To(BeNil())
		})

		Context("when moving over any element fails", func() {
			It("should retun an error", func() {
				client.MoveToCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to 'CSS: #selector': some error"))
			})
		})

		It("should double-click on each element", func() {
			selection.DoubleClick()
			Expect(client.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when the double-clicking any element fails", func() {
			It("should return an error", func() {
				client.DoubleClickCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to double-click on 'CSS: #selector': some error"))
			})
		})

		Context("when the double-clicking all elements succeeds", func() {
			It("should return nil", func() {
				Expect(selection.DoubleClick()).To(BeNil())
			})
		})
	})

	Describe("#Fill", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Fill("some text")
		})

		Context("when clearing any element fails", func() {
			It("should return an error", func() {
				secondElement.ClearCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to clear 'CSS: #selector': some error"))
			})
		})

		Context("when entering text into any element fails", func() {
			It("should return an error", func() {
				secondElement.ValueCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to enter text into 'CSS: #selector': some error"))
			})
		})

		Context("when the fill succeeds", func() {
			It("should clear each element", func() {
				selection.Fill("some text")
				Expect(firstElement.ClearCall.Called).To(BeTrue())
				Expect(secondElement.ClearCall.Called).To(BeTrue())
			})

			It("should fill each element with the provided text", func() {
				selection.Fill("some text")
				Expect(firstElement.ValueCall.Text).To(Equal("some text"))
				Expect(secondElement.ValueCall.Text).To(Equal("some text"))
			})

			It("should return nil", func() {
				Expect(selection.Fill("some text")).To(BeNil())
			})
		})
	})

	Describe("#Check", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Check()
		})

		It("should check the type of each checkbox", func() {
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			selection.Check()
			Expect(firstElement.GetAttributeCall.Attribute).To(Equal("type"))
			Expect(secondElement.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when any element fails to retrieve the 'type' attribute", func() {
			It("should return an error", func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.Err = errors.New("some error")
				Expect(selection.Check()).To(MatchError("failed to retrieve type of 'CSS: #selector': some error"))
			})
		})

		Context("when any element is not a checkbox", func() {
			It("should return an error", func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.ReturnValue = "banana"
				Expect(selection.Check()).To(MatchError("'CSS: #selector' does not refer to a checkbox"))
			})
		})

		Context("when all elements are checkboxes", func() {
			BeforeEach(func() {
				firstElement.GetAttributeCall.ReturnValue = "checkbox"
				secondElement.GetAttributeCall.ReturnValue = "checkbox"
			})

			Context("when the determining the selected status of any element fails", func() {
				It("should return an error", func() {
					secondElement.IsSelectedCall.Err = errors.New("some error")
					Expect(selection.Check()).To(MatchError("failed to retrieve state of 'CSS: #selector': some error"))
				})
			})

			Context("when clicking succeeds", func() {
				BeforeEach(func() {
					firstElement.IsSelectedCall.ReturnSelected = true
					secondElement.IsSelectedCall.ReturnSelected = false
					selection.Check()
				})

				It("should not click on the checked checkbox", func() {
					Expect(firstElement.ClickCall.Called).To(BeFalse())
				})

				It("should click on the unchecked checkboxes", func() {
					Expect(secondElement.ClickCall.Called).To(BeTrue())
				})
			})

			Context("when clicking on the checkbox fails", func() {
				BeforeEach(func() {
					secondElement.ClickCall.Err = errors.New("some error")
				})

				It("should return an error", func() {
					Expect(selection.Check()).To(MatchError("failed to click on 'CSS: #selector': some error"))
				})
			})
		})
	})

	Describe("#Uncheck", func() {
		It("should click on an checked checkbox", func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.IsSelectedCall.ReturnSelected = true
			selection.Uncheck()
			Expect(firstElement.ClickCall.Called).To(BeFalse())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})
	})

	Describe("#Select", func() {
		var (
			firstOptions  []*mocks.Element
			secondOptions []*mocks.Element
		)

		BeforeEach(func() {
			firstOptions = []*mocks.Element{&mocks.Element{}, &mocks.Element{}}
			secondOptions = []*mocks.Element{&mocks.Element{}, &mocks.Element{}}
			firstElement.GetElementsCall.ReturnElements = []types.Element{firstOptions[0], firstOptions[1]}
			secondElement.GetElementsCall.ReturnElements = []types.Element{secondOptions[0], secondOptions[1]}
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Select("some text")
		})

		It("should retrieve the options with matching text for each selected element", func() {
			selection.Select("some text")
			Expect(firstElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(firstElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space(text())="some text"]`))
			Expect(secondElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(secondElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space(text())="some text"]`))
		})

		Context("when we fail to retrieve any option", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError("failed to select specified option for some 'CSS: #selector': some error"))
			})
		})

		Context("when any of the elements has no options with matching text", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.ReturnElements = []types.Element{}
				Expect(selection.Select("some text")).To(MatchError(`no options with text "some text" found for some 'CSS: #selector'`))
			})
		})

		It("should click on all options with matching text", func() {
			selection.Select("some text")
			Expect(firstOptions[0].ClickCall.Called).To(BeTrue())
			Expect(firstOptions[1].ClickCall.Called).To(BeTrue())
			Expect(secondOptions[0].ClickCall.Called).To(BeTrue())
			Expect(secondOptions[1].ClickCall.Called).To(BeTrue())
		})

		Context("when the click fails for any of the options", func() {
			It("should return an error", func() {
				secondOptions[1].ClickCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError(`failed to click on option with text "some text" for some 'CSS: #selector': some error`))
			})
		})

		Context("when all clicks succeed", func() {
			It("should not return an error", func() {
				err := selection.Select("some text")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Submit()
		})

		It("should submit all selected elements", func() {
			selection.Submit()
			Expect(firstElement.SubmitCall.Called).To(BeTrue())
			Expect(secondElement.SubmitCall.Called).To(BeTrue())
		})

		Context("when any submit fails", func() {
			It("should return an error", func() {
				secondElement.SubmitCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to submit 'CSS: #selector': some error"))
			})
		})

		Context("when all submits succeed", func() {
			It("should return nil", func() {
				Expect(selection.Submit()).To(BeNil())
			})
		})
	})
})
