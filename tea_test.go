package tea_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/marconi/tea"
	"github.com/stretchr/testify/assert"
)

func TestTeaGet(t *testing.T) {
	tt := tea.New()
	tt.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Test")
	})

	ts := httptest.NewServer(tt)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	assert.Nil(t, err)
	assert.Equal(t, "Test", string(result))
}
