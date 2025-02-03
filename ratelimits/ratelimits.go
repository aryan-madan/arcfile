package ratelimits

import (
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/ulule/limiter/v3"
	mGin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func SetupLimiters() map[string]gin.HandlerFunc {
	store := memory.NewStore()

	limits := map[string]limiter.Rate{
		"getFile":         {Period: 5 * time.Minute, Limit: 120}, // 120 requests per 5 minutes
		"postFile":        {Period: 10 * time.Minute, Limit: 20}, // 20 requests per 10 minutes
		"deleteFile":      {Period: 10 * time.Minute, Limit: 20}, // 20 requests per 10 minutes
		"getDownloadFile": {Period: 1 * time.Hour, Limit: 50},    // 50 requests per hour
	}

	return map[string]gin.HandlerFunc{
		"getFile":         mGin.NewMiddleware(limiter.New(store, limits["getFile"])),
		"postFile":        mGin.NewMiddleware(limiter.New(store, limits["postFile"])),
		"deleteFile":      mGin.NewMiddleware(limiter.New(store, limits["deleteFile"])),
		"getDownloadFile": mGin.NewMiddleware(limiter.New(store, limits["getDownloadFile"])),
	}
}
