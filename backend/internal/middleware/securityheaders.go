package middleware

import "net/http"

func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Prevent site framing
		w.Header().Set("X-Frame-Options", "DENY")

		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Enable XSS Protection in older browsers
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Basic referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Basic content security policy
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'; frame-ancestors 'none'; base-uri 'self'")

		next.ServeHTTP(w, r)
	})
}
