package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Actions", func() {
	var (
		selection     *Selection
		client        *mocks.Client
		firstElement  *mocks.Element
		secondElement *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		firstElement = &mocks.Element{}
		secondElement = &mocks.Element{}
		emptySelection := &Selection{Client: client}
		selection = emptySelection.AppendCSS("#selector")
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

		It("should successfully click on all selected elements", func() {
			Expect(selection.Click()).To(Succeed())
			Expect(firstElement.ClickCall.Called).To(BeTrue())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})

		Context("when any click fails", func() {
			It("should return an error", func() {
				secondElement.ClickCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to click on 'CSS: #selector': some error"))
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

		It("should successfully move the mouse to the middle of each selected element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(client.MoveToCall.Element).To(Equal(secondElement))
			Expect(client.MoveToCall.Point).To(BeNil())
		})

		Context("when moving over any element fails", func() {
			It("should retun an error", func() {
				client.MoveToCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to 'CSS: #selector': some error"))
			})
		})

		It("should successfully double-click on each element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(client.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when the double-clicking any element fails", func() {
			It("should return an error", func() {
				client.DoubleClickCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to double-click on 'CSS: #selector': some error"))
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

		It("should successfully clear each element", func() {
			Expect(selection.Fill("some text")).To(Succeed())
			Expect(firstElement.ClearCall.Called).To(BeTrue())
			Expect(secondElement.ClearCall.Called).To(BeTrue())
		})

		Context("when clearing any element fails", func() {
			It("should return an error", func() {
				secondElement.ClearCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to clear 'CSS: #selector': some error"))
			})
		})

		It("should successfully fill each element with the provided text", func() {
			Expect(selection.Fill("some text")).To(Succeed())
			Expect(firstElement.ValueCall.Text).To(Equal("some text"))
			Expect(secondElement.ValueCall.Text).To(Equal("some text"))
		})

		Context("when entering text into any element fails", func() {
			It("should return an error", func() {
				secondElement.ValueCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to enter text into 'CSS: #selector': some error"))
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

		It("should successfully check the type of each checkbox", func() {
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			Expect(selection.Check()).To(Succeed())
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

			It("should not click on the checked checkbox successfully", func() {
				firstElement.IsSelectedCall.ReturnSelected = true
				Expect(selection.Check()).To(Succeed())
				Expect(firstElement.ClickCall.Called).To(BeFalse())
			})

			It("should click on the unchecked checkboxes successfully", func() {
				secondElement.IsSelectedCall.ReturnSelected = false
				Expect(selection.Check()).To(Succeed())
				Expect(secondElement.ClickCall.Called).To(BeTrue())
			})

			Context("when the determining the selected status of any element fails", func() {
				It("should return an error", func() {
					secondElement.IsSelectedCall.Err = errors.New("some error")
					Expect(selection.Check()).To(MatchError("failed to retrieve state of 'CSS: #selector': some error"))
				})
			})

			Context("when clicking on the checkbox fails", func() {
				It("should return an error", func() {
					secondElement.ClickCall.Err = errors.New("some error")
					Expect(selection.Check()).To(MatchError("failed to click on 'CSS: #selector': some error"))
				})
			})
		})
	})

	Describe("#Uncheck", func() {
		It("should successfully click on a checked checkbox", func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.IsSelectedCall.ReturnSelected = true
			Expect(selection.Uncheck()).To(Succeed())
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

		It("should successfully retrieve the options with matching text for each selected element", func() {
			Expect(selection.Select("some text")).To(Succeed())
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

		It("should successfully click on all options with matching text", func() {
			Expect(selection.Select("some text")).To(Succeed())
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
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{firstElement, secondElement}
		})

		ItShouldEnsureAtLeastOneElement(func() error {
			return selection.Submit()
		})

		It("should successfully submit all selected elements", func() {
			Expect(selection.Submit()).To(Succeed())
			Expect(firstElement.SubmitCall.Called).To(BeTrue())
			Expect(secondElement.SubmitCall.Called).To(BeTrue())
		})

		Context("when any submit fails", func() {
			It("should return an error", func() {
				secondElement.SubmitCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to submit 'CSS: #selector': some error"))
			})
		})
	})
})
