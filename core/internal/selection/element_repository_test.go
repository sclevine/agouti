package selection_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti/core/internal/api"
	"github.com/sclevine/agouti/core/internal/mocks"
	. "github.com/sclevine/agouti/core/internal/selection"
	"github.com/sclevine/agouti/core/internal/types"
)

var _ = Describe("Elements", func() {
	var (
		client     *mocks.Client
		repository *ElementRepository
	)

	BeforeEach(func() {
		client = &mocks.Client{}
		repository = &ElementRepository{Client: client}
	})

	Describe("#GetAtLeastOne", func() {
		Context("when the client retrieves zero elements", func() {
			It("should fail with an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.GetAtLeastOne([]types.Selector{types.Selector{}})
				Expect(err).To(MatchError("no elements found"))
			})
		})
	})

	Describe("#GetExactlyOne", func() {
		Context("when the client retrieves zero elements", func() {
			It("should return an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.GetExactlyOne([]types.Selector{types.Selector{}})
				Expect(err).To(MatchError("no elements found"))
			})
		})

		Context("when the client retrieves more than one element", func() {
			It("should return an error", func() {
				client.GetElementsCall.ReturnElements = []*api.Element{&api.Element{}, &api.Element{}}
				_, err := repository.GetExactlyOne([]types.Selector{types.Selector{}})
				Expect(err).To(MatchError("method does not support multiple elements (2)"))
			})
		})
	})

	Describe("#Get", func() {
		var (
			firstParentSession  *mocks.Session
			firstParent         *api.Element
			secondParentSession *mocks.Session
			secondParent        *api.Element
			children            []Element
			parentSelector      types.Selector
			parentSelectorJSON  string
			childSelector       types.Selector
			childSelectorJSON   string
		)

		BeforeEach(func() {
			firstParentSession = &mocks.Session{}
			firstParent = &api.Element{Session: firstParentSession}
			secondParentSession = &mocks.Session{}
			secondParent = &api.Element{Session: secondParentSession}
			children = []Element{
				Element(&api.Element{ID: "first child", Session: firstParentSession}),
				Element(&api.Element{ID: "second child", Session: firstParentSession}),
				Element(&api.Element{ID: "third child", Session: secondParentSession}),
				Element(&api.Element{ID: "fourth child", Session: secondParentSession}),
			}
			firstParentSession.ExecuteCall.Result = `[{"ELEMENT": "first child"}, {"ELEMENT": "second child"}]`
			secondParentSession.ExecuteCall.Result = `[{"ELEMENT": "third child"}, {"ELEMENT": "fourth child"}]`
			client.GetElementsCall.ReturnElements = []*api.Element{firstParent, secondParent}
			parentSelector = types.Selector{Using: "css selector", Value: "parents"}
			parentSelectorJSON = `{"using": "css selector", "value": "parents"}`
			childSelector = types.Selector{Using: "xpath", Value: "children"}
			childSelectorJSON = `{"using": "xpath", "value": "children"}`
		})

		Context("when all elements are successfully retrieved", func() {
			It("should retrieve the parent elements using the client", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector))
			})

			It("should retrieve the child elements of the parent selector", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(firstParentSession.ExecuteCall.BodyJSON).To(MatchJSON(childSelectorJSON))
				Expect(secondParentSession.ExecuteCall.BodyJSON).To(MatchJSON(childSelectorJSON))

			})

			It("should successfully return all of the children", func() {
				selectedChildren, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(selectedChildren).To(Equal(children))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when all non-zero-indexed elements are successfully retrieved", func() {
			BeforeEach(func() {
				parentSelector.Index = 1
				parentSelector.Indexed = true
				childSelector.Index = 1
				childSelector.Indexed = true
			})

			It("should retrieve the parent elements using the client", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector))
			})

			It("should retrieve the child elements of the second parent selector", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(firstParentSession.ExecuteCall.BodyJSON).To(BeEmpty())
				Expect(secondParentSession.ExecuteCall.BodyJSON).To(MatchJSON(childSelectorJSON))
			})

			It("should return only the selected child elements", func() {
				selectedChildren, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(selectedChildren).To(Equal([]Element{children[3]}))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when all zero-indexed elements are successfully retrieved", func() {
			BeforeEach(func() {
				firstParentSession.ExecuteCall.Result = `{"ELEMENT": "first child"}`
				client.GetElementCall.ReturnElement = firstParent
				parentSelector.Index = 0
				parentSelector.Indexed = true
				childSelector.Index = 0
				childSelector.Indexed = true
			})

			It("should retrieve the first parent element using the client", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(client.GetElementCall.Selector).To(Equal(parentSelector))
			})

			It("should retrieve the first child element of the parent selector", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(firstParentSession.ExecuteCall.BodyJSON).To(MatchJSON(childSelectorJSON))
			})

			It("should return only the selected child element", func() {
				selectedChildren, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(selectedChildren).To(Equal([]Element{children[0]}))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when single-element-only elements are successfully retrieved", func() {
			BeforeEach(func() {
				firstParentSession.ExecuteCall.Result = `[{"ELEMENT": "first child"}]`
				client.GetElementsCall.ReturnElements = []*api.Element{firstParent}
				parentSelector.Single = true
				childSelector.Single = true
			})

			It("should retrieve the parent element using the client", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(client.GetElementsCall.Selector).To(Equal(parentSelector))
			})

			It("should retrieve the child element of the parent selector", func() {
				repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(firstParentSession.ExecuteCall.BodyJSON).To(MatchJSON(childSelectorJSON))
			})

			It("should return only the selected child element", func() {
				selectedChildren, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(selectedChildren).To(Equal([]Element{children[0]}))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is no selection", func() {
			It("should return an error", func() {
				_, err := repository.Get([]types.Selector{})
				Expect(err).To(MatchError("empty selection"))
			})
		})

		Context("when retrieving the parent elements fails", func() {
			It("should return an error", func() {
				client.GetElementsCall.Err = errors.New("some error")
				_, err := repository.Get([]types.Selector{parentSelector})
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when retrieving any of the child elements fails", func() {
			It("should return an error", func() {
				secondParentSession.ExecuteCall.Err = errors.New("some error")
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when a single-element-only parent selection refers to multiple parents", func() {
			It("should return an error", func() {
				parentSelector.Single = true
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("ambiguous find"))
			})
		})

		Context("when a single-element-only parent selection refers to no parents", func() {
			It("should return an error", func() {
				parentSelector.Single = true
				client.GetElementsCall.ReturnElements = []*api.Element{}
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("element not found"))
			})
		})

		Context("when any single-element-only child selection refers to multiple child elements", func() {
			It("should return an error", func() {
				childSelector.Single = true
				firstParentSession.ExecuteCall.Result = `[{"ELEMENT": "first child"}]`
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("ambiguous find"))
			})
		})

		Context("when any single-element-only child selection refers to no child elements", func() {
			It("should return an error", func() {
				childSelector.Single = true
				firstParentSession.ExecuteCall.Result = `[]`
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("element not found"))
			})
		})

		Context("when the parent selection index is out of range", func() {
			It("should return an error", func() {
				parentSelector.Index = 2
				parentSelector.Indexed = true
				_, err := repository.Get([]types.Selector{parentSelector})
				Expect(err).To(MatchError("element index out of range"))
			})
		})

		Context("when child selection indices are out of range", func() {
			It("should return an error", func() {
				parentSelector.Index = 1
				parentSelector.Indexed = true
				childSelector.Index = 2
				childSelector.Indexed = true
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("element index out of range"))
			})
		})

		Context("when a zero-indexed parent selection element does not exist", func() {
			It("should return an error", func() {
				client.GetElementCall.Err = errors.New("some error")
				parentSelector.Index = 0
				parentSelector.Indexed = true
				_, err := repository.Get([]types.Selector{parentSelector})
				Expect(err).To(MatchError("some error"))
			})
		})

		Context("when a zero-indexed child selection element does not exist", func() {
			It("should return an error", func() {
				firstParentSession.ExecuteCall.Err = errors.New("some error")
				client.GetElementCall.ReturnElement = firstParent
				parentSelector.Index = 0
				parentSelector.Indexed = true
				childSelector.Index = 0
				childSelector.Indexed = true
				_, err := repository.Get([]types.Selector{parentSelector, childSelector})
				Expect(err).To(MatchError("some error"))
			})
		})
	})
})
