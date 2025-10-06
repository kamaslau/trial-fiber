package utils

import (
	"log"
	"time"
)

// 示例：在路由处理函数中使用 JWT
func ExampleUseJWTInHandler() {
	// 由于 init() 函数已经验证了 JWT_SECRET
	// 我们可以安全地使用 JWT 功能

	// 基本用法
	token, _, err := IssueJWT("user-123")
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return
	}

	log.Printf("Generated token: %s", token)

	// 验证 token
	userID, err := VerifyJWT(token)
	if err != nil {
		log.Printf("Failed to verify JWT: %v", err)
		return
	}

	log.Printf("Verified user: %v", userID)
}

// 示例：使用自定义选项生成 JWT
func ExampleIssueJWTWithOptions() {
	// 生成带有自定义选项的 JWT
	token, _, err := IssueJWT("admin-user",
		WithDuration(24*time.Hour), // 24小时有效期
		WithIssuer("my-app"),       // 设置发行者
	)
	if err != nil {
		log.Printf("Failed to generate JWT with options: %v", err)
		return
	}

	log.Printf("Generated token with options: %s", token)

	// 解析 token 获取完整信息
	claims, err := ParseJWT(token)
	if err != nil {
		log.Printf("Failed to parse JWT: %v", err)
		return
	}

	log.Printf("Token details:")
	log.Printf("  User ID: %s", claims.UserID)
	log.Printf("  Issuer: %s", claims.Issuer)
	log.Printf("  Issued At: %v", time.Unix(claims.IssuedAt, 0))
	log.Printf("  Expires At: %v", time.Unix(claims.ExpiresAt, 0))
}

// 示例：检查 token 过期状态
func ExampleCheckTokenExpiration() {
	// 生成一个短期 token
	token, _, err := IssueJWT("test-user", WithDuration(1*time.Hour))
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		return
	}

	// 检查是否过期
	expired, err := IsTokenExpired(token)
	if err != nil {
		log.Printf("Failed to check expiration: %v", err)
		return
	}

	if expired {
		log.Println("Token is expired")
	} else {
		log.Println("Token is still valid")
	}

	// 获取过期时间
	expiration, err := GetTokenExpiration(token)
	if err != nil {
		log.Printf("Failed to get expiration time: %v", err)
		return
	}

	log.Printf("Token expires at: %v", expiration)
	log.Printf("Time until expiration: %v", expiration.Sub(time.Now()))
}

// 示例：错误处理
func ExampleErrorHandling() {
	// 测试空 userID
	_, _, err := IssueJWT("")
	if err != nil {
		log.Printf("Expected error for empty userID: %v", err)
	}

	// 测试无效持续时间
	_, _, err = IssueJWT("user", WithDuration(-1*time.Hour))
	if err != nil {
		log.Printf("Expected error for negative duration: %v", err)
	}

	// 测试无效 token
	_, err = VerifyJWT("invalid.token.here")
	if err != nil {
		log.Printf("Expected error for invalid token: %v", err)
	}

	// 测试空 token
	_, err = VerifyJWT("")
	if err != nil {
		log.Printf("Expected error for empty token: %v", err)
	}
}

// 示例：setSecret 函数的重复调用行为
func ExampleSetSecretMultipleCalls() {
	// 第一次调用 setSecret - 会初始化 jwtSecret
	setSecret()
	log.Println("First call to setSecret: jwtSecret initialized")

	// 第二次调用 setSecret - 由于 jwtSecret 已经初始化，什么都不做
	setSecret()
	log.Println("Second call to setSecret: nothing happens (already initialized)")

	// 第三次调用 setSecret - 同样什么都不做
	setSecret()
	log.Println("Third call to setSecret: nothing happens (already initialized)")

	// 这种设计确保了：
	// 1. jwtSecret 只会被初始化一次
	// 2. 重复调用不会产生副作用
	// 3. 提高了性能，避免了不必要的环境变量读取
}
