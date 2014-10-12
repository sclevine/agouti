package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sclevine/agouti/page"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var server *httptest.Server
var submitted bool

var _ = BeforeSuite(func() {
	StartPhantom()
	server = httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		submitted = request.Method == "POST"
		html, _ := ioutil.ReadFile("test_page.html")
		response.Write(html)
	}))
})

var _ = AfterSuite(func() {
	server.Close()
	StopPhantom()
})
