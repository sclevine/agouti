package integration

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var (
	Submitted = false
	handler   = func(response http.ResponseWriter, request *http.Request) {
		if request.Method == "POST" {
			Submitted = true
		}
		html, _ := ioutil.ReadFile("test_page.html")
		response.Write(html)
	}

	Server = httptest.NewUnstartedServer(http.HandlerFunc(handler))
)
