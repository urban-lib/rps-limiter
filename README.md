```bash
go get github.com/urban-lib/rps-limiter
```

```go

func MaxAllowedMiddleware(cfg *limiter.Config) gin.HandlerFunc {
	l := limiter.NewRateLimiter(cfg)
	go l.CleanupVisitors()

	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			logger.Errorf(err.Error())
			NewResponse(c, http.StatusInternalServerError, "Server error")
			//c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if !l.GetVisitor(ip).Allow() {
			NewResponse(c, http.StatusTooManyRequests, "To many requests")
			//c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
```