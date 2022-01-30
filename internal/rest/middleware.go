package rest

import (
	"fmt"
	"log"
	"net/http"
)

func logMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("Request to %v", r.URL.Path))
		next.ServeHTTP(w, r)
		log.Println(fmt.Sprintf("End of request %v", r.URL.Path))
	}

	return http.HandlerFunc(fn)
}
