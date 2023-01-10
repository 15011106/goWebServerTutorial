package myapp

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestIndexPathHander(t *testing.T) {

// 	assert := assert.New(t)

// 	res := httptest.NewRecorder()
// 	req := httptest.NewRequest("Get", "/handleB", nil)

// 	mux := NewHttpHander()
// 	mux.ServeHTTP(res, req)

// 	assert.Equal(http.StatusOK, res.Code)
// 	data, _ := io.ReadAll(res.Body)

// 	assert.Contains(string(data), "Welcome handleB")
// }

// func TestBHanderWithoutName(t *testing.T) {
// 	assert := assert.New(t)

// 	res := httptest.NewRecorder()
// 	req := httptest.NewRequest("Get", "/handleB", nil)

// 	mux := NewHttpHander()
// 	mux.ServeHTTP(res, req)

// 	assert.Equal(http.StatusOK, res.Code)
// 	data, _ := io.ReadAll(res.Body)
// 	assert.Contains(string(data), "Welcome handleB")

// }

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

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHander())
	defer ts.Close()

	res, err := http.Get(ts.URL)

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
}

func HandleB(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHander())
	defer ts.Close()

	res, err := http.Get(ts.URL + "/handleB/30")

	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	data, _ := io.ReadAll(res.Body)
	assert.Contains(string(data), "30")
}

func TestCreate(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHttpHander())
	defer ts.Close()

	res, err := http.Post(ts.URL+"/handleB", "application/json",
		strings.NewReader(`{"first_name":"MinHyeok", "last_name":"Kwon", "email":"rkqhwk1@naver.com"}`))

	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)

	user := new(User)
	err = json.NewDecoder(res.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.Id)

	id := user.Id
	res, err = http.Get(ts.URL + "/handleB/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	user2 := new(User)
	err = json.NewDecoder(res.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.Id, user2.Id)
	assert.Equal(user.FirstName, user2.FirstName)
}
