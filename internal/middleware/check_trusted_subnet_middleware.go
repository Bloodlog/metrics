package middleware

import (
	"net"
	"net/http"

	"go.uber.org/zap"
)

func CheckTrustedSubnetMiddleware(
	logger *zap.SugaredLogger,
	trustedSubnet *net.IPNet,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if trustedSubnet == nil {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			realIP := r.Header.Get("X-Real-IP")
			ip := net.ParseIP(realIP)
			if ip == nil {
				logger.Infow("invalid or missing X-Real-IP", "ip", realIP)
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if !trustedSubnet.Contains(ip) {
				logger.Infow("unauthorized IP", "ip", ip.String())
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
