package headers

import (
	"net/http"

	cfg "tmpnotes/internal/config"
)

func AddStandardHeaders(h http.Header) {
	h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self' https://cdnjs.cloudflare.com; style-src 'self' https://cdn.jsdelivr.net")
	h.Set("X-Frame-Options", "DENY")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("X-XSS-Protection", "1; mode=block")
	if cfg.Config.EnableHsts {
		h.Set("Strict-Transport-Security", "max-age=15552000")
	}
}
