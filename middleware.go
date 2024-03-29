package toolib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"io/ioutil"
	"net/http"
	"time"
)

type MiddlewareRespHandle func(*gin.Context, string, error)

func MiddlewareCacheByRedis(red *redis.Client, isCookie bool, dataExpiration, lockExpiration, updateExpiration time.Duration, respHandle MiddlewareRespHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("url path:", c.Request.URL.Path)
		if red == nil {
			return
		}
		key := getCacheKeyByGet(c, isCookie)
		if c.Request.Method == http.MethodPost {
			key = getCacheKeyByPost(c, isCookie)
		}
		cacheHandle := func() (string, error) {
			blw := &bodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
			c.Next()
			statusCode := c.Writer.Status()
			// no cache failed req
			if statusCode != http.StatusOK {
				return "", fmt.Errorf("status code [%d]", statusCode)
			}
			if blw.body.String() == "" {
				return "", fmt.Errorf("body is nil")
			}
			return blw.body.String(), nil
		}
		res, err := CacheByRedis(red, key, dataExpiration, lockExpiration, updateExpiration, cacheHandle)
		if respHandle != nil {
			respHandle(c, res, err)
		} else {
			if err != nil {
				fmt.Println("CacheByRedis err:", err.Error())
				c.AbortWithStatusJSON(http.StatusOK, err.Error())
			} else if res != "" {
				var respMap map[string]interface{}
				_ = json.Unmarshal([]byte(res), &respMap)
				c.AbortWithStatusJSON(http.StatusOK, respMap)
			}
		}
	}
}

func CacheByRedis(red *redis.Client, key string, dataExpiration, lockExpiration, updateExpiration time.Duration, cacheHandle func() (string, error)) (string, error) {
	updateExpirationKey := fmt.Sprintf("uek:%s", key)
	lockExpirationKey := fmt.Sprintf("lek:%s", key)
	// check cache is exist
	if dataStr, err := red.Get(key).Result(); err == nil { // check update time is expired
		if exi, err := red.Exists(updateExpirationKey).Result(); err != nil {
			return "", err
		} else if exi == 0 { // check lock
			fmt.Println("CacheByRedis is expired:", key)
			if ok, err := red.SetNX(lockExpirationKey, "", lockExpiration).Result(); err != nil {
				fmt.Println("SetNX 1:", ok, err.Error())
				return dataStr, nil
			} else if !ok {
				fmt.Println("SetNX 2:", ok)
				return dataStr, nil
			} else {
				if dataStr, err = cacheHandle(); err != nil {
					return "", err
				} else if err = red.Set(key, dataStr, dataExpiration).Err(); err != nil {
					return "", err
				} else {
					_ = red.Set(updateExpirationKey, "", updateExpiration).Err()
					_ = red.Expire(lockExpirationKey, time.Second*5).Err()
					return "", nil
				}
			}
		} else { // not expired
			fmt.Println("CacheByRedis OK:", key)
			return dataStr, nil
		}
	} else if err == redis.Nil { // not exist
		fmt.Println("CacheByRedis is nil:", key)
		if dataStr, err = cacheHandle(); err != nil {
			return "", err
		} else if err = red.Set(key, dataStr, dataExpiration).Err(); err != nil {
			return "", err
		} else {
			_ = red.Set(updateExpirationKey, "", updateExpiration).Err()
			return "", nil
		}
	} else {
		return "", err
	}
}

func getCacheKeyByGet(c *gin.Context, isCookie bool) string {
	if isCookie {
		cook, _ := json.Marshal(c.Request.Cookies()) //加入cookie的部分
		urlBytes := append([]byte(c.Request.URL.String()), cook...)
		return Md5Hash(urlBytes)
	}
	return Md5Hash([]byte(c.Request.URL.String()))
}

func getCacheKeyByPost(c *gin.Context, isCookie bool) string {
	if isCookie {
		bodyBytes, _ := c.GetRawData()
		cook, _ := json.Marshal(c.Request.Cookies())
		urlBytes := append([]byte(c.Request.URL.String()), cook...)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
		return Md5Hash(append(urlBytes, bodyBytes...))
	}
	bodyBytes, _ := c.GetRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
	return Md5Hash(append([]byte(c.Request.URL.String()), bodyBytes...))

}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (b bodyWriter) Write(bys []byte) (int, error) {
	b.body.Write(bys)
	return b.ResponseWriter.Write(bys)
}

// ============================ JWT

func MiddlewareJwtCheck(JwtAuthorization, JwtKey string, claimsCheck JwtClaimsCheck, respHandle MiddlewareRespHandle) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie(JwtAuthorization)
		if err != nil {
			respHandle(ctx, "", err)
			return
		}
		token := cookie.Value
		if claims, err := JwtVerify(token, JwtKey); err != nil {
			respHandle(ctx, "", err)
			return
		} else if claimsCheck != nil {
			if err = claimsCheck(ctx, claims, token); err != nil {
				respHandle(ctx, "", err)
				return
			}
		}
		ctx.Next()
	}
}

type JwtClaimsCheck func(*gin.Context, jwt.MapClaims, string) error
