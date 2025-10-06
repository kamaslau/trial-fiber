package utils

import (
	"os"
	"testing"
	"time"
)

func TestJWTGenerateAndVerify(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-jwt-tests")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 测试生成 JWT（使用默认配置）
	userID := "test-user-123"
	token, _, err := IssueJWT(userID)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	t.Logf("Generated token: %s", token)

	// 测试验证 JWT
	verified, err := VerifyJWT(token)
	if err != nil {
		t.Fatalf("Failed to verify JWT: %v", err)
	}

	if verified.UserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, verified.UserID)
	}

	t.Logf("Successfully verified token for user: %s", verified.UserID)
}

func TestJWTWithOptions(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-options-tests")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 测试使用自定义选项生成 JWT
	userID := "admin-user"
	token, _, err := IssueJWT(userID,
		WithDuration(1*time.Hour),
		WithIssuer("test-issuer"),
	)
	if err != nil {
		t.Fatalf("Failed to generate JWT with options: %v", err)
	}

	// 解析 token 并验证所有字段
	claims, err := ParseJWT(token)
	if err != nil {
		t.Fatalf("Failed to parse JWT: %v", err)
	}

	// 验证字段
	if claims.UserID != userID {
		t.Errorf("Expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Issuer != "test-issuer" {
		t.Errorf("Expected issuer 'test-issuer', got %s", claims.Issuer)
	}

	t.Logf("Successfully generated and parsed JWT with custom options")
}

func TestJWTValidationErrors(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-error-tests")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 测试空 userID
	_, _, err := IssueJWT("")
	if err == nil {
		t.Error("Expected error for empty userID, but got none")
	}
	t.Logf("Correctly rejected empty userID: %v", err)

	// 测试无效 token
	invalidToken := "invalid.token.here"
	_, err = VerifyJWT(invalidToken)
	if err == nil {
		t.Error("Expected error for invalid token, but got none")
	}
	t.Logf("Correctly rejected invalid token: %v", err)

	// 测试空 token
	_, err = VerifyJWT("")
	if err == nil {
		t.Error("Expected error for empty token, but got none")
	}
	t.Logf("Correctly rejected empty token: %v", err)
}

func TestJWTExpiration(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-expiration-tests")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 生成一个短期 token（1秒）
	token, _, err := IssueJWT("test-user", WithDuration(1*time.Second))
	if err != nil {
		t.Fatalf("Failed to generate short-lived JWT: %v", err)
	}

	// 立即检查过期状态
	expired, err := IsTokenExpired(token)
	if err != nil {
		t.Fatalf("Failed to check token expiration: %v", err)
	}
	if expired {
		t.Error("Token should not be expired immediately after generation")
	}

	// 获取过期时间
	expiration, err := GetTokenExpiration(token)
	if err != nil {
		t.Fatalf("Failed to get token expiration: %v", err)
	}

	now := time.Now()
	if expiration.Before(now) {
		t.Error("Token expiration should be in the future")
	}

	t.Logf("Token expires at: %v", expiration)
	t.Logf("Current time: %v", now)
	t.Logf("Time until expiration: %v", expiration.Sub(now))
}

func TestJWTInvalidDuration(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-duration-tests")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 测试负持续时间
	_, _, err := IssueJWT("test-user", WithDuration(-1*time.Hour))
	if err == nil {
		t.Error("Expected error for negative duration, but got none")
	}
	t.Logf("Correctly rejected negative duration: %v", err)

	// 测试零持续时间
	_, _, err = IssueJWT("test-user", WithDuration(0))
	if err == nil {
		t.Error("Expected error for zero duration, but got none")
	}
	t.Logf("Correctly rejected zero duration: %v", err)
}

func TestSetSecretMultipleCalls(t *testing.T) {
	// 设置测试环境变量
	originalSecret := os.Getenv("JWT_SECRET")
	os.Setenv("JWT_SECRET", "test-secret-key-for-multiple-calls")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
	}()

	// 第一次调用 setSecret
	setSecret()

	// 验证 jwtSecret 已经被设置
	if len(jwtSecret) == 0 {
		t.Fatal("jwtSecret should be set after first call to setSecret")
	}

	originalSecretLength := len(jwtSecret)

	// 第二次调用 setSecret - 应该什么都不做
	setSecret()

	// 验证 jwtSecret 没有被重新设置
	if len(jwtSecret) != originalSecretLength {
		t.Error("jwtSecret should not be modified on subsequent calls to setSecret")
	}

	t.Log("setSecret correctly handles multiple calls without re-initialization")
}
