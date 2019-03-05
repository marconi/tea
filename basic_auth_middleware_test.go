package tea_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marconi/tea"
	"github.com/stretchr/testify/suite"
)

type BasicAuthMiddlewareTestSuite struct {
	suite.Suite
	t  *tea.Tea
	ts *httptest.Server
}

func (suite *BasicAuthMiddlewareTestSuite) SetupTest() {
	suite.t = tea.New()
	suite.t.Use(tea.NewBasicAuthMiddleware(
		"Restricted",
		func(userId string, password string) bool {
			if userId == "admin" && password == "admin" {
				return true
			}
			return false
		},
		nil,
	))
	suite.t.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Test")
	})

	suite.ts = httptest.NewServer(suite.t)
}

func (suite *BasicAuthMiddlewareTestSuite) TearDownTest() {
	suite.ts.Close()
}

func (suite *BasicAuthMiddlewareTestSuite) TestInitialReq() {
	res, err := http.Get(suite.ts.URL)
	suite.Nil(err)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	suite.Nil(err)
	suite.Equal("Unauthorized\n", string(result))
	suite.Equal(res.Header.Get("WWW-Authenticate"), "Basic realm=Restricted")
}

func (suite *BasicAuthMiddlewareTestSuite) TestInvalidAuthHeader() {
	r, err := http.NewRequest("GET", suite.ts.URL, nil)
	suite.Nil(err)
	r.Header.Set("Authorization", "Something else")

	res, err := http.DefaultClient.Do(r)
	suite.Nil(err)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	suite.Nil(err)
	suite.Equal("Invalid authentication\n", string(result))
}

func (suite *BasicAuthMiddlewareTestSuite) TestInvalidAuthCreds() {
	encoded := base64.StdEncoding.EncodeToString([]byte("foo:bar"))
	r, err := http.NewRequest("GET", suite.ts.URL, nil)
	suite.Nil(err)
	r.Header.Set("Authorization", fmt.Sprintf("Basic: %s", encoded))

	res, err := http.DefaultClient.Do(r)
	suite.Nil(err)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	suite.Nil(err)
	suite.Equal("Invalid authentication\n", string(result))
}

func (suite *BasicAuthMiddlewareTestSuite) TestValidAuthCreds() {
	encoded := base64.StdEncoding.EncodeToString([]byte("admin:admin"))
	r, err := http.NewRequest("GET", suite.ts.URL, nil)
	suite.Nil(err)
	r.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	res, err := http.DefaultClient.Do(r)
	suite.Nil(err)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	suite.Nil(err)
	suite.Equal("Test", string(result))
}

func TestBasicAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(BasicAuthMiddlewareTestSuite))
}
