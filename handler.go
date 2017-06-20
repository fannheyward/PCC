package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/patrickmn/go-cache"
)

func insertAtFirst(origin []map[string]string, item map[string]string) []map[string]string {
	newList := make([]map[string]string, 0)
	if len(origin) == 0 {
		newList = append(newList, item)
	} else {
		newList = append(newList, item)
		newList = append(newList, origin...)
	}

	return newList
}

func likeHandler(c *gin.Context, oid, uid string) {
	score := cli.ZScore(fmt.Sprintf("like:%s", oid), uid).Val()
	if score > 0 {
		c.JSON(200, errInfo("already liked"))
		return
	}

	err := cli.ZAdd(fmt.Sprintf("like:%s", oid), redis.Z{Score: float64(time.Now().UnixNano()), Member: uid}).Err()
	if err != nil {
		c.JSON(200, errInfo("3:"+err.Error()))
		return
	}

	userList := make([]map[string]string, 0)
	users := cli.ZRevRangeByScore(fmt.Sprintf("like:%s", oid), redis.ZRangeBy{Min: "-inf", Max: "+inf", Count: 20}).Val()
	for _, u := range users {
		name := cli.Get(fmt.Sprintf("user:%s", u)).String()
		user := map[string]string{
			u: name,
		}

		// is oid and u friend?
		score := cli.ZScore(fmt.Sprintf("friend:%s", uid), u).Val()
		if score == 0 {
			userList = append(userList, user)
		} else {
			userList = insertAtFirst(userList, user)
		}
	}

	c.JSON(200, map[string]interface{}{
		"oid":       strToInt(oid),
		"uid":       strToInt(uid),
		"like_list": userList,
	})
}

func isLikeHandler(c *gin.Context, oid, uid string) {
	isLike := 0
	cacheKey := fmt.Sprintf("cache_like:%s:%s", oid, uid)
	_, exist := mc.Get(cacheKey)
	if exist {
		isLike = 1
	} else {
		score := cli.ZScore(fmt.Sprintf("like:%s", oid), uid).Val()
		if score > 0 {
			isLike = 1
			mc.Set(cacheKey, 1, cache.NoExpiration)
		}
	}

	c.JSON(200, map[string]interface{}{
		"oid":     strToInt(oid),
		"uid":     strToInt(uid),
		"is_like": isLike,
	})
}

func countHandler(c *gin.Context, oid, uid string) {
	count := cli.ZCount(fmt.Sprintf("like:%s", oid), "-inf", "+inf").Val()
	c.JSON(200, map[string]interface{}{
		"oid":   strToInt(oid),
		"count": count,
	})
}

func getLikedUsers(oid string, cursor, pageSize int64) []string {
	key := fmt.Sprintf("like:%s", oid)
	opt := redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Count:  pageSize,
		Offset: cursor,
	}

	return cli.ZRevRangeByScore(key, opt).Val()
}

func listHandler(c *gin.Context, oid, uid string, cursor, pageSize int64, friendOnly bool) {
	var nextCursor int64
	userList := make([]map[string]string, 0)

	if friendOnly {
		start := cursor
		userIDs := make([]string, 0)

	L:
		for int64(len(userIDs)) < pageSize {
			users := getLikedUsers(oid, start, pageSize)
			if len(users) == 0 {
				break L
			}

			for i, u := range users {
				score := cli.ZScore(fmt.Sprintf("friend:%s", uid), u).Val()
				if score == 0 {
					continue
				}

				userIDs = append(userIDs, u)
				if int64(len(userIDs)) == pageSize {
					start = start + int64(i) + 1
					break L
				}
			}
			start = start + pageSize
		}

		nextCursor = start
		for _, u := range userIDs {
			name := cli.Get(fmt.Sprintf("user:%s", u)).String()
			user := map[string]string{
				u: name,
			}

			userList = append(userList, user)
		}
	} else {
		users := getLikedUsers(oid, cursor, pageSize)
		for _, u := range users {
			name := cli.Get(fmt.Sprintf("user:%s", u)).String()
			user := map[string]string{
				u: name,
			}

			userList = append(userList, user)
		}
		nextCursor = (cursor + 1) * pageSize
	}

	c.JSON(200, map[string]interface{}{
		"oid":         strToInt(oid),
		"like_list":   userList,
		"next_cursor": nextCursor,
	})
}
