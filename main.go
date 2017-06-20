package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

var (
	cli *redis.Client
	mc  *cache.Cache
)

func init() {
	cli = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6079",
		PoolSize: 100,
	})
	err := cli.Ping().Err()
	if err != nil {
		panic("redis failed")
	}

	mc = cache.New(cache.NoExpiration, time.Minute)
}

func handler(c *gin.Context) {
	action := c.Query("action")
	oid := c.Query("oid")
	uid := c.Query("uid")

	if action == "" || oid == "" || uid == "" {
		c.JSON(503, errInfo("action || oid || uid is missing"))
		return
	}

	switch action {
	case "like":
		likeHandler(c, oid, uid)
	case "is_like":
		isLikeHandler(c, oid, uid)
	case "count":
		countHandler(c, oid, uid)
	case "list":
		cursor := c.Query("cursor")
		pageSize := c.Query("page_size")
		isFriend := c.Query("is_friend")
		if cursor == "" || pageSize == "" || isFriend == "" {
			c.JSON(503, errInfo("cursor || page_size || is_friend is missing"))
			return
		}

		friendOnly := false
		if isFriend == "1" {
			friendOnly = true
		}

		listHandler(c, oid, uid, strToInt(cursor), strToInt(pageSize), friendOnly)
	}
}

func main() {
	debug := false
	gin.SetMode(gin.ReleaseMode)

	mux := gin.New()
	if debug {
		gin.SetMode(gin.DebugMode)
		mux.Use(gin.Logger())
	}

	mux.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// /pcc?action=like|is_like|count|list&oid=xxx&uid=xxx
	mux.GET("/pcc", handler)
	mux.POST("/test/gen", testDataGenHandler)
	mux.Run("127.0.0.1:9009")
}
