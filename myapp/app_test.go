package myapp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPathHander(t *testing.T) {

	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("Get", "/handleB", nil)

	mux := NewHttpHander()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := io.ReadAll(res.Body)

	assert.Equal("Welcome handleB", string(data))
}

func TestBHanderWithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("Get", "/handleB", nil)

	mux := NewHttpHander()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := io.ReadAll(res.Body)
	assert.Equal("Welcome handleB", string(data))

}

func TestHanderWithName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("Get", "/?name=kwon", nil)

	mux := NewHttpHander()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)
	data, _ := io.ReadAll(res.Body)
	assert.Equal("Hello kwon!", string(data))

}

func TestAHanderWithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/handleA",
		strings.NewReader(`{"first_name":"MinHyeok", "last_name":"Kwon", "email":"rkqhwk1@naver.com"}`))

	mux := NewHttpHander()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)

	assert.Equal(nil, err)
	assert.Equal("Kwon", user.LastName)
	assert.Equal("MinHyeok", user.FirstName)

}
