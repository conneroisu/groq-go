package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
)

// ToolhouseTestServer creates a test server for emulating the toolhouse api.
func (ts *ServerTest) ToolhouseTestServer() *httptest.Server {
	return httptest.NewUnstartedServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf(
				"received a %s request at path %q\n",
				r.Method,
				r.URL.Path,
			)

			// check auth
			if r.Header.Get("Authorization") != "Bearer "+GetTestToken() &&
				r.Header.Get("api-key") != GetTestToken() {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Handle /path/* routes.
			// Note: the * is converted to a .* in register handler for proper regex handling
			for route, handler := range ts.handlers {
				// Adding ^ and $ to make path matching deterministic since go map iteration isn't ordered
				pattern, _ := regexp.Compile("^" + route + "$")
				if pattern.MatchString(r.URL.Path) {
					handler(w, r)
					return
				}
			}
			http.Error(
				w,
				"the resource path doesn't exist",
				http.StatusNotFound,
			)
		}),
	)
}
