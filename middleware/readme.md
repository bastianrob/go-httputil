# HTTP Middleware

A middleware to chain HTTP request handler

## Example

```go
step1 := func(number int) HTTPMiddleware {
    return func(h http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("number", strconv.Itoa(number))
            h.ServeHTTP(w, r)
        }
    }
}

step2 := func(word string) HTTPMiddleware {
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

pipeline := httpmiddleware.NewPipeline().
        Do(step1(100)).
        Do(step2("I am a word")).
        For(api)

router := http.NewServeMux()
router.Handle("/", pipeline)

req := httptest.NewRequest("GET", "/", nil)
rec := httptest.NewRecorder()

router.ServeHTTP(rec, req)
log.Println(rec.Code, rec.Header().Get("number"), rec.Header().Get("word"))
//200 100 I am a word
```
