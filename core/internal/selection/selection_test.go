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
		driver    *mocks.Driver
		element   *mocks.Element
	)

	BeforeEach(func() {
		driver = &mocks.Driver{}
		element = &mocks.Element{}
		selection = &Selection{Driver: driver}
		selection = selection.Find("#selector")
	})

	ItShouldEnsureASingleElement := func(matcher func() error) {
		Context("ensures a single element is returned", func() {
			It("returns an error with the number of elements", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{element, element}
				Expect(matcher()).To(MatchError("failed to retrieve element with 'CSS: #selector': mutiple elements (2) were selected"))
			})
		})
	}

	Describe("most methods: retrieving elements", func() {
		var (
			parentOne *mocks.Element
			parentTwo *mocks.Element
			count     int
		)

		BeforeEach(func() {
			selection = selection.FindXPath("children")
			parentOne = &mocks.Element{}
			parentTwo = &mocks.Element{}
			parentOne.GetElementsCall.ReturnElements = []types.Element{&mocks.Element{}, &mocks.Element{}}
			parentTwo.GetElementsCall.ReturnElements = []types.Element{&mocks.Element{}, &mocks.Element{}}
			driver.GetElementsCall.ReturnElements = []types.Element{parentOne, parentTwo}
			count, _ = selection.Count()
		})

		Context("when successful", func() {
			It("retrieves the parent elements using the driver", func() {
				Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{"css selector", "#selector"}))
			})

			It("retrieves the child elements of the parent selector", func() {
				Expect(parentOne.GetElementsCall.Selector).To(Equal(types.Selector{"xpath", "children"}))
				Expect(parentTwo.GetElementsCall.Selector).To(Equal(types.Selector{"xpath", "children"}))
			})

			It("returns all child elements of the terminal selector", func() {
				Expect(count).To(Equal(4))
			})
		})

		Context("when there is no selection", func() {
			BeforeEach(func() {
				selection = &Selection{Driver: driver}
			})

			It("returns an error", func() {
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to retrieve elements for '': empty selection"))
			})
		})

		Context("when retrieving the parent elements fails", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("returns the error", func() {
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to retrieve elements for 'CSS: #selector | XPath: children': some error"))
			})
		})

		Context("when retrieving any of the child elements fails", func() {
			BeforeEach(func() {
				parentTwo.GetElementsCall.Err = errors.New("some error")
			})

			It("returns the error", func() {
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to retrieve elements for 'CSS: #selector | XPath: children': some error"))
			})
		})
	})

	Describe("most methods: retrieving a single element", func() {
		It("requests an element from the driver using the element's selector", func() {
			selection.Click()
			Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{"css selector", "#selector"}))
		})

		Context("when the driver fails to retrieve any elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("returns error from the driver", func() {
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': some error"))
			})
		})

		Context("when the driver retrieves more than one element", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []types.Element{element, element}
			})

			It("returns an error with the number of elements", func() {
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': mutiple elements (2) were selected"))
			})
		})

		Context("when the driver retrieves zero elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.ReturnElements = []types.Element{}
			})

			It("fails with an error indicating there were no elements", func() {
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': no element found"))
			})
		})
	})

	Describe("#Find", func() {
		Context("when there is no selection", func() {
			It("adds a new css selector to the selection", func() {
				selection := &Selection{Driver: driver}
				Expect(selection.Find("#selector").String()).To(Equal("CSS: #selector"))
			})
		})

		Context("when the selection ends with an xpath selector", func() {
			It("adds a new css selector to the selection", func() {
				xpath := selection.FindXPath("//subselector")
				Expect(xpath.Find("#subselector").String()).To(Equal("CSS: #selector | XPath: //subselector | CSS: #subselector"))
			})
		})

		Context("when the selection ends with a CSS selector", func() {
			It("modifies the terminal css selector to include the new selector", func() {
				Expect(selection.Find("#subselector").String()).To(Equal("CSS: #selector #subselector"))
			})
		})
	})

	Describe("#FindXPath", func() {
		It("adds a new XPath selector to the selection", func() {
			Expect(selection.FindXPath("//subselector").String()).To(Equal("CSS: #selector | XPath: //subselector"))
		})
	})

	Describe("#FindByLabel", func() {
		It("adds an XPath selector for finding by label", func() {
			Expect(selection.FindByLabel("label name").String()).To(Equal(`CSS: #selector | XPath: //input[@id=(//label[text()="label name"]/@for)] | //label[text()="label name"]/input`))
		})
	})

	Describe("#String", func() {
		It("returns the separated selectors", func() {
			Expect(selection.FindXPath("//subselector").String()).To(Equal("CSS: #selector | XPath: //subselector"))
		})
	})

	Describe("#Count", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element, element}
		})

		It("requests elements from the driver using the provided selector", func() {
			selection.Count()
			Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{"css selector", "#selector"}))
		})

		Context("when the driver succeeds in retrieving the elements", func() {
			It("returns the text", func() {
				count, _ := selection.Count()
				Expect(count).To(Equal(2))
			})

			It("does not return an error", func() {
				_, err := selection.Count()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the the driver fails to retrieve the elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to retrieve elements for 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Click", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
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
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.DoubleClick()
		})

		It("moves the mouse to the middle of the selected element", func() {
			selection.DoubleClick()
			Expect(driver.MoveToCall.Element).To(Equal(element))
			Expect(driver.MoveToCall.Point).To(BeNil())
		})

		Context("when moving over the element fails", func() {
			BeforeEach(func() {
				driver.MoveToCall.Err = errors.New("some error")
			})

			It("retuns an error", func() {
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to 'CSS: #selector': some error"))
			})
		})

		It("double-clicks on an element", func() {
			selection.DoubleClick()
			Expect(driver.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when the double-clicking the element fails", func() {
			BeforeEach(func() {
				driver.DoubleClickCall.Err = errors.New("some error")
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
			driver.GetElementsCall.ReturnElements = []types.Element{element}
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
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			return selection.Check()
		})

		It("checks the type of the checkbox", func() {
			selection.Check()
			Expect(element.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when the the driver fails to retrieve the 'type' attribute", func() {
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
			driver.GetElementsCall.ReturnElements = []types.Element{element}
			element.GetAttributeCall.ReturnValue = "checkbox"
			element.IsSelectedCall.ReturnSelected = true
		})

		It("clicks on an checked checkbox", func() {
			selection.Uncheck()
			Expect(element.ClickCall.Called).To(BeTrue())
		})
	})

	Describe("#Text", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		Context("if the the driver fails to retrieve the element text", func() {
			BeforeEach(func() {
				element.GetTextCall.Err = errors.New("some error")
			})

			It("returns an error", func() {
				_, err := selection.Text()
				Expect(err).To(MatchError("failed to retrieve text for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element text", func() {
			BeforeEach(func() {
				element.GetTextCall.ReturnText = "some text"
			})

			It("returns the text", func() {
				text, _ := selection.Text()
				Expect(text).To(Equal("some text"))
			})

			It("does not return an error", func() {
				_, err := selection.Text()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Attribute", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		It("requests the attribute value using the attribute name", func() {
			selection.Attribute("some-attribute")
			Expect(element.GetAttributeCall.Attribute).To(Equal("some-attribute"))
		})

		Context("if the the driver fails to retrieve the requested element attribute", func() {
			It("returns an error", func() {
				element.GetAttributeCall.Err = errors.New("some error")
				_, err := selection.Attribute("some-attribute")
				Expect(err).To(MatchError("failed to retrieve attribute value for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the requested element attribute", func() {
			BeforeEach(func() {
				element.GetAttributeCall.ReturnValue = "some value"
			})

			It("returns the attribute value", func() {
				value, _ := selection.Attribute("some-attribute")
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				_, err := selection.Attribute("some-attribute")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#CSS", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Text()
			return err
		})

		It("requests the CSS property value using the property name", func() {
			selection.CSS("some-property")
			Expect(element.GetCSSCall.Property).To(Equal("some-property"))
		})

		Context("if the the driver fails to retrieve the requested element CSS property", func() {
			It("returns an error", func() {
				element.GetCSSCall.Err = errors.New("some error")
				_, err := selection.CSS("some-property")
				Expect(err).To(MatchError("failed to retrieve CSS property for 'CSS: #selector': some error"))
			})
		})

		Context("if the driver succeeds in retrieving the requested element CSS property", func() {
			BeforeEach(func() {
				element.GetCSSCall.ReturnValue = "some value"
			})

			It("returns the property value", func() {
				value, _ := selection.CSS("some-property")
				Expect(value).To(Equal("some value"))
			})

			It("does not return an error", func() {
				_, err := selection.CSS("some-property")
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Selected", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Selected()
			return err
		})

		Context("if the the driver fails to retrieve the element's selected status", func() {
			It("returns an error", func() {
				element.IsSelectedCall.Err = errors.New("some error")
				_, err := selection.Selected()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is selected: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's selected status", func() {
			It("returns the selected status when selected", func() {
				element.IsSelectedCall.ReturnSelected = true
				value, _ := selection.Selected()
				Expect(value).To(BeTrue())
			})

			It("returns the selected status when not selected", func() {
				element.IsSelectedCall.ReturnSelected = false
				value, _ := selection.Selected()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Selected()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Visible", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Visible()
			return err
		})

		Context("if the the driver fails to retrieve the element's visible status", func() {
			It("returns an error", func() {
				element.IsDisplayedCall.Err = errors.New("some error")
				_, err := selection.Visible()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is visible: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's visible status", func() {
			It("returns the visible status when visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = true
				value, _ := selection.Visible()
				Expect(value).To(BeTrue())
			})

			It("returns the visible status when not visible", func() {
				element.IsDisplayedCall.ReturnDisplayed = false
				value, _ := selection.Visible()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Visible()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("#Enabled", func() {
		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.Enabled()
			return err
		})

		Context("if the the driver fails to retrieve the element's enabled status", func() {
			It("returns an error", func() {
				element.IsEnabledCall.Err = errors.New("some error")
				_, err := selection.Enabled()
				Expect(err).To(MatchError("failed to determine whether 'CSS: #selector' is enabled: some error"))
			})
		})

		Context("if the driver succeeds in retrieving the element's enabled status", func() {
			It("returns the enabled status when enabled", func() {
				element.IsEnabledCall.ReturnEnabled = true
				value, _ := selection.Enabled()
				Expect(value).To(BeTrue())
			})

			It("returns the enabled status when not enabled", func() {
				element.IsEnabledCall.ReturnEnabled = false
				value, _ := selection.Enabled()
				Expect(value).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.Enabled()
				Expect(err).NotTo(HaveOccurred())
			})
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
			driver.GetElementsCall.ReturnElements = []types.Element{optionOne, optionTwo, optionThree}
		})

		It("request child option elements from the driver", func() {
			selection.Select("some text")
			Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{"css selector", "#selector option"}))
		})

		Context("when the driver fails to retrieve any elements", func() {
			BeforeEach(func() {
				driver.GetElementsCall.Err = errors.New("some error")
			})

			It("returns error from the driver", func() {
				Expect(selection.Select("some text")).To(MatchError("failed to retrieve options for 'CSS: #selector': some error"))
			})
		})

		Context("when the driver fails to retrieve text for an element", func() {
			BeforeEach(func() {
				optionOne.GetTextCall.Err = errors.New("some error")
			})

			It("returns error from the driver", func() {
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
			driver.GetElementsCall.ReturnElements = []types.Element{element}
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

	Describe("#EqualsElement", func() {
		var (
			otherDriver    *mocks.Driver
			otherSelection types.Selection
			otherElement   *mocks.Element
		)

		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
			otherDriver = &mocks.Driver{}
			otherSelection = (&Selection{Driver: otherDriver}).Find("#other_selector")
			otherElement = &mocks.Element{}
			otherDriver.GetElementsCall.ReturnElements = []types.Element{otherElement}
		})

		ItShouldEnsureASingleElement(func() error {
			_, err := selection.EqualsElement(otherSelection)
			return err
		})

		It("ensures that the other selection is a single element", func() {
			otherDriver.GetElementsCall.ReturnElements = []types.Element{element, element}
			_, err := selection.EqualsElement(otherSelection)
			Expect(err).To(MatchError("failed to retrieve element with 'CSS: #other_selector': mutiple elements (2) were selected"))
		})

		It("compares the selection elements for equality", func() {
			selection.EqualsElement(otherSelection)
			Expect(element.IsEqualToCall.Element).To(Equal(otherElement))
		})

		Context("if the the driver fails to compare the elements", func() {
			It("returns an error", func() {
				element.IsEqualToCall.Err = errors.New("some error")
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).To(MatchError("failed to compare 'CSS: #selector' to 'CSS: #other_selector': some error"))
			})
		})

		Context("if the driver succeeds in comparing the elements", func() {
			It("returns true if they are equal", func() {
				element.IsEqualToCall.ReturnEquals = true
				equal, _ := selection.EqualsElement(otherSelection)
				Expect(equal).To(BeTrue())
			})

			It("returns false if they are not equal", func() {
				element.IsEqualToCall.ReturnEquals = false
				equal, _ := selection.EqualsElement(otherSelection)
				Expect(equal).To(BeFalse())
			})

			It("does not return an error", func() {
				_, err := selection.EqualsElement(otherSelection)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
