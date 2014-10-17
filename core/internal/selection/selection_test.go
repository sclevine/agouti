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
				Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector"}))
			})

			It("retrieves the child elements of the parent selector", func() {
				Expect(parentOne.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children"}))
				Expect(parentTwo.GetElementsCall.Selector).To(Equal(types.Selector{Using: "xpath", Value: "children"}))
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

	Describe("#At & most methods: retrieving the selected element", func() {
		It("requests an element from the driver using the element's selector", func() {
			selection.Click()
			Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector"}))
		})

		Context("when the driver fails to retrieve any elements", func() {
			It("returns error from the driver", func() {
				driver.GetElementsCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': some error"))
			})
		})

		Context("when the driver retrieves zero elements", func() {
			It("fails with an error indicating there were no elements", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{}
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': no element found"))
			})
		})

		Context("when the driver retrieves more than one element and indexing is disabled", func() {
			It("returns an error with the number of elements", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{element, element}
				Expect(selection.Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': mutiple elements (2) were selected"))
			})
		})

		Context("when the selection index is out of range", func() {
			It("returns an error with the index and total number of elements", func() {
				driver.GetElementsCall.ReturnElements = []types.Element{element, element}
				Expect(selection.At(2).Click()).To(MatchError("failed to retrieve element with 'CSS: #selector': element index (2) out of range (>1)"))
			})
		})

		Context("when the index refers to a selected element", func() {
			It("processes only that element", func() {
				selectedElement := &mocks.Element{}
				driver.GetElementsCall.ReturnElements = []types.Element{element, element, selectedElement}
				selection.At(2).Click()
				Expect(element.ClickCall.Called).To(BeFalse())
				Expect(selectedElement.ClickCall.Called).To(BeTrue())
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
			Expect(selection.FindByLabel("label name").String()).To(Equal(`CSS: #selector | XPath: //input[@id=(//label[normalize-space(text())="label name"]/@for)] | //label[normalize-space(text())="label name"]/input`))
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
			Expect(driver.GetElementsCall.Selector).To(Equal(types.Selector{Using: "css selector", Value: "#selector"}))
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

	Describe("#EqualsElement", func() {
		var (
			otherDriver    *mocks.Driver
			otherSelection types.Selection
			otherElement   *mocks.Element
		)

		BeforeEach(func() {
			driver.GetElementsCall.ReturnElements = []types.Element{element}
			otherDriver = &mocks.Driver{}
			otherSelection = &Selection{Driver: otherDriver}
			otherSelection = otherSelection.Find("#other_selector")
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

		Context("if the provided element is not a *Selection", func() {
			It("returns an error", func() {
				_, err := selection.EqualsElement("not a selection")
				Expect(err).To(MatchError("provided object is not a selection"))
			})
		})

		Context("if the driver fails to compare the elements", func() {
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
