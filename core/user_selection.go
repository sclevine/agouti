package core

import (
	"errors"

	"github.com/sclevine/agouti/core/internal/selection"
)

type userSelection struct {
	*selection.Selection
}

func (u *userSelection) EqualsElement(comparable interface{}) (bool, error) {
	other, ok := comparable.(*userSelection)
	if !ok {
		return false, errors.New("provided object is not a Selection")
	}
	return u.Selection.EqualsElement(other.Selection)
}

func (u *userSelection) At(index int) Selection {
	return &userSelection{u.Selection.At(index)}
}

func (u *userSelection) Find(selector string) Selection {
	return &userSelection{u.AppendCSS(selector).Single()}
}

func (u *userSelection) FindByXPath(selector string) Selection {
	return &userSelection{u.AppendXPath(selector).Single()}
}

func (u *userSelection) FindByLink(text string) Selection {
	return &userSelection{u.AppendLink(text).Single()}
}

func (u *userSelection) FindByLabel(text string) Selection {
	return &userSelection{u.AppendLabeled(text).Single()}
}

func (u *userSelection) FindByButton(text string) Selection {
	return &userSelection{u.AppendButton(text).Single()}
}

func (u *userSelection) First(selector string) Selection {
	return &userSelection{u.AppendCSS(selector).At(0)}
}

func (u *userSelection) FirstByXPath(selector string) Selection {
	return &userSelection{u.AppendXPath(selector).At(0)}
}

func (u *userSelection) FirstByLink(text string) Selection {
	return &userSelection{u.AppendLink(text).At(0)}
}

func (u *userSelection) FirstByLabel(text string) Selection {
	return &userSelection{u.AppendLabeled(text).At(0)}
}

func (u *userSelection) FirstByButton(text string) Selection {
	return &userSelection{u.AppendButton(text).At(0)}
}

func (u *userSelection) All(selector string) MultiSelection {
	return &userSelection{u.AppendCSS(selector)}
}

func (u *userSelection) AllByXPath(selector string) MultiSelection {
	return &userSelection{u.AppendXPath(selector)}
}

func (u *userSelection) AllByLink(text string) MultiSelection {
	return &userSelection{u.AppendLink(text)}
}

func (u *userSelection) AllByLabel(text string) MultiSelection {
	return &userSelection{u.AppendLabeled(text)}
}

func (u *userSelection) AllByButton(text string) MultiSelection {
	return &userSelection{u.AppendButton(text)}
}
