package tea_test

import (
	"bytes"
	"fmt"
	"io"
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

	buf := bytes.NewBufferString("")
	_, err = io.Copy(buf, res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "Test", buf.String())
}
