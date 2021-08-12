package toolib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"io/ioutil"
	"net/http"
	"time"
)

func MiddlewareCacheByRedis(red *redis.Client, dataExpiration, lockExpiration, updateExpiration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := getCacheKeyByGet(c)
		if c.Request.Method == http.MethodPost {
			key = getCacheKeyByPost(c)
		}
		cacheHandle := func() (string, error) {
			blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
			c.Next()
			statusCode := c.Writer.Status()
			// 不缓存失败的请求
			if statusCode != http.StatusOK {
				return "", fmt.Errorf("status code not ok")
			}
			if blw.body.String() == "" {
				return "", fmt.Errorf("body is nil")
			}
			return blw.body.String(), nil
		}
		res, err := CacheByRedis(red, key, dataExpiration, lockExpiration, updateExpiration, cacheHandle)
		if err != nil {
			fmt.Println("CacheByRedis err:", err.Error())
			c.AbortWithStatusJSON(http.StatusOK, err.Error())
		} else {
			var respMap map[string]interface{}
			_ = json.Unmarshal([]byte(res), &respMap)
			c.AbortWithStatusJSON(http.StatusOK, respMap)
		}
	}
}

func CacheByRedis(red *redis.Client, key string, dataExpiration, lockExpiration, updateExpiration time.Duration, cacheHandle func() (string, error)) (string, error) {
	updateExpirationKey := fmt.Sprintf("uek:%s", key)
	lockExpirationKey := fmt.Sprintf("lek:%s", key)
	// 查询缓存是否存在
	if dataStr, err := red.Get(key).Result(); err == nil { // 存在，判断更新时间是否过期
		if exi, err := red.Exists(updateExpirationKey).Result(); err != nil {
			return "", err
		} else if exi == 0 { //过期判断当前分布式锁是否被占用
			if ok, err := red.SetNX(lockExpirationKey, "", lockExpiration).Result(); err != nil {
				return dataStr, nil
			} else if !ok {
				return dataStr, nil
			} else {
				if dataStr, err = cacheHandle(); err != nil {
					return "", err
				} else if err = red.Set(key, dataStr, dataExpiration).Err(); err != nil {
					return "", err
				} else {
					_ = red.Set(updateExpirationKey, "", updateExpiration).Err()
					_ = red.Expire(lockExpirationKey, time.Second*5).Err()
					return dataStr, nil
				}
			}
		} else { //没过期返回数据
			fmt.Println("cacheByRedis OK:", key)
			return dataStr, nil
		}
	} else if err == redis.Nil { // 不存在查询数据库，写缓存
		if dataStr, err = cacheHandle(); err != nil {
			return "", err
		} else if err = red.Set(key, dataStr, dataExpiration).Err(); err != nil {
			return "", err
		} else {
			_ = red.Set(updateExpirationKey, "", updateExpiration).Err()
			return dataStr, nil
		}
	} else {
		return "", err
	}
}

func getCacheKeyByGet(c *gin.Context) string {
	cook, _ := json.Marshal(c.Request.Cookies()) //加入cookie的部分
	urlBytes := append([]byte(c.Request.URL.String()), cook...)
	return Md5Hash(urlBytes)
}

func getCacheKeyByPost(c *gin.Context) string {
	bodyBytes, _ := c.GetRawData()
	cook, _ := json.Marshal(c.Request.Cookies())
	urlBytes := append([]byte(c.Request.URL.String()), cook...)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
	return Md5Hash(append(urlBytes, bodyBytes...))
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyWriter) Write(bys []byte) (int, error) {
	return b.body.Write(bys)
}
