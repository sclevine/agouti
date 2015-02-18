package core

import (
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/core/internal/page"
)

type userPage struct {
	*page.Page
	*userSelection
}

func (u *userPage) SetCookie(cookie agouti.Cookie) error {
	return u.Page.SetCookie(cookie)
}

func (u *userPage) ReadLogs(logType string, all ...bool) ([]Log, error) {
	logs, err := u.Page.ReadLogs(logType, all...)
	if err != nil {
		return nil, err
	}

	var copiedLogs []Log
	for _, log := range logs {
		copiedLogs = append(copiedLogs, Log(log))
	}

	return copiedLogs, nil
}
