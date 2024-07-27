package middleware

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Sherrira/rateLimiter/internal/infra/api/auth"
	"github.com/gorilla/mux"
)

type RateLimiterConfiguration struct {
	IpRateLimit   int
	TokenLimits   map[string]int
	BlockDuration time.Duration
}

type RateLimiter interface {
	IsRateLimited(ctx context.Context, key string, limit int) bool
}

func RateLimitMiddleware(rl RateLimiter, ipRateLimit int, authorizer auth.Authorizer) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Invalid remote address", http.StatusInternalServerError)
				return
			}
			apiKey := r.Header.Get("API_KEY")

			if apiKey != "" {
				if token := authorizer.Authorize(apiKey); token != nil {
					if rl.IsRateLimited(ctx, "token:"+token.Token, token.Limit) {
						http.Error(w, "your token have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
						return
					}
				}
			} else if rl.IsRateLimited(ctx, "ip:"+ip, ipRateLimit) {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
