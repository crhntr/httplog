# httplog

**This was initially pulled out of a side project then we used the repo as an acceptance test stub at Pivotal/VMware (note the silly commits)**

`httplog` is a http request logger to wrap your http.Handler.

[![GoDoc](https://godoc.org/github.com/crhntr/httplog?status.svg)](https://godoc.org/github.com/crhntr/httplog)

## Example
```go
mux := http.NewServeMux()
mux.HandleFunc("/greeting", func(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  fmt.Fprint(w, "Hello, world!")
})
logMux := httplog.Wrap(mux)
w := httptest.NewRecorder()
r := httptest.NewRequest(http.MethodGet, "/greeting", nil)
logMux.ServeHTTP(w, r)
// Output:
// {"type": "HTTP_REQUEST", "method": "GET", "path": "/greeting", "duration": "48.572µs", "status": 200}
```
