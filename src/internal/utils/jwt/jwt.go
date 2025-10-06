package utils

import (
	"app/src/internal/utils/uuid"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5" // https://golang-jwt.github.io/jwt/
)

var jwtSecret []byte

// Setup JWT Secret
func setSecret() {
	// 如果 jwtSecret 已经初始化，直接返回
	if len(jwtSecret) > 0 {
		return
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET environment variable is not set")
	}
	jwtSecret = []byte(secret)
}

// JWTClaims 定义 JWT 的声明结构
type JWTClaims struct {
	ID        string `json:"id"`
	UserID    string `json:"sub"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss,omitempty"`
}

// IssueJWT
func IssueJWT(userID string, options ...JWTOption) (string, jwt.MapClaims, error) {
	// 验证输入参数
	if userID == "" {
		return "", nil, fmt.Errorf("userID cannot be empty")
	}

	setSecret()

	// 设置默认值
	config := &JWTConfig{
		Duration: 72 * time.Hour, // 默认72小时
	}

	// 应用自定义选项
	for _, option := range options {
		option(config)
	}

	// 验证配置
	if config.Duration <= 0 {
		return "", nil, fmt.Errorf("duration must be positive")
	}

	// token 元数据
	now := time.Now()
	claims := jwt.MapClaims{
		"iat": now.Unix(),
		"exp": now.Add(config.Duration).Unix(),
		"iss": config.Issuer,
		"id":  uuid.NewString(),
		"sub": userID,
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名 token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, claims, nil
}

// JWTConfig 定义 JWT 配置
type JWTConfig struct {
	Duration time.Duration
	Issuer   string
}

// JWTOption 定义 JWT 选项函数类型
type JWTOption func(*JWTConfig)

// WithDuration 设置 token 有效期
func WithDuration(duration time.Duration) JWTOption {
	return func(config *JWTConfig) {
		config.Duration = duration
	}
}

// WithIssuer 设置发行者
func WithIssuer(issuer string) JWTOption {
	return func(config *JWTConfig) {
		config.Issuer = issuer
	}
}

type VerifyJWTPayload struct {
	UserID    string
	ExpiresAt int64
	TokenID   string
}

// VerifyJWT Veriry JWT and return certain fields
func VerifyJWT(tokenString string) (VerifyJWTPayload, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return VerifyJWTPayload{}, err
	}

	// return claims.UserID, nil
	return VerifyJWTPayload{
		UserID:    claims.UserID,
		ExpiresAt: claims.ExpiresAt,
		TokenID:   claims.ID,
	}, nil
}

// ParseJWT 解析 JWT token 并返回完整的声明信息
func ParseJWT(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token string cannot be empty")
	}

	setSecret()

	// 创建一个安全的密钥获取函数
	getSigningKey := func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// 返回密钥的副本，避免直接暴露原始密钥
		secretCopy := make([]byte, len(jwtSecret))
		copy(secretCopy, jwtSecret)
		return secretCopy, nil
	}

	// 使用安全的密钥获取函数解析 token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, getSigningKey)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 验证 token 是否有效
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		// 提取用户ID
		userID, ok := (*claims)["sub"].(string)
		if !ok || userID == "" {
			return nil, fmt.Errorf("invalid or missing user ID in token")
		}

		// 构建 JWTClaims 结构
		jwtClaims := &JWTClaims{
			UserID: userID,
		}

		// 提取可选字段
		if issuedAt, ok := (*claims)["iat"].(float64); ok {
			jwtClaims.IssuedAt = int64(issuedAt)
		}
		if expiresAt, ok := (*claims)["exp"].(float64); ok {
			jwtClaims.ExpiresAt = int64(expiresAt)
		}
		if issuer, ok := (*claims)["iss"].(string); ok {
			jwtClaims.Issuer = issuer
		}
		if id, ok := (*claims)["id"].(string); ok {
			jwtClaims.ID = id
		}

		return jwtClaims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// IsTokenExpired 检查 token 是否已过期
func IsTokenExpired(tokenString string) (bool, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return true, err
	}

	now := time.Now().Unix()
	return claims.ExpiresAt < now, nil
}

// GetTokenExpiration 获取 token 的过期时间
func GetTokenExpiration(tokenString string) (time.Time, error) {
	claims, err := ParseJWT(tokenString)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(claims.ExpiresAt, 0), nil
}
