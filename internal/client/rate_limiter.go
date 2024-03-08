package client

import (
	"sync"
	"time"

	"go.uber.org/ratelimit"
)

const ApplicationDefaultRateLimit int = 5
const OrganizationDefaultRateLimit int = 10

var lock = &sync.Mutex{}

type RateLimiter struct {
	applicationRateLimiter  ratelimit.Limiter
	organisationRateLimiter ratelimit.Limiter
}

var rateLimiter *RateLimiter

// GetRateLimiter returns a singleton instance of RateLimiter.
// It will prevent multiple instances of RateLimiter from being created.
func GetRateLimiter(applicationRateLimit, organizationRateLimit int) *RateLimiter {
	if rateLimiter == nil {
		lock.Lock()
		defer lock.Unlock()
		if rateLimiter == nil {
			rateLimiter = &RateLimiter{
				applicationRateLimiter: ratelimit.New(
					// TODO: instead of reducing the limit by one, should implement a retry mechanism in the client.
					applicationRateLimit-1, ratelimit.Per(time.Second), ratelimit.WithSlack(0)),
				organisationRateLimiter: ratelimit.New(
					organizationRateLimit-1, ratelimit.Per(time.Minute), ratelimit.WithSlack(0)),
			}
		}
	}
	return rateLimiter
}
