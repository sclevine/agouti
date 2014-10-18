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
		selection types.Selection
		client    *mocks.Client
		element   *mocks.Element
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		element = &mocks.Element{}
		selection = &Selection{Client: client}
		selection = selection.Find("#selector")
	})

	ItShouldEnsureASingleElement := func(matcher func() error) {
		Context("ensures a single element is returned", func() {
			It("returns an error with the number of elements", func() {
				client.GetElementsCall.ReturnElements = []types.Element{element, element}
				Expect(matcher()).To(MatchError("failed to retrieve element with 'CSS: #selector': multiple elements (2) were selected"))
			})
		})
	}

	Describe("#Click", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Click()
		})

		It("clicks on an element", func() {
			selection.Click()
			Expect(element.ClickCall.Called).To(BeTrue())
		})

		Context("if the click fails", func() {
			BeforeEach(func() {
				element.ClickCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Click()).To(MatchError("failed to click on 'CSS: #selector': some error"))
			})
		})

		Context("if the click succeeds", func() {
			It("returns nil", func() {
				Expect(selection.Click()).To(BeNil())
			})
		})
	})

	Describe("#DoubleClick", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.DoubleClick()
		})

		It("moves the mouse to the middle of the selected element", func() {
			selection.DoubleClick()
			Expect(client.MoveToCall.Element).To(Equal(element))
			Expect(client.MoveToCall.Point).To(BeNil())
		})

		Context("when moving over the element fails", func() {
			BeforeEach(func() {
				client.MoveToCall.Err = errors.New("some error")
			})

			It("retuns an error", func() {
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to 'CSS: #selector': some error"))
			})
		})

		It("double-clicks on an element", func() {
			selection.DoubleClick()
			Expect(client.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when the double-clicking the element fails", func() {
			BeforeEach(func() {
				client.DoubleClickCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.DoubleClick()).To(MatchError("failed to double-click on 'CSS: #selector': some error"))
			})
		})

		Context("when the double-clicking the element succeeds", func() {
			It("returns nil", func() {
				Expect(selection.DoubleClick()).To(BeNil())
			})
		})
	})

	Describe("#Fill", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Fill("some text")
		})

		Context("if clearing the element fails", func() {
			BeforeEach(func() {
				element.ClearCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Fill("some text")).To(MatchError("failed to clear 'CSS: #selector': some error"))
			})
		})

		Context("if entering text into the element fails", func() {
			BeforeEach(func() {
				element.ValueCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Fill("some text")).To(MatchError("failed to enter text into 'CSS: #selector': some error"))
			})
		})

		Context("if the fill succeeds", func() {
			It("clears the element", func() {
				selection.Fill("some text")
				Expect(element.ClearCall.Called).To(BeTrue())
			})

			It("fills the element with the provided text", func() {
				selection.Fill("some text")
				Expect(element.ValueCall.Text).To(Equal("some text"))
			})

			It("returns nil", func() {
				Expect(selection.Fill("some text")).To(BeNil())
			})
		})
	})

	Describe("#Check", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Check()
		})

		It("checks the type of the checkbox", func() {
			selection.Check()
			Expect(element.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when the the client fails to retrieve the 'type' attribute", func() {
			BeforeEach(func() {
				element.GetAttributeCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Check()).To(MatchError("failed to retrieve type of 'CSS: #selector': some error"))
			})
		})

		Context("when the selection is not a checkbox", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "banana"
			})

			It("returns an error", func() {
				Expect(selection.Check()).To(MatchError("'CSS: #selector' does not refer to a checkbox"))
			})
		})

		Context("when the selection is a checkbox", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "checkbox"
			})

			Context("when the determining the selected status of the element fails", func() {
				BeforeEach(func() {
					element.IsSelectedCall.Err = errors.New("some error")
				})

				It("returns an error", func() {
					Expect(selection.Check()).To(MatchError("failed to retrieve state of 'CSS: #selector': some error"))
				})
			})

			Context("when the box is already checked", func() {
				BeforeEach(func() {
					element.IsSelectedCall.ReturnSelected = true
				})

				It("does not click on the checkbox", func() {
					selection.Check()
					Expect(element.ClickCall.Called).To(BeFalse())
				})
			})

			Context("when the box is not checked", func() {
				BeforeEach(func() {
					element.IsSelectedCall.ReturnSelected = false
				})

				It("clicks on the checkbox", func() {
					selection.Check()
					Expect(element.ClickCall.Called).To(BeTrue())
				})

				Context("when clicking on the checkbox fails", func() {
					BeforeEach(func() {
						element.ClickCall.Err = errors.New("some error")
					})

					It("returns an error", func() {
						Expect(selection.Check()).To(MatchError("failed to click on 'CSS: #selector': some error"))
					})
				})
			})
		})
	})

	Describe("#Uncheck", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
			element.GetAttributeCall.ReturnValue = "checkbox"
			element.IsSelectedCall.ReturnSelected = true
		})

		It("clicks on an checked checkbox", func() {
			selection.Uncheck()
			Expect(element.ClickCall.Called).To(BeTrue())
		})
	})

	Describe("#Select", func() {
		var (
			optionOne   *mocks.Element
			optionTwo   *mocks.Element
			optionThree *mocks.Element
		)

		BeforeEach(func() {
			optionOne = &mocks.Element{}
			optionTwo = &mocks.Element{}
			optionThree = &mocks.Element{}
			client.GetElementsCall.ReturnElements = []types.Element{optionOne, optionTwo, optionThree}
		})

		It("request child option elements from the client", func() {
			selection.Select("some text")
			Expect(client.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector option"}))
		})

		Context("when the client fails to retrieve any elements", func() {
			BeforeEach(func() {
				client.GetElementsCall.Err = errors.New("some error")
			})

			It("returns error from the client", func() {
				Expect(selection.Select("some text")).To(MatchError("failed to retrieve options for 'CSS: #selector': some error"))
			})
		})

		Context("when the client fails to retrieve text for an element", func() {
			BeforeEach(func() {
				optionOne.GetTextCall.Err = errors.New("some error")
			})

			It("returns error from the client", func() {
				Expect(selection.Select("some text")).To(MatchError("failed to retrieve option text for 'CSS: #selector': some error"))
			})
		})

		Context("when at least one element has matching text", func() {
			BeforeEach(func() {
				optionOne.GetTextCall.ReturnText = "some other text"
				optionTwo.GetTextCall.ReturnText = "some text"
				optionThree.GetTextCall.ReturnText = "some text"
			})

			It("clicks on the first matching element", func() {
				selection.Select("some text")
				Expect(optionOne.ClickCall.Called).To(BeFalse())
				Expect(optionTwo.ClickCall.Called).To(BeTrue())
				Expect(optionThree.ClickCall.Called).To(BeFalse())
			})

			It("does not return an error", func() {
				err := selection.Select("some text")
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the click fails", func() {
				BeforeEach(func() {
					optionTwo.ClickCall.Err = errors.New("some error")
				})

				It("return an error indicating that it failed to click on the element", func() {
					err := selection.Select("some text")
					Expect(err).To(MatchError(`failed to click on option with text "some text" for 'CSS: #selector': some error`))
				})
			})
		})

		Context("when no elements have matching text", func() {
			BeforeEach(func() {
				optionOne.GetTextCall.ReturnText = "some other text"
				optionTwo.GetTextCall.ReturnText = "some different text"
				optionThree.GetTextCall.ReturnText = "some other different text"
			})

			It("returns an error indicating that no options could be selected", func() {
				err := selection.Select("some text")
				Expect(err).To(MatchError(`no options with text "some text" found for 'CSS: #selector'`))
			})
		})
	})

	Describe("#Submit", func() {
		BeforeEach(func() {
			client.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Submit()
		})

		Context("when submitting fails", func() {
			BeforeEach(func() {
				element.SubmitCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				Expect(selection.Submit()).To(MatchError("failed to submit 'CSS: #selector': some error"))
			})
		})

		Context("when submitting succeeds", func() {
			It("submits the element", func() {
				selection.Submit()
				Expect(element.SubmitCall.Called).To(BeTrue())
			})

			It("returns nil", func() {
				Expect(selection.Submit()).To(BeNil())
			})
		})
	})
})
