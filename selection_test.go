package agouti_test

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"github.com/sclevine/agouti/internal/element"
	. "github.com/sclevine/agouti/internal/matchers"
	"github.com/sclevine/agouti/internal/mocks"
)

var minimalpng = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a,
	0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52, 0x00,
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x06, 0x00,
	0x00, 0x00, 0x1f, 0x15, 0xc4, 0x89, 0x00, 0x00, 0x00, 0x11,
	0x49, 0x44, 0x41, 0x54, 0x78, 0x9c, 0x62, 0x62, 0x60, 0x60,
	0x60, 0x00, 0x04, 0x00, 0x00, 0xff, 0xff, 0x00, 0x0f, 0x00,
	0x03, 0xfe, 0x8f, 0xeb, 0xcf, 0x00, 0x00, 0x00, 0x00, 0x49,
	0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82}

var _ = Describe("Selection", func() {
	var (
		firstElement  *mocks.Element
		secondElement *api.Element
	)

	BeforeEach(func() {
		firstElement = &mocks.Element{}
		secondElement = &api.Element{}
	})

	Describe("#String", func() {
		It("should return a string representation of the selection", func() {
			selection := NewTestMultiSelection(nil, nil, "#selector", nil)
			Expect(selection.AllByXPath("#subselector").String()).To(Equal("selection 'CSS: #selector | XPath: #subselector'"))
		})
	})

	Describe("#Elements", func() {
		var (
			selection         *Selection
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			selection = NewTestSelection(nil, elementRepository, "#selector", nil)
		})

		It("should return a []*api.Elements retrieved from the element repository", func() {
			elements := []*api.Element{{ID: "first"}, {ID: "second"}}
			elementRepository.GetCall.ReturnElements = []element.Element{elements[0], elements[1]}
			Expect(selection.Elements()).To(Equal(elements))
		})

		Context("when retrieving the elements fails", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Elements()
				Expect(err).To(MatchError("some error"))
			})
		})
	})

	Describe("#Count", func() {
		var (
			selection         *MultiSelection
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			selection = NewTestMultiSelection(nil, elementRepository, "#selector", nil)
			elementRepository.GetCall.ReturnElements = []element.Element{firstElement, secondElement}
		})

		It("should successfully return the number of elements", func() {
			Expect(selection.Count()).To(Equal(2))
		})

		Context("when the the session fails to retrieve the elements", func() {
			It("should return an error", func() {
				elementRepository.GetCall.Err = errors.New("some error")
				_, err := selection.Count()
				Expect(err).To(MatchError("failed to select elements from selection 'CSS: #selector': some error"))
			})
		})
	})

	Describe("#EqualsElement", func() {
		var (
			firstSelection          *Selection
			secondSelection         *Selection
			firstElementRepository  *mocks.ElementRepository
			secondElementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			firstElementRepository = &mocks.ElementRepository{}
			firstElementRepository.GetExactlyOneCall.ReturnElement = firstElement
			firstSelection = NewTestSelection(nil, firstElementRepository, "#first_selector", nil)

			secondElementRepository = &mocks.ElementRepository{}
			secondElementRepository.GetExactlyOneCall.ReturnElement = secondElement
			secondSelection = NewTestSelection(nil, secondElementRepository, "#second_selector", nil)
		})

		It("should compare the selection elements for equality", func() {
			firstSelection.EqualsElement(secondSelection)
			Expect(firstElement.IsEqualToCall.Element).To(ExactlyEqual(secondElement))
		})

		It("should successfully return true if they are equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = true
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeTrue())
		})

		It("should successfully return false if they are not equal", func() {
			firstElement.IsEqualToCall.ReturnEquals = false
			Expect(firstSelection.EqualsElement(secondSelection)).To(BeFalse())
		})

		Context("when the provided object is a *MultiSelection", func() {
			It("should not fail", func() {
				multiSelection := NewTestMultiSelection(nil, secondElementRepository, "#multi_selector", nil)
				Expect(firstSelection.EqualsElement(multiSelection)).To(BeFalse())
				Expect(firstElement.IsEqualToCall.Element).To(ExactlyEqual(secondElement))
			})
		})

		Context("when the provided object is not a type of selection", func() {
			It("should return an error", func() {
				_, err := firstSelection.EqualsElement("not a selection")
				Expect(err).To(MatchError("must be *Selection or *MultiSelection"))
			})
		})

		Context("when there is an error retrieving elements from the selection", func() {
			It("should return an error", func() {
				firstElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #first_selector [single]': some error"))
			})
		})

		Context("when there is an error retrieving elements from the other selection", func() {
			It("should return an error", func() {
				secondElementRepository.GetExactlyOneCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #second_selector [single]': some error"))
			})
		})

		Context("when the session fails to compare the elements", func() {
			It("should return an error", func() {
				firstElement.IsEqualToCall.Err = errors.New("some error")
				_, err := firstSelection.EqualsElement(secondSelection)
				Expect(err).To(MatchError("failed to compare selection 'CSS: #first_selector [single]' to selection 'CSS: #second_selector [single]': some error"))
			})
		})
	})

	Describe("#MouseToElement", func() {
		var (
			selection         *Selection
			session           *mocks.Session
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			elementRepository = &mocks.ElementRepository{}
			elementRepository.GetExactlyOneCall.ReturnElement = secondElement
			session = &mocks.Session{}
			selection = NewTestSelection(session, elementRepository, "#selector", nil)
		})

		It("should successfully instruct the session to move the mouse over the selection", func() {
			Expect(selection.MouseToElement()).To(Succeed())
			Expect(session.MoveToCall.Element).To(Equal(secondElement))
			Expect(session.MoveToCall.Offset).To(BeNil())
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				err := selection.MouseToElement()
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector [single]': some error"))
			})
		})

		Context("when the session fails to move the mouse to the element", func() {
			It("should return an error", func() {
				session.MoveToCall.Err = errors.New("some error")
				err := selection.MouseToElement()
				Expect(err).To(MatchError("failed to move mouse to element for selection 'CSS: #selector [single]': some error"))
			})
		})
	})

	Describe("#Screenshot", func() {
		var (
			selection         *Selection
			session           *mocks.Session
			cropper           *mocks.Cropper
			firstElement      *mocks.Element
			elementRepository *mocks.ElementRepository
		)

		BeforeEach(func() {
			firstElement = &mocks.Element{}
			elementRepository = &mocks.ElementRepository{}
			elementRepository.GetExactlyOneCall.ReturnElement = firstElement
			session = &mocks.Session{}
			cropper = &mocks.Cropper{}
			selection = NewTestSelection(session, elementRepository, "#selector", cropper)
		})

		It("should successfully return the screenshot", func() {
			firstElement.GetScreenshotCall.ReturnImage = []byte("some-image")
			filename, _ := filepath.Abs(".test.screenshot.png")
			Expect(selection.Screenshot(".test.screenshot.png")).To(Succeed())
			defer os.Remove(filename)
			result, _ := ioutil.ReadFile(filename)
			Expect(string(result)).To(Equal("some-image"))
		})

		Context("when a new screenshot file cannot be saved", func() {
			It("should return an error", func() {
				err := selection.Screenshot("")
				Expect(err.Error()).To(ContainSubstring("failed to save screenshot: open"))
			})
		})

		Context("when the element repository fails to return exactly one element", func() {
			It("should return an error", func() {
				elementRepository.GetExactlyOneCall.Err = errors.New("some error")
				err := selection.Screenshot(".test.screenshot.png")
				Expect(err).To(MatchError("failed to select element from selection 'CSS: #selector [single]': some error"))
			})
		})

		Context("when the selection fails to retrieve a screenshot", func() {
			BeforeEach(func() {
				firstElement.GetScreenshotCall.Err = errors.New("some error")
			})

			It("should fall back to using the session screenshot and cropping", func() {
				session.GetScreenshotCall.ReturnImage = minimalpng
				cropper.ReturnImage, _ = png.Decode(bytes.NewBuffer(minimalpng))
				filename, _ := filepath.Abs(".test.screenshot.png")
				Expect(selection.Screenshot(".test.screenshot.png")).To(Succeed())
				defer os.Remove(filename)
				result, _ := ioutil.ReadFile(filename)
				Expect(result).To(Equal(minimalpng))
			})

			Context("and the session fails to retrieve a screenshot", func() {
				It("should return an error", func() {
					session.GetScreenshotCall.Err = errors.New("some error")
					err := selection.Screenshot(".test.screenshot.png")
					Expect(err).To(MatchError("failed to retrieve screenshot: some error"))
				})
			})

			Context("and the session screenshot cannot be decoded", func() {
				It("should return an error", func() {
					session.GetScreenshotCall.ReturnImage = []byte("some-image")
					err := selection.Screenshot(".test.screenshot.png")
					Expect(err).To(MatchError("failed to decode screenshot: png: invalid format: not a PNG file"))
				})
			})

			Context("and the selections bounding rectangle cannot be retrieved", func() {
				It("should return an error", func() {
					session.GetScreenshotCall.ReturnImage = minimalpng
					firstElement.GetRectCall.Err = errors.New("some error")
					err := selection.Screenshot(".test.screenshot.png")
					Expect(err).To(MatchError("failed to retrieve bounds for selection: some error"))
				})
			})

			Context("and the image cannot be cropped", func() {
				It("should return an error", func() {
					session.GetScreenshotCall.ReturnImage = minimalpng
					cropper.Err = errors.New("some error")
					err := selection.Screenshot(".test.screenshot.png")
					Expect(err).To(MatchError("failed to crop screenshot: some error"))
				})
			})

			Context("and the image cannot be encoded", func() {
				It("should return an error", func() {
					session.GetScreenshotCall.ReturnImage = minimalpng
					cropper.ReturnImage = image.Rectangle{}
					err := selection.Screenshot(".test.screenshot.png")
					Expect(err).To(MatchError("failed to encode screenshot: png: invalid format: invalid image size: 0x0"))
				})
			})
		})
	})
})
