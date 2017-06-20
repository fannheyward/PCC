package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const nameLength = 4

func getFakeName() string {
	b := make([]byte, nameLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func genKV(k string, count int64) error {
	for i := int64(0); i < count; i++ {
		key := fmt.Sprintf("%s:%d", k, i)
		value := fmt.Sprintf("%s-%d-%s", k, i, getFakeName())
		err := cli.Set(key, value, 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func genFriend(count int64) error {
	top := int(float64(count) * 0.02)
	mid := int(float64(count) * 0.05)
	bottom := int(float64(count) * 0.1)

	topUsers := map[string]bool{}
	midUsers := map[string]bool{}
	bottomUsers := map[string]bool{}

	rand.Seed(time.Now().Unix())
	for {
		user := fmt.Sprintf("friend:%d", rand.Int63n(count))
		if topUsers[user] && midUsers[user] && bottomUsers[user] {
			continue
		}
		topUsers[user] = true
		if len(topUsers) == top {
			break
		}
	}
	for {
		user := fmt.Sprintf("friend:%d", rand.Int63n(count))
		if topUsers[user] && midUsers[user] && bottomUsers[user] {
			continue
		}
		midUsers[user] = true

		if len(midUsers) == mid {
			break
		}
	}
	for {
		user := fmt.Sprintf("friend:%d", rand.Int63n(count))
		if topUsers[user] && midUsers[user] && bottomUsers[user] {
			continue
		}
		bottomUsers[user] = true

		if len(bottomUsers) == bottom {
			break
		}
	}

	for user := range topUsers {
		for {
			target := fmt.Sprintf("%d", rand.Int63n(count))
			if fmt.Sprintf("friend:%s", target) == user {
				continue
			}

			// make user friend with target
			err := cli.ZAdd(user, redis.Z{Score: float64(time.Now().UnixNano()), Member: target}).Err()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			length, err := cli.ZCount(user, "-inf", "+inf").Result()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if length >= 50 {
				break
			}
		}
	}

	for user := range midUsers {
		for {
			target := fmt.Sprintf("%d", rand.Int63n(count))
			if fmt.Sprintf("friend:%s", target) == user {
				continue
			}

			// make user friend with target
			err := cli.ZAdd(user, redis.Z{Score: float64(time.Now().UnixNano()), Member: target}).Err()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			length, err := cli.ZCount(user, "-inf", "+inf").Result()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if length >= 20 {
				break
			}
		}
	}

	for user := range bottomUsers {
		for {
			target := fmt.Sprintf("%d", rand.Int63n(count))
			if fmt.Sprintf("friend:%s", target) == user {
				continue
			}

			// make user friend with target
			err := cli.ZAdd(user, redis.Z{Score: float64(time.Now().UnixNano()), Member: target}).Err()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			length, err := cli.ZCount(user, "-inf", "+inf").Result()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			if length >= 10 {
				break
			}
		}
	}

	return nil
}

func testDataGenHandler(c *gin.Context) {
	userCount := c.Query("user")
	objCount := c.Query("object")

	key := fmt.Sprintf("user:%d", strToInt(userCount)-1)
	value := cli.Get(key).Val()
	if value != "" {
		c.JSON(200, map[string]interface{}{
			"status": true,
			"info":   "already init test data",
		})
		return
	}

	err := genKV("user", strToInt(userCount))
	if err != nil {
		c.JSON(200, errInfo("user"+err.Error()))
		return
	}

	err = genKV("obj", strToInt(objCount))
	if err != nil {
		c.JSON(200, errInfo("object"+err.Error()))
		return
	}

	err = genFriend(strToInt(userCount))
	if err != nil {
		c.JSON(200, errInfo("genFriend"+err.Error()))
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": true,
	})
}
