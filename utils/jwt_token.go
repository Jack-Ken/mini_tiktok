package utils

import (
	"errors"
	"mini_tiktok/internal/dao/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 24

var CustomSecret = []byte("pbkdf2-sha512$%s$%s")

type CustomClaims struct {
	// 可根据需要自行添加字段
	UserName             string `json:"username"`
	UserID               int64  `json:"user_id"`
	jwt.RegisteredClaims        // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(username string, user_id int64) (string, error) {
	// 创建一个我们自己的声明
	claims := CustomClaims{
		username, // 自定义字段
		user_id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "potatoes", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(CustomSecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析token
	//var mc = new(CustomClaims)
	// 如果是自定义Claim结构体则需要使用 ParseWithClaims 方法
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		// 直接使用标准的Claim则可以直接使用Parse方法
		//token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		return CustomSecret, nil
	})
	if err != nil {
		return nil, err
	}
	// 对token对象中的Claim进行类型断言
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		// 获取参数和参数校验
		tokenStr := c.Query("token")
		if tokenStr == "" {
			tokenStr = c.PostForm("token")
		}
		//用户不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, models.Response{StatusCode: 401, StatusMsg: "用户不存在"})
			c.Abort() //阻止执行
			return
		}

		mc, err := ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, models.RegisterResponse{
				Response: models.Response{
					StatusCode: -1,
					StatusMsg:  "无效的Token/Token验证错误/Token过期",
				},
			})
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.UserName)
		c.Set("userid", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
//func JWTAuthMiddleware() func(c *gin.Context) {
//	return func(c *gin.Context) {
//		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
//		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
//		// 这里的具体实现方式要依据你的实际业务情况决定
//		authHeader := c.Request.Header.Get("Authorization")
//		if authHeader == "" {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2003,
//				"msg":  "请求头中auth为空",
//			})
//			c.Abort()
//			return
//		}
//		// 按空格分割
//		parts := strings.SplitN(authHeader, " ", 2)
//		if !(len(parts) == 2 && parts[0] == "Bearer") {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2004,
//				"msg":  "请求头中auth格式有误",
//			})
//			c.Abort()
//			return
//		}
//		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
//		mc, err := ParseToken(parts[1])
//		if err != nil {
//			c.JSON(http.StatusOK, gin.H{
//				"code": 2005,
//				"msg":  "无效的Token",
//			})
//			c.Abort()
//			return
//		}
//		// 将当前请求的username信息保存到请求的上下文c上
//		c.Set("username", mc.UserName)
//		c.Set("userid", mc.UserID)
//		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
//	}
//}
