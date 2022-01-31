package log

import (
	"fmt"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (l *Logger) LogMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         200,
		}
		start := time.Now()
		next.ServeHTTP(recorder, r)
		stop := time.Now()
		duration := stop.Sub(start)
		l.Request(r.Method, r.URL.Path, duration, fmt.Sprint(recorder.status))
	}

	return http.HandlerFunc(fn)
}
