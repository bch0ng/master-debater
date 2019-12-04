package handlers

import "net/http"

// CORSMiddleware stores the HTTP handler
type CORSMiddleware struct {
	handler http.Handler
}

// ServeHTTP serves HTTP with CORS enabled.
func (c *CORSMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	c.handler.ServeHTTP(w, r)
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	(*w).Header().Set("Access-Control-Expose-Headers", "Authorization")
	(*w).Header().Set("Access-Control-Max-Age", "600")
	(*w).Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, access-control-expose-headers, Access-Control-Allow-Headers, access-control-allow-origin, X-Requested-With, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

}

// NewCORSMiddleware initializes a new CORSMiddleware struct with the given
// HTTP handler.
func NewCORSMiddleware(handler http.Handler) *CORSMiddleware {
	return &CORSMiddleware{handler}
}
