package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func performRequest(router *gin.Engine, method, url string) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, url, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	return w
}

func performPutRequest(router *gin.Engine, method, url string, body io.Reader) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, url, body)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, request)
	return w
}

func TestAlbumList(t *testing.T) {
	router := GetRouter()
	w := performRequest(router, "GET", "/albums")
	if w.Code != http.StatusOK {
		t.Fatal("Status not ok")
	}
	var response []album
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}
	for i := range response {
		if response[i] != albums[i] {
			t.Fatal("Status not ok")
		}
	}
}

func TestAlbumById(t *testing.T) {
	router := GetRouter()
	for i := 1; i < 4; i++ {
		w := performRequest(router, "GET", "/albums/"+strconv.Itoa(i))
		if w.Code != http.StatusOK {
			t.Fatal("status not ok")
		}
		var response album
		err := json.Unmarshal([]byte(w.Body.String()), &response)
		if err != nil {
			t.Fatal(err)
		}
		if response != albums[i-1] {
			t.Fatal("Not equal Album")
		}
	}

	for i := 4; i < 10; i++ {
		w := performRequest(router, "GET", "/albums/"+strconv.Itoa(i))
		if w.Code != http.StatusNotFound {
			t.Fatal("status should be 404")
		}
	}
}

func TestDeleteAlbum(t *testing.T) {
	router := GetRouter()
	w := performRequest(router, "DELETE", "/albums/1")
	if w.Code != http.StatusNoContent {
		t.Fatal("status must be 204")
	}
	for _, a := range albums {
		if a.ID == "1" {
			t.Fatal("It should be deleted")
		}
	}
}

func TestUpdateAlbum(t *testing.T) {
	router := GetRouter()
	w := performPutRequest(router, "PUT", "/albums/2", strings.NewReader(`{"title": "test"}`))
	if w.Code != http.StatusOK {
		t.Fatal("status must be OK")
	}

	var response album

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	for _, a := range albums {
		if a.ID == "2" && response != a {
			t.Fatal("It should be updated")
		}
	}
}

func TestPostNewAlbum(t *testing.T) {
	router := GetRouter()
	w := performPutRequest(router, "POST", "/albums", strings.NewReader(`{"id": "4","title": "The Moder Sound","artist": "Molodets","price": 29.9}`))
	if w.Code != http.StatusCreated {
		t.Fatal("status must be Created")
	}

	var response album

	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if err != nil {
		t.Fatal(err)
	}

	if albums[len(albums)-1] != response {
		t.Fatal("It should be added")
	}

	w = performPutRequest(router, "POST", "/albums", strings.NewReader(`{"title": "The Moder Sound","artist": "Molodets","price": 29.9}`))
	if w.Code != http.StatusBadRequest {
		t.Fatal("status must Bad Request")
	}

}
