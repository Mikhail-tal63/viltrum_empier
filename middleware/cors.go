package middleware

import "net/http"

func IsAllowedOrigin(origin string) bool {
	_, ok := allowedOrigins[origin]
	return ok
}

var allowedOrigins = map[string]struct{}{
	"http://localhost:5173": {},
	"http://localhost:8081": {},
	"http://localhost:8082": {},
	"http://localhost:3000": {},
	"http://127.0.0.1:5173": {},
	"http://127.0.0.1:8081": {},
	"http://127.0.0.1:8082": {},
	"http://127.0.0.1:3000": {},
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if _, ok := allowedOrigins[origin]; ok {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
