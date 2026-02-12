package middleware

import (
	"context"
	"strconv"
	"time"

	"reelcut/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var limitsByTier = map[string]int{
	"free":       2000,  // per hour per IP; 100 was too low for normal dev/frontend refetches
	"pro":        1000,
	"enterprise": 10000,
}

func RateLimit(rdb *redis.Client, getTier func(c *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := "ratelimit:" + c.ClientIP()
		tier := "free"
		if getTier != nil {
			tier = getTier(c)
		}
		limit := limitsByTier[tier]
		if limit <= 0 {
			limit = 2000
		}
		window := time.Hour
		ctx := context.Background()
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.Next()
			return
		}
		if count == 1 {
			rdb.Expire(ctx, key, window)
		}
		ttl, _ := rdb.TTL(ctx, key).Result()
		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		remaining := limit - int(count)
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
		if int(count) > limit {
			utils.Error(c, 429, "RATE_LIMIT_EXCEEDED", "Too many requests", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}

func TierFromUser(c *gin.Context) string {
	u := GetUser(c)
	if u != nil {
		return u.SubscriptionTier
	}
	return "free"
}
