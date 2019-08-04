package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPMiddleware(t *testing.T) {
	adapter1 := func(number int) HTTPMiddleware {
		return func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("number", strconv.Itoa(number))
				h.ServeHTTP(w, r)
			}
		}
	}

	adapter2 := func(word string) HTTPMiddleware {
		return func(h http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("word", word)
				h.ServeHTTP(w, r)
			}
		}
	}

	api := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	router := http.NewServeMux()
	router.Handle("/", NewPipeline().
		Do(adapter1(100)).
		Do(adapter2("I am a word")).
		For(api))

	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "100", rec.Header().Get("number"))
	assert.Equal(t, "I am a word", rec.Header().Get("word"))
}
