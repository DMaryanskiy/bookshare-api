package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	Redis *redis.Client
	Rules map[string]RateLimitRule
}

type RateLimitRule struct {
	Limit   int
	Window  time.Duration
}

func NewRateLimiter(redis *redis.Client, rules map[string]RateLimitRule) *RateLimiter {
	return &RateLimiter{
		Redis: redis,
		Rules: rules,
	}
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		rule, exists := r.Rules[path]
		if !exists {
			c.Next()
			return
		}

		userID := c.GetString("user_id")
		if userID == "" {
			userID = c.ClientIP()
		}

		limit := rule.Limit

		key := fmt.Sprintf("rate:%s:%s", path, userID)
		now := time.Now()

		pipe := r.Redis.TxPipeline()
		count := pipe.Incr(c, key)
		pipe.Expire(c, key, rule.Window)
		_, err := pipe.Exec(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			return
		}

		current := count.Val()
		remaining := limit - int(current)

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(rule.Window).Unix()))

		if remaining < 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}

		c.Next()
	}
}
