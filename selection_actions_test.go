package agouti_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	"github.com/sclevine/agouti/internal/mocks"
)

var _ = Describe("Selection Actions", func() {
	var (
		selection         *MultiSelection
		session           *mocks.Session
		elementRepository *mocks.ElementRepository
		firstElement      *mocks.Element
		secondElement     *mocks.Element
	)

	BeforeEach(func() {
		session = &mocks.Session{}
		firstElement = &mocks.Element{}
		secondElement = &mocks.Element{}
		elementRepository = &mocks.ElementRepository{}
		selection = NewTestMultiSelection(elementRepository, session, "#selector")
		elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{firstElement, secondElement}
	})

	Describe("#Click", func() {
		It("should successfully click on all selected elements", func() {
			Expect(selection.Click()).To(Succeed())
			Expect(firstElement.ClickCall.Called).To(BeTrue())
			Expect(secondElement.ClickCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Click()).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
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
		var apiElement *api.Element

		BeforeEach(func() {
			apiElement = &api.Element{}
			elementRepository.GetAtLeastOneCall.ReturnElements = []element.Element{&api.Element{}, apiElement}
		})

		It("should successfully move the mouse to the middle of each selected element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(session.MoveToCall.Element).To(Equal(apiElement))
			Expect(session.MoveToCall.Offset).To(BeNil())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})

		Context("when moving over any element fails", func() {
			It("should retun an error", func() {
				session.MoveToCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to move mouse to 'CSS: #selector': some error"))
			})
		})

		It("should successfully double-click on each element", func() {
			Expect(selection.DoubleClick()).To(Succeed())
			Expect(session.DoubleClickCall.Called).To(BeTrue())
		})

		Context("when the double-clicking any element fails", func() {
			It("should return an error", func() {
				session.DoubleClickCall.Err = errors.New("some error")
				Expect(selection.DoubleClick()).To(MatchError("failed to double-click on 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#Fill", func() {
		It("should successfully clear each element", func() {
			Expect(selection.Fill("some text")).To(Succeed())
			Expect(firstElement.ClearCall.Called).To(BeTrue())
			Expect(secondElement.ClearCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Fill("some text")).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
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
		It("should successfully check the type of each checkbox", func() {
			firstElement.GetAttributeCall.ReturnValue = "checkbox"
			secondElement.GetAttributeCall.ReturnValue = "checkbox"
			Expect(selection.Check()).To(Succeed())
			Expect(firstElement.GetAttributeCall.Attribute).To(Equal("type"))
			Expect(secondElement.GetAttributeCall.Attribute).To(Equal("type"))
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Check()).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
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
			firstOptionBuses  []*mocks.Bus
			secondOptionBuses []*mocks.Bus
			firstOptions      []*api.Element
			secondOptions     []*api.Element
		)

		BeforeEach(func() {
			firstOptionBuses = []*mocks.Bus{&mocks.Bus{}, &mocks.Bus{}}
			secondOptionBuses = []*mocks.Bus{&mocks.Bus{}, &mocks.Bus{}}
			firstOptions = []*api.Element{
				&api.Element{ID: "one", Session: &api.Session{firstOptionBuses[0]}},
				&api.Element{ID: "two", Session: &api.Session{firstOptionBuses[1]}},
			}
			secondOptions = []*api.Element{
				&api.Element{ID: "three", Session: &api.Session{secondOptionBuses[0]}},
				&api.Element{ID: "four", Session: &api.Session{secondOptionBuses[1]}},
			}
			firstElement.GetElementsCall.ReturnElements = []*api.Element{firstOptions[0], firstOptions[1]}
			secondElement.GetElementsCall.ReturnElements = []*api.Element{secondOptions[0], secondOptions[1]}
		})

		It("should successfully retrieve the options with matching text for each selected element", func() {
			Expect(selection.Select("some text")).To(Succeed())
			Expect(firstElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(firstElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space()="some text"]`))
			Expect(secondElement.GetElementsCall.Selector.Using).To(Equal("xpath"))
			Expect(secondElement.GetElementsCall.Selector.Value).To(Equal(`./option[normalize-space()="some text"]`))
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})

		Context("when we fail to retrieve any option", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError("failed to select specified option for some 'CSS: #selector': some error"))
			})
		})

		Context("when any of the elements has no options with matching text", func() {
			It("should return an error", func() {
				secondElement.GetElementsCall.ReturnElements = []*api.Element{}
				Expect(selection.Select("some text")).To(MatchError(`no options with text "some text" found for some 'CSS: #selector'`))
			})
		})

		It("should successfully click on all options with matching text", func() {
			Expect(selection.Select("some text")).To(Succeed())
			Expect(firstOptionBuses[0].SendCall.Endpoint).To(Equal("element/one/click"))
			Expect(firstOptionBuses[1].SendCall.Endpoint).To(Equal("element/two/click"))
			Expect(secondOptionBuses[0].SendCall.Endpoint).To(Equal("element/three/click"))
			Expect(secondOptionBuses[1].SendCall.Endpoint).To(Equal("element/four/click"))
		})

		Context("when the click fails for any of the options", func() {
			It("should return an error", func() {
				secondOptionBuses[1].SendCall.Err = errors.New("some error")
				Expect(selection.Select("some text")).To(MatchError(`failed to click on option with text "some text" for some 'CSS: #selector': some error`))
			})
		})
	})

	Describe("#Submit", func() {
		It("should successfully submit all selected elements", func() {
			Expect(selection.Submit()).To(Succeed())
			Expect(firstElement.SubmitCall.Called).To(BeTrue())
			Expect(secondElement.SubmitCall.Called).To(BeTrue())
		})

		Context("when zero elements are returned", func() {
			It("should return an error", func() {
				elementRepository.GetAtLeastOneCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to select 'CSS: #selector': some error"))
			})
		})

		Context("when any submit fails", func() {
			It("should return an error", func() {
				secondElement.SubmitCall.Err = errors.New("some error")
				Expect(selection.Submit()).To(MatchError("failed to submit 'CSS: #selector': some error"))
			})
		})
	})
})
